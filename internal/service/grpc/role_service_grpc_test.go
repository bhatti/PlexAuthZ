package grpc

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/service"
	"github.com/stretchr/testify/require"
	"testing"
)

func testCRUDRoles(
	ctx context.Context,
	t *testing.T,
	authService service.AuthAdminService,
	org *types.Organization,
) {
	role, err := domain.NewRoleBuilder().
		WithNamespace(org.Namespaces[0]).
		WithName("name-1").Build()
	require.NoError(t, err)

	// WHEN Creating role
	role, err = authService.CreateRole(ctx, org.Id, role)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN Updating role
	err = authService.UpdateRole(ctx, org.Id, role)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN finding a role
	role, err = authService.GetRole(ctx, org.Id, org.Namespaces[0], role.Id)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, "name-1", role.Name)

	// WHEN querying roles
	res, _, err := authService.GetRoles(ctx, org.Id, org.Namespaces[0], map[string]string{"name": "name-1"}, "", 0)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, 1, len(res))

	for _, next := range res {
		require.Equal(t, role.Name, next.Name)
		// WHEN Deleting a role
		err := authService.DeleteRole(ctx, org.Id, org.Namespaces[0], next.Id)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func testCRUDRolesWithPermissions(
	ctx context.Context,
	t *testing.T,
	authService service.AuthAdminService,
	org *types.Organization,
) {
	role, err := domain.NewRoleBuilder().
		WithNamespace(org.Namespaces[0]).
		WithName("role1").Build()
	require.NoError(t, err)

	role, err = authService.CreateRole(ctx, org.Id, role)
	require.NoError(t, err)
	err = authService.AddPermissionsToRole(ctx, org.Id, org.Namespaces[0], role.Id, "perm1", "perm2")
	require.NoError(t, err)

	res, _, err := authService.GetRoles(ctx, org.Id, org.Namespaces[0], map[string]string{"name": "role1"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	for _, next := range res {
		require.Equal(t, role.Name, next.Name)
		require.Equal(t, 2, len(next.PermissionIds))
		err = authService.DeletePermissionsToRole(ctx, org.Id, org.Namespaces[0], next.Id, "perm1", "perm2")
		require.NoError(t, err)
		err := authService.DeleteRole(ctx, org.Id, org.Namespaces[0], next.Id)
		require.NoError(t, err)
		// Delete should fail without org-id, namespace or id
		require.Error(t, authService.DeleteRole(ctx, "", org.Namespaces[0], next.Id))
		require.Error(t, authService.DeleteRole(ctx, org.Id, "", next.Id))
		require.Error(t, authService.DeleteRole(ctx, org.Id, org.Namespaces[0], ""))
	}

}
