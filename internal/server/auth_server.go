package server

import (
	"context"
	api "github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/internal/authz"
	"github.com/bhatti/PlexAuthZ/internal/service"
	"time"
)

type authServer struct {
	api.AuthZServiceServer
	authAdminService service.AuthAdminService
	authorizer       authz.Authorizer
}

// NewAuthServer constructor.
func NewAuthServer(
	authAdminService service.AuthAdminService,
	authorizer authz.Authorizer,
) (api.AuthZServiceServer, error) {
	return &authServer{
		authAdminService: authAdminService,
		authorizer:       authorizer,
	}, nil
}

// Authorize request for access.
func (s *authServer) Authorize(
	ctx context.Context,
	req *api.AuthRequest,
) (*api.AuthResponse, error) {
	return s.authorizer.Authorize(ctx, req)
}

// Allocate Resources
func (s *authServer) Allocate(
	ctx context.Context,
	req *api.AllocateResourceRequest,
) (*api.AllocateResourceResponse, error) {
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
	var expiry time.Duration
	if req.Expiry != nil {
		expiry = req.Expiry.AsDuration()
	}
	err := s.authAdminService.AllocateResourceInstance(
		ctx,
		req.OrganizationId,
		req.Namespace,
		req.ResourceId,
		req.PrincipalId,
		req.Constraints,
		expiry,
		req.Context,
	)
	if err != nil {
		return nil, err
	}
	return &api.AllocateResourceResponse{}, nil
}

// Deallocate Resources
func (s *authServer) Deallocate(
	ctx context.Context,
	req *api.DeallocateResourceRequest,
) (*api.DeallocateResourceResponse, error) {
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

	err := s.authAdminService.DeallocateResourceInstance(
		ctx,
		req.OrganizationId,
		req.Namespace,
		req.ResourceId,
		req.PrincipalId,
	)
	if err != nil {
		return nil, err
	}
	return &api.DeallocateResourceResponse{}, nil
}
