package authz

import (
	"context"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
	"testing"
)

func Test_AuthConstraintsForGRPC(t *testing.T) {
	// GIVEN auth-authAdminService and organization
	ctx := context.TODO()
	authAdminService, cfg, err := newAuthService()
	require.NoError(t, err)

	org, err := domain.NewOrganizationBuilder().
		WithId("test-org-"+uuid.NewV4().String()).
		WithName("org-name").
		WithUrl("org-url").
		WithNamespaces("finance", "loan").Build()
	org, err = authAdminService.CreateOrganization(ctx, org)
	require.NoError(t, err)

	// AND with following principals
	tom, err := domain.NewPrincipalBuilder().
		WithOrganizationId(org.Id).
		WithNamespaces(org.Namespaces...).
		WithAttribute("Region", "Midwest").
		WithName("Tom").
		WithUsername("tom").Build()
	require.NoError(t, err) //
	require.NotNil(t, tom)

	_, err = NewGrpcAuth(cfg)
	require.Error(t, err) //

	//_, err = authorizer.Check(ctx, &services.CheckConstraintsRequest{
	//	OrganizationId: org.Id,
	//	Namespace:      org.Namespaces[0],
	//	PrincipalId:    tom.Id,
	//	Constraints:    `.CurrentLocation "Chicago"`,
	//	Context:        map[string]string{"CurrentLocation": "Chicago"},
	//})
	//require.Error(t, err) // constraints not implemented
}

func Test_PermissionsForDepositAccount(t *testing.T) {
	ctx := context.TODO()
	_, cfg, err := newAuthService()
	require.NoError(t, err)

	_ = Subject(ctx)

	_, err = NewGrpcAuth(cfg)
	require.Error(t, err)

	//req := &services.AuthRequest{
	//	PrincipalId: "root",
	//	Action:      "query",
	//	Resource:    "*",
	//}
	//_, err = authorizer.Authorize(ctx, req)
	//require.NoError(t, err)
}
