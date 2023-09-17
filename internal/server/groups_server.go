package server

import (
	"context"
	api "github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/authz"
	"github.com/bhatti/PlexAuthZ/internal/service"
)

type groupsServer struct {
	api.GroupsServiceServer
	authAdminService service.AuthAdminService
	authorizer       authz.Authorizer
}

// NewGroupsServer constructor
func NewGroupsServer(
	authAdminService service.AuthAdminService,
	authorizer authz.Authorizer,
) (api.GroupsServiceServer, error) {
	return &groupsServer{
		authAdminService: authAdminService,
		authorizer:       authorizer,
	}, nil
}

// Create Group
func (s *groupsServer) Create(
	ctx context.Context,
	req *api.CreateGroupRequest,
) (*api.CreateGroupResponse, error) {
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
	group := &types.Group{
		Namespace: req.Namespace,
		Name:      req.Name,
		ParentIds: req.ParentIds,
	}
	group, err := s.authAdminService.CreateGroup(ctx, req.OrganizationId, group)
	if err != nil {
		return nil, err
	}
	return &api.CreateGroupResponse{
		Id: group.Id,
	}, nil
}

// Update Group
func (s *groupsServer) Update(
	ctx context.Context,
	req *api.UpdateGroupRequest,
) (*api.UpdateGroupResponse, error) {
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
	group := &types.Group{
		Id:        req.Id,
		Namespace: req.Namespace,
		Name:      req.Name,
		ParentIds: req.ParentIds,
	}
	err := s.authAdminService.UpdateGroup(ctx, req.OrganizationId, group)
	if err != nil {
		return nil, err
	}
	return &api.UpdateGroupResponse{}, nil
}

// Query Group
//
// Responses:
// 200: queryGroupResponse
// 400	Bad Request
// 401	Not Authorized
// 500	Internal Error
func (s *groupsServer) Query(
	req *api.QueryGroupRequest,
	sender api.GroupsService_QueryServer,
) error {
	if _, err := s.authorizer.Authorize(
		sender.Context(),
		&api.AuthRequest{
			PrincipalId: authz.Subject(sender.Context()),
			Resource:    objectWildcard,
			Action:      queryAction,
		},
	); err != nil {
		return err
	}
	res, nextOffset, err := s.authAdminService.GetGroups(
		sender.Context(),
		req.OrganizationId,
		req.Namespace,
		req.Predicates,
		req.Offset,
		req.Limit)
	if err != nil {
		return err
	}
	for _, group := range res {
		err = sender.Send(
			&api.QueryGroupResponse{
				Id:         group.Id,
				Version:    group.Version,
				Name:       group.Name,
				RoleIds:    group.RoleIds,
				ParentIds:  group.ParentIds,
				Created:    group.Created,
				Updated:    group.Updated,
				NextOffset: nextOffset,
			})
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete Group
func (s *groupsServer) Delete(
	ctx context.Context,
	req *api.DeleteGroupRequest,
) (*api.DeleteGroupResponse, error) {
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
	err := s.authAdminService.DeleteGroup(ctx, req.OrganizationId, req.Namespace, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.DeleteGroupResponse{}, nil
}

// AddRoles Group
func (s *groupsServer) AddRoles(
	ctx context.Context,
	req *api.AddRolesToGroupRequest,
) (*api.AddRolesToGroupResponse, error) {
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
	if err := s.authAdminService.AddRolesToGroup(
		ctx,
		req.OrganizationId,
		req.Namespace,
		req.GroupId,
		req.RoleIds...); err != nil {
		return nil, err
	}
	return &api.AddRolesToGroupResponse{}, nil
}

// DeleteRoles Group
func (s *groupsServer) DeleteRoles(
	ctx context.Context,
	req *api.DeleteRolesToGroupRequest,
) (*api.DeleteRolesToGroupResponse, error) {
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
	if err := s.authAdminService.DeleteRolesToGroup(
		ctx,
		req.OrganizationId,
		req.Namespace,
		req.GroupId,
		req.RoleIds...); err != nil {
		return nil, err
	}
	return &api.DeleteRolesToGroupResponse{}, nil
}
