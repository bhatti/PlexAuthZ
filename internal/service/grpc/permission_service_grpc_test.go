package grpc

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/service"
	"github.com/stretchr/testify/require"
	"testing"
)

func testCRUDPermissions(
	ctx context.Context,
	t *testing.T,
	authService service.AuthAdminService,
	org *types.Organization,
) {
	resource, err := authService.CreateResource(ctx, org.Id, &types.Resource{
		Namespace:      org.Namespaces[0],
		Name:           "resource1",
		AllowedActions: []string{"read", "write"},
	})
	require.NoError(t, err)
	permission, err := domain.NewPermissionBuilder().
		WithNamespace(org.Namespaces[0]).
		WithScope(fmt.Sprintf("scope_%d", 1)).
		WithActions("read", "write").
		WithResourceId(resource.Id).
		WithEffect(types.Effect_PERMITTED).
		WithConstraints("time > 10").Build()
	require.NoError(t, err)

	// WHEN creating a permission
	savedPermission, err := authService.CreatePermission(ctx, org.Id, permission)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN updating a permission
	err = authService.UpdatePermission(ctx, org.Id, permission)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN getting a permission
	loaded, err := authService.GetPermission(ctx, org.Id, permission.Namespace, permission.Id)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, permission.Id, loaded.Id)

	// WHEN querying an organization by id
	res, _, err := authService.GetPermissions(
		ctx,
		org.Id,
		org.Namespaces[0],
		map[string]string{"id": savedPermission.Id}, "", 0)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, 1, len(res))

	// WHEN querying an organization by scope
	res, _, err = authService.GetPermissions(
		ctx,
		org.Id,
		org.Namespaces[0],
		map[string]string{"scope": "scope_1"}, "", 0)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, 1, len(res))

	// Iterating permissions
	for _, next := range res {
		require.Equal(t, permission.Namespace, next.Namespace)
		require.Equal(t, permission.Scope, next.Scope)
		require.Equal(t, resource.Id, next.ResourceId)
		require.Equal(t, types.Effect_PERMITTED, next.Effect)
		require.Equal(t, 2, len(next.Actions))
		// WHEN deleting a permission
		err := authService.DeletePermission(ctx, org.Id, org.Namespaces[0], next.Id)
		// THEN it should not fail
		require.NoError(t, err)
		// WHEN deleting a permission without org-id, namespace or id should fail
		require.Error(t, authService.DeletePermission(ctx, "", org.Namespaces[0], next.Id))
		require.Error(t, authService.DeletePermission(ctx, org.Id, "", next.Id))
		require.Error(t, authService.DeletePermission(ctx, org.Id, org.Namespaces[0], ""))
	}
}
