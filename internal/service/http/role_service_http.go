package http

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/web"
)

// RoleServiceHTTP - manages persistence of roles data
type RoleServiceHTTP struct {
	*baseHTTPClient
}

// NewRoleServiceHTTP manages persistence of roles data
func NewRoleServiceHTTP(
	client web.HTTPClient,
	baseURL string,
) *RoleServiceHTTP {
	return &RoleServiceHTTP{
		baseHTTPClient: &baseHTTPClient{
			client:  client,
			baseURL: baseURL,
		},
	}
}

// CreateRole - creates a new role
func (h *RoleServiceHTTP) CreateRole(
	ctx context.Context,
	organizationID string,
	role *types.Role) (*types.Role, error) {
	if organizationID == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	req := &services.CreateRoleRequest{
		OrganizationId: organizationID,
		Namespace:      role.Namespace,
		Name:           role.Name,
		ParentIds:      role.ParentIds,
	}
	res := &services.CreateRoleResponse{}
	_, _, err := h.post(ctx,
		fmt.Sprintf("/api/v1/%s/%s/roles", organizationID, role.Namespace),
		req,
		res,
	)
	if err != nil {
		return nil, err
	}
	role.Id = res.Id
	return role, nil
}

// UpdateRole - updates an existing role
func (h *RoleServiceHTTP) UpdateRole(
	ctx context.Context,
	organizationID string,
	role *types.Role) error {
	xRole := domain.NewRoleExt(role)
	if err := xRole.Validate(); err != nil {
		return err
	}
	if role.Id == "" {
		return domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	req := &services.UpdateRoleRequest{
		OrganizationId: organizationID,
		Namespace:      role.Namespace,
		Id:             role.Id,
		Name:           role.Name,
		ParentIds:      role.ParentIds,
	}
	res := &services.UpdateRoleResponse{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/%s/%s/roles/%s", organizationID, role.Namespace, role.Id),
		req,
		res,
	)
	return err
}

// DeleteRole removes role
func (h *RoleServiceHTTP) DeleteRole(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	if namespace == "" {
		return domain.NewValidationError(
			fmt.Sprintf("namespace is not defined"))
	}
	if id == "" {
		return domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	_, _, err := h.del(ctx,
		fmt.Sprintf("/api/v1/%s/%s/roles/%s", organizationID, namespace, id),
	)
	return err
}

// GetRole - finds role
func (h *RoleServiceHTTP) GetRole(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string,
) (*types.Role, error) {
	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	roles, _, err := h.GetRoles(
		ctx,
		organizationID,
		namespace,
		map[string]string{"id": id},
		"",
		1,
	)
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return nil, domain.NewNotFoundError(fmt.Sprintf("role with id  %s is not found", id))
	}
	return roles[0], nil
}

// GetRoles - queries roles
func (h *RoleServiceHTTP) GetRoles(
	ctx context.Context,
	organizationID string,
	namespace string,
	predicates map[string]string,
	offset string,
	limit int64) (arr []*types.Role, nextOffset string, err error) {
	if organizationID == "" {
		return nil, "", domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	if namespace == "" {
		return nil, "", domain.NewValidationError(
			fmt.Sprintf("namespace is not defined"))
	}
	if predicates == nil {
		predicates = make(map[string]string)
	}
	res := &[]services.QueryRoleResponse{}
	predicates["offset"] = offset
	predicates["limit"] = fmt.Sprintf("%d", limit)
	_, resHeaders, err := h.get(
		ctx,
		fmt.Sprintf("/api/v1/%s/%s/roles", organizationID, namespace),
		predicates,
		res,
	)
	if err != nil {
		return nil, "", err
	}
	for _, roleRes := range *res {
		arr = append(arr, &types.Role{
			Id:            roleRes.Id,
			Version:       roleRes.Version,
			Namespace:     roleRes.Namespace,
			Name:          roleRes.Name,
			PermissionIds: roleRes.PermissionIds,
			ParentIds:     roleRes.ParentIds,
			Created:       roleRes.Created,
			Updated:       roleRes.Updated,
		})
	}
	nextOffset = resHeaders[domain.NextOffsetHeader]
	return
}

// AddPermissionsToRole helper
func (h *RoleServiceHTTP) AddPermissionsToRole(
	ctx context.Context,
	organizationID string,
	namespace string,
	roleID string,
	permissionIds ...string,
) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	if namespace == "" {
		return domain.NewValidationError(
			fmt.Sprintf("namespace is not defined"))
	}
	if roleID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("role-id is not defined"))
	}
	if len(permissionIds) == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("permission-ids is not defined"))
	}
	req := &services.AddPermissionsToRoleRequest{
		OrganizationId: organizationID,
		Namespace:      namespace,
		RoleId:         roleID,
		PermissionIds:  permissionIds,
	}
	res := &services.AddPermissionsToRoleResponse{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/%s/%s/roles/%s/permissions/add", organizationID, namespace, roleID),
		req,
		res,
	)
	return err
}

// DeletePermissionsToRole helper
func (h *RoleServiceHTTP) DeletePermissionsToRole(
	ctx context.Context,
	organizationID string,
	namespace string,
	roleID string,
	permissionIds ...string,
) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	if namespace == "" {
		return domain.NewValidationError(
			fmt.Sprintf("namespace is not defined"))
	}
	if roleID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("role-id is not defined"))
	}
	if len(permissionIds) == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("permission-ids is not defined"))
	}
	req := &services.DeletePermissionsToRoleRequest{
		OrganizationId: organizationID,
		Namespace:      namespace,
		RoleId:         roleID,
		PermissionIds:  permissionIds,
	}
	res := &services.DeletePermissionsToRoleResponse{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/%s/%s/roles/%s/permissions/delete", organizationID, namespace, roleID),
		req,
		res,
	)
	return err
}
