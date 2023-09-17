package db

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/metrics"
	"github.com/bhatti/PlexAuthZ/internal/repository"
	"github.com/twinj/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// RelationshipServiceDB - manages persistence of relationship data
type RelationshipServiceDB struct {
	metricsRegistry        *metrics.Registry
	orgService             *OrganizationServiceDB
	relationshipRepository repository.Repository[types.Relationship]
	hashRepository         repository.Repository[domain.HashIndex]
}

// NewRelationshipServiceDB manages persistence of relationship data
func NewRelationshipServiceDB(
	metricsRegistry *metrics.Registry,
	orgService *OrganizationServiceDB,
	relationshipRepository repository.Repository[types.Relationship],
	hashRepository repository.Repository[domain.HashIndex],
) *RelationshipServiceDB {
	return &RelationshipServiceDB{
		metricsRegistry:        metricsRegistry,
		orgService:             orgService,
		relationshipRepository: relationshipRepository,
		hashRepository:         hashRepository,
	}
}

// CreateRelationship - creates a new relationship
func (s *RelationshipServiceDB) CreateRelationship(
	ctx context.Context,
	organizationID string,
	relationship *types.Relationship) (*types.Relationship, error) {
	defer s.metricsRegistry.Elapsed("relationships_svc_create", "org", organizationID)()
	xRelationship := domain.NewRelationshipExt(relationship)
	if err := xRelationship.Validate(); err != nil {
		return nil, err
	}
	hash := xRelationship.Hash()
	hashIndex, _ := s.hashRepository.GetByID(ctx, organizationID, relationship.Namespace, hash)
	if hashIndex != nil {
		return nil, domain.NewDuplicateError(
			fmt.Sprintf("relation %s already exists with id %v",
				relationship.Relation, hashIndex.Ids))
	}

	relationship.Id = uuid.NewV4().String()
	relationship.Version = 1
	relationship.Created = timestamppb.Now()
	relationship.Updated = timestamppb.Now()
	err := s.updateRelationship(
		ctx,
		organizationID,
		0, // first version
		xRelationship,
	)
	if err != nil {
		return nil, err
	}
	return relationship, nil
}

// UpdateRelationship - updates an existing relationship
func (s *RelationshipServiceDB) UpdateRelationship(
	ctx context.Context,
	organizationID string,
	relationship *types.Relationship) error {
	defer s.metricsRegistry.Elapsed("relationships_svc_update", "org", organizationID)()
	xRelationship := domain.NewRelationshipExt(relationship)
	if err := xRelationship.Validate(); err != nil {
		return err
	}
	if relationship.Id == "" {
		return domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}

	existing, err := s.relationshipRepository.GetByID(
		ctx,
		organizationID,
		relationship.Namespace,
		relationship.Id)
	if err != nil {
		return err
	}
	version := relationship.Version
	if version == 0 {
		version = existing.Version
	}
	relationship.Version = version + 1
	relationship.Updated = timestamppb.Now()
	return s.updateRelationship(
		ctx,
		organizationID,
		version,
		xRelationship,
	)
}

// DeleteRelationship removes relationship
func (s *RelationshipServiceDB) DeleteRelationship(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string) error {
	defer s.metricsRegistry.Elapsed("relationships_svc_delete", "org", organizationID)()
	if id == "" {
		return domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, namespace); err != nil {
		return err
	}
	return s.relationshipRepository.Delete(ctx, organizationID, namespace, id)
}

// GetRelationship - finds relationship
func (s *RelationshipServiceDB) GetRelationship(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string,
) (*types.Relationship, error) {
	defer s.metricsRegistry.Elapsed("relationships_svc_get", "org", organizationID)()
	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	if _, err := s.orgService.verifyOrganizationNamespace(ctx, organizationID, namespace); err != nil {
		return nil, err
	}
	return s.relationshipRepository.GetByID(
		ctx,
		organizationID,
		namespace,
		id,
	)
}

// GetRelationships - queries relationships
func (s *RelationshipServiceDB) GetRelationships(
	ctx context.Context,
	organizationID string,
	namespace string,
	predicate map[string]string,
	offset string,
	limit int64) (res []*types.Relationship, nextOffset string, err error) {
	defer s.metricsRegistry.Elapsed("relationships_svc_query", "org", organizationID)()
	if predicate["id"] != "" {
		relationship, err := s.GetRelationship(ctx, organizationID, namespace, predicate["id"])
		if err != nil {
			return nil, "", err
		}
		return []*types.Relationship{relationship}, "", nil
	}
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, namespace); err != nil {
		return res, "", err
	}
	return s.relationshipRepository.Query(
		ctx,
		organizationID,
		namespace,
		predicate,
		offset,
		limit)
}

func (s *RelationshipServiceDB) updateRelationship(
	ctx context.Context,
	organizationID string,
	version int64,
	xRelationship *domain.RelationshipExt) (err error) {
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, xRelationship.Delegate.Namespace); err != nil {
		return err
	}
	if version == 0 {
		err = s.relationshipRepository.Create(
			ctx,
			organizationID,
			xRelationship.Delegate.Namespace,
			xRelationship.Delegate.Id,
			xRelationship.Delegate,
			time.Duration(0),
		)
	} else {
		err = s.relationshipRepository.Update(
			ctx,
			organizationID,
			xRelationship.Delegate.Namespace,
			xRelationship.Delegate.Id,
			version,
			xRelationship.Delegate,
			time.Duration(0),
		)
	}
	if err != nil {
		return err
	}
	hash := xRelationship.Hash()
	return s.hashRepository.Update(
		ctx,
		organizationID,
		xRelationship.Delegate.Namespace,
		hash,
		-1, // no version
		domain.NewHashIndex(hash, []string{xRelationship.Delegate.Id}),
		time.Duration(0),
	)
}
