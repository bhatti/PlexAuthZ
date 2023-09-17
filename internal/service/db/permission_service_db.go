package db

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/metrics"
	"github.com/bhatti/PlexAuthZ/internal/repository"
	"github.com/bhatti/PlexAuthZ/internal/utils"
	"github.com/twinj/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// PermissionServiceDB - manages persistence of permission data
type PermissionServiceDB struct {
	metricsRegistry      *metrics.Registry
	orgService           *OrganizationServiceDB
	groupRepository      repository.Repository[types.Group]
	resourceRepository   repository.Repository[types.Resource]
	permissionRepository repository.Repository[types.Permission]
	hashRepository       repository.Repository[domain.HashIndex]
}

// NewPermissionServiceDB manages persistence of permission data
func NewPermissionServiceDB(
	metricsRegistry *metrics.Registry,
	orgService *OrganizationServiceDB,
	resourceRepository repository.Repository[types.Resource],
	permissionRepository repository.Repository[types.Permission],
	hashRepository repository.Repository[domain.HashIndex],
) *PermissionServiceDB {
	return &PermissionServiceDB{
		metricsRegistry:      metricsRegistry,
		orgService:           orgService,
		resourceRepository:   resourceRepository,
		permissionRepository: permissionRepository,
		hashRepository:       hashRepository,
	}
}

// CreatePermission - creates a new permission
func (s *PermissionServiceDB) CreatePermission(
	ctx context.Context,
	organizationID string,
	permission *types.Permission) (*types.Permission, error) {
	defer s.metricsRegistry.Elapsed("permissions_svc_create", "org", organizationID)()
	xPermission := domain.NewPermissionExt(permission)
	if err := xPermission.Validate(); err != nil {
		return nil, err
	}

	hash := xPermission.Hash()
	hashIndex, _ := s.hashRepository.GetByID(ctx, organizationID, permission.Namespace, hash)
	if hashIndex != nil {
		return nil, domain.NewDuplicateError(
			fmt.Sprintf("similar permission %v already exists with id %v",
				permission, hashIndex.Ids))
	}
	resource, err := s.resourceRepository.GetByID(ctx, organizationID, permission.Namespace, permission.ResourceId)
	if err != nil {
		return nil, err
	}
	for _, action := range permission.Actions {
		if action != "*" && !utils.Includes(resource.AllowedActions, action) {
			return nil, domain.NewValidationError(fmt.Sprintf("action %s is not allowed for resource %s (%v)",
				action, resource.Name, resource.AllowedActions))
		}
	}
	permission.Id = uuid.NewV4().String()
	permission.Version = 1
	permission.Created = timestamppb.Now()
	permission.Updated = timestamppb.Now()

	err = s.updatePermission(
		ctx,
		organizationID,
		0,
		xPermission)
	if err != nil {
		return nil, err
	}
	return permission, nil
}

// UpdatePermission - updates an existing permission
func (s *PermissionServiceDB) UpdatePermission(
	ctx context.Context,
	organizationID string,
	permission *types.Permission) error {
	defer s.metricsRegistry.Elapsed("permissions_svc_update", "org", organizationID)()
	xPermission := domain.NewPermissionExt(permission)
	if err := xPermission.Validate(); err != nil {
		return err
	}
	if permission.Id == "" {
		return domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, permission.Namespace); err != nil {
		return err
	}

	existing, err := s.permissionRepository.GetByID(
		ctx,
		organizationID,
		permission.Namespace,
		permission.Id)
	if err != nil {
		return err
	}
	version := permission.Version
	if version == 0 {
		version = existing.Version
	}
	permission.Created = existing.Created
	permission.Version = version + 1
	permission.Updated = timestamppb.Now()

	return s.updatePermission(
		ctx,
		organizationID,
		version,
		xPermission)
}

// DeletePermission removes permission
func (s *PermissionServiceDB) DeletePermission(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string) error {
	defer s.metricsRegistry.Elapsed("permissions_svc_delete", "org", organizationID)()
	if id == "" {
		return domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, namespace); err != nil {
		return err
	}
	return s.permissionRepository.Delete(
		ctx,
		organizationID,
		namespace,
		id)
}

// GetPermission - finds permission
func (s *PermissionServiceDB) GetPermission(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string,
) (*types.Permission, error) {
	defer s.metricsRegistry.Elapsed("permissions_svc_get", "org", organizationID)()
	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, namespace); err != nil {
		return nil, err
	}
	return s.permissionRepository.GetByID(
		ctx,
		organizationID,
		namespace,
		id)
}

// GetPermissions - queries permissions
func (s *PermissionServiceDB) GetPermissions(
	ctx context.Context,
	organizationID string,
	namespace string,
	predicate map[string]string,
	offset string,
	limit int64) (res []*types.Permission, nextToken string, err error) {
	defer s.metricsRegistry.Elapsed("permissions_svc_query", "org", organizationID)()
	if predicate["id"] != "" {
		perm, err := s.GetPermission(ctx, organizationID, namespace, predicate["id"])
		if err != nil {
			return nil, "", err
		}
		return []*types.Permission{perm}, "", nil
	}
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, namespace); err != nil {
		return res, "", err
	}
	return s.permissionRepository.Query(
		ctx,
		organizationID,
		namespace,
		predicate,
		offset,
		limit)
}

func (s *PermissionServiceDB) updatePermission(
	ctx context.Context,
	organizationID string,
	version int64,
	xPermission *domain.PermissionExt) (err error) {
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, xPermission.Delegate.Namespace); err != nil {
		return err
	}

	if version == 0 {
		err = s.permissionRepository.Create(
			ctx,
			organizationID,
			xPermission.Delegate.Namespace,
			xPermission.Delegate.Id,
			xPermission.Delegate,
			time.Duration(0),
		)
	} else {
		err = s.permissionRepository.Update(
			ctx,
			organizationID,
			xPermission.Delegate.Namespace,
			xPermission.Delegate.Id,
			version,
			xPermission.Delegate,
			time.Duration(0))
	}
	if err != nil {
		return err
	}

	// update mapping between permission-hash and id
	hash := xPermission.Hash()
	return s.hashRepository.Update(
		ctx,
		organizationID,
		xPermission.Delegate.Namespace,
		hash,
		-1,
		domain.NewHashIndex(hash, []string{xPermission.Delegate.Id}),
		time.Duration(0),
	)
}
