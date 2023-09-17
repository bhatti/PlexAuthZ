package domain

import (
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/utils"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
	"testing"
)

func Test_ShouldCreateHashIndex(t *testing.T) {
	hash := NewHashIndex("h", nil)
	require.Error(t, hash.Validate())
	hash.Ids = []string{"id"}
	require.NoError(t, hash.Validate())
}

func Test_ShouldCreateOrganization(t *testing.T) {
	org := createTestOrg()
	require.Equal(t, "org-name", org.Name)
	require.Equal(t, int64(0), org.Version)
	require.Equal(t, 0, len(org.ParentIds))

	xOrg := NewOrganizationExt(&org)
	require.NoError(t, xOrg.Validate())
	org.Namespaces = nil
	require.Error(t, xOrg.Validate())
	org.Namespaces = []string{"ns"}
	org.Name = ""
	require.Error(t, xOrg.Validate())
	org.Name = "name"
	xOrg.Delegate = nil
	require.Error(t, xOrg.Validate())
}

func Test_ShouldCreateGroup(t *testing.T) {
	group := createTestGroup(555)
	require.Equal(t, "name_555", group.Name)
	require.Equal(t, int64(0), group.Version)
	require.Equal(t, 0, len(group.ParentIds))

	xGroup := NewGroupExt(&group)
	require.NoError(t, xGroup.Validate())
	xGroup.Delegate.Namespace = ""
	require.Error(t, xGroup.Validate())
	xGroup.Delegate.Namespace = "ns"
	xGroup.Delegate.Name = ""
	require.Error(t, xGroup.Validate())
	require.NotEqual(t, "", xGroup.Hash())
	require.NotEqual(t, "", xGroup.String())
}

func Test_ShouldCreateRelationships(t *testing.T) {
	relationship := createTestRelationship(555)
	require.Equal(t, "rel_555", relationship.Relation)
	require.Equal(t, int64(0), relationship.Version)

	xRelationship := NewRelationshipExt(&relationship)
	require.NoError(t, xRelationship.Validate())
	xRelationship.Delegate.Namespace = ""
	require.Error(t, xRelationship.Validate())
	xRelationship.Delegate.Namespace = "ns"
	xRelationship.Delegate.Relation = ""
	require.Error(t, xRelationship.Validate())
	xRelationship.Delegate.Relation = "rel"
	xRelationship.Delegate.PrincipalId = ""
	require.Error(t, xRelationship.Validate())
	xRelationship.Delegate.PrincipalId = "id"
	xRelationship.Delegate.ResourceId = ""
	require.Error(t, xRelationship.Validate())
	require.NotEqual(t, "", xRelationship.Hash())
	require.NotEqual(t, "", xRelationship.String())
}

func Test_ShouldCreateRole(t *testing.T) {
	role := createTestRole(555)
	require.Equal(t, "/file/555", role.Name)
	require.Equal(t, 2, len(role.PermissionIds))

	xPerm := NewRoleExt(&role)
	require.NoError(t, xPerm.Validate())
	xPerm.Delegate.Namespace = ""
	require.Error(t, xPerm.Validate())
	xPerm.Delegate.Namespace = "ns"
	xPerm.Delegate.Name = ""
	require.Error(t, xPerm.Validate())
	require.NotEqual(t, "", xPerm.Hash())
	require.NotEqual(t, "", xPerm.String())
}

func Test_ShouldCreateResource(t *testing.T) {
	resource := createTestResource(555)
	require.Equal(t, "/file/555", resource.Name)
	require.Equal(t, 2, len(resource.AllowedActions))

	xResource := NewResourceExt(&resource)
	require.NoError(t, xResource.Validate())
	resource.Namespace = ""
	require.Error(t, xResource.Validate())
	resource.Namespace = "ns"
	resource.Name = ""
	require.Error(t, xResource.Validate())
	resource.Name = "name"
	resource.AllowedActions = nil
	require.Error(t, xResource.Validate())

	resource.AllowedActions = []string{"read"}
	xResource.Delegate = nil
	require.Error(t, xResource.Validate())

	xResource.Delegate = &resource
	require.NotEqual(t, "", xResource.Hash())

	instance := NewResourceInstanceExt(resource.Namespace, "res-id", "user-id")
	require.NoError(t, instance.Validate())
	instance.Delegate.Namespace = ""
	require.Error(t, instance.Validate())
	instance.Delegate.Namespace = "ns"
	instance.Delegate.ResourceId = ""
	require.Error(t, instance.Validate())
	instance.Delegate.ResourceId = "res-id"
	instance.Delegate.PrincipalId = ""
	require.Error(t, instance.Validate())
	require.NotEqual(t, "", instance.Hash())
	require.NotEqual(t, "", instance.String())
}

func Test_ShouldConvertUnsafeCaseInt32ToBytes(t *testing.T) {
	require.True(t, len(unsafeCaseInt32ToBytes(111)) > 0)
}

func Test_ShouldMatchResourceName(t *testing.T) {
	require.True(t, doesResourceNameMatches("test*", "test123"))
}

func Test_ShouldCreatePermission(t *testing.T) {
	permission := createTestPermission(555)
	require.Equal(t, "ns_555", permission.Namespace)
	require.Equal(t, "scope_555", permission.Scope)
	require.Equal(t, "2101", permission.ResourceId)
	require.Equal(t, types.Effect_PERMITTED, permission.Effect)
	require.Equal(t, 2, len(permission.Actions))

	xPerm := NewPermissionExt(&permission)
	require.NoError(t, xPerm.Validate())
	xPerm.Delegate.Namespace = ""
	require.Error(t, xPerm.Validate())
	xPerm.Delegate.Namespace = "ns"
	xPerm.Delegate.ResourceId = ""
	require.Error(t, xPerm.Validate())
	xPerm.Delegate.ResourceId = "res-id"
	xPerm.Delegate.Actions = nil
	require.Error(t, xPerm.Validate())
	require.NotEqual(t, "", xPerm.Hash())
	require.NotEqual(t, "", xPerm.String())
}

func Test_ShouldCreatePrincipal(t *testing.T) {
	org := createTestOrg()
	principal := createTestPrincipal(777)
	principal.OrganizationId = org.Id
	principal.Username = "p123"

	var permissionIds1 []string
	var permissionIds2 []string
	var relationshipIds []string
	for i := 0; i < 10; i++ {
		resource := createTestResource(i)
		for j := 0; j < 10; j++ {
			permission := createTestPermission(i*10 + j)
			permission.ResourceId = resource.Id
			if i%2 == 0 {
				permissionIds1 = utils.AddSlice(permissionIds1, permission.Id)
			} else {
				permissionIds2 = utils.AddSlice(permissionIds2, permission.Id)
			}
		}
		for i := 0; i < 10; i++ {
			relationship := createTestRelationship(i)
			relationship.ResourceId = resource.Id
			relationship.PrincipalId = principal.Id
			relationshipIds = utils.AddSlice(relationshipIds, relationship.Id)
		}
	}

	var roleIDs1 []string
	var roleIDs2 []string
	for i := 0; i < 10; i++ {
		role := createTestRole(i)
		if i%2 == 0 {
			role.PermissionIds = permissionIds1
		} else {
			role.PermissionIds = permissionIds2
		}
		if i%2 == 0 {
			roleIDs1 = utils.AddSlice(roleIDs1, role.Id)
		} else {
			roleIDs2 = utils.AddSlice(roleIDs1, role.Id)
		}
	}
	var groupIDs []string
	for i := 0; i < 10; i++ {
		group := createTestGroup(i)
		if i%2 == 0 {
			group.RoleIds = roleIDs1
		} else {
			group.RoleIds = roleIDs2
		}
		groupIDs = utils.AddSlice(groupIDs, group.Id)
	}
	principal.RoleIds = roleIDs2
	principal.GroupIds = groupIDs
	principal.PermissionIds = permissionIds1
	principal.RelationIds = relationshipIds

	require.Equal(t, "p123", principal.Username)
	require.Equal(t, org.Id, principal.OrganizationId)
	require.Equal(t, len(groupIDs), len(principal.GroupIds))
	require.Equal(t, len(roleIDs2), len(principal.RoleIds))
	require.Equal(t, len(permissionIds1), len(principal.PermissionIds))
	require.Equal(t, len(relationshipIds), len(principal.RelationIds))
}

func Test_ShouldCreatePrincipalExt(t *testing.T) {
	org := createTestOrg()
	principal := createTestPrincipal(777)
	principal.OrganizationId = org.Id
	principal.Username = "p123"
	principal.Namespaces = org.Namespaces
	xPrincipal := NewPrincipalExt(&principal)

	var permissionIds1 []string
	var permissionIds2 []string
	var relationshipIds []string
	for i := 0; i < 10; i++ {
		resource := createTestResource(i)
		xPrincipal.ResourcesById[resource.Id] = &resource
		for j := 0; j < 10; j++ {
			permission := createTestPermission(i*10 + j)
			permission.ResourceId = resource.Id
			err := xPrincipal.AddPermission(&permission)
			require.NoError(t, err)
			permission.ResourceId = resource.Id
			if i%2 == 0 {
				permissionIds1 = utils.AddSlice(permissionIds1, permission.Id)
			} else {
				permissionIds2 = utils.AddSlice(permissionIds2, permission.Id)
			}
		}
		for i := 0; i < 10; i++ {
			relationship := createTestRelationship(i)
			xPrincipal.RelationsById[relationship.Id] = &relationship
			relationship.ResourceId = resource.Id
			relationship.PrincipalId = principal.Id
			relationshipIds = utils.AddSlice(relationshipIds, relationship.Id)
		}
	}

	var roleIDs1 []string
	var roleIDs2 []string
	for i := 0; i < 10; i++ {
		role := createTestRole(i)
		xPrincipal.RolesByName[role.Name] = &role
		if i%2 == 0 {
			role.PermissionIds = permissionIds1
		} else {
			role.PermissionIds = permissionIds2
		}
		if i%2 == 0 {
			roleIDs1 = utils.AddSlice(roleIDs1, role.Id)
		} else {
			roleIDs2 = utils.AddSlice(roleIDs1, role.Id)
		}
	}
	var groupIDs []string
	for i := 0; i < 10; i++ {
		group := createTestGroup(i)
		xPrincipal.GroupsByName[group.Name] = &group
		if i%2 == 0 {
			group.RoleIds = roleIDs1
		} else {
			group.RoleIds = roleIDs2
		}
		groupIDs = utils.AddSlice(groupIDs, group.Id)
	}
	principal.RoleIds = roleIDs2
	principal.GroupIds = groupIDs
	principal.PermissionIds = permissionIds1
	principal.RelationIds = relationshipIds

	require.Equal(t, "p123", principal.Username)
	require.Equal(t, org.Id, principal.OrganizationId)
	require.Equal(t, len(groupIDs), len(principal.GroupIds))
	require.Equal(t, len(roleIDs2), len(principal.RoleIds))
	require.Equal(t, len(permissionIds1), len(principal.PermissionIds))
	require.Equal(t, len(relationshipIds), len(principal.RelationIds))

	require.Equal(t, len(groupIDs), len(xPrincipal.GroupsByName))
	require.Equal(t, 10, len(xPrincipal.RolesByName))
	require.Equal(t, len(permissionIds1)+len(permissionIds2), len(xPrincipal.AllPermissions()))
	require.Equal(t, 10, len(xPrincipal.ResourcesById))
	require.Equal(t, 100, len(xPrincipal.RelationsById))
	for id := range xPrincipal.ResourcesById {
		require.Equal(t, 10, len(xPrincipal.RelationNames(id)))
	}

	response := xPrincipal.ToGetPrincipalResponse()
	clone := NewPrincipalExtFromResponse(response)
	require.NotNil(t, clone)
	m := xPrincipal.ToMap(&services.AuthRequest{
		OrganizationId: principal.OrganizationId,
		Namespace:      "ns",
		PrincipalId:    "user-id",
		Action:         "action",
		Resource:       "res",
		Context:        map[string]string{"k": "v"},
	},
		&types.Resource{
			Name:       "name",
			Attributes: map[string]string{"k": "v"},
		})
	require.NotNil(t, m)

	require.NoError(t, xPrincipal.Validate())
	xPrincipal.Delegate = nil
	require.Error(t, xPrincipal.Validate())
	xPrincipal.Delegate = &principal
	require.NoError(t, xPrincipal.Validate())
	xPrincipal.Delegate.Username = ""
	require.Error(t, xPrincipal.Validate())
	xPrincipal.Delegate.Username = "user"
	xPrincipal.Delegate.OrganizationId = ""
	require.Error(t, xPrincipal.Validate())
	xPrincipal.Delegate.OrganizationId = "id"
	xPrincipal.Delegate.Namespaces = nil
	require.Error(t, xPrincipal.Validate())

	require.True(t, len(xPrincipal.Roles()) > 0)
	require.True(t, len(xPrincipal.RoleNames()) > 0)
	require.True(t, len(xPrincipal.GroupNames()) > 0)
	require.True(t, len(xPrincipal.Groups()) > 0)
	require.True(t, len(xPrincipal.ResourceNames()) > 0)
	require.True(t, len(xPrincipal.Hash()) > 0)
	require.True(t, len(xPrincipal.GroupHashIndex()) > 0)
	require.True(t, len(xPrincipal.RoleHashIndex()) > 0)
	require.True(t, len(xPrincipal.String()) > 0)
	require.Nil(t, xPrincipal.LatestGroupRoleDate())
	_, err := xPrincipal.CheckPermission(&services.AuthRequest{
		OrganizationId: principal.OrganizationId,
		Namespace:      "ns",
		PrincipalId:    "user-id",
		Action:         "action",
		Resource:       "res",
		Context:        map[string]string{"k": "v"},
	})
	require.Error(t, err)
	res, _, err := xPrincipal.CheckConstraints(&services.AuthRequest{
		OrganizationId: principal.OrganizationId,
		Namespace:      "ns",
		PrincipalId:    "user-id",
		Action:         "action",
		Resource:       "res",
		Context:        map[string]string{"k": "v"},
	}, &types.Resource{}, "eq 11 12")
	require.NoError(t, err)
	require.False(t, res)
}

func createTestResource(i int) types.Resource {
	return types.Resource{
		Id:             "test-resource-" + uuid.NewV4().String(),
		Namespace:      "ns",
		Name:           fmt.Sprintf("/file/%d", i),
		Capacity:       int32(i + 1),
		Attributes:     map[string]string{"k": "v"},
		AllowedActions: []string{"read", "write"},
	}
}

func createTestGroup(i int) types.Group {
	return types.Group{
		Namespace: "ns",
		Id:        "test-group-" + uuid.NewV4().String(),
		Name:      fmt.Sprintf("name_%d", i),
		RoleIds:   []string{"group-role1", "group-role2"},
	}
}

func createTestOrg() types.Organization {
	return types.Organization{
		Id:         "test-org-" + uuid.NewV4().String(),
		Name:       "org-name",
		Url:        "org-url",
		Namespaces: []string{"admin", "finance", "engineering"},
	}
}

func createTestRelationship(i int) types.Relationship {
	return types.Relationship{
		Namespace:   "ns",
		Id:          "test-relation-" + uuid.NewV4().String(),
		Relation:    fmt.Sprintf("rel_%d", i),
		PrincipalId: fmt.Sprintf("user_%d", i),
		ResourceId:  fmt.Sprintf("res_%d", i),
		Attributes:  map[string]string{"k": "v"},
	}
}

func createTestRole(i int) types.Role {
	return types.Role{
		Namespace:     "ns",
		Id:            "test-role-" + uuid.NewV4().String(),
		Name:          fmt.Sprintf("/file/%d", i),
		PermissionIds: []string{"perm1", "perm2"},
	}
}

func createTestPermission(i int) types.Permission {
	return types.Permission{
		Id:          "test-permission-" + uuid.NewV4().String(),
		Namespace:   fmt.Sprintf("ns_%d", i),
		Scope:       fmt.Sprintf("scope_%d", i),
		Actions:     []string{"read", "write"},
		ResourceId:  "2101",
		Effect:      types.Effect_PERMITTED,
		Constraints: "time > 10",
	}
}

func createTestPrincipal(i int) types.Principal {
	return types.Principal{
		Id:             "test-principal-" + uuid.NewV4().String(),
		Username:       fmt.Sprintf("name_%d", i),
		OrganizationId: fmt.Sprintf("org_%d", i),
		Attributes:     map[string]string{"k": "v"},
	}
}
