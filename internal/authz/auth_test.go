package authz

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/metrics"
	"github.com/bhatti/PlexAuthZ/internal/service"
	"github.com/bhatti/PlexAuthZ/internal/service/db"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
	"testing"
)

func Test_AuthConstraintsForDefaultAuthorizer(t *testing.T) {
	// GIVEN auth-authService and organization
	ctx := context.TODO()
	authService, cfg, err := newAuthService()
	require.NoError(t, err)

	org, err := domain.NewOrganizationBuilder().
		WithId("test-org-"+uuid.NewV4().String()).
		WithName("org-name").
		WithUrl("org-url").
		WithNamespaces("finance", "loan").Build()
	org, err = authService.CreateOrganization(ctx, org)
	require.NoError(t, err)

	// AND with following principals
	tom, err := domain.NewPrincipalBuilder().
		WithOrganizationId(org.Id).
		WithNamespaces(org.Namespaces...).
		WithAttribute("Region", "Midwest").
		WithName("Tom").
		WithUsername("tom").Build()
	require.NoError(t, err)
	tom, err = authService.CreatePrincipal(ctx, tom)
	require.NoError(t, err)

	// WHEN creating a resource
	depositAccount, err := domain.NewResourceBuilder().
		WithNamespace(org.Namespaces[0]).
		WithName("DepositAccount").
		WithAttribute("AccountType", "Checking").
		WithAllowedActions("balance", "withdraw", "deposit", "open", "close").Build()
	require.NoError(t, err)
	depositAccount, err = authService.CreateResource(ctx, org.Id, depositAccount)
	// THEN we should be able to save resource in the database
	require.NoError(t, err)

	// AND with roles for:
	employee := assertCreateRole(t, authService, org, "Employee")
	teller := assertCreateRole(t, authService, org, "Teller", employee.Id)

	namespace := org.Namespaces[0]

	// WHEN assigning creating roles
	require.NoError(t, authService.AddRolesToPrincipal(ctx, org.Id, namespace, tom.Id, teller.Id))

	authorizer, err := CreateAuthorizer(DefaultAuthorizerKind, cfg, authService)
	require.NoError(t, err)

	res, err := authorizer.Check(ctx, &services.CheckConstraintsRequest{
		OrganizationId: org.Id,
		Namespace:      namespace,
		PrincipalId:    tom.Id,
	})
	require.Error(t, err) // check without constraints should fail

	res, err = authorizer.Check(ctx, &services.CheckConstraintsRequest{
		OrganizationId: org.Id,
		Namespace:      namespace,
		PrincipalId:    "blah",
		Constraints:    `eq .CurrentLocation "Seattle"`,
		Context:        map[string]string{"CurrentLocation": "Seattle"},
	})
	require.Error(t, err) // bad principal id should fail

	res, err = authorizer.Check(ctx, &services.CheckConstraintsRequest{
		OrganizationId: org.Id,
		Namespace:      namespace,
		PrincipalId:    tom.Id,
		Constraints:    `eq .CurrentLocation "Seattle"`,
		Context:        map[string]string{"CurrentLocation": "Seattle"},
	})
	require.NoError(t, err)
	require.True(t, res.Matched)

	res, err = authorizer.Check(ctx, &services.CheckConstraintsRequest{
		OrganizationId: org.Id,
		Namespace:      namespace,
		PrincipalId:    tom.Id,
		Constraints:    `eq .CurrentLocation "Chicago"`,
		Context:        map[string]string{"CurrentLocation": "Seattle"},
	})
	require.Error(t, err) // constraints should not match
}

func Test_PermissionsForDepositAccountForDefaultAuthorizer(t *testing.T) {
	// GIVEN auth-authService and organization
	ctx := context.TODO()
	authService, cfg, err := newAuthService()
	require.NoError(t, err)

	org, err := domain.NewOrganizationBuilder().
		WithId("test-org-"+uuid.NewV4().String()).
		WithName("org-name").
		WithUrl("org-url").
		WithNamespaces("finance", "loan").Build()
	org, err = authService.CreateOrganization(ctx, org)
	require.NoError(t, err)

	// AND with following principals
	tom, err := domain.NewPrincipalBuilder().
		WithOrganizationId(org.Id).
		WithNamespaces(org.Namespaces...).
		WithAttribute("Region", "Midwest").
		WithName("Tom").
		WithUsername("tom").Build()
	require.NoError(t, err)
	tom, err = authService.CreatePrincipal(ctx, tom)
	require.NoError(t, err)

	// WHEN creating a resource
	depositAccount, err := domain.NewResourceBuilder().
		WithNamespace(org.Namespaces[0]).
		WithName("DepositAccount").
		WithAttribute("AccountType", "Checking").
		WithAllowedActions("balance", "withdraw", "deposit", "open", "close").Build()
	require.NoError(t, err)
	depositAccount, err = authService.CreateResource(ctx, org.Id, depositAccount)
	// THEN we should be able to save resource in the database
	require.NoError(t, err)

	// AND with roles for:
	employee := assertCreateRole(t, authService, org, "Employee")
	teller := assertCreateRole(t, authService, org, "Teller", employee.Id)

	// AND with following permissions
	balancePerm := assertCreatePermission(t, depositAccount, authService, org,
		"", "balance")
	depositPerm := assertCreatePermission(t, depositAccount, authService, org,
		`and (eq .Principal.Region "Midwest") (eq .CurrentLocation "Chicago")`,
		"deposit", "withdraw")

	namespace := org.Namespaces[0]

	// assigning permission to roles
	require.NoError(t, authService.AddPermissionsToRole(ctx, org.Id, namespace, employee.Id, balancePerm.Id))
	require.NoError(t, authService.AddPermissionsToRole(ctx, org.Id, namespace, teller.Id, depositPerm.Id))

	// WHEN assigning creating roles
	require.NoError(t, authService.AddRolesToPrincipal(ctx, org.Id, namespace, tom.Id, teller.Id))

	authorizer, err := CreateAuthorizer(DefaultAuthorizerKind, cfg, authService)
	require.NoError(t, err)

	// Test for DepositAccount
	for i, action := range []string{"balance", "deposit", "withdraw"} {
		// WHEN checking for action and permission
		req := &services.AuthRequest{
			OrganizationId: org.Id,
			Namespace:      namespace,
			PrincipalId:    tom.Id,
			Action:         action,
			Resource:       depositAccount.Name,
			Context:        map[string]string{"CurrentLocation": "Seattle"},
		}
		res, err := authorizer.Authorize(ctx, req)
		if i > 0 {
			// Should fail without Location equal to Chicago
			require.Error(t, err)
			req.Context["CurrentLocation"] = "Chicago"
			res, err = authorizer.Authorize(ctx, req)
		}
		// BUT it should succeed with PERMITTED when using Chicago
		require.NoError(t, err, fmt.Sprintf("i %d, action %s, principal %v", i, action, tom))
		require.Equal(t, types.Effect_PERMITTED, res.Effect)
	}
	for _, action := range []string{"close"} {
		// WHEN checking for action and permission
		req := &services.AuthRequest{
			OrganizationId: org.Id,
			Namespace:      namespace,
			PrincipalId:    tom.Id,
			Action:         action,
			Resource:       depositAccount.Name,
			Context:        map[string]string{"CurrentLocation": "Chicago"},
		}
		_, err := authorizer.Authorize(ctx, req)
		// tom should not be able to invoke `close` action
		_, err = authorizer.Authorize(ctx, req)
		require.Error(t, err)
	}
}

func assertCreateRole(
	t *testing.T,
	authAdminService service.AuthAdminService,
	org *types.Organization,
	name string,
	parentIds ...string) *types.Role {
	ctx := context.TODO()
	role, err := domain.NewRoleBuilder().
		WithNamespace(org.Namespaces[0]).
		WithName(name).
		WithParentIds(parentIds...).
		Build()
	require.NoError(t, err)
	role, err = authAdminService.CreateRole(ctx, org.Id, role)
	require.NoError(t, err)
	return role
}

func assertCreatePermission(
	t *testing.T,
	resource *types.Resource,
	authAdminService service.AuthAdminService,
	org *types.Organization,
	constraints string,
	actions ...string) *types.Permission {
	ctx := context.TODO()
	// WHEN creating a permission for the resource
	perm, err := domain.NewPermissionBuilder().
		WithNamespace(org.Namespaces[0]).
		WithActions(actions...).
		WithResourceId(resource.Id).
		WithEffect(types.Effect_PERMITTED).
		WithScope("").
		WithConstraints(constraints).Build()
	require.NoError(t, err)
	perm, err = authAdminService.CreatePermission(ctx, org.Id, perm)
	// THEN we should be able to save permission in the database
	require.NoError(t, err)
	return perm
}

func newAuthService() (service.AuthAdminService, *domain.Config, error) {
	cfg, err := domain.NewConfig("")
	if err != nil {
		return nil, nil, err
	}
	authService, _, err := db.CreateDatabaseAuthService(cfg, metrics.New())
	if err != nil {
		return nil, nil, err
	}
	return authService, cfg, err
}
