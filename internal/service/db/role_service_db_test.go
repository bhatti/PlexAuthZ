package db

import (
	"context"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Should_CRUD_Roles(t *testing.T) {
	// GIVEN auth-service
	ctx := context.TODO()
	store, org, err := newAuthServiceAndOrg()
	require.NoError(t, err)

	role, err := domain.NewRoleBuilder().
		WithNamespace(org.Namespaces[0]).
		WithName("name-1").Build()
	require.NoError(t, err)

	// WHEN Creating role
	role, err = store.CreateRole(ctx, org.Id, role)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN Updating role
	err = store.UpdateRole(ctx, org.Id, role)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN finding a role
	role, err = store.GetRole(ctx, org.Id, org.Namespaces[0], role.Id)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, "name-1", role.Name)

	// WHEN querying roles
	res, _, err := store.GetRoles(ctx, org.Id, org.Namespaces[0], map[string]string{"name": "name-1"}, "", 0)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, 1, len(res))

	for _, next := range res {
		require.Equal(t, role.Name, next.Name)
		// WHEN Deleting a role
		err := store.DeleteRole(ctx, org.Id, org.Namespaces[0], next.Id)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func Test_ShouldCRUDRoleWithPermissions(t *testing.T) {
	// GIVEN auth-service
	ctx := context.TODO()
	store, org, err := newAuthServiceAndOrg()
	require.NoError(t, err)
	role, err := domain.NewRoleBuilder().
		WithNamespace(org.Namespaces[0]).
		WithName("role1").Build()
	require.NoError(t, err)

	role, err = store.CreateRole(ctx, org.Id, role)
	require.NoError(t, err)
	err = store.AddPermissionsToRole(ctx, org.Id, org.Namespaces[0], role.Id, "perm1", "perm2")
	require.NoError(t, err)

	res, _, err := store.GetRoles(ctx, org.Id, org.Namespaces[0], map[string]string{"name": "role1"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	for _, next := range res {
		require.Equal(t, role.Name, next.Name)
		require.Equal(t, 2, len(next.PermissionIds))
		err = store.DeletePermissionsToRole(ctx, org.Id, org.Namespaces[0], next.Id, "perm1", "perm2")
		require.NoError(t, err)
		err := store.DeleteRole(ctx, org.Id, org.Namespaces[0], next.Id)
		require.NoError(t, err)
	}

}
