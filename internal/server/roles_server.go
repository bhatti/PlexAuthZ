package server

import (
	"context"
	api "github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/authz"
	"github.com/bhatti/PlexAuthZ/internal/service"
)

type rolesServer struct {
	api.RolesServiceServer
	authAdminService service.AuthAdminService
	authorizer       authz.Authorizer
}

// NewRolesServer constructor
func NewRolesServer(
	authAdminService service.AuthAdminService,
	authorizer authz.Authorizer,
) (api.RolesServiceServer, error) {
	return &rolesServer{
		authAdminService: authAdminService,
		authorizer:       authorizer,
	}, nil
}

// Create Role
func (s *rolesServer) Create(
	ctx context.Context,
	req *api.CreateRoleRequest,
) (*api.CreateRoleResponse, error) {
	if _, err := s.authorizer.Authorize(
		ctx,
		&api.AuthRequest{
			PrincipalId: authz.Subject(ctx),
			Resource:    objectWildcard,
			Action:      updateAction,
		},
	); err != nil {
		return nil, err
	}

	role := &types.Role{
		Namespace: req.Namespace,
		Name:      req.Name,
		ParentIds: req.ParentIds,
	}
	role, err := s.authAdminService.CreateRole(ctx, req.OrganizationId, role)
	if err != nil {
		return nil, err
	}
	return &api.CreateRoleResponse{
		Id: role.Id,
	}, nil
}

// Update Role
func (s *rolesServer) Update(
	ctx context.Context,
	req *api.UpdateRoleRequest,
) (*api.UpdateRoleResponse, error) {
	if _, err := s.authorizer.Authorize(
		ctx,
		&api.AuthRequest{
			PrincipalId: authz.Subject(ctx),
			Resource:    objectWildcard,
			Action:      updateAction,
		},
	); err != nil {
		return nil, err
	}

	role := &types.Role{
		Id:        req.Id,
		Namespace: req.Namespace,
		Name:      req.Name,
		ParentIds: req.ParentIds,
	}
	if err := s.authAdminService.UpdateRole(ctx, req.OrganizationId, role); err != nil {
		return nil, err
	}
	return &api.UpdateRoleResponse{}, nil
}

// Query Role swagger:route GET /api/{organization_id}/{namespace}/roles/{id} roles queryRoleRequest
//
// Responses:
// 200: queryRoleResponse
// 400	Bad Request
// 401	Not Authorized
// 500	Internal Error
func (s *rolesServer) Query(
	req *api.QueryRoleRequest,
	sender api.RolesService_QueryServer,
) error {
	if _, err := s.authorizer.Authorize(
		sender.Context(),
		&api.AuthRequest{
			PrincipalId: authz.Subject(sender.Context()),
			Resource:    objectWildcard,
			Action:      updateAction,
		},
	); err != nil {
		return err
	}

	res, nextOffset, err := s.authAdminService.GetRoles(
		sender.Context(),
		req.OrganizationId,
		req.Namespace,
		req.Predicates,
		req.Offset,
		req.Limit)
	if err != nil {
		return err
	}
	for _, role := range res {
		err = sender.Send(
			&api.QueryRoleResponse{
				Id:            role.Id,
				Version:       role.Version,
				Name:          role.Name,
				Namespace:     role.Namespace,
				PermissionIds: role.PermissionIds,
				ParentIds:     role.ParentIds,
				Created:       role.Created,
				Updated:       role.Updated,
				NextOffset:    nextOffset,
			})
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete Role
func (s *rolesServer) Delete(
	ctx context.Context,
	req *api.DeleteRoleRequest,
) (*api.DeleteRoleResponse, error) {
	if _, err := s.authorizer.Authorize(
		ctx,
		&api.AuthRequest{
			PrincipalId: authz.Subject(ctx),
			Resource:    objectWildcard,
			Action:      deleteAction,
		},
	); err != nil {
		return nil, err
	}

	err := s.authAdminService.DeleteRole(ctx, req.OrganizationId, req.Namespace, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.DeleteRoleResponse{}, nil
}

// AddPermissions Role
func (s *rolesServer) AddPermissions(
	ctx context.Context,
	req *api.AddPermissionsToRoleRequest,
) (*api.AddPermissionsToRoleResponse, error) {
	if _, err := s.authorizer.Authorize(
		ctx,
		&api.AuthRequest{
			PrincipalId: authz.Subject(ctx),
			Resource:    objectWildcard,
			Action:      updateAction,
		},
	); err != nil {
		return nil, err
	}

	err := s.authAdminService.AddPermissionsToRole(
		ctx,
		req.OrganizationId,
		req.Namespace,
		req.RoleId,
		req.PermissionIds...)
	if err != nil {
		return nil, err
	}
	return &api.AddPermissionsToRoleResponse{}, nil
}

// DeletePermissions Role
func (s *rolesServer) DeletePermissions(
	ctx context.Context,
	req *api.DeletePermissionsToRoleRequest,
) (*api.DeletePermissionsToRoleResponse, error) {

	if _, err := s.authorizer.Authorize(
		ctx,
		&api.AuthRequest{
			PrincipalId: authz.Subject(ctx),
			Resource:    objectWildcard,
			Action:      updateAction,
		},
	); err != nil {
		return nil, err
	}

	err := s.authAdminService.DeletePermissionsToRole(
		ctx,
		req.OrganizationId,
		req.Namespace,
		req.RoleId,
		req.PermissionIds...)
	if err != nil {
		return nil, err
	}
	return &api.DeletePermissionsToRoleResponse{}, nil
}
