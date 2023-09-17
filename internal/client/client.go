package client

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/authz"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/service"
	"github.com/bhatti/PlexAuthZ/internal/utils"
	"time"
)

// AuthAdapter adapter for auth-service.
type AuthAdapter struct {
	authorizer       authz.Authorizer
	authAdminService service.AuthAdminService
}

// New constructor
func New(
	authorizer authz.Authorizer,
	authAdminService service.AuthAdminService,
) *AuthAdapter {
	return &AuthAdapter{authorizer: authorizer, authAdminService: authAdminService}
}

// CreateOrganization adapter
func (c *AuthAdapter) CreateOrganization(
	org *types.Organization) (*OrganizationAdapter, error) {
	org, err := c.authAdminService.CreateOrganization(context.Background(), org)
	if err != nil {
		return nil, err
	}
	return &OrganizationAdapter{
		authorizer:       c.authorizer,
		authAdminService: c.authAdminService,
		Organization:     org,
	}, nil
}

// GetOrganization adapter
func (c *AuthAdapter) GetOrganization(
	orgID string) (*OrganizationAdapter, error) {
	org, err := c.authAdminService.GetOrganization(context.Background(), orgID)
	if err != nil {
		return nil, err
	}
	return &OrganizationAdapter{
		authorizer:       c.authorizer,
		authAdminService: c.authAdminService,
		Organization:     org,
	}, nil
}

// OrganizationAdapter for managing organizations.
type OrganizationAdapter struct {
	authorizer       authz.Authorizer
	authAdminService service.AuthAdminService
	Organization     *types.Organization
}

// Update adapter for updating organization.
func (c *OrganizationAdapter) Update() error {
	return c.authAdminService.UpdateOrganization(context.Background(), c.Organization)
}

// Delete adapter for deleting organization.
func (c *OrganizationAdapter) Delete() error {
	return c.authAdminService.DeleteOrganization(context.Background(), c.Organization.Id)
}

// Principals adapter for managing principals.
func (c *OrganizationAdapter) Principals() *PrincipalAdapter {
	return &PrincipalAdapter{
		authAdminService: c.authAdminService,
		authorizer:       c.authorizer,
		Principal: &types.Principal{
			OrganizationId: c.Organization.Id,
			Namespaces:     c.Organization.Namespaces,
		},
	}
}

// Groups adapter for managing groups.
func (c *OrganizationAdapter) Groups(namespace string) *GroupAdapter {
	return &GroupAdapter{
		authAdminService: c.authAdminService,
		orgID:            c.Organization.Id,
		Group: &types.Group{
			Namespace: namespace,
		},
	}
}

// Permissions adapter for managing permissions.
func (c *OrganizationAdapter) Permissions(namespace string) *PermissionAdapter {
	return &PermissionAdapter{
		authAdminService: c.authAdminService,
		orgID:            c.Organization.Id,
		Permission: &types.Permission{
			Namespace: namespace,
		},
	}
}

// Resources adapter for managing resources.
func (c *OrganizationAdapter) Resources(namespace string) *ResourceAdapter {
	return &ResourceAdapter{
		authAdminService: c.authAdminService,
		orgID:            c.Organization.Id,
		Resource: &types.Resource{
			Namespace: namespace,
		},
		context: make(map[string]string),
	}
}

// Roles adapter for managing roles.
func (c *OrganizationAdapter) Roles(namespace string) *RoleAdapter {
	return &RoleAdapter{
		authAdminService: c.authAdminService,
		orgID:            c.Organization.Id,
		Role: &types.Role{
			Namespace: namespace,
		},
	}
}

// AuthorizerAdapter for authorization request.
type AuthorizerAdapter struct {
	authorizer  authz.Authorizer
	Principal   *types.Principal
	namespace   string
	action      string
	resource    string
	scope       string
	context     map[string]string
	constraints string
	LastMessage string
}

// WithAction setter.
func (c *AuthorizerAdapter) WithAction(action string) *AuthorizerAdapter {
	c.action = action
	return c
}

// WithConstraints setter.
func (c *AuthorizerAdapter) WithConstraints(constraints string) *AuthorizerAdapter {
	c.constraints = constraints
	return c
}

// WithResource setter.
func (c *AuthorizerAdapter) WithResource(resource *types.Resource) *AuthorizerAdapter {
	return c.WithResourceName(resource.Name)
}

// WithResourceName setter.
func (c *AuthorizerAdapter) WithResourceName(resource string) *AuthorizerAdapter {
	c.resource = resource
	return c
}

// WithScope setter.
func (c *AuthorizerAdapter) WithScope(scope string) *AuthorizerAdapter {
	c.scope = scope
	return c
}

// WithContext setter.
func (c *AuthorizerAdapter) WithContext(ctx ...string) *AuthorizerAdapter {
	c.context = utils.ArrayToMap(ctx...)
	return c
}

// Check checks for authorization access.
func (c *AuthorizerAdapter) Check() error {
	if c.constraints != "" {
		req := &services.CheckConstraintsRequest{
			OrganizationId: c.Principal.OrganizationId,
			Namespace:      c.namespace,
			PrincipalId:    c.Principal.Id,
			Constraints:    c.constraints,
			Context:        c.context,
		}
		res, err := c.authorizer.Check(context.Background(), req)
		if err != nil {
			return err
		}
		if !res.Matched {
			return domain.NewAuthError(fmt.Sprintf("faied to check for constraints '%s' for principal %s",
				c.constraints, c.Principal.Username))
		}
		c.LastMessage = res.Output
	} else {
		req := &services.AuthRequest{
			OrganizationId: c.Principal.OrganizationId,
			Namespace:      c.namespace,
			PrincipalId:    c.Principal.Id,
			Action:         c.action,
			Resource:       c.resource,
			Scope:          c.scope,
			Context:        c.context,
		}
		res, err := c.authorizer.Authorize(context.Background(), req)
		if err != nil {
			return err
		}
		c.LastMessage = res.Message
		if res.Effect != types.Effect_PERMITTED {
			return domain.NewAuthError(fmt.Sprintf("principal %s cannot access %s for %s %s",
				c.Principal.Username, c.resource, c.action, res.Message))
		}
	}
	return nil
}

// PrincipalAdapter for managing principals.
type PrincipalAdapter struct {
	authorizer       authz.Authorizer
	authAdminService service.AuthAdminService
	Principal        *types.Principal
}

// WithUsername setter.
func (c *PrincipalAdapter) WithUsername(username string) *PrincipalAdapter {
	c.Principal.Username = username
	return c
}

// WithName setter.
func (c *PrincipalAdapter) WithName(name string) *PrincipalAdapter {
	c.Principal.Name = name
	return c
}

// WithEmail setter.
func (c *PrincipalAdapter) WithEmail(email string) *PrincipalAdapter {
	c.Principal.Email = email
	return c
}

// WithAttributes setter.
func (c *PrincipalAdapter) WithAttributes(
	kv ...string,
) *PrincipalAdapter {
	c.Principal.Attributes = utils.ArrayToMap(kv...)
	return c
}

// Create adds Principal in the database.
func (c *PrincipalAdapter) Create() (*PrincipalAdapter, error) {
	var err error
	c.Principal, err = c.authAdminService.CreatePrincipal(context.Background(), c.Principal)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Get finds Principal from the database.
func (c *PrincipalAdapter) Get(principalID string) error {
	principal, err := c.authAdminService.GetPrincipal(context.Background(), c.Principal.OrganizationId, principalID)
	if err != nil {
		return err
	}
	c.Principal = principal
	return nil
}

// Update updates Principal.
func (c *PrincipalAdapter) Update() error {
	return c.authAdminService.UpdatePrincipal(context.Background(), c.Principal)
}

// Delete removes principal.
func (c *PrincipalAdapter) Delete() error {
	return c.authAdminService.DeletePrincipal(context.Background(), c.Principal.OrganizationId, c.Principal.Id)
}

// AddGroups adds groups to principal.
func (c *PrincipalAdapter) AddGroups(groups ...*types.Group) error {
	var groupIDs []string
	var namespace string
	for _, group := range groups {
		groupIDs = utils.AddSlice(groupIDs, group.Id)
		namespace = group.Namespace
	}
	return c.authAdminService.AddGroupsToPrincipal(
		context.Background(),
		c.Principal.OrganizationId,
		namespace,
		c.Principal.Id,
		groupIDs...)
}

// DeleteGroups removes groups from the principal.
func (c *PrincipalAdapter) DeleteGroups(groups ...*types.Group) error {
	var groupIDs []string
	var namespace string
	for _, group := range groups {
		groupIDs = utils.AddSlice(groupIDs, group.Id)
		namespace = group.Namespace
	}
	return c.authAdminService.DeleteGroupsToPrincipal(
		context.Background(),
		c.Principal.OrganizationId,
		namespace,
		c.Principal.Id,
		groupIDs...)
}

// AddRoles adds roles to principal.
func (c *PrincipalAdapter) AddRoles(roles ...*types.Role) error {
	var roleIDs []string
	var namespace string
	for _, role := range roles {
		roleIDs = utils.AddSlice(roleIDs, role.Id)
		namespace = role.Namespace
	}
	return c.authAdminService.AddRolesToPrincipal(
		context.Background(),
		c.Principal.OrganizationId,
		namespace,
		c.Principal.Id,
		roleIDs...)
}

// DeleteRoles removes roles to principal.
func (c *PrincipalAdapter) DeleteRoles(roles ...*types.Role) error {
	var roleIDs []string
	var namespace string
	for _, role := range roles {
		roleIDs = utils.AddSlice(roleIDs, role.Id)
		namespace = role.Namespace
	}
	return c.authAdminService.DeleteRolesToPrincipal(
		context.Background(),
		c.Principal.OrganizationId,
		namespace,
		c.Principal.Id,
		roleIDs...)
}

// AddPermissions adds permissions to principal.
func (c *PrincipalAdapter) AddPermissions(permissions ...*types.Permission) error {
	var permissionIds []string
	var namespace string
	for _, permission := range permissions {
		permissionIds = utils.AddSlice(permissionIds, permission.Id)
		namespace = permission.Namespace
	}
	return c.authAdminService.AddPermissionsToPrincipal(
		context.Background(),
		c.Principal.OrganizationId,
		namespace,
		c.Principal.Id,
		permissionIds...)
}

// DeletePermissions removes permissions to principal.
func (c *PrincipalAdapter) DeletePermissions(permissions ...*types.Permission) error {
	var permissionIds []string
	var namespace string
	for _, permission := range permissions {
		permissionIds = utils.AddSlice(permissionIds, permission.Id)
		namespace = permission.Namespace
	}
	return c.authAdminService.DeletePermissionsToPrincipal(
		context.Background(),
		c.Principal.OrganizationId,
		namespace,
		c.Principal.Id,
		permissionIds...)
}

// AddRelations adds relations to principal.
func (c *PrincipalAdapter) AddRelations(relations ...*types.Relationship) error {
	var relationIds []string
	var namespace string
	for _, relation := range relations {
		relationIds = utils.AddSlice(relationIds, relation.Id)
		namespace = relation.Namespace
	}
	return c.authAdminService.AddRelationshipsToPrincipal(
		context.Background(),
		c.Principal.OrganizationId,
		namespace,
		c.Principal.Id,
		relationIds...)
}

// DeleteRelations removes relations to principal.
func (c *PrincipalAdapter) DeleteRelations(relations ...*types.Relationship) error {
	var relationIds []string
	var namespace string
	for _, relation := range relations {
		relationIds = utils.AddSlice(relationIds, relation.Id)
		namespace = relation.Namespace
	}
	return c.authAdminService.DeleteRelationshipsToPrincipal(
		context.Background(),
		c.Principal.OrganizationId,
		namespace,
		c.Principal.Id,
		relationIds...)
}

// Relationships manages relationships.
func (c *PrincipalAdapter) Relationships(namespace string) *RelationshipAdapter {
	return &RelationshipAdapter{
		authAdminService: c.authAdminService,
		orgID:            c.Principal.OrganizationId,
		Relationship: &types.Relationship{
			Namespace:   namespace,
			PrincipalId: c.Principal.Id,
		},
		principalAdapter: c,
	}
}

// Authorizer builds authorizer adapter.
func (c *PrincipalAdapter) Authorizer(
	namespace string,
) *AuthorizerAdapter {
	return &AuthorizerAdapter{
		authorizer: c.authorizer,
		Principal:  c.Principal,
		namespace:  namespace,
	}
}

// GroupAdapter adapter for managing groups.
type GroupAdapter struct {
	authAdminService service.AuthAdminService
	orgID            string
	Group            *types.Group
}

// WithName setter.
func (c *GroupAdapter) WithName(
	name string,
) *GroupAdapter {
	c.Group.Name = name
	return c
}

// WithParents setter.
func (c *GroupAdapter) WithParents(
	parents ...*types.Group,
) *GroupAdapter {
	var parentIds []string
	for _, parent := range parents {
		parentIds = utils.AddSlice(parentIds, parent.Id)
	}
	c.Group.ParentIds = parentIds
	return c
}

// Create adapter.
func (c *GroupAdapter) Create() (*GroupAdapter, error) {
	var err error
	c.Group, err = c.authAdminService.CreateGroup(context.Background(), c.orgID, c.Group)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Update adapter.
func (c *GroupAdapter) Update() error {
	return c.authAdminService.UpdateGroup(context.Background(), c.orgID, c.Group)
}

// Get adapter.
func (c *GroupAdapter) Get(id string) error {
	group, err := c.authAdminService.GetGroup(context.Background(), c.orgID, c.Group.Namespace, id)
	if err != nil {
		return err
	}
	c.Group = group
	return nil
}

// AddRoles adapter for adding roles.
func (c *GroupAdapter) AddRoles(roles ...*types.Role) error {
	var roleIDs []string
	var namespace string
	for _, role := range roles {
		roleIDs = utils.AddSlice(roleIDs, role.Id)
		namespace = role.Namespace
	}
	return c.authAdminService.AddRolesToGroup(
		context.Background(),
		c.orgID,
		namespace,
		c.Group.Id,
		roleIDs...)
}

// DeleteRoles adapter for deleting roles.
func (c *GroupAdapter) DeleteRoles(roles ...*types.Role) error {
	var roleIDs []string
	var namespace string
	for _, role := range roles {
		roleIDs = utils.AddSlice(roleIDs, role.Id)
		namespace = role.Namespace
	}
	return c.authAdminService.DeleteRolesToGroup(
		context.Background(),
		c.orgID,
		namespace,
		c.Group.Id,
		roleIDs...)
}

// PermissionAdapter adapter for permissions.
type PermissionAdapter struct {
	authAdminService service.AuthAdminService
	orgID            string
	Permission       *types.Permission
}

// WithResource setter.
func (c *PermissionAdapter) WithResource(
	resource *types.Resource,
) *PermissionAdapter {
	c.Permission.ResourceId = resource.Id
	return c
}

// WithEffect setter.
func (c *PermissionAdapter) WithEffect(
	effect types.Effect,
) *PermissionAdapter {
	c.Permission.Effect = effect
	return c
}

// WithScope setter.
func (c *PermissionAdapter) WithScope(
	scope string,
) *PermissionAdapter {
	c.Permission.Scope = scope
	return c
}

// WithConstraints setter.
func (c *PermissionAdapter) WithConstraints(
	constraints string,
) *PermissionAdapter {
	c.Permission.Constraints = constraints
	return c
}

// WithActions setter.
func (c *PermissionAdapter) WithActions(
	actions ...string) *PermissionAdapter {
	c.Permission.Actions = actions
	return c
}

// Create adapter.
func (c *PermissionAdapter) Create() (*PermissionAdapter, error) {
	var err error
	c.Permission, err = c.authAdminService.CreatePermission(context.Background(), c.orgID, c.Permission)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Update adapter for updating permissions.
func (c *PermissionAdapter) Update() error {
	return c.authAdminService.UpdatePermission(context.Background(), c.orgID, c.Permission)
}

// Get adapter for fetching permission.
func (c *PermissionAdapter) Get(id string) error {
	permission, err := c.authAdminService.GetPermission(context.Background(), c.orgID, c.Permission.Namespace, id)
	if err != nil {
		return err
	}
	c.Permission = permission
	return nil
}

// RelationshipAdapter adapter for managing relationships.
type RelationshipAdapter struct {
	authAdminService service.AuthAdminService
	orgID            string
	Relationship     *types.Relationship
	principalAdapter *PrincipalAdapter
}

// WithRelation setter.
func (c *RelationshipAdapter) WithRelation(relation string) *RelationshipAdapter {
	c.Relationship.Relation = relation
	return c
}

// WithResource setter.
func (c *RelationshipAdapter) WithResource(resource *types.Resource) *RelationshipAdapter {
	c.Relationship.ResourceId = resource.Id
	return c
}

// WithAttributes setter.
func (c *RelationshipAdapter) WithAttributes(
	kv ...string,
) *RelationshipAdapter {
	c.Relationship.Attributes = utils.ArrayToMap(kv...)
	return c
}

// Create adapter for adding relationship.
func (c *RelationshipAdapter) Create() (*RelationshipAdapter, error) {
	var err error
	c.Relationship, err = c.authAdminService.CreateRelationship(context.Background(), c.orgID, c.Relationship)
	if err != nil {
		return nil, err
	}
	err = c.principalAdapter.AddRelations(c.Relationship)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Update adapter for updating relationship.
func (c *RelationshipAdapter) Update() error {
	return c.authAdminService.UpdateRelationship(context.Background(), c.orgID, c.Relationship)
}

// Get adapter for initializing relationship from the database.
func (c *RelationshipAdapter) Get(id string) error {
	relationship, err := c.authAdminService.GetRelationship(context.Background(), c.orgID, c.Relationship.Namespace, id)
	if err != nil {
		return err
	}
	c.Relationship = relationship
	return nil
}

// ResourceAdapter adapter for managing resources.
type ResourceAdapter struct {
	authAdminService service.AuthAdminService
	orgID            string
	Resource         *types.Resource
	constraints      string
	context          map[string]string
	expiry           time.Duration
}

// WithName setter.
func (c *ResourceAdapter) WithName(
	name string,
) *ResourceAdapter {
	c.Resource.Name = name
	return c
}

// WithCapacity setter.
func (c *ResourceAdapter) WithCapacity(
	capacity int,
) *ResourceAdapter {
	c.Resource.Capacity = int32(capacity)
	return c
}

// WithExpiration setter.
func (c *ResourceAdapter) WithExpiration(
	expiry time.Duration,
) *ResourceAdapter {
	c.expiry = expiry
	return c
}

// WithActions setter.
func (c *ResourceAdapter) WithActions(
	actions ...string) *ResourceAdapter {
	c.Resource.AllowedActions = actions
	return c
}

// WithConstraints setter.
func (c *ResourceAdapter) WithConstraints(
	constraints string) *ResourceAdapter {
	c.constraints = constraints
	return c
}

// WithContext setter.
func (c *ResourceAdapter) WithContext(
	kv ...string,
) *ResourceAdapter {
	c.context = utils.ArrayToMap(kv...)
	return c
}

// WithAttributes setter.
func (c *ResourceAdapter) WithAttributes(
	kv ...string,
) *ResourceAdapter {
	c.Resource.Attributes = utils.ArrayToMap(kv...)
	return c
}

// Create adapter for adding resource.
func (c *ResourceAdapter) Create() (*ResourceAdapter, error) {
	var err error
	c.Resource, err = c.authAdminService.CreateResource(context.Background(), c.orgID, c.Resource)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Update adapter for updating resource.
func (c *ResourceAdapter) Update() error {
	return c.authAdminService.UpdateResource(context.Background(), c.orgID, c.Resource)
}

// Allocate adapter for allocating resource.
func (c *ResourceAdapter) Allocate(principal *types.Principal) error {
	return c.authAdminService.AllocateResourceInstance(
		context.Background(),
		c.orgID,
		c.Resource.Namespace,
		c.Resource.Id,
		principal.Id,
		c.constraints,
		c.expiry,
		c.context,
	)
}

// Deallocate adapter for deallocating resource.
func (c *ResourceAdapter) Deallocate(principal *types.Principal) error {
	return c.authAdminService.DeallocateResourceInstance(
		context.Background(),
		c.orgID,
		c.Resource.Namespace,
		c.Resource.Id,
		principal.Id)
}

// AllocatedCount adapter for counting allocated resources.
func (c *ResourceAdapter) AllocatedCount() (allocated int32, err error) {
	_, allocated, err = c.authAdminService.CountResourceInstances(
		context.Background(),
		c.orgID,
		c.Resource.Namespace,
		c.Resource.Id,
	)
	return
}

// Get adapter for initializing resource.
func (c *ResourceAdapter) Get(id string) error {
	resource, err := c.authAdminService.GetResource(context.Background(), c.orgID, c.Resource.Namespace, id)
	if err != nil {
		return err
	}
	c.Resource = resource
	return nil
}

// RoleAdapter adapter for managing roles.
type RoleAdapter struct {
	authAdminService service.AuthAdminService
	orgID            string
	Role             *types.Role
}

// WithName setter.
func (c *RoleAdapter) WithName(
	name string,
) *RoleAdapter {
	c.Role.Name = name
	return c
}

// WithParents setter.
func (c *RoleAdapter) WithParents(
	parents ...*types.Role) *RoleAdapter {
	var parentIds []string
	for _, parent := range parents {
		parentIds = utils.AddSlice(parentIds, parent.Id)
	}
	c.Role.ParentIds = parentIds
	return c
}

// Create adapter.
func (c *RoleAdapter) Create() (*RoleAdapter, error) {
	var err error
	c.Role, err = c.authAdminService.CreateRole(context.Background(), c.orgID, c.Role)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Update adapter.
func (c *RoleAdapter) Update() error {
	return c.authAdminService.UpdateRole(context.Background(), c.orgID, c.Role)
}

// Get adapter.
func (c *RoleAdapter) Get(id string) error {
	role, err := c.authAdminService.GetRole(context.Background(), c.orgID, c.Role.Namespace, id)
	if err != nil {
		return err
	}
	c.Role = role
	return nil
}

// AddPermissions adapter for adding permissions.
func (c *RoleAdapter) AddPermissions(permissions ...*types.Permission) error {
	var permissionIds []string
	var namespace string
	for _, permission := range permissions {
		permissionIds = utils.AddSlice(permissionIds, permission.Id)
		namespace = permission.Namespace
	}
	return c.authAdminService.AddPermissionsToRole(
		context.Background(),
		c.orgID,
		namespace,
		c.Role.Id,
		permissionIds...)
}

// DeletePermissions adapter for deleting permissions.
func (c *RoleAdapter) DeletePermissions(permissions ...*types.Permission) error {
	var permissionIds []string
	var namespace string
	for _, permission := range permissions {
		permissionIds = utils.AddSlice(permissionIds, permission.Id)
		namespace = permission.Namespace
	}
	return c.authAdminService.DeletePermissionsToRole(
		context.Background(),
		c.orgID,
		namespace,
		c.Role.Id,
		permissionIds...)
}
