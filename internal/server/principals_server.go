package server

import (
	"context"
	api "github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/authz"
	"github.com/bhatti/PlexAuthZ/internal/service"
)

type principalsServer struct {
	api.PrincipalsServiceServer
	authAdminService service.AuthAdminService
	authorizer       authz.Authorizer
}

// NewPrincipalsServer constructor
func NewPrincipalsServer(
	authAdminService service.AuthAdminService,
	authorizer authz.Authorizer,
) (api.PrincipalsServiceServer, error) {
	return &principalsServer{
		authAdminService: authAdminService,
		authorizer:       authorizer,
	}, nil
}

// Create Principal
func (s *principalsServer) Create(
	ctx context.Context,
	req *api.CreatePrincipalRequest,
) (*api.CreatePrincipalResponse, error) {
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
	principal := &types.Principal{
		OrganizationId: req.OrganizationId,
		Namespaces:     req.Namespaces,
		Username:       req.Username,
		Name:           req.Name,
		Attributes:     req.Attributes,
	}
	principal, err := s.authAdminService.CreatePrincipal(ctx, principal)
	if err != nil {
		return nil, err
	}
	return &api.CreatePrincipalResponse{
		Id: principal.Id,
	}, nil
}

// Update Principal
func (s *principalsServer) Update(
	ctx context.Context,
	req *api.UpdatePrincipalRequest,
) (*api.UpdatePrincipalResponse, error) {
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
	principal := &types.Principal{
		Id:             req.Id,
		OrganizationId: req.OrganizationId,
		Namespaces:     req.Namespaces,
		Username:       req.Username,
		Name:           req.Name,
		Attributes:     req.Attributes,
	}
	if err := s.authAdminService.UpdatePrincipal(ctx, principal); err != nil {
		return nil, err
	}
	return &api.UpdatePrincipalResponse{}, nil
}

// Get Principal
func (s *principalsServer) Get(
	ctx context.Context,
	req *api.GetPrincipalRequest,
) (*api.GetPrincipalResponse, error) {
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
	xPrincipal, err := s.authAdminService.GetPrincipalExt(
		ctx,
		req.OrganizationId,
		req.Namespace,
		req.Id)
	if err != nil {
		return nil, err
	}
	res := &api.GetPrincipalResponse{
		Id:             xPrincipal.Delegate.Id,
		Version:        xPrincipal.Delegate.Version,
		OrganizationId: xPrincipal.Delegate.OrganizationId,
		Namespaces:     xPrincipal.Delegate.Namespaces,
		Username:       xPrincipal.Delegate.Username,
		Name:           xPrincipal.Delegate.Name,
		Email:          xPrincipal.Delegate.Email,
		Attributes:     xPrincipal.Delegate.Attributes,
		Groups:         xPrincipal.Groups(),
		Roles:          xPrincipal.Roles(),
		Resources:      xPrincipal.Resources(),
		Permissions:    xPrincipal.AllPermissions(),
		Relations:      xPrincipal.Relations(),
		GroupIds:       xPrincipal.Delegate.GroupIds,
		RoleIds:        xPrincipal.Delegate.RoleIds,
		PermissionIds:  xPrincipal.Delegate.PermissionIds,
		Created:        xPrincipal.Delegate.Created,
		Updated:        xPrincipal.Delegate.Updated,
	}
	return res, nil
}

// Query Principal
func (s *principalsServer) Query(
	req *api.QueryPrincipalRequest,
	sender api.PrincipalsService_QueryServer,
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
	res, nextOffset, err := s.authAdminService.GetPrincipals(
		sender.Context(),
		req.OrganizationId,
		req.Predicates,
		req.Offset,
		req.Limit)
	if err != nil {
		return err
	}
	for _, principal := range res {
		err = sender.Send(
			&api.QueryPrincipalResponse{
				Id:             principal.Id,
				Version:        principal.Version,
				OrganizationId: principal.OrganizationId,
				Username:       principal.Username,
				Name:           principal.Name,
				Email:          principal.Email,
				Attributes:     principal.Attributes,
				GroupIds:       principal.GroupIds,
				RoleIds:        principal.RoleIds,
				PermissionIds:  principal.PermissionIds,
				RelationIds:    principal.RelationIds,
				Created:        principal.Created,
				Updated:        principal.Updated,
				NextOffset:     nextOffset,
			})
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete Principal
func (s *principalsServer) Delete(
	ctx context.Context,
	req *api.DeletePrincipalRequest,
) (*api.DeletePrincipalResponse, error) {
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
	err := s.authAdminService.DeletePrincipal(ctx, req.OrganizationId, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.DeletePrincipalResponse{}, nil
}

// AddGroups Principal
func (s *principalsServer) AddGroups(
	ctx context.Context,
	req *api.AddGroupsToPrincipalRequest,
) (*api.AddGroupsToPrincipalResponse, error) {
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
	err := s.authAdminService.AddGroupsToPrincipal(
		ctx,
		req.OrganizationId,
		req.Namespace,
		req.PrincipalId,
		req.GroupIds...)
	if err != nil {
		return nil, err
	}
	return &api.AddGroupsToPrincipalResponse{}, nil
}

// DeleteGroups Principal
func (s *principalsServer) DeleteGroups(
	ctx context.Context,
	req *api.DeleteGroupsToPrincipalRequest,
) (*api.DeleteGroupsToPrincipalResponse, error) {
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
	err := s.authAdminService.DeleteGroupsToPrincipal(
		ctx,
		req.OrganizationId,
		req.Namespace,
		req.PrincipalId,
		req.GroupIds...)
	if err != nil {
		return nil, err
	}
	return &api.DeleteGroupsToPrincipalResponse{}, nil
}

// AddRoles Principal
func (s *principalsServer) AddRoles(
	ctx context.Context,
	req *api.AddRolesToPrincipalRequest,
) (*api.AddRolesToPrincipalResponse, error) {
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
	err := s.authAdminService.AddRolesToPrincipal(
		ctx,
		req.OrganizationId,
		req.Namespace,
		req.PrincipalId,
		req.RoleIds...)
	if err != nil {
		return nil, err
	}
	return &api.AddRolesToPrincipalResponse{}, nil
}

// DeleteRoles Principal
func (s *principalsServer) DeleteRoles(
	ctx context.Context,
	req *api.DeleteRolesToPrincipalRequest,
) (*api.DeleteRolesToPrincipalResponse, error) {
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
	err := s.authAdminService.DeleteRolesToPrincipal(
		ctx,
		req.OrganizationId,
		req.Namespace,
		req.PrincipalId,
		req.RoleIds...)
	if err != nil {
		return nil, err
	}
	return &api.DeleteRolesToPrincipalResponse{}, nil
}

// AddPermissions Principal
func (s *principalsServer) AddPermissions(
	ctx context.Context,
	req *api.AddPermissionsToPrincipalRequest,
) (*api.AddPermissionsToPrincipalResponse, error) {
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
	err := s.authAdminService.AddPermissionsToPrincipal(
		ctx,
		req.OrganizationId,
		req.Namespace,
		req.PrincipalId,
		req.PermissionIds...)
	if err != nil {
		return nil, err
	}
	return &api.AddPermissionsToPrincipalResponse{}, nil
}

// DeletePermissions Principal
func (s *principalsServer) DeletePermissions(
	ctx context.Context,
	req *api.DeletePermissionsToPrincipalRequest,
) (*api.DeletePermissionsToPrincipalResponse, error) {
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
	err := s.authAdminService.DeletePermissionsToPrincipal(
		ctx,
		req.OrganizationId,
		req.Namespace,
		req.PrincipalId,
		req.PermissionIds...)
	if err != nil {
		return nil, err
	}
	return &api.DeletePermissionsToPrincipalResponse{}, nil
}

// AddRelationships Principal
func (s *principalsServer) AddRelationships(
	ctx context.Context,
	req *api.AddRelationshipsToPrincipalRequest,
) (*api.AddRelationshipsToPrincipalResponse, error) {
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
	err := s.authAdminService.AddRelationshipsToPrincipal(
		ctx,
		req.OrganizationId,
		req.Namespace,
		req.PrincipalId,
		req.RelationshipIds...)
	if err != nil {
		return nil, err
	}
	return &api.AddRelationshipsToPrincipalResponse{}, nil
}

// DeleteRelationships Principal
func (s *principalsServer) DeleteRelationships(
	ctx context.Context,
	req *api.DeleteRelationshipsToPrincipalRequest,
) (*api.DeleteRelationshipsToPrincipalResponse, error) {
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
	err := s.authAdminService.DeleteRelationshipsToPrincipal(
		ctx,
		req.OrganizationId,
		req.Namespace,
		req.PrincipalId,
		req.RelationshipIds...)
	if err != nil {
		return nil, err
	}
	return &api.DeleteRelationshipsToPrincipalResponse{}, nil
}
