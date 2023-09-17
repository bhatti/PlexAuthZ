package db

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Should_CRUD_Groups(t *testing.T) {
	// GIVEN auth-service
	ctx := context.TODO()
	store, org, err := newAuthServiceAndOrg()
	require.NoError(t, err)

	group, err := domain.NewGroupBuilder().
		WithNamespace(org.Namespaces[0]).
		WithName(fmt.Sprintf("name_%d", 1)).Build()
	require.NoError(t, err)
	// WHEN Creating group
	group, err = store.CreateGroup(ctx, org.Id, group)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN Updating group
	err = store.UpdateGroup(ctx, org.Id, group)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN adding roles to a group
	err = store.AddRolesToGroup(ctx, org.Id, org.Namespaces[0], group.Id, "group-role1")
	require.NoError(t, err)
	err = store.AddRolesToGroup(ctx, org.Id, org.Namespaces[0], group.Id, "group-role2")
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN finding a group
	group, err = store.GetGroup(ctx, org.Id, org.Namespaces[0], group.Id)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, 2, len(group.RoleIds))

	// WHEN finding a group without id
	_, err = store.GetGroup(ctx, org.Id, org.Namespaces[0], "")
	// THEN it should fail
	require.Error(t, err)

	// WHEN querying groups
	res, _, err := store.GetGroups(ctx, org.Id, org.Namespaces[0], map[string]string{"name": "name_1"}, "", 0)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, 1, len(res))

	for _, next := range res {
		require.Equal(t, group.Name, next.Name)
		require.Equal(t, 2, len(next.RoleIds))

		// WHEN deleting roles to a group
		err = store.DeleteRolesToGroup(ctx, org.Id, org.Namespaces[0], next.Id, "group-role1")
		require.NoError(t, err)
		err = store.DeleteRolesToGroup(ctx, org.Id, org.Namespaces[0], next.Id, "group-role2")
		// THEN it should not fail
		require.NoError(t, err)

		// WHEN Deleting a group
		err := store.DeleteGroup(ctx, org.Id, org.Namespaces[0], next.Id)
		// THEN it should not fail
		require.NoError(t, err)

		// WHEN deleting a group without id
		err = store.DeleteGroup(ctx, org.Id, org.Namespaces[0], "")
		// THEN it should fail
		require.Error(t, err)
	}

}
