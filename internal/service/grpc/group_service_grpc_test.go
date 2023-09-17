package grpc

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/service"
	"github.com/stretchr/testify/require"
	"testing"
)

func testCRUDGroups(
	ctx context.Context,
	t *testing.T,
	authService service.AuthAdminService,
	org *types.Organization,
) {
	group, err := domain.NewGroupBuilder().
		WithNamespace(org.Namespaces[0]).
		WithName("name_1").Build()
	require.NoError(t, err)

	// WHEN Creating group
	group, err = authService.CreateGroup(ctx, org.Id, group)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN Updating group
	err = authService.UpdateGroup(ctx, org.Id, group)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN adding roles to a group
	err = authService.AddRolesToGroup(ctx, org.Id, org.Namespaces[0], group.Id, "group-role1")
	require.NoError(t, err)
	err = authService.AddRolesToGroup(ctx, org.Id, org.Namespaces[0], group.Id, "group-role2")
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN finding a group
	group, err = authService.GetGroup(ctx, org.Id, org.Namespaces[0], group.Id)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, 2, len(group.RoleIds))

	// WHEN querying groups
	res, _, err := authService.GetGroups(ctx, org.Id, org.Namespaces[0], map[string]string{"name": "name_1"}, "", 0)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, 1, len(res))

	for _, next := range res {
		require.Equal(t, group.Name, next.Name)
		require.Equal(t, 2, len(next.RoleIds))
		// WHEN Deleting a group
		err := authService.DeleteGroup(ctx, org.Id, org.Namespaces[0], next.Id)
		// THEN it should not fail
		require.NoError(t, err)
	}
}
