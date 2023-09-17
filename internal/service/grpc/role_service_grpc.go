package grpc

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/server"
)

// RoleServiceGrpc - manages persistence of roles data
type RoleServiceGrpc struct {
	clients server.Clients
}

// NewRoleServiceGrpc manages persistence of roles data
func NewRoleServiceGrpc(
	clients server.Clients,
) *RoleServiceGrpc {
	return &RoleServiceGrpc{
		clients: clients,
	}
}

// CreateRole - creates a new role
func (s *RoleServiceGrpc) CreateRole(
	ctx context.Context,
	organizationID string,
	role *types.Role) (*types.Role, error) {
	if organizationID == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	res, err := s.clients.RolesClient.Create(
		ctx,
		&services.CreateRoleRequest{
			OrganizationId: organizationID,
			Namespace:      role.Namespace,
			Name:           role.Name,
			ParentIds:      role.ParentIds,
		})
	if err != nil {
		return nil, err
	}
	role.Id = res.Id
	return role, nil
}

// UpdateRole - updates an existing role
func (s *RoleServiceGrpc) UpdateRole(
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
	_, err := s.clients.RolesClient.Update(
		ctx,
		&services.UpdateRoleRequest{
			OrganizationId: organizationID,
			Namespace:      role.Namespace,
			Id:             role.Id,
			Name:           role.Name,
			ParentIds:      role.ParentIds,
		})
	return err
}

// DeleteRole removes role
func (s *RoleServiceGrpc) DeleteRole(
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
	_, err := s.clients.RolesClient.Delete(
		ctx,
		&services.DeleteRoleRequest{
			OrganizationId: organizationID,
			Namespace:      namespace,
			Id:             id,
		})
	return err
}

// GetRole - finds role
func (s *RoleServiceGrpc) GetRole(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string,
) (*types.Role, error) {
	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	roles, _, err := s.GetRoles(
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
func (s *RoleServiceGrpc) GetRoles(
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
	res, err := s.clients.RolesClient.Query(
		ctx,
		&services.QueryRoleRequest{
			OrganizationId: organizationID,
			Namespace:      namespace,
			Predicates:     predicates,
			Offset:         offset,
			Limit:          limit,
		})
	if err != nil {
		return nil, "", err
	}
	for {
		roleRes, err := res.Recv()
		if err != nil {
			break
		}
		nextOffset = roleRes.NextOffset
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
	return
}

// AddPermissionsToRole helper
func (s *RoleServiceGrpc) AddPermissionsToRole(
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
	_, err := s.clients.RolesClient.AddPermissions(ctx, &services.AddPermissionsToRoleRequest{
		OrganizationId: organizationID,
		Namespace:      namespace,
		RoleId:         roleID,
		PermissionIds:  permissionIds,
	})
	return err
}

// DeletePermissionsToRole helper
func (s *RoleServiceGrpc) DeletePermissionsToRole(
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
	_, err := s.clients.RolesClient.DeletePermissions(ctx, &services.DeletePermissionsToRoleRequest{
		OrganizationId: organizationID,
		Namespace:      namespace,
		RoleId:         roleID,
		PermissionIds:  permissionIds,
	})
	return err
}
