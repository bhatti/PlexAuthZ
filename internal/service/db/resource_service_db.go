package db

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/metrics"
	"github.com/bhatti/PlexAuthZ/internal/repository"
	"github.com/twinj/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"time"
)

// ResourceServiceDB - manages persistence of resources and resource instances
type ResourceServiceDB struct {
	metricsRegistry                   *metrics.Registry
	orgService                        *OrganizationServiceDB
	principalService                  *PrincipalServiceDB
	resourceRepository                repository.Repository[types.Resource]
	resourceInstanceRepositoryFactory repository.ResourceInstanceRepositoryFactory
	hashRepository                    repository.Repository[domain.HashIndex]
}

// NewResourceServiceDB manages persistence of resources
func NewResourceServiceDB(
	metricsRegistry *metrics.Registry,
	orgService *OrganizationServiceDB,
	principalService *PrincipalServiceDB,
	resourceRepository repository.Repository[types.Resource],
	resourceInstanceRepositoryFactory repository.ResourceInstanceRepositoryFactory,
	hashRepository repository.Repository[domain.HashIndex],
) *ResourceServiceDB {
	return &ResourceServiceDB{
		metricsRegistry:                   metricsRegistry,
		orgService:                        orgService,
		principalService:                  principalService,
		resourceRepository:                resourceRepository,
		resourceInstanceRepositoryFactory: resourceInstanceRepositoryFactory,
		hashRepository:                    hashRepository,
	}
}

// CreateResource - creates a new instance of resource
func (s *ResourceServiceDB) CreateResource(
	ctx context.Context,
	organizationID string,
	resource *types.Resource) (*types.Resource, error) {
	defer s.metricsRegistry.Elapsed("resources_svc_create", "org", organizationID)()
	xResource := domain.NewResourceExt(resource)
	if err := xResource.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, resource.Namespace); err != nil {
		return nil, err
	}
	hash := xResource.Hash()
	hashIndex, _ := s.hashRepository.GetByID(ctx, organizationID, resource.Namespace, hash)
	if hashIndex != nil {
		return nil, domain.NewDuplicateError(
			fmt.Sprintf("resource with name %s already exists with id %v",
				resource.Name, hashIndex.Ids))
	}

	resource.Id = uuid.NewV4().String()
	resource.Version = 1
	resource.Wildcard = strings.Contains(resource.Name, "*")
	resource.Created = timestamppb.Now()
	resource.Updated = timestamppb.Now()

	err := s.updateResource(
		ctx,
		organizationID,
		0, // first version
		xResource)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

// UpdateResource - updates an existing resource
func (s *ResourceServiceDB) UpdateResource(
	ctx context.Context,
	organizationID string,
	resource *types.Resource) error {
	defer s.metricsRegistry.Elapsed("resources_svc_update", "org", organizationID)()
	xResource := domain.NewResourceExt(resource)
	if err := xResource.Validate(); err != nil {
		return err
	}
	if resource.Id == "" {
		return domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}

	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, resource.Namespace); err != nil {
		return err
	}

	existing, err := s.resourceRepository.GetByID(
		ctx,
		organizationID,
		resource.Namespace,
		resource.Id)
	if err != nil {
		return err
	}
	version := resource.Version
	if version == 0 {
		version = existing.Version
	}
	resource.Wildcard = strings.Contains(resource.Name, "*")
	resource.Version = version + 1
	resource.Created = existing.Created
	resource.Updated = timestamppb.Now()

	return s.updateResource(
		ctx,
		organizationID,
		version,
		xResource)
}

// DeleteResource removes resource
func (s *ResourceServiceDB) DeleteResource(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string) error {
	defer s.metricsRegistry.Elapsed("resources_svc_delete", "org", organizationID)()
	if id == "" {
		return domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, namespace); err != nil {
		return err
	}
	existing, err := s.resourceRepository.GetByID(ctx, organizationID, namespace, id)
	if err != nil {
		return err
	}
	err = s.resourceRepository.Delete(ctx, organizationID, namespace, id)
	if err != nil {
		return err
	}
	return s.hashRepository.Delete(
		ctx,
		organizationID,
		namespace,
		domain.NewResourceExt(existing).Hash())
}

// GetResource - finds resource
func (s *ResourceServiceDB) GetResource(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string,
) (*types.Resource, error) {
	defer s.metricsRegistry.Elapsed("resources_svc_get", "org", organizationID)()
	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, namespace); err != nil {
		return nil, err
	}
	return s.resourceRepository.GetByID(
		ctx,
		organizationID,
		namespace,
		id,
	)
}

// QueryResources - queries resources by predicates.
func (s *ResourceServiceDB) QueryResources(
	ctx context.Context,
	organizationID string,
	namespace string,
	predicate map[string]string,
	offset string,
	limit int64) (res []*types.Resource, nextOffset string, err error) {
	defer s.metricsRegistry.Elapsed("resources_svc_query", "org", organizationID)()
	if predicate["id"] != "" {
		resource, err := s.GetResource(ctx, organizationID, namespace, predicate["id"])
		if err != nil {
			return nil, "", err
		}
		return []*types.Resource{resource}, "", nil
	}
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, namespace); err != nil {
		return nil, "", err
	}
	return s.resourceRepository.Query(
		ctx,
		organizationID,
		namespace,
		predicate,
		offset, limit)
}

// QueryResourceInstances - queries resource-instances
func (s *ResourceServiceDB) QueryResourceInstances(
	ctx context.Context,
	organizationID string,
	namespace string,
	resourceID string,
	predicate map[string]string,
	offset string,
	limit int64) (res []*types.ResourceInstance, nextOffset string, err error) {
	defer s.metricsRegistry.Elapsed("resources_svc_query_instances", "org", organizationID)()
	if resourceID == "" {
		return res, "", domain.NewValidationError(
			fmt.Sprintf("resource_id is not defined"))
	}
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, namespace); err != nil {
		return nil, "", err
	}
	instanceRepository, err := s.resourceInstanceRepositoryFactory.CreateResourceInstanceRepository(resourceID)
	if err != nil {
		return res, "", err
	}

	res, _, err = instanceRepository.Query(
		ctx,
		organizationID,
		namespace,
		predicate,
		offset,
		limit)
	return instanceRepository.Query(
		ctx,
		organizationID,
		namespace,
		predicate,
		offset,
		limit)
}

// AllocateResourceInstance - allocates resource-instance
func (s *ResourceServiceDB) AllocateResourceInstance(
	ctx context.Context,
	organizationID string,
	namespace string,
	resourceID string,
	principalID string,
	constraints string,
	expiry time.Duration,
	context map[string]string,
) error {
	defer s.metricsRegistry.Elapsed("resources_svc_allocate", "org", organizationID)()
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, namespace); err != nil {
		return err
	}
	if constraints != "" {
		resource, err := s.GetResource(ctx, organizationID, namespace, resourceID)
		if err != nil {
			return err
		}
		xPrincipal, err := s.principalService.GetPrincipalExt(
			ctx,
			organizationID,
			namespace,
			principalID)
		if err != nil {
			return err
		}
		matched, _, err := xPrincipal.CheckConstraints(
			&services.AuthRequest{
				OrganizationId: organizationID,
				Namespace:      namespace,
				PrincipalId:    principalID,
				Context:        context,
			},
			resource,
			constraints,
		)
		if err != nil {
			return err
		}
		if !matched {
			return domain.NewAuthError(fmt.Sprintf("constraints '%s' didn't match", constraints))
		}
	} else {
		if _, err := s.principalService.GetPrincipal(ctx, organizationID, principalID); err != nil {
			return err
		}
	}

	xInstance := domain.NewResourceInstanceExt(namespace, resourceID, principalID)
	if err := xInstance.Validate(); err != nil {
		return err
	}
	capacity, allocated, err := s.CountResourceInstances(ctx, organizationID, namespace, resourceID)
	if err != nil {
		return err
	}
	if capacity == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("capacity is not defined"))
	}
	if allocated >= capacity {
		return domain.NewValidationError(
			fmt.Sprintf("usage %d exceeds capacity %d", allocated, capacity))
	}
	instanceRepository, err := s.resourceInstanceRepositoryFactory.CreateResourceInstanceRepository(resourceID)
	if err != nil {
		return err
	}
	existing, _ := instanceRepository.GetByID(ctx, organizationID, namespace, xInstance.Delegate.Id)
	var version int64 = 0
	xInstance.Delegate.Updated = timestamppb.Now()
	if existing != nil {
		version = existing.Version
		xInstance.Delegate.Version = existing.Version + 1
		xInstance.Delegate.Created = existing.Created
		return instanceRepository.Update(
			ctx,
			organizationID,
			namespace,
			xInstance.Delegate.Id,
			version,
			xInstance.Delegate,
			expiry)
	}
	return instanceRepository.Create(
		ctx,
		organizationID,
		namespace,
		xInstance.Delegate.Id,
		xInstance.Delegate,
		expiry)
}

// DeallocateResourceInstance - deallocates resource-instance
func (s *ResourceServiceDB) DeallocateResourceInstance(
	ctx context.Context,
	organizationID string,
	namespace string,
	resourceID string,
	principalID string,
) error {
	defer s.metricsRegistry.Elapsed("resources_svc_deallocate", "org", organizationID)()
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, namespace); err != nil {
		return err
	}
	xInstance := domain.NewResourceInstanceExt(namespace, resourceID, principalID)
	if err := xInstance.Validate(); err != nil {
		return err
	}
	instanceRepository, err := s.resourceInstanceRepositoryFactory.CreateResourceInstanceRepository(resourceID)
	if err != nil {
		return err
	}
	return instanceRepository.Delete(
		ctx,
		organizationID,
		namespace,
		xInstance.Delegate.Id)
}

// CountResourceInstances - size of total and allocated resource-instances
func (s *ResourceServiceDB) CountResourceInstances(
	ctx context.Context,
	organizationID string,
	namespace string,
	resourceID string,
) (capacity int32, allocated int32, err error) {
	defer s.metricsRegistry.Elapsed("resources_svc_count_instances", "org", organizationID)()
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, namespace); err != nil {
		return 0, 0, err
	}
	if resourceID == "" {
		return 0, 0,
			domain.NewValidationError(
				fmt.Sprintf("resource_id is not defined"))
	}
	resource, err := s.resourceRepository.GetByID(ctx, organizationID, namespace, resourceID)
	if err != nil {
		return 0, 0, err
	}
	instanceRepository, err := s.resourceInstanceRepositoryFactory.CreateResourceInstanceRepository(resourceID)
	if err != nil {
		return 0, 0, err
	}
	capacity = resource.Capacity
	count, err := instanceRepository.Size(ctx, organizationID, namespace)
	if err != nil {
		return 0, 0, err
	}
	allocated = int32(count)
	return
}

// updateResource - save resource
func (s *ResourceServiceDB) updateResource(
	ctx context.Context,
	organizationID string,
	version int64,
	xResource *domain.ResourceExt) (err error) {

	if version == 0 {
		err = s.resourceRepository.Create(
			ctx,
			organizationID,
			xResource.Delegate.Namespace,
			xResource.Delegate.Id,
			xResource.Delegate,
			time.Duration(0))
	} else {
		err = s.resourceRepository.Update(
			ctx,
			organizationID,
			xResource.Delegate.Namespace,
			xResource.Delegate.Id,
			version,
			xResource.Delegate,
			time.Duration(0))
	}
	if err != nil {
		return err
	}

	// update mapping between resource-name and id
	hash := xResource.Hash()
	return s.hashRepository.Update(
		ctx,
		organizationID,
		xResource.Delegate.Namespace,
		hash,
		-1, // no version
		domain.NewHashIndex(hash, []string{xResource.Delegate.Id}),
		time.Duration(0))
}
