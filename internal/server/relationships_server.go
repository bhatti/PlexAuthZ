package server

import (
	"context"
	api "github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/authz"
	"github.com/bhatti/PlexAuthZ/internal/service"
)

type relationshipsServer struct {
	api.RelationshipsServiceServer
	authAdminService service.AuthAdminService
	authorizer       authz.Authorizer
}

// NewRelationshipsServer constructor
func NewRelationshipsServer(
	authAdminService service.AuthAdminService,
	authorizer authz.Authorizer,
) (api.RelationshipsServiceServer, error) {
	return &relationshipsServer{
		authAdminService: authAdminService,
		authorizer:       authorizer,
	}, nil
}

// Create Relationship
func (s *relationshipsServer) Create(
	ctx context.Context,
	req *api.CreateRelationshipRequest,
) (*api.CreateRelationshipResponse, error) {
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
	relationship := &types.Relationship{
		Namespace:   req.Namespace,
		Relation:    req.Relation,
		PrincipalId: req.PrincipalId,
		ResourceId:  req.ResourceId,
		Attributes:  req.Attributes,
	}
	relationship, err := s.authAdminService.CreateRelationship(ctx, req.OrganizationId, relationship)
	if err != nil {
		return nil, err
	}
	return &api.CreateRelationshipResponse{
		Id: relationship.Id,
	}, nil
}

// Update Relationship
func (s *relationshipsServer) Update(
	ctx context.Context,
	req *api.UpdateRelationshipRequest,
) (*api.UpdateRelationshipResponse, error) {
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
	relationship := &types.Relationship{
		Id:          req.Id,
		Namespace:   req.Namespace,
		Relation:    req.Relation,
		PrincipalId: req.PrincipalId,
		ResourceId:  req.ResourceId,
		Attributes:  req.Attributes,
	}
	if err := s.authAdminService.UpdateRelationship(ctx, req.OrganizationId, relationship); err != nil {
		return nil, err
	}
	return &api.UpdateRelationshipResponse{}, nil
}

// Query Relationship swagger:route GET /api/{organization_id}/{namespace}/relationships/{id} relationships queryRelationshipRequest
//
// Responses:
// 200: queryRelationshipResponse
// 400	Bad Request
// 401	Not Authorized
// 500	Internal Error
func (s *relationshipsServer) Query(
	req *api.QueryRelationshipRequest,
	sender api.RelationshipsService_QueryServer,
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
	res, nextOffset, err := s.authAdminService.GetRelationships(
		sender.Context(),
		req.OrganizationId,
		req.Namespace,
		req.Predicates,
		req.Offset,
		req.Limit)
	if err != nil {
		return err
	}
	for _, relationship := range res {
		err = sender.Send(
			&api.QueryRelationshipResponse{
				Id:          relationship.Id,
				Namespace:   relationship.Namespace,
				Version:     relationship.Version,
				Relation:    relationship.Relation,
				PrincipalId: relationship.PrincipalId,
				ResourceId:  relationship.ResourceId,
				Created:     relationship.Created,
				Updated:     relationship.Updated,
				NextOffset:  nextOffset,
			})
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete Relationship
func (s *relationshipsServer) Delete(
	ctx context.Context,
	req *api.DeleteRelationshipRequest,
) (*api.DeleteRelationshipResponse, error) {
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
	err := s.authAdminService.DeleteRelationship(ctx, req.OrganizationId, req.Namespace, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.DeleteRelationshipResponse{}, nil
}
