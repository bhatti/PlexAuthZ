package service

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
)

// RelationshipService - admin APIs for relationships data
type RelationshipService interface {
	// CreateRelationship - creates a new relationship
	CreateRelationship(
		ctx context.Context,
		organizationID string,
		relationship *types.Relationship) (*types.Relationship, error)

	// UpdateRelationship - creates a new relationship
	UpdateRelationship(
		ctx context.Context,
		organizationID string,
		relationship *types.Relationship) error

	// DeleteRelationship removes relationship
	DeleteRelationship(
		ctx context.Context,
		organizationID string,
		namespace string,
		id string) error

	// GetRelationship - finds relationship
	GetRelationship(
		ctx context.Context,
		organizationID string,
		namespace string,
		id string,
	) (*types.Relationship, error)

	// GetRelationships - queries relationships
	GetRelationships(
		ctx context.Context,
		organizationID string,
		namespace string,
		predicate map[string]string,
		offset string,
		limit int64) (res []*types.Relationship, nextOffset string, err error)
}
