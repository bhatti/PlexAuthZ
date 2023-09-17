package domain

import (
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/utils"
	"github.com/twinj/uuid"
)

// OrganizationBuilder that owns roles, groups, relations, and principals for a given namespace.
type OrganizationBuilder struct {
	// ID unique identifier assigned to this organization.
	Id string
	// Name of organization.
	Name string
	// Allowed Namespaces for organization.
	Namespaces []string
	// url for organization.
	Url string
}

// NewOrganizationBuilder constructor
func NewOrganizationBuilder() *OrganizationBuilder {
	return &OrganizationBuilder{}
}

// WithId setter
func (b *OrganizationBuilder) WithId(id string) *OrganizationBuilder {
	b.Id = id
	return b
}

// WithName setter
func (b *OrganizationBuilder) WithName(name string) *OrganizationBuilder {
	b.Name = name
	return b
}

// WithNamespaces setter
func (b *OrganizationBuilder) WithNamespaces(namespaces ...string) *OrganizationBuilder {
	b.Namespaces = namespaces
	return b
}

// WithUrl setter
func (b *OrganizationBuilder) WithUrl(url string) *OrganizationBuilder {
	b.Url = url
	return b
}

// Build helper
func (b *OrganizationBuilder) Build() (*types.Organization, error) {
	org := &types.Organization{
		Id:         b.Id,
		Name:       b.Name,
		Namespaces: b.Namespaces,
		Url:        b.Url,
	}
	if err := NewOrganizationExt(org).Validate(); err != nil {
		return nil, err
	}
	return org, nil
}

// ResourceBuilder - The object that the principal wants to access (e.g., a file, a database record).
type ResourceBuilder struct {
	// Namespace of resource.
	Namespace string
	// Name of the resource.
	Name string
	// capacity of resource.
	Capacity int32
	// Attributes of resource.
	Attributes map[string]string
	// AllowedActions that can be performed.
	AllowedActions []string
}

// NewResourceBuilder constructor
func NewResourceBuilder() *ResourceBuilder {
	return &ResourceBuilder{Attributes: make(map[string]string)}
}

// WithNamespace setter
func (b *ResourceBuilder) WithNamespace(namespace string) *ResourceBuilder {
	b.Namespace = namespace
	return b
}

// WithName setter
func (b *ResourceBuilder) WithName(name string) *ResourceBuilder {
	b.Name = name
	return b
}

// WithCapacity setter
func (b *ResourceBuilder) WithCapacity(capacity int) *ResourceBuilder {
	b.Capacity = int32(capacity)
	return b
}

// WithAttribute setter
func (b *ResourceBuilder) WithAttribute(name string, val string) *ResourceBuilder {
	b.Attributes[name] = val
	return b
}

// WithAllowedActions setter
func (b *ResourceBuilder) WithAllowedActions(actions ...string) *ResourceBuilder {
	b.AllowedActions = utils.AddSlice(b.AllowedActions, actions...)
	return b
}

// Build helper
func (b *ResourceBuilder) Build() (*types.Resource, error) {
	resource := &types.Resource{
		Id:             uuid.NewV4().String(),
		Namespace:      b.Namespace,
		Name:           b.Name,
		Capacity:       b.Capacity,
		Attributes:     b.Attributes,
		AllowedActions: b.AllowedActions,
	}
	if err := NewResourceExt(resource).Validate(); err != nil {
		return nil, err
	}
	return resource, nil
}

// PermissionBuilder - An action that a principal is allowed to perform on a particular resource.
type PermissionBuilder struct {
	// Namespace of permission
	Namespace string
	// Scope for permission.
	Scope string
	// Actions that can be performed.
	Actions []string
	// Resource for the action.
	ResourceId string
	// Effect Permitted or Denied
	Effect types.Effect
	// Constraints expression with dynamic properties.
	Constraints string
}

// NewPermissionBuilder constructor
func NewPermissionBuilder() *PermissionBuilder {
	return &PermissionBuilder{Effect: types.Effect_PERMITTED}
}

// WithNamespace setter
func (b *PermissionBuilder) WithNamespace(namespace string) *PermissionBuilder {
	b.Namespace = namespace
	return b
}

// WithScope setter
func (b *PermissionBuilder) WithScope(scope string) *PermissionBuilder {
	b.Scope = scope
	return b
}

// WithActions setter
func (b *PermissionBuilder) WithActions(actions ...string) *PermissionBuilder {
	b.Actions = utils.AddSlice(b.Actions, actions...)
	return b
}

// WithResourceId setter
func (b *PermissionBuilder) WithResourceId(resourceID string) *PermissionBuilder {
	b.ResourceId = resourceID
	return b
}

// WithEffect setter
func (b *PermissionBuilder) WithEffect(effect types.Effect) *PermissionBuilder {
	b.Effect = effect
	return b
}

// WithConstraints setter
func (b *PermissionBuilder) WithConstraints(constraints string) *PermissionBuilder {
	b.Constraints = constraints
	return b
}

// Build helper
func (b *PermissionBuilder) Build() (*types.Permission, error) {
	permission := &types.Permission{
		Id:          uuid.NewV4().String(),
		Namespace:   b.Namespace,
		Scope:       b.Scope,
		Actions:     b.Actions,
		ResourceId:  b.ResourceId,
		Effect:      b.Effect,
		Constraints: b.Constraints,
	}
	if err := NewPermissionExt(permission).Validate(); err != nil {
		return nil, err
	}
	return permission, nil
}

// RoleBuilder - A named collection of permissions that can be assigned to a principal.
type RoleBuilder struct {
	// Namespace of role.
	Namespace string
	// Name of role
	Name string
	// PermissionIDs that can be performed.
	PermissionIds []string
	// Optional parent ids
	ParentIds []string
}

// NewRoleBuilder constructor
func NewRoleBuilder() *RoleBuilder {
	return &RoleBuilder{}
}

// WithNamespace setter
func (b *RoleBuilder) WithNamespace(namespace string) *RoleBuilder {
	b.Namespace = namespace
	return b
}

// WithName setter
func (b *RoleBuilder) WithName(name string) *RoleBuilder {
	b.Name = name
	return b
}

// WithParentIds setter
func (b *RoleBuilder) WithParentIds(ids ...string) *RoleBuilder {
	b.ParentIds = utils.AddSlice(b.ParentIds, ids...)
	return b
}

// Build helper
func (b *RoleBuilder) Build() (*types.Role, error) {
	role := &types.Role{
		Id:            uuid.NewV4().String(),
		Namespace:     b.Namespace,
		Name:          b.Name,
		PermissionIds: b.PermissionIds,
		ParentIds:     b.ParentIds,
	}
	if err := NewRoleExt(role).Validate(); err != nil {
		return nil, err
	}
	return role, nil
}

// GroupBuilder - A collection of principals that are treated as a single unit for the purpose of granting permissions.
type GroupBuilder struct {
	// Namespace of group.
	Namespace string
	// Name of the group.
	Name string
	// RoleIDs that are associated.
	RoleIds []string
	// Optional parent ids.
	ParentIds []string
}

// NewGroupBuilder constructor
func NewGroupBuilder() *GroupBuilder {
	return &GroupBuilder{}
}

// WithNamespace setter
func (b *GroupBuilder) WithNamespace(namespace string) *GroupBuilder {
	b.Namespace = namespace
	return b
}

// WithName setter
func (b *GroupBuilder) WithName(name string) *GroupBuilder {
	b.Name = name
	return b
}

// WithParentIds setter
func (b *GroupBuilder) WithParentIds(ids ...string) *GroupBuilder {
	b.ParentIds = utils.AddSlice(b.ParentIds, ids...)
	return b
}

// Build helper
func (b *GroupBuilder) Build() (*types.Group, error) {
	group := &types.Group{
		Id:        uuid.NewV4().String(),
		Namespace: b.Namespace,
		Name:      b.Name,
		RoleIds:   b.RoleIds,
		ParentIds: b.ParentIds,
	}
	if err := NewGroupExt(group).Validate(); err != nil {
		return nil, err
	}
	return group, nil
}

// RelationshipBuilder - represents a relationship between a resource and a principal.
type RelationshipBuilder struct {
	// Namespace of relationship.
	Namespace string
	// Relation name.
	Relation string
	// PrincipalID for relationship.
	PrincipalId string
	// ResourceID for relationship.
	ResourceId string
	// Attributes of relationship.
	Attributes map[string]string
}

// NewRelationshipBuilder constructor
func NewRelationshipBuilder() *RelationshipBuilder {
	return &RelationshipBuilder{
		Attributes: make(map[string]string),
	}
}

// WithNamespace setter
func (b *RelationshipBuilder) WithNamespace(namespace string) *RelationshipBuilder {
	b.Namespace = namespace
	return b
}

// WithRelation setter
func (b *RelationshipBuilder) WithRelation(relation string) *RelationshipBuilder {
	b.Relation = relation
	return b
}

// WithPrincipalId setter
func (b *RelationshipBuilder) WithPrincipalId(id string) *RelationshipBuilder {
	b.PrincipalId = id
	return b
}

// WithResourceId setter
func (b *RelationshipBuilder) WithResourceId(id string) *RelationshipBuilder {
	b.ResourceId = id
	return b
}

// WithAttribute setter
func (b *RelationshipBuilder) WithAttribute(name string, val string) *RelationshipBuilder {
	b.Attributes[name] = val
	return b
}

// Build helper
func (b *RelationshipBuilder) Build() (*types.Relationship, error) {
	relation := &types.Relationship{
		Id:          uuid.NewV4().String(),
		Namespace:   b.Namespace,
		Relation:    b.Relation,
		PrincipalId: b.PrincipalId,
		ResourceId:  b.ResourceId,
		Attributes:  b.Attributes,
	}
	if err := NewRelationshipExt(relation).Validate(); err != nil {
		return nil, err
	}
	return relation, nil
}

// PrincipalBuilder - The entity (which could be a user, system, or another service) that is making the request.
type PrincipalBuilder struct {
	// OrganizationId of the principal user.
	OrganizationId string
	// Allowed Namespaces for organization.
	Namespaces []string
	// Username of the principal user.
	Username string
	// Name of the principal user.
	Name string
	// Email of the principal user.
	Email string
	// Attributes of principal
	Attributes map[string]string
	// Groups that the principal belongs to.
	GroupIds []string
	// Roles that the principal belongs to.
	RoleIds []string
	// Permissions that the principal belongs to.
	PermissionIds []string
	// Relationships that the principal belongs to.
	RelationIds []string
}

// NewPrincipalBuilder constructor
func NewPrincipalBuilder() *PrincipalBuilder {
	return &PrincipalBuilder{Attributes: make(map[string]string)}
}

// WithOrganizationId setter
func (b *PrincipalBuilder) WithOrganizationId(id string) *PrincipalBuilder {
	b.OrganizationId = id
	return b
}

// WithUsername setter
func (b *PrincipalBuilder) WithUsername(username string) *PrincipalBuilder {
	b.Username = username
	return b
}

// WithEmail setter
func (b *PrincipalBuilder) WithEmail(email string) *PrincipalBuilder {
	b.Email = email
	return b
}

// WithName setter
func (b *PrincipalBuilder) WithName(name string) *PrincipalBuilder {
	b.Name = name
	return b
}

// WithAttribute setter
func (b *PrincipalBuilder) WithAttribute(name string, val string) *PrincipalBuilder {
	b.Attributes[name] = val
	return b
}

// WithNamespaces setter
func (b *PrincipalBuilder) WithNamespaces(namespaces ...string) *PrincipalBuilder {
	b.Namespaces = namespaces
	return b
}

// Build helper
func (b *PrincipalBuilder) Build() (*types.Principal, error) {
	principal := &types.Principal{
		Id:             uuid.NewV4().String(),
		OrganizationId: b.OrganizationId,
		Namespaces:     b.Namespaces,
		Username:       b.Username,
		Name:           b.Name,
		Email:          b.Email,
		Attributes:     b.Attributes,
		GroupIds:       b.GroupIds,
		RoleIds:        b.RoleIds,
		PermissionIds:  b.PermissionIds,
		RelationIds:    b.RelationIds,
	}
	if err := NewPrincipalExt(principal).Validate(); err != nil {
		return nil, err
	}
	return principal, nil
}
