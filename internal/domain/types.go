package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"unsafe"
)

// ClientType alias
type ClientType string

const (
	// RootClientType admin access
	RootClientType = ClientType("root")

	// NobodyClientType without any access
	NobodyClientType = ClientType("nobody")

	// DefaultClientType with specified access
	DefaultClientType = ClientType("client")
)

// HashIndex for indexing
type HashIndex struct {
	Hash    string                 `json:"hash,omitempty"`
	Ids     []string               `json:"ids,omitempty"`
	Updated *timestamppb.Timestamp `json:"updated,omitempty"`
}

// NewHashIndex constructor
func NewHashIndex(hash string, ids []string) *HashIndex {
	return &HashIndex{
		Hash:    hash,
		Ids:     ids,
		Updated: timestamppb.Now(),
	}
}

// Validate helper
func (x *HashIndex) Validate() error {
	if len(x.Ids) == 0 {
		return NewValidationError(fmt.Sprintf("ids are not defined"))
	}
	return nil
}

// OrganizationExt that owns roles, groups, relations, and principals for a given namespace.
type OrganizationExt struct {
	Delegate *types.Organization
}

// NewOrganizationExt constructor
func NewOrganizationExt(delegate *types.Organization) *OrganizationExt {
	return &OrganizationExt{Delegate: delegate}
}

// Validate helper
func (x *OrganizationExt) Validate() error {
	if x.Delegate == nil {
		return NewValidationError(fmt.Sprintf("organization delegate is not defined"))
	}
	if x.Delegate.Name == "" {
		return NewValidationError(fmt.Sprintf("name is not defined"))
	}
	if len(x.Delegate.Namespaces) == 0 {
		return NewValidationError(fmt.Sprintf("namespaces are not defined"))
	}
	return nil
}

func (x *OrganizationExt) String() string {
	return x.Delegate.String()
}

// ResourceExt - The object that the principal wants to access (e.g., a file, a database record).
type ResourceExt struct {
	Delegate *types.Resource
}

// NewResourceExt constructor
func NewResourceExt(delegate *types.Resource) *ResourceExt {
	return &ResourceExt{Delegate: delegate}
}

// Validate helper
func (x *ResourceExt) Validate() error {
	if x.Delegate == nil {
		return NewValidationError(fmt.Sprintf("resource delegate is not defined"))
	}
	if x.Delegate.Namespace == "" {
		return NewValidationError(fmt.Sprintf("namespace is not defined"))
	}
	if x.Delegate.Name == "" {
		return NewValidationError(fmt.Sprintf("name is not defined"))
	}
	if len(x.Delegate.AllowedActions) == 0 {
		return NewValidationError(fmt.Sprintf("allowed_actions are not defined"))
	}

	return nil
}

// Hash calculator
func (x *ResourceExt) Hash() string {
	sort.Slice(x.Delegate.AllowedActions, func(i, j int) bool {
		return x.Delegate.AllowedActions[i] < x.Delegate.AllowedActions[j]
	})
	hasher := sha256.New()
	hasher.Write([]byte("resource"))
	hasher.Write([]byte(strings.ToLower(x.Delegate.Name)))
	for _, action := range x.Delegate.AllowedActions {
		hasher.Write([]byte(strings.ToLower(action)))
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func (x *ResourceExt) String() string {
	return x.Delegate.String()
}

// ResourceInstanceExt - instance of the resource for tracking quota of resource.
type ResourceInstanceExt struct {
	Delegate *types.ResourceInstance
}

// NewResourceInstanceExt constructor
func NewResourceInstanceExt(
	namespace string,
	resourceID string,
	principalID string) *ResourceInstanceExt {
	delegate := &types.ResourceInstance{
		Namespace:   namespace,
		ResourceId:  resourceID,
		PrincipalId: principalID,
		State:       types.ResourceState_ALLOCATED,
		Created:     timestamppb.Now(),
	}
	xInstance := &ResourceInstanceExt{Delegate: delegate}
	delegate.Id = xInstance.Hash()
	return xInstance
}

// Validate helper
func (x *ResourceInstanceExt) Validate() error {
	if x.Delegate == nil {
		return NewValidationError(fmt.Sprintf("resource-instance delegate is not defined"))
	}
	if x.Delegate.Namespace == "" {
		return NewValidationError(fmt.Sprintf("namespace is not defined"))
	}
	if x.Delegate.ResourceId == "" {
		return NewValidationError(fmt.Sprintf("resource_id is not defined"))
	}
	if x.Delegate.PrincipalId == "" {
		return NewValidationError(fmt.Sprintf("principal_id is not defined"))
	}

	return nil
}

// Hash calculator
func (x *ResourceInstanceExt) Hash() string {
	hasher := sha256.New()
	hasher.Write([]byte("resource_instance"))
	hasher.Write([]byte(x.Delegate.ResourceId))
	hasher.Write([]byte(x.Delegate.PrincipalId))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (x *ResourceInstanceExt) String() string {
	return x.Delegate.String()
}

// PermissionExt - An action that a principal is allowed to perform on a particular resource.
type PermissionExt struct {
	Delegate *types.Permission
}

// NewPermissionExt constructor
func NewPermissionExt(delegate *types.Permission) *PermissionExt {
	return &PermissionExt{Delegate: delegate}
}

// Validate helper
func (x *PermissionExt) Validate() error {
	if x.Delegate == nil {
		return NewValidationError(fmt.Sprintf("permission delegate is not defined"))
	}
	if x.Delegate.Namespace == "" {
		return NewValidationError(fmt.Sprintf("namespace is not defined"))
	}
	if x.Delegate.ResourceId == "" {
		return NewValidationError(fmt.Sprintf("resource_id is not defined"))
	}
	if len(x.Delegate.Actions) == 0 {
		return NewValidationError(fmt.Sprintf("actions are not defined"))
	}

	return nil
}

// Hash calculator
func (x *PermissionExt) Hash() string {
	sort.Slice(x.Delegate.Actions, func(i, j int) bool {
		return x.Delegate.Actions[i] < x.Delegate.Actions[j]
	})
	hasher := sha256.New()
	hasher.Write([]byte("permission"))
	hasher.Write([]byte(strings.ToLower(x.Delegate.Namespace)))
	hasher.Write([]byte(strings.ToLower(x.Delegate.Scope)))
	hasher.Write([]byte(x.Delegate.ResourceId))
	hasher.Write(unsafeCaseInt32ToBytes(int32(x.Delegate.Effect)))
	for _, action := range x.Delegate.Actions {
		hasher.Write([]byte(strings.ToLower(action)))
	}
	hasher.Write([]byte(x.Delegate.Constraints))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (x *PermissionExt) String() string {
	return x.Delegate.String()
}

// RoleExt - A named collection of permissions that can be assigned to a principal.
type RoleExt struct {
	Delegate *types.Role
}

// NewRoleExt constructor
func NewRoleExt(delegate *types.Role) *RoleExt {
	return &RoleExt{Delegate: delegate}
}

// Validate helper
func (x *RoleExt) Validate() error {
	if x.Delegate == nil {
		return NewValidationError(fmt.Sprintf("role delegate is not defined"))
	}
	if x.Delegate.Namespace == "" {
		return NewValidationError(fmt.Sprintf("namespace is not defined"))
	}
	if x.Delegate.Name == "" {
		return NewValidationError(fmt.Sprintf("name is not defined"))
	}

	return nil
}

// Hash calculator
func (x *RoleExt) Hash() string {
	hasher := sha256.New()
	hasher.Write([]byte("role"))
	hasher.Write([]byte(strings.ToLower(x.Delegate.Name)))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (x *RoleExt) String() string {
	return x.Delegate.String()
}

// GroupExt - A collection of principals that are treated as a single unit for the purpose of granting permissions.
type GroupExt struct {
	Delegate *types.Group
}

// NewGroupExt constructor
func NewGroupExt(delegate *types.Group) *GroupExt {
	return &GroupExt{Delegate: delegate}
}

// Validate helper
func (x *GroupExt) Validate() error {
	if x.Delegate == nil {
		return NewValidationError(fmt.Sprintf("group delegate is not defined"))
	}
	if x.Delegate.Namespace == "" {
		return NewValidationError(fmt.Sprintf("namespace is not defined"))
	}
	if x.Delegate.Name == "" {
		return NewValidationError(fmt.Sprintf("name is not defined"))
	}

	return nil
}

// Hash calculator
func (x *GroupExt) Hash() string {
	hasher := sha256.New()
	hasher.Write([]byte("group"))
	hasher.Write([]byte(strings.ToLower(x.Delegate.Name)))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (x *GroupExt) String() string {
	return x.Delegate.String()
}

// RelationshipExt - represents a relationship between a resource and a principal.
type RelationshipExt struct {
	Delegate *types.Relationship
}

// NewRelationshipExt constructor
func NewRelationshipExt(delegate *types.Relationship) *RelationshipExt {
	return &RelationshipExt{Delegate: delegate}
}

// Validate helper
func (x *RelationshipExt) Validate() error {
	if x.Delegate == nil {
		return NewValidationError(fmt.Sprintf("relationship delegate is not defined"))
	}
	if x.Delegate.Namespace == "" {
		return NewValidationError(fmt.Sprintf("namespace is not defined"))
	}
	if x.Delegate.Relation == "" {
		return NewValidationError(fmt.Sprintf("relation is not defined"))
	}
	if x.Delegate.PrincipalId == "" {
		return NewValidationError(fmt.Sprintf("principal_id is not defined"))
	}
	if x.Delegate.ResourceId == "" {
		return NewValidationError(fmt.Sprintf("resource_id is not defined"))
	}

	return nil
}

// Hash calculator
func (x *RelationshipExt) Hash() string {
	hasher := sha256.New()
	hasher.Write([]byte("relationship"))
	hasher.Write([]byte(x.Delegate.ResourceId))
	hasher.Write([]byte(x.Delegate.PrincipalId))
	hasher.Write([]byte(strings.ToLower(x.Delegate.Relation)))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (x *RelationshipExt) String() string {
	return x.Delegate.String()
}

// PrincipalExt - The entity (which could be a user, system, or another service) that is making the request.
type PrincipalExt struct {
	Delegate                  *types.Principal
	Organization              *types.Organization
	GroupsByName              map[string]*types.Group
	RolesByName               map[string]*types.Role
	RelationsById             map[string]*types.Relationship
	ResourcesById             map[string]*types.Resource
	PermissionsByResourceName map[string]map[string]*types.Permission
}

// NewPrincipalExt constructor
func NewPrincipalExt(delegate *types.Principal) *PrincipalExt {
	return &PrincipalExt{
		Delegate:                  delegate,
		GroupsByName:              make(map[string]*types.Group),
		RolesByName:               make(map[string]*types.Role),
		RelationsById:             make(map[string]*types.Relationship),
		ResourcesById:             make(map[string]*types.Resource),
		PermissionsByResourceName: make(map[string]map[string]*types.Permission)}
}

func NewPrincipalExtFromResponse(
	res *services.GetPrincipalResponse,
) *PrincipalExt {
	principal := &types.Principal{
		Id:             res.Id,
		Version:        res.Version,
		Namespaces:     res.Namespaces,
		OrganizationId: res.OrganizationId,
		Username:       res.Username,
		Email:          res.Email,
		Name:           res.Name,
		Attributes:     res.Attributes,
		GroupIds:       res.GroupIds,
		RoleIds:        res.RoleIds,
		PermissionIds:  res.PermissionIds,
		Created:        res.Created,
		Updated:        res.Updated,
	}
	principalExt := NewPrincipalExt(principal)
	principalExt.Organization = &types.Organization{Id: res.OrganizationId, Namespaces: res.Namespaces}

	for _, group := range res.Groups {
		principalExt.GroupsByName[group.Name] = group
	}
	for _, role := range res.Roles {
		principalExt.RolesByName[role.Name] = role
	}
	for _, relation := range res.Relations {
		principal.RelationIds = append(principal.RelationIds, relation.Id)
		principalExt.RelationsById[relation.Id] = relation
	}
	for _, resource := range res.Resources {
		principalExt.ResourcesById[resource.Id] = resource
		permsForResource := make(map[string]*types.Permission)
		for _, perm := range res.Permissions {
			if perm.ResourceId == resource.Id {
				permsForResource[perm.Id] = perm
			}
		}
		principalExt.PermissionsByResourceName[resource.Name] = permsForResource
	}
	return principalExt
}

// ToMap helper
func (x *PrincipalExt) ToMap(
	req *services.AuthRequest,
	resource *types.Resource,
) (res map[string]any) {
	res = make(map[string]any)
	principalMap := make(map[string]any)
	principalMap["Id"] = x.Delegate.Id
	principalMap["OrganizationId"] = x.Delegate.OrganizationId
	principalMap["Username"] = x.Delegate.Username
	principalMap["Name"] = x.Delegate.Name
	principalMap["Email"] = x.Delegate.Email
	principalMap["Groups"] = x.GroupNames()
	principalMap["Roles"] = x.RoleNames()
	principalMap["Resource"] = resource.Name
	resourceMap := make(map[string]any)
	for k, v := range resource.Attributes {
		resourceMap[k] = v
	}
	resourceMap["Name"] = resource.Name
	resourceMap["Capacity"] = resource.Capacity
	resourceMap["AllowedActions"] = resource.AllowedActions
	res["Resource"] = resourceMap

	principalMap["Action"] = req.Action
	principalMap["Scope"] = req.Scope
	for k, v := range x.Delegate.Attributes {
		principalMap[k] = v
	}
	res["Principal"] = principalMap
	for k, v := range req.Context {
		res[k] = v
	}
	relations := make(map[string]map[string]string)
	for _, relation := range x.RelationsByResource(resource.Id) {
		attrs := make(map[string]string)
		attrs["Name"] = relation.Relation
		for k, v := range relation.Attributes {
			attrs[k] = v
		}
		relations[relation.Relation] = attrs
	}
	principalMap["Relations"] = relations
	res["Relations"] = relations
	return
}

// Validate helper
func (x *PrincipalExt) Validate() error {
	if x.Delegate == nil {
		return NewValidationError(fmt.Sprintf("principal delegate is not defined"))
	}
	if x.Delegate.Username == "" {
		return NewValidationError(fmt.Sprintf("username is not defined"))
	}
	if x.Delegate.OrganizationId == "" {
		return NewValidationError(fmt.Sprintf("organization_id is not defined"))
	}
	if len(x.Delegate.Namespaces) == 0 {
		return NewValidationError(fmt.Sprintf("namespaces are not defined"))
	}
	if len(x.Delegate.Attributes) > 255 {
		return NewValidationError(fmt.Sprintf("too many attributes are defined"))
	}

	return nil
}

// AddPermission helper
func (x *PrincipalExt) AddPermission(perm *types.Permission) error {
	resource := x.ResourcesById[perm.ResourceId]
	if resource == nil {
		return NewNotFoundError(
			fmt.Sprintf("failed to add permission, resource %s not found for permission %s",
				perm.ResourceId, perm.Id))
	}
	permsForResource := x.PermissionsByResourceName[resource.Name]
	if permsForResource == nil {
		permsForResource = make(map[string]*types.Permission)
		x.PermissionsByResourceName[resource.Name] = permsForResource
	}
	permsForResource[perm.Id] = perm
	return nil
}

// Roles Getter
func (x *PrincipalExt) Roles() (res []*types.Role) {
	for _, r := range x.RolesByName {
		res = append(res, r)
	}
	return
}

// RoleNames Getter
func (x *PrincipalExt) RoleNames() (res []string) {
	for k := range x.RolesByName {
		res = append(res, k)
	}
	return
}

// GroupNames Getter
func (x *PrincipalExt) GroupNames() (res []string) {
	for k := range x.GroupsByName {
		res = append(res, k)
	}
	return
}

// Groups Getter
func (x *PrincipalExt) Groups() (res []*types.Group) {
	for _, g := range x.GroupsByName {
		res = append(res, g)
	}
	return
}

// Resources Getter
func (x *PrincipalExt) Resources() (res []*types.Resource) {
	for _, r := range x.ResourcesById {
		res = append(res, r)
	}
	return
}

// ResourceNames Getter
func (x *PrincipalExt) ResourceNames() (names []string) {
	for _, res := range x.ResourcesById {
		names = append(names, res.Name)
	}
	return
}

// ResourcesByPartialNameAndAction Getter
func (x *PrincipalExt) ResourcesByPartialNameAndAction(resourceName string, action string) (arr []*types.Resource) {
	for _, res := range x.ResourcesById {
		if !utils.Includes(res.AllowedActions, action) {
			continue
		}
		if res.Name == resourceName || (res.Wildcard && doesResourceNameMatches(res.Name, resourceName)) {
			arr = append(arr, res)
		}
	}
	return arr
}

// ResourceByName Getter
func (x *PrincipalExt) ResourceByName(resourceName string) *types.Resource {
	for _, res := range x.ResourcesById {
		if res.Name == resourceName {
			return res
		}
	}
	return nil
}

// RelationNamesByResourceName Getter
func (x *PrincipalExt) RelationNamesByResourceName(resourceName string) (res []string) {
	resource := x.ResourceByName(resourceName)
	if resource != nil {
		for _, rel := range x.RelationsById {
			if rel.ResourceId == resource.Id {
				res = append(res, rel.Relation)
			}
		}
	}
	return
}

// Relations Getter
func (x *PrincipalExt) Relations() (res []*types.Relationship) {
	for _, rel := range x.RelationsById {
		res = append(res, rel)
	}
	return
}

// RelationsByResource Getter
func (x *PrincipalExt) RelationsByResource(resourceID string) (res []*types.Relationship) {
	for _, rel := range x.RelationsById {
		if rel.ResourceId == resourceID {
			res = append(res, rel)
		}
	}
	return
}

// RelationNames Getter
func (x *PrincipalExt) RelationNames(resourceID string) (res []string) {
	for _, rel := range x.RelationsById {
		if rel.ResourceId == resourceID {
			res = append(res, rel.Relation)
		}
	}
	return
}

// AllPermissions Getter
func (x *PrincipalExt) AllPermissions() (res []*types.Permission) {
	for _, permMap := range x.PermissionsByResourceName {
		for _, perm := range permMap {
			res = append(res, perm)
		}
	}
	return
}

// LatestGroupRoleDate helper
func (x *PrincipalExt) LatestGroupRoleDate() (latestGroupRoleDate *timestamppb.Timestamp) {
	for _, group := range x.GroupsByName {
		if latestGroupRoleDate == nil || group.Updated.Seconds > latestGroupRoleDate.Seconds {
			latestGroupRoleDate = group.Updated
		}
	}

	for _, role := range x.RolesByName {
		if latestGroupRoleDate == nil || role.Updated.Seconds > latestGroupRoleDate.Seconds {
			latestGroupRoleDate = role.Updated
		}
	}
	return
}

// ToGetPrincipalResponse helper
func (x *PrincipalExt) ToGetPrincipalResponse() *services.GetPrincipalResponse {
	return &services.GetPrincipalResponse{
		Id:             x.Delegate.Id,
		Version:        x.Delegate.Version,
		OrganizationId: x.Delegate.OrganizationId,
		Namespaces:     x.Delegate.Namespaces,
		Username:       x.Delegate.Username,
		Email:          x.Delegate.Email,
		Name:           x.Delegate.Name,
		Attributes:     x.Delegate.Attributes,
		Groups:         x.Groups(),
		Roles:          x.Roles(),
		Resources:      x.Resources(),
		Permissions:    x.AllPermissions(),
		Relations:      x.Relations(),
		GroupIds:       x.Delegate.GroupIds,
		RoleIds:        x.Delegate.RoleIds,
		PermissionIds:  x.Delegate.PermissionIds,
		Created:        x.Delegate.Created,
		Updated:        x.Delegate.Updated,
	}
}

// Hash calculator
func (x *PrincipalExt) Hash() string {
	hasher := sha256.New()
	hasher.Write([]byte("principal"))
	hasher.Write([]byte(x.Delegate.OrganizationId))
	hasher.Write([]byte(strings.ToLower(x.Delegate.Username)))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GroupHashIndex calculator
func (x *PrincipalExt) GroupHashIndex() string {
	hasher := sha256.New()
	hasher.Write([]byte("groups-index"))
	hasher.Write([]byte(strings.ToLower(x.Delegate.Id)))
	return hex.EncodeToString(hasher.Sum(nil))
}

// RoleHashIndex calculator
func (x *PrincipalExt) RoleHashIndex() string {
	hasher := sha256.New()
	hasher.Write([]byte("roles-index"))
	hasher.Write([]byte(strings.ToLower(x.Delegate.Id)))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (x *PrincipalExt) String() string {
	return x.Delegate.String()
}

func (x *PrincipalExt) CheckConstraints(
	req *services.AuthRequest,
	resource *types.Resource,
	constraints string) (bool, string, error) {

	tf, err := ParseTemplate(constraints, x, resource, req)
	if err != nil {
		return false, "", err
	}
	output := string(tf)
	return output == "true", output, nil
}

func (x *PrincipalExt) CheckPermission(
	req *services.AuthRequest,
) (res *services.AuthResponse, err error) {
	resources := x.ResourcesByPartialNameAndAction(req.Resource, req.Action)
	if len(resources) == 0 {
		return nil, NewAuthError(
			fmt.Sprintf("resource %s not found with action %s, available resources %v",
				req.Resource, req.Action, x.ResourceNames()))
	}
	actionMatched := false
	constraintsFailed := false
	effects := make(map[types.Effect]int)
	var permissionIds []string
	for _, resource := range resources {
		res = &services.AuthResponse{}
		perms := x.PermissionsByResourceName[resource.Name]
		for _, perm := range perms {
			matched := false
			for _, permAction := range perm.Actions {
				if (perm.Scope == "*" || perm.Scope == req.Scope) &&
					(permAction == "*" || permAction == req.Action) {
					matched = true
					actionMatched = true
					break
				}
			}
			if matched {
				if perm.Constraints == "" {
					res.Effect = perm.Effect
					effects[perm.Effect]++
					permissionIds = append(permissionIds, perm.Id)
				} else {
					matched, tf, err := x.CheckConstraints(
						req, resource, perm.Constraints)
					if err != nil {
						return res, err
					}
					if matched {
						res.Effect = perm.Effect
						effects[perm.Effect]++
						permissionIds = append(permissionIds, perm.Id)
					} else {
						if log.IsLevelEnabled(log.DebugLevel) {
							log.WithFields(log.Fields{
								"Component":  "PrincipalExt",
								"Output":     tf,
								"Permission": perm,
								"Resource":   resource,
								"Context":    req.Context,
							}).
								Infof("template failed to match")
						}
						constraintsFailed = true
					}
				}
			}
		}
	}
	// Checking matched permissions
	if len(effects) == 0 {
		return res, NewAuthError(fmt.Sprintf("no permissions[%d/%v/%v] matched for %s %s",
			len(permissionIds), actionMatched, constraintsFailed, req.Resource, req.Action))
	} else if len(effects) == 1 {
		if len(permissionIds) > 1 {
			res.Message = fmt.Sprintf("multiple permissions matched %v [%s]",
				permissionIds, MultiplePermissionsMatchedCode)
		}
	} else {
		res.Effect = types.Effect_DENIED // if both permit and deny found then treat as denied
		res.Message = fmt.Sprintf("conflicting permissions [%d %v] found for %s %s [%s]",
			len(permissionIds), effects, req.Resource, req.Action, ConflictingPermissionsCode)
	}
	return res, nil
}

// BytesInInt32 constant
const BytesInInt32 = 4

func unsafeCaseInt32ToBytes(val int32) []byte {
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&val)), Len: BytesInInt32, Cap: BytesInInt32}
	return *(*[]byte)(unsafe.Pointer(&hdr))
}

const NextOffsetHeader = "X-Next-Offset"

// doesResourceNameMatches uses regex to match resource
func doesResourceNameMatches(patternResourceName string, reqResourceName string) bool {
	patternResourceName = strings.ReplaceAll(patternResourceName, "*", ".*")
	re, err := regexp.Compile(patternResourceName)
	if err != nil {
		return false
	}
	return re.MatchString(reqResourceName)
}
