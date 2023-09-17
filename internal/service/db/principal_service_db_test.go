package db

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/utils"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
	"testing"
)

func Test_Should_CRUD_Principals(t *testing.T) {
	// GIVEN auth-service
	ctx := context.TODO()
	store, org, err := newAuthServiceAndOrg()
	require.NoError(t, err)

	principal, err := domain.NewPrincipalBuilder().
		WithNamespaces(org.Namespaces...).
		WithUsername("user1").
		WithName("john").
		WithOrganizationId(org.Id).Build()
	require.NoError(t, err)
	// WHEN Creating principal
	principal, err = store.CreatePrincipal(ctx, principal)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN Updating principal
	err = store.UpdatePrincipal(ctx, principal)
	// THEN it should not fail
	require.NoError(t, err)

	role1, err := store.CreateRole(ctx, org.Id, &types.Role{
		Name:      "name-1",
		Namespace: org.Namespaces[0],
	})
	require.NoError(t, err)
	role2, err := store.CreateRole(ctx, org.Id, &types.Role{
		Name:      "name-2",
		Namespace: org.Namespaces[0],
	})
	require.NoError(t, err)
	// WHEN adding roles to a principal
	err = store.AddRolesToPrincipal(ctx, org.Id, org.Namespaces[0], principal.Id, role1.Id)
	require.NoError(t, err)
	err = store.AddRolesToPrincipal(ctx, org.Id, org.Namespaces[0], principal.Id, role2.Id)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN finding a principal
	principal, err = store.GetPrincipal(ctx, org.Id, principal.Id)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, 2, len(principal.RoleIds))

	// WHEN querying principals
	res, _, err := store.GetPrincipals(ctx, org.Id, map[string]string{"username": "user1"}, "", 0)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, 1, len(res))

	for _, next := range res {
		require.Equal(t, 2, len(next.RoleIds))
		// WHEN Deleting a principal
		err := store.DeletePrincipal(ctx, org.Id, next.Id)
		// THEN it should not fail
		require.NoError(t, err)

		// WHEN Deleting a principal without id
		err = store.DeletePrincipal(ctx, org.Id, "")
		// THEN it should fail
		require.Error(t, err)
	}
}

func Test_ShouldCRUDPrincipalWithPermissions(t *testing.T) {
	// GIVEN auth-service
	ctx := context.TODO()
	store, org, err := newAuthServiceAndOrg()
	require.NoError(t, err)
	principal, err := domain.NewPrincipalBuilder().
		WithOrganizationId(org.Id).
		WithNamespaces(org.Namespaces...).
		WithName("john").
		WithUsername(uuid.NewV4().String()).Build()
	require.NoError(t, err)

	savedPrincipal, err := store.CreatePrincipal(ctx, principal)
	require.NoError(t, err)

	var permissionIds1 []string
	var permissionIds2 []string
	var relationshipIds []string
	for i := 0; i < 5; i++ {
		resource, _ := domain.NewResourceBuilder().
			WithNamespace(org.Namespaces[0]).
			WithName(fmt.Sprintf("/file/%d", i)).
			WithCapacity(1).
			WithAllowedActions("read", "write").Build()
		savedRes, err := store.CreateResource(ctx, org.Id, resource)
		require.NoError(t, err)

		for j := 0; j < 2; j++ {
			k := i*5 + j
			permission, _ := domain.NewPermissionBuilder().
				WithNamespace(org.Namespaces[0]).
				WithScope(fmt.Sprintf("scope_%d", k)).
				WithActions("read", "write").
				WithResourceId(savedRes.Id).
				WithEffect(types.Effect_PERMITTED).
				WithConstraints("time > 10").Build()
			savedPermission, err := store.CreatePermission(ctx, org.Id, permission)
			require.NoError(t, err)
			if i%2 == 0 {
				permissionIds1 = utils.AddSlice(permissionIds1, savedPermission.Id)
			} else {
				permissionIds2 = utils.AddSlice(permissionIds2, savedPermission.Id)
			}
		}
		relationship, err := domain.NewRelationshipBuilder().
			WithNamespace(org.Namespaces[0]).
			WithResourceId(savedRes.Id).
			WithPrincipalId(savedPrincipal.Id).
			WithRelation(fmt.Sprintf("rel_%d", i)).Build()
		require.NoError(t, err)
		savedRelation, err := store.CreateRelationship(ctx, org.Id, relationship)
		require.NoError(t, err)
		relationshipIds = utils.AddSlice(relationshipIds, savedRelation.Id)
	}

	var roleIDs1 []string
	var roleIDs2 []string
	for i := 0; i < 10; i++ {
		var rolePermIds []string
		if i%2 == 0 {
			rolePermIds = permissionIds1
		} else {
			rolePermIds = permissionIds2
		}
		role, err := domain.NewRoleBuilder().
			WithNamespace(org.Namespaces[0]).
			WithName(fmt.Sprintf("role_%d", i)).Build()
		require.NoError(t, err)
		role, err = store.CreateRole(ctx, org.Id, role)
		require.NoError(t, err)
		err = store.AddPermissionsToRole(ctx, org.Id, org.Namespaces[0], role.Id, rolePermIds...)
		require.NoError(t, err)

		if i%2 == 0 {
			roleIDs1 = utils.AddSlice(roleIDs1, role.Id)
		} else {
			roleIDs2 = utils.AddSlice(roleIDs1, role.Id)
		}
	}
	var groupIDs []string
	for i := 0; i < 10; i++ {
		var groupRoleIds []string
		if i%2 == 0 {
			groupRoleIds = roleIDs1
		} else {
			groupRoleIds = roleIDs2
		}
		group, err := domain.NewGroupBuilder().
			WithNamespace(org.Namespaces[0]).
			WithName(fmt.Sprintf("name_%d", i)).Build()
		require.NoError(t, err)
		group, err = store.CreateGroup(ctx, org.Id, group)
		require.NoError(t, err)
		err = store.AddRolesToGroup(ctx, org.Id, org.Namespaces[0], group.Id, groupRoleIds...)
		require.NoError(t, err)
		groupIDs = utils.AddSlice(groupIDs, group.Id)
	}
	for _, roleID := range roleIDs2 {
		err = store.AddRolesToPrincipal(ctx, savedPrincipal.OrganizationId, org.Namespaces[0], savedPrincipal.Id, roleID)
		require.NoError(t, err)
	}
	for _, groupID := range groupIDs {
		err = store.AddGroupsToPrincipal(ctx, savedPrincipal.OrganizationId, org.Namespaces[0], savedPrincipal.Id, groupID)
		require.NoError(t, err)
	}
	for _, permID := range permissionIds1 {
		err = store.AddPermissionsToPrincipal(ctx, savedPrincipal.OrganizationId, org.Namespaces[0], savedPrincipal.Id, permID)
		require.NoError(t, err)
	}
	for _, relationID := range relationshipIds {
		err = store.AddRelationshipsToPrincipal(ctx, savedPrincipal.OrganizationId, org.Namespaces[0], savedPrincipal.Id, relationID)
		require.NoError(t, err)
	}

	loadedPrincipal, err := store.GetPrincipalExt(ctx, org.Id, org.Namespaces[0], savedPrincipal.Id)
	require.NoError(t, err)
	require.Equal(t, savedPrincipal.Username, loadedPrincipal.Delegate.Username)
	require.Equal(t, savedPrincipal.Name, loadedPrincipal.Delegate.Name)
	require.Equal(t, savedPrincipal.OrganizationId, loadedPrincipal.Delegate.OrganizationId)
	require.Equal(t, len(roleIDs2), len(loadedPrincipal.Delegate.RoleIds))
	require.Equal(t, len(groupIDs), len(loadedPrincipal.Delegate.GroupIds))
	require.Equal(t, len(permissionIds1), len(loadedPrincipal.Delegate.PermissionIds))
	require.Equal(t, len(relationshipIds), len(loadedPrincipal.Delegate.RelationIds))

	require.Equal(t, len(groupIDs), len(loadedPrincipal.GroupsByName))
	require.Equal(t, len(roleIDs2), len(loadedPrincipal.RolesByName))
	require.Equal(t, len(permissionIds1)+len(permissionIds2), len(loadedPrincipal.AllPermissions()))
	require.Equal(t, 5, len(loadedPrincipal.ResourcesById))
	require.Equal(t, 5, len(loadedPrincipal.RelationsById))
	for id := range loadedPrincipal.ResourcesById {
		require.Equal(t, 1, len(loadedPrincipal.RelationNames(id)))
	}

	loadedPrincipals, _, err := store.GetPrincipals(ctx, org.Id, nil, "", 0)
	require.NoError(t, err)
	require.True(t, len(loadedPrincipals) > 0)

	// WHEN deleting groups, permissions and relationships
	for _, groupID := range groupIDs {
		err = store.DeleteGroupsToPrincipal(ctx, savedPrincipal.OrganizationId, org.Namespaces[0], savedPrincipal.Id, groupID)
		// THEN should not fail
		require.NoError(t, err)
	}
	for _, permID := range permissionIds1 {
		err = store.DeletePermissionsToPrincipal(ctx, savedPrincipal.OrganizationId, org.Namespaces[0], savedPrincipal.Id, permID)
		// THEN should not fail
		require.NoError(t, err)
	}
	for _, relationID := range relationshipIds {
		err = store.DeleteRelationshipsToPrincipal(ctx, savedPrincipal.OrganizationId, org.Namespaces[0], savedPrincipal.Id, relationID)
		// THEN should not fail
		require.NoError(t, err)
	}

}
