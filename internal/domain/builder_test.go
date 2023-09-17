package domain

import (
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/utils"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
	"testing"
)

func Test_ShouldBuildOrganization(t *testing.T) {
	org, err := NewOrganizationBuilder().
		WithId("test-org-"+uuid.NewV4().String()).
		WithName("org-name").
		WithUrl("org-url").
		WithNamespaces("admin", "finance", "engineering").Build()
	require.NoError(t, err)
	require.Equal(t, "org-name", org.Name)
	require.Equal(t, "org-url", org.Url)
	require.Equal(t, int64(0), org.Version)
	require.Equal(t, 3, len(org.Namespaces))
	require.Equal(t, 0, len(org.ParentIds))
}

func Test_ShouldBuildGroup(t *testing.T) {
	group, err := NewGroupBuilder().
		WithNamespace("ns").
		WithName("name1").Build()
	require.NoError(t, err)
	require.Equal(t, "name1", group.Name)
	require.Equal(t, 0, len(group.ParentIds))
}

func Test_ShouldRelationshipsWithoutPrincipal(t *testing.T) {
	_, err := NewRelationshipBuilder().
		WithRelation(fmt.Sprintf("rel-%d", 1)).
		WithAttribute("k", "v").
		WithResourceId(fmt.Sprintf("res-%d", 1)).Build()
	require.Error(t, err)
}

func Test_ShouldRelationships(t *testing.T) {
	relationship, err := NewRelationshipBuilder().
		WithNamespace("ns").
		WithRelation(fmt.Sprintf("rel-%d", 1)).
		WithPrincipalId(fmt.Sprintf("user-%d", 1)).
		WithResourceId(fmt.Sprintf("res-%d", 1)).Build()
	require.NoError(t, err)
	require.Equal(t, "rel-1", relationship.Relation)
	require.Equal(t, int64(0), relationship.Version)
}

func Test_ShouldBuildRole(t *testing.T) {
	role, err := NewRoleBuilder().
		WithNamespace("ns").
		WithName("role1").Build()
	require.NoError(t, err)
	require.Equal(t, "role1", role.Name)
}

func Test_ShouldBuildResource(t *testing.T) {
	resource, err := NewResourceBuilder().
		WithNamespace("ns").
		WithName("one").
		WithCapacity(1).
		WithAttribute("k1", "v1").
		WithAllowedActions("read", "write").Build()
	require.NoError(t, err)
	require.Equal(t, "one", resource.Name)
	require.Equal(t, int32(1), resource.Capacity)
	require.Equal(t, 2, len(resource.AllowedActions))
}

func Test_ShouldBuildPermission(t *testing.T) {
	permission, err := NewPermissionBuilder().
		WithNamespace(fmt.Sprintf("ns_%d", 1)).
		WithScope(fmt.Sprintf("scope_%d", 1)).
		WithActions("read", "write").
		WithResourceId("2101").
		WithEffect(types.Effect_PERMITTED).
		WithConstraints("time > 10").Build()
	require.NoError(t, err)
	require.Equal(t, "ns_1", permission.Namespace)
	require.Equal(t, "scope_1", permission.Scope)
	require.Equal(t, "2101", permission.ResourceId)
	require.Equal(t, types.Effect_PERMITTED, permission.Effect)
	require.Equal(t, 2, len(permission.Actions))
}

func Test_ShouldBuildPrincipal(t *testing.T) {
	org, err := NewOrganizationBuilder().
		WithId("test-org-"+uuid.NewV4().String()).
		WithName("org-name").
		WithUrl("org-url").
		WithNamespaces("admin", "finance", "engineering").Build()
	require.NoError(t, err)
	principal, err := NewPrincipalBuilder().
		WithOrganizationId(org.Id).
		WithNamespaces(org.Namespaces...).
		WithEmail("email").
		WithAttribute("k", "v").
		WithName("john").
		WithUsername("p123").Build()
	require.NoError(t, err)

	var permissionIds1 []string
	var permissionIds2 []string
	var relationshipIds []string
	for i := 0; i < 10; i++ {
		resource, _ := NewResourceBuilder().
			WithNamespace(org.Namespaces[0]).
			WithName(fmt.Sprintf("/file/%d", i)).
			WithCapacity(1).
			WithAllowedActions("read", "write").Build()
		for j := 0; j < 10; j++ {
			k := i*10 + j
			permission, _ := NewPermissionBuilder().
				WithNamespace(org.Namespaces[0]).
				WithScope(fmt.Sprintf("scope_%d", k)).
				WithActions("read", "write").
				WithResourceId("2101").
				WithEffect(types.Effect_PERMITTED).
				WithConstraints("time > 10").Build()

			permission.ResourceId = resource.Id
			if i%2 == 0 {
				permissionIds1 = utils.AddSlice(permissionIds1, permission.Id)
			} else {
				permissionIds2 = utils.AddSlice(permissionIds2, permission.Id)
			}
		}
		for i := 0; i < 10; i++ {
			relationship, _ := NewRelationshipBuilder().
				WithNamespace(org.Namespaces[0]).
				WithRelation(fmt.Sprintf("rel-%d", i)).
				WithPrincipalId(principal.Id).
				WithResourceId(resource.Id).Build()
			relationshipIds = utils.AddSlice(relationshipIds, relationship.Id)
		}
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
		role, err := NewRoleBuilder().
			WithNamespace(org.Namespaces[0]).
			WithParentIds("id").
			WithName(fmt.Sprintf("role_%d", i)).Build()
		require.NoError(t, err)
		role.PermissionIds = rolePermIds
	}
	var groupIDs []string
	for i := 0; i < 10; i++ {
		var groupRoleIds []string
		if i%2 == 0 {
			groupRoleIds = roleIDs1
		} else {
			groupRoleIds = roleIDs2
		}
		group, err := NewGroupBuilder().
			WithNamespace(org.Namespaces[0]).
			WithParentIds("id").
			WithName(fmt.Sprintf("test-group-%d", i)).Build()
		require.NoError(t, err)
		group.RoleIds = groupRoleIds
		groupIDs = utils.AddSlice(groupIDs, group.Id)
	}

	principal.RoleIds = roleIDs2
	principal.GroupIds = groupIDs
	principal.PermissionIds = permissionIds1
	principal.RelationIds = relationshipIds

	require.Equal(t, "p123", principal.Username)
	require.Equal(t, "john", principal.Name)
	require.Equal(t, org.Id, principal.OrganizationId)
	require.Equal(t, len(groupIDs), len(principal.GroupIds))
	require.Equal(t, len(roleIDs2), len(principal.RoleIds))
	require.Equal(t, len(permissionIds1), len(principal.PermissionIds))
	require.Equal(t, len(relationshipIds), len(principal.RelationIds))
}
