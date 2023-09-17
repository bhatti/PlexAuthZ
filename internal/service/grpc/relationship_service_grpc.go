package grpc

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/server"
)

// RelationshipServiceGrpc - manages persistence of relationship data
type RelationshipServiceGrpc struct {
	clients server.Clients
}

// NewRelationshipServiceGrpc manages persistence of relationship data
func NewRelationshipServiceGrpc(
	clients server.Clients,
) *RelationshipServiceGrpc {
	return &RelationshipServiceGrpc{
		clients: clients,
	}
}

// CreateRelationship - creates a new relationship
func (s *RelationshipServiceGrpc) CreateRelationship(
	ctx context.Context,
	organizationID string,
	relationship *types.Relationship) (*types.Relationship, error) {
	if organizationID == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	res, err := s.clients.RelationshipsClient.Create(
		ctx,
		&services.CreateRelationshipRequest{
			OrganizationId: organizationID,
			Namespace:      relationship.Namespace,
			Relation:       relationship.Relation,
			PrincipalId:    relationship.PrincipalId,
			ResourceId:     relationship.ResourceId,
			Attributes:     relationship.Attributes,
		})
	if err != nil {
		return nil, err
	}
	relationship.Id = res.Id
	return relationship, nil
}

// UpdateRelationship - updates an existing relationship
func (s *RelationshipServiceGrpc) UpdateRelationship(
	ctx context.Context,
	organizationID string,
	relationship *types.Relationship) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	_, err := s.clients.RelationshipsClient.Update(
		ctx,
		&services.UpdateRelationshipRequest{
			OrganizationId: organizationID,
			Namespace:      relationship.Namespace,
			Id:             relationship.Id,
			Relation:       relationship.Relation,
			PrincipalId:    relationship.PrincipalId,
			ResourceId:     relationship.ResourceId,
			Attributes:     relationship.Attributes,
		})
	return err
}

// DeleteRelationship removes relationship
func (s *RelationshipServiceGrpc) DeleteRelationship(
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
	_, err := s.clients.RelationshipsClient.Delete(
		ctx,
		&services.DeleteRelationshipRequest{
			OrganizationId: organizationID,
			Namespace:      namespace,
			Id:             id,
		})
	return err
}

// GetRelationship - finds relationship
func (s *RelationshipServiceGrpc) GetRelationship(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string,
) (*types.Relationship, error) {
	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	relations, _, err := s.GetRelationships(
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
	if len(relations) == 0 {
		return nil, domain.NewNotFoundError(fmt.Sprintf("relationship %s is not found", id))
	}
	return relations[0], nil
}

// GetRelationships - queries relationships
func (s *RelationshipServiceGrpc) GetRelationships(
	ctx context.Context,
	organizationID string,
	namespace string,
	predicates map[string]string,
	offset string,
	limit int64) (arr []*types.Relationship, nextOffset string, err error) {
	if organizationID == "" {
		return nil, "", domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	if namespace == "" {
		return nil, "", domain.NewValidationError(
			fmt.Sprintf("namespace is not defined"))
	}
	res, err := s.clients.RelationshipsClient.Query(
		ctx,
		&services.QueryRelationshipRequest{
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
		relationRes, err := res.Recv()
		if err != nil {
			break
		}
		nextOffset = relationRes.NextOffset
		arr = append(arr, &types.Relationship{
			Id:          relationRes.Id,
			Version:     relationRes.Version,
			Namespace:   relationRes.Namespace,
			Relation:    relationRes.Relation,
			PrincipalId: relationRes.PrincipalId,
			ResourceId:  relationRes.ResourceId,
			Attributes:  relationRes.Attributes,
			Created:     relationRes.Created,
			Updated:     relationRes.Updated,
		})
	}
	return
}
