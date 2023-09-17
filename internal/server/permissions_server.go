package server

import (
	"context"
	api "github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/authz"
	"github.com/bhatti/PlexAuthZ/internal/service"
)

type permissionsServer struct {
	api.PermissionsServiceServer
	authAdminService service.AuthAdminService
	authorizer       authz.Authorizer
}

// NewPermissionsServer constructor
func NewPermissionsServer(
	authAdminService service.AuthAdminService,
	authorizer authz.Authorizer,
) (api.PermissionsServiceServer, error) {
	return &permissionsServer{
		authAdminService: authAdminService,
		authorizer:       authorizer,
	}, nil
}

// Create Permission
func (s *permissionsServer) Create(
	ctx context.Context,
	req *api.CreatePermissionRequest,
) (*api.CreatePermissionResponse, error) {
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
	permission := &types.Permission{
		Namespace:   req.Namespace,
		Scope:       req.Scope,
		Actions:     req.Actions,
		ResourceId:  req.ResourceId,
		Effect:      req.Effect,
		Constraints: req.Constraints,
	}
	permission, err := s.authAdminService.CreatePermission(ctx, req.OrganizationId, permission)
	if err != nil {
		return nil, err
	}
	return &api.CreatePermissionResponse{
		Id: permission.Id,
	}, nil
}

// Update Permission
func (s *permissionsServer) Update(
	ctx context.Context,
	req *api.UpdatePermissionRequest,
) (*api.UpdatePermissionResponse, error) {
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
	permission := &types.Permission{
		Id:          req.Id,
		Namespace:   req.Namespace,
		Scope:       req.Scope,
		Actions:     req.Actions,
		ResourceId:  req.ResourceId,
		Effect:      req.Effect,
		Constraints: req.Constraints,
	}
	if err := s.authAdminService.UpdatePermission(ctx, req.OrganizationId, permission); err != nil {
		return nil, err
	}
	return &api.UpdatePermissionResponse{}, nil
}

// Query Permission swagger:route GET /api/{organization_id}/{namespace}/permissions/{id} permissions queryPermissionRequest
//
// Responses:
// 200: queryPermissionResponse
// 400	Bad Request
// 401	Not Authorized
// 500	Internal Error
func (s *permissionsServer) Query(
	req *api.QueryPermissionRequest,
	sender api.PermissionsService_QueryServer,
) error {
	if _, err := s.authorizer.Authorize(
		sender.Context(),
		&api.AuthRequest{
			PrincipalId: authz.Subject(sender.Context()),
			Resource:    objectWildcard,
			Action:      deleteAction,
		},
	); err != nil {
		return err
	}
	res, nextOffset, err := s.authAdminService.GetPermissions(
		sender.Context(),
		req.OrganizationId,
		req.Namespace,
		req.Predicates,
		req.Offset,
		req.Limit)
	if err != nil {
		return err
	}
	for _, permission := range res {
		err = sender.Send(
			&api.QueryPermissionResponse{
				Id:          permission.Id,
				Version:     permission.Version,
				Namespace:   permission.Namespace,
				Scope:       permission.Scope,
				Actions:     permission.Actions,
				ResourceId:  permission.ResourceId,
				Effect:      permission.Effect,
				Constraints: permission.Constraints,
				Created:     permission.Created,
				Updated:     permission.Updated,
				NextOffset:  nextOffset,
			})
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete Permission
func (s *permissionsServer) Delete(
	ctx context.Context,
	req *api.DeletePermissionRequest,
) (*api.DeletePermissionResponse, error) {
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
	err := s.authAdminService.DeletePermission(ctx, req.OrganizationId, req.Namespace, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.DeletePermissionResponse{}, nil
}
