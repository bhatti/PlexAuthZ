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

// RoleServiceDB - manages persistence of roles data
type RoleServiceDB struct {
	metricsRegistry *metrics.Registry
	orgService      *OrganizationServiceDB
	roleRepository  repository.Repository[types.Role]
	hashRepository  repository.Repository[domain.HashIndex]
}

// NewRoleServiceDB manages persistence of roles data
func NewRoleServiceDB(
	metricsRegistry *metrics.Registry,
	orgService *OrganizationServiceDB,
	roleRepository repository.Repository[types.Role],
	hashRepository repository.Repository[domain.HashIndex],
) *RoleServiceDB {
	return &RoleServiceDB{
		metricsRegistry: metricsRegistry,
		orgService:      orgService,
		roleRepository:  roleRepository,
		hashRepository:  hashRepository,
	}
}

// CreateRole - creates a new role
func (s *RoleServiceDB) CreateRole(
	ctx context.Context,
	organizationID string,
	role *types.Role) (*types.Role, error) {
	defer s.metricsRegistry.Elapsed("roles_svc_create", "org", organizationID)()
	xRole := domain.NewRoleExt(role)
	if err := xRole.Validate(); err != nil {
		return nil, err
	}
	hash := xRole.Hash()
	hashIndex, _ := s.hashRepository.GetByID(ctx, organizationID, role.Namespace, hash)
	if hashIndex != nil {
		return nil, domain.NewDuplicateError(
			fmt.Sprintf("role with name %s already exists with id %v",
				role.Name, hashIndex.Ids))
	}

	role.Id = uuid.NewV4().String()
	role.Version = 1
	role.Created = timestamppb.Now()
	role.Updated = timestamppb.Now()
	err := s.updateRole(
		ctx,
		organizationID,
		0, // first version
		xRole)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// UpdateRole - updates an existing role
func (s *RoleServiceDB) UpdateRole(
	ctx context.Context,
	organizationID string,
	role *types.Role) error {
	xRole := domain.NewRoleExt(role)
	defer s.metricsRegistry.Elapsed("roles_svc_update", "org", organizationID)()
	if err := xRole.Validate(); err != nil {
		return err
	}
	if role.Id == "" {
		return domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}

	existing, err := s.roleRepository.GetByID(
		ctx,
		organizationID,
		role.Namespace,
		role.Id)
	if err != nil {
		return err
	}
	version := role.Version
	if version == 0 {
		version = existing.Version
	}
	role.Version = version + 1
	role.Updated = timestamppb.Now()
	return s.updateRole(
		ctx,
		organizationID,
		version,
		xRole)
}

// DeleteRole removes role
func (s *RoleServiceDB) DeleteRole(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string) error {
	defer s.metricsRegistry.Elapsed("roles_svc_delete", "org", organizationID)()
	if id == "" {
		return domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, namespace); err != nil {
		return err
	}
	return s.roleRepository.Delete(
		ctx,
		organizationID,
		namespace,
		id)
}

// GetRole - finds role
func (s *RoleServiceDB) GetRole(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string,
) (*types.Role, error) {
	defer s.metricsRegistry.Elapsed("roles_svc_get", "org", organizationID)()
	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, namespace); err != nil {
		return nil, err
	}
	return s.roleRepository.GetByID(ctx, organizationID, namespace, id)
}

// GetRoles - queries roles
func (s *RoleServiceDB) GetRoles(
	ctx context.Context,
	organizationID string,
	namespace string,
	predicate map[string]string,
	offset string,
	limit int64) (res []*types.Role, nextOffset string, err error) {
	defer s.metricsRegistry.Elapsed("roles_svc_query", "org", organizationID)()
	if predicate["id"] != "" {
		role, err := s.GetRole(ctx, organizationID, namespace, predicate["id"])
		if err != nil {
			return nil, "", err
		}
		return []*types.Role{role}, "", nil
	}
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, namespace); err != nil {
		return nil, "", err
	}
	return s.roleRepository.Query(
		ctx,
		organizationID,
		namespace,
		predicate,
		offset,
		limit)
}

// AddPermissionsToRole helper
func (s *RoleServiceDB) AddPermissionsToRole(
	ctx context.Context,
	organizationID string,
	namespace string,
	roleID string,
	permissionIds ...string,
) error {
	defer s.metricsRegistry.Elapsed("roles_svc_add_permissions", "org", organizationID)()
	role, err := s.roleRepository.GetByID(
		ctx,
		organizationID,
		namespace,
		roleID)
	if err != nil {
		return err
	}
	version := role.Version
	role.PermissionIds = utils.AddSlice(role.PermissionIds, permissionIds...)
	role.Version++

	// update role
	return s.updateRole(ctx, organizationID, version, domain.NewRoleExt(role))
}

// DeletePermissionsToRole helper
func (s *RoleServiceDB) DeletePermissionsToRole(
	ctx context.Context,
	organizationID string,
	namespace string,
	roleID string,
	permissionIds ...string,
) error {
	defer s.metricsRegistry.Elapsed("roles_svc_delete_permissions", "org", organizationID)()
	role, err := s.roleRepository.GetByID(
		ctx,
		organizationID,
		namespace,
		roleID)
	if err != nil {
		return err
	}
	version := role.Version
	role.PermissionIds = utils.RemoveSlice(role.PermissionIds, permissionIds...)
	role.Version++

	// update role
	return s.updateRole(ctx, organizationID, version, domain.NewRoleExt(role))
}

func (s *RoleServiceDB) updateRole(
	ctx context.Context,
	organizationID string,
	version int64,
	xRole *domain.RoleExt) (err error) {
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, xRole.Delegate.Namespace); err != nil {
		return err
	}
	if version == 0 {
		err = s.roleRepository.Create(
			ctx,
			organizationID,
			xRole.Delegate.Namespace,
			xRole.Delegate.Id,
			xRole.Delegate,
			time.Duration(0),
		)
	} else {
		err = s.roleRepository.Update(
			ctx,
			organizationID,
			xRole.Delegate.Namespace,
			xRole.Delegate.Id,
			version,
			xRole.Delegate,
			time.Duration(0),
		)
	}
	if err != nil {
		return err
	}
	hash := xRole.Hash()
	return s.hashRepository.Update(
		ctx,
		organizationID,
		xRole.Delegate.Namespace,
		hash,
		-1,
		domain.NewHashIndex(hash, []string{xRole.Delegate.Id}),
		time.Duration(0),
	)
}
