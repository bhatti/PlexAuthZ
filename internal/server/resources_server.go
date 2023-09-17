package server

import (
	"context"
	api "github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/authz"
	"github.com/bhatti/PlexAuthZ/internal/service"
)

type resourcesServer struct {
	api.ResourcesServiceServer
	authAdminService service.AuthAdminService
	authorizer       authz.Authorizer
}

// NewResourcesServer constructor
func NewResourcesServer(
	authAdminService service.AuthAdminService,
	authorizer authz.Authorizer,
) (api.ResourcesServiceServer, error) {
	return &resourcesServer{
		authAdminService: authAdminService,
		authorizer:       authorizer,
	}, nil
}

// Create Resource
func (s *resourcesServer) Create(
	ctx context.Context,
	req *api.CreateResourceRequest,
) (*api.CreateResourceResponse, error) {
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
	resource := &types.Resource{
		Namespace:      req.Namespace,
		Name:           req.Name,
		Capacity:       req.Capacity,
		Attributes:     req.Attributes,
		AllowedActions: req.AllowedActions,
	}
	resource, err := s.authAdminService.CreateResource(ctx, req.OrganizationId, resource)
	if err != nil {
		return nil, err
	}
	return &api.CreateResourceResponse{
		Id: resource.Id,
	}, nil
}

// Update Resource
func (s *resourcesServer) Update(
	ctx context.Context,
	req *api.UpdateResourceRequest,
) (*api.UpdateResourceResponse, error) {
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
	resource := &types.Resource{
		Id:             req.Id,
		Namespace:      req.Namespace,
		Name:           req.Name,
		Capacity:       req.Capacity,
		Attributes:     req.Attributes,
		AllowedActions: req.AllowedActions,
	}
	if err := s.authAdminService.UpdateResource(ctx, req.OrganizationId, resource); err != nil {
		return nil, err
	}
	return &api.UpdateResourceResponse{}, nil
}

// Query Resource swagger:route GET /api/{organization_id}/{namespace}/resources/{id} resources queryResourceRequest
//
// Responses:
// 200: queryResourceResponse
// 400	Bad Request
// 401	Not Authorized
// 500	Internal Error
func (s *resourcesServer) Query(
	req *api.QueryResourceRequest,
	sender api.ResourcesService_QueryServer,
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

	res, nextOffset, err := s.authAdminService.QueryResources(
		sender.Context(),
		req.OrganizationId,
		req.Namespace,
		req.Predicates,
		req.Offset,
		req.Limit)
	if err != nil {
		return err
	}
	for _, resource := range res {
		err = sender.Send(
			&api.QueryResourceResponse{
				Id:             resource.Id,
				Version:        resource.Version,
				Namespace:      resource.Namespace,
				Name:           resource.Name,
				Capacity:       resource.Capacity,
				Attributes:     resource.Attributes,
				AllowedActions: resource.AllowedActions,
				Created:        resource.Created,
				Updated:        resource.Updated,
				NextOffset:     nextOffset,
			})
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete Resource
func (s *resourcesServer) Delete(
	ctx context.Context,
	req *api.DeleteResourceRequest,
) (*api.DeleteResourceResponse, error) {
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

	err := s.authAdminService.DeleteResource(
		ctx,
		req.OrganizationId,
		req.Namespace,
		req.Id)
	if err != nil {
		return nil, err
	}
	return &api.DeleteResourceResponse{}, nil
}

// CountResourceInstances Resources
func (s *resourcesServer) CountResourceInstances(
	ctx context.Context,
	req *api.CountResourceInstancesRequest,
) (*api.CountResourceInstancesResponse, error) {
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

	capacity, allocated, err := s.authAdminService.CountResourceInstances(
		ctx,
		req.OrganizationId,
		req.Namespace,
		req.ResourceId,
	)
	if err != nil {
		return nil, err
	}
	return &api.CountResourceInstancesResponse{
		Capacity:  capacity,
		Allocated: allocated,
	}, nil
}

// QueryResourceInstances Resources
func (s *resourcesServer) QueryResourceInstances(
	req *api.QueryResourceInstanceRequest,
	sender api.ResourcesService_QueryResourceInstancesServer) error {
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
	res, nextToken, err := s.authAdminService.QueryResourceInstances(
		sender.Context(),
		req.OrganizationId,
		req.Namespace,
		req.ResourceId,
		req.Predicates,
		req.Offset,
		req.Limit,
	)
	if err != nil {
		return err
	}
	for _, instance := range res {
		err = sender.Send(&api.QueryResourceInstanceResponse{
			Id:          instance.Id,
			Version:     instance.Version,
			Namespace:   instance.Namespace,
			ResourceId:  instance.ResourceId,
			PrincipalId: instance.PrincipalId,
			State:       instance.State,
			NextOffset:  nextToken,
			Created:     instance.Created,
			Updated:     instance.Updated,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
