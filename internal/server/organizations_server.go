package server

import (
	"context"
	api "github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/authz"
	"github.com/bhatti/PlexAuthZ/internal/service"
	log "github.com/sirupsen/logrus"
)

type organizationsServer struct {
	api.OrganizationsServiceServer
	authAdminService service.AuthAdminService
	authorizer       authz.Authorizer
}

// NewOrganizationsServer constructor
func NewOrganizationsServer(
	authAdminService service.AuthAdminService,
	authorizer authz.Authorizer,
) (api.OrganizationsServiceServer, error) {
	return &organizationsServer{
		authAdminService: authAdminService,
		authorizer:       authorizer,
	}, nil
}

// Create Organization
func (s *organizationsServer) Create(
	ctx context.Context,
	req *api.CreateOrganizationRequest,
) (*api.CreateOrganizationResponse, error) {
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
	organization := &types.Organization{
		Name:       req.Name,
		Url:        req.Url,
		Namespaces: req.Namespaces,
		ParentIds:  req.ParentIds,
	}
	log.WithFields(log.Fields{
		"Component": "OrganizationsServer",
		"Request":   req,
	}).
		Debugf("creating organization")
	savedOrganization, err := s.authAdminService.CreateOrganization(ctx, organization)
	if err != nil {
		return nil, err
	}
	return &api.CreateOrganizationResponse{
		Id: savedOrganization.Id,
	}, nil
}

// Update Organizations swagger:route PUT /api/{organization_id}/{namespace}/organizations/{id} organizations updateOrganizationRequest
//
// Responses:
// 200: updateOrganizationResponse
// 400	Bad Request
// 401	Not Authorized
// 500	Internal Error
func (s *organizationsServer) Update(
	ctx context.Context,
	req *api.UpdateOrganizationRequest,
) (*api.UpdateOrganizationResponse, error) {
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
	organization := &types.Organization{
		Id:         req.Id,
		Name:       req.Name,
		Url:        req.Url,
		Namespaces: req.Namespaces,
		ParentIds:  req.ParentIds,
	}
	log.WithFields(log.Fields{
		"Component": "OrganizationsServer",
		"Request":   req,
	}).
		Debugf("updating organization")
	err := s.authAdminService.UpdateOrganization(ctx, organization)
	if err != nil {
		return nil, err
	}
	return &api.UpdateOrganizationResponse{}, nil
}

// Get Organization swagger:route GET /api/organizations/{id} organizations getOrganizationRequest
func (s *organizationsServer) Get(
	ctx context.Context, req *api.GetOrganizationRequest,
) (*api.GetOrganizationResponse, error) {
	if _, err := s.authorizer.Authorize(
		ctx,
		&api.AuthRequest{
			PrincipalId: authz.Subject(ctx),
			Resource:    objectWildcard,
			Action:      queryAction,
		},
	); err != nil {
		return nil, err
	}
	organization, err := s.authAdminService.GetOrganization(
		ctx,
		req.Id,
	)
	if err != nil {
		return nil, err
	}

	return &api.GetOrganizationResponse{
		Id:         organization.Id,
		Version:    organization.Version,
		Name:       organization.Name,
		Url:        organization.Url,
		Namespaces: organization.Namespaces,
		ParentIds:  organization.ParentIds,
		Created:    organization.Created,
		Updated:    organization.Updated,
	}, nil
}

// Query Organization
func (s *organizationsServer) Query(
	req *api.QueryOrganizationRequest,
	sender api.OrganizationsService_QueryServer,
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
	res, nextOffset, err := s.authAdminService.GetOrganizations(
		sender.Context(),
		req.Predicates,
		req.Offset,
		req.Limit)
	if err != nil {
		return err
	}

	for _, organization := range res {
		err = sender.Send(
			&api.QueryOrganizationResponse{
				Id:         organization.Id,
				Version:    organization.Version,
				Name:       organization.Name,
				Url:        organization.Url,
				Namespaces: organization.Namespaces,
				ParentIds:  organization.ParentIds,
				Created:    organization.Created,
				Updated:    organization.Updated,
				NextOffset: nextOffset,
			})
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete Organization
func (s *organizationsServer) Delete(
	ctx context.Context,
	req *api.DeleteOrganizationRequest,
) (*api.DeleteOrganizationResponse, error) {
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
	log.WithFields(log.Fields{
		"Component": "OrganizationsServer",
		"Request":   req,
	}).
		Debugf("deleting organization")
	err := s.authAdminService.DeleteOrganization(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.DeleteOrganizationResponse{}, nil
}
