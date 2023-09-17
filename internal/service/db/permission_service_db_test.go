package db

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Should_CRUD_Permissions(t *testing.T) {
	// GIVEN auth-service
	ctx := context.TODO()
	store, org, err := newAuthServiceAndOrg()
	require.NoError(t, err)

	// WHEN creating a resource
	resource, err := store.CreateResource(ctx, org.Id, &types.Resource{
		Namespace:      org.Namespaces[0],
		Name:           "resource1",
		AllowedActions: []string{"read", "write"},
	})
	// THEN it should not fail
	require.NoError(t, err)

	// AND permission
	permission, err := domain.NewPermissionBuilder().
		WithNamespace(org.Namespaces[0]).
		WithScope(fmt.Sprintf("scope_%d", 1)).
		WithActions("read", "write").
		WithResourceId(resource.Id).
		WithEffect(types.Effect_PERMITTED).
		WithConstraints("time > 10").Build()
	require.NoError(t, err)

	// WHEN creating a permission
	savedPermission, err := store.CreatePermission(ctx, org.Id, permission)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN updating a permission
	err = store.UpdatePermission(ctx, org.Id, permission)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN querying an organization by id
	res, _, err := store.GetPermissions(ctx, org.Id, org.Namespaces[0], map[string]string{"id": savedPermission.Id}, "", 0)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, 1, len(res))

	// WHEN querying an organization by scope
	res, _, err = store.GetPermissions(ctx, org.Id, org.Namespaces[0], map[string]string{"scope": "scope_1"}, "", 0)
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
		err := store.DeletePermission(ctx, org.Id, org.Namespaces[0], next.Id)
		// THEN it should not fail
		require.NoError(t, err)

		// WHEN deleting a permission without id
		err = store.DeletePermission(ctx, org.Id, org.Namespaces[0], "")
		// THEN it should fail
		require.Error(t, err)
	}
}
