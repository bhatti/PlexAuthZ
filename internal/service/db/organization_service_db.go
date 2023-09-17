package db

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/metrics"
	"github.com/bhatti/PlexAuthZ/internal/repository"
	"github.com/bhatti/PlexAuthZ/internal/utils"
	"github.com/hashicorp/golang-lru/v2/expirable"
	log "github.com/sirupsen/logrus"
	"github.com/twinj/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// OrganizationServiceDB - manages persistence of AuthZ data
type OrganizationServiceDB struct {
	config          *domain.Config
	metricsRegistry *metrics.Registry
	orgRepository   repository.Repository[types.Organization]
	orgCache        *expirable.LRU[string, *types.Organization]
}

// NewOrganizationServiceDB manages persistence of organization
func NewOrganizationServiceDB(
	metricsRegistry *metrics.Registry,
	orgRepository repository.Repository[types.Organization],
	maxCacheSize int,
	cacheExpirationMillis int,
) *OrganizationServiceDB {
	return &OrganizationServiceDB{
		metricsRegistry: metricsRegistry,
		orgRepository:   orgRepository,
		orgCache: expirable.NewLRU[string, *types.Organization](
			maxCacheSize,
			nil,
			time.Millisecond*time.Duration(cacheExpirationMillis)),
	}
}

// GetOrganization finds organization
func (s *OrganizationServiceDB) GetOrganization(
	ctx context.Context,
	id string) (org *types.Organization, err error) {
	defer s.metricsRegistry.Elapsed("orgs_svc_get", "org", id)()
	org, _ = s.orgCache.Get(id)
	if org != nil {
		return org, nil
	}
	org, err = s.orgRepository.GetByID(ctx, id, "", id)
	if err != nil {
		return nil, err
	}
	s.orgCache.Add(id, org)
	return
}

// GetOrganizations - queries organizations
func (s *OrganizationServiceDB) GetOrganizations(
	ctx context.Context,
	predicate map[string]string,
	offset string,
	limit int64) (res []*types.Organization, nextToken string, err error) {
	defer s.metricsRegistry.Elapsed("orgs_svc_query")
	return s.orgRepository.Query(
		ctx,
		"",
		"",
		predicate,
		offset,
		limit)
}

// CreateOrganization - adds an organization
func (s *OrganizationServiceDB) CreateOrganization(
	ctx context.Context,
	org *types.Organization) (*types.Organization, error) {
	defer s.metricsRegistry.Elapsed("orgs_svc_create")
	if err := domain.NewOrganizationExt(org).Validate(); err != nil {
		return nil, err
	}
	now := timestamppb.Now()
	org.Id = uuid.NewV4().String()
	org.Created = now
	org.Updated = now
	org.Version = 1
	err := s.orgRepository.Create(ctx, org.Id, "", org.Id, org, time.Duration(0))
	if err != nil {
		return nil, err
	}
	_ = s.orgCache.Add(org.Id, org)
	if log.IsLevelEnabled(log.DebugLevel) {
		log.WithFields(log.Fields{
			"Component":    "OrganizationServiceDB",
			"Organization": org.Id,
		}).
			Debugf("created org")
	}
	return org, nil
}

// UpdateOrganization - updates organization
func (s *OrganizationServiceDB) UpdateOrganization(
	ctx context.Context,
	org *types.Organization) error {
	defer s.metricsRegistry.Elapsed("orgs_svc_update", "org", org.Id)()
	if err := domain.NewOrganizationExt(org).Validate(); err != nil {
		return err
	}
	if org.Id == "" {
		return domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	existing, err := s.orgRepository.GetByID(ctx, org.Id, "", org.Id)
	if err != nil {
		return err
	}
	version := org.Version
	if version == 0 {
		version = existing.Version
	}
	org.Version = existing.Version + 1
	org.Updated = timestamppb.Now()
	err = s.orgRepository.Update(ctx, org.Id, "", org.Id, version, org, time.Duration(0))
	if err != nil {
		return err
	}
	_ = s.orgCache.Add(org.Id, org)
	log.WithFields(log.Fields{
		"Component":    "OrganizationServiceDB",
		"Organization": org.Id,
	}).
		Infof("updated org")
	return nil
}

func (s *OrganizationServiceDB) DeleteOrganization(
	ctx context.Context,
	id string) error {
	defer s.metricsRegistry.Elapsed("orgs_svc_delete", "org", id)()
	if id == "" {
		return domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	err := s.orgRepository.Delete(ctx, id, "", id)
	if err != nil {
		return err
	}
	_ = s.orgCache.Remove(id)
	log.WithFields(log.Fields{
		"Component":      "OrganizationServiceDB",
		"OrganizationId": id,
	}).
		Infof("deleted org")
	return nil
}

func (s *OrganizationServiceDB) verifyOrganizationNamespace(
	ctx context.Context,
	organizationID string,
	namespace string) (*types.Organization, error) {
	if organizationID == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("organization_id is not defined"))
	}
	if namespace == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("namespace is not defined"))
	}
	org, err := s.GetOrganization(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	if !utils.Includes(org.Namespaces, namespace) {
		return nil, domain.NewValidationError(
			fmt.Sprintf("namespace %s is not allowed %v", namespace, org.Namespaces))
	}
	return org, nil
}

func toKey(organizationID string, namespace string, id string) string {
	return fmt.Sprintf("%s_%s_%s", organizationID, namespace, id)
}
