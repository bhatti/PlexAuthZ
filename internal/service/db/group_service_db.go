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

// GroupServiceDB - manages persistence of groups data
type GroupServiceDB struct {
	metricsRegistry *metrics.Registry
	orgService      *OrganizationServiceDB
	groupRepository repository.Repository[types.Group]
	hashRepository  repository.Repository[domain.HashIndex]
}

// NewGroupServiceDB manages persistence of groups data
func NewGroupServiceDB(
	metricsRegistry *metrics.Registry,
	orgService *OrganizationServiceDB,
	groupsRepository repository.Repository[types.Group],
	hashRepository repository.Repository[domain.HashIndex],
) *GroupServiceDB {
	return &GroupServiceDB{
		metricsRegistry: metricsRegistry,
		orgService:      orgService,
		groupRepository: groupsRepository,
		hashRepository:  hashRepository,
	}
}

// CreateGroup - creates a new group
// Note: Redis doesn't allow optimistic concurrency check based on version, so it's possible to have duplicate groups records.
func (s *GroupServiceDB) CreateGroup(
	ctx context.Context,
	organizationID string,
	group *types.Group) (*types.Group, error) {
	xGroup := domain.NewGroupExt(group)
	defer s.metricsRegistry.Elapsed("groups_svc_create", "org", organizationID)()
	if err := xGroup.Validate(); err != nil {
		return nil, err
	}
	hash := xGroup.Hash()
	hashIndex, _ := s.hashRepository.GetByID(ctx, organizationID, group.Namespace, hash)
	if hashIndex != nil {
		return nil, domain.NewDuplicateError(
			fmt.Sprintf("group with name %s already exists with id %v",
				group.Name, hashIndex.Ids))
	}

	group.Id = uuid.NewV4().String()
	group.Version = 1
	group.RoleIds = []string{}
	group.Created = timestamppb.Now()
	group.Updated = timestamppb.Now()
	err := s.updateGroup(ctx, organizationID, 0, xGroup)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// UpdateGroup - updates an existing group
// Note: Redis doesn't allow optimistic concurrency check based on version, so it's possible to have duplicate groups records.
func (s *GroupServiceDB) UpdateGroup(
	ctx context.Context,
	organizationID string,
	group *types.Group) error {
	defer s.metricsRegistry.Elapsed("groups_svc_update", "org", organizationID)()
	xGroup := domain.NewGroupExt(group)
	if err := xGroup.Validate(); err != nil {
		return err
	}
	if group.Id == "" {
		return domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}

	existing, err := s.groupRepository.GetByID(
		ctx,
		organizationID,
		group.Namespace,
		group.Id)
	if err != nil {
		return err
	}
	version := group.Version
	if version == 0 {
		version = existing.Version
	}
	group.Version = version + 1
	group.RoleIds = existing.RoleIds
	group.Updated = timestamppb.Now()
	return s.updateGroup(ctx, organizationID, version, xGroup)
}

// DeleteGroup removes group
func (s *GroupServiceDB) DeleteGroup(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string) error {
	defer s.metricsRegistry.Elapsed("groups_svc_delete", "org", organizationID)()
	if id == "" {
		return domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, namespace); err != nil {
		return err
	}
	return s.groupRepository.Delete(ctx, organizationID, namespace, id)
}

// GetGroup - finds group
func (s *GroupServiceDB) GetGroup(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string,
) (*types.Group, error) {
	defer s.metricsRegistry.Elapsed("groups_svc_get", "org", organizationID)()
	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, namespace); err != nil {
		return nil, err
	}
	return s.groupRepository.GetByID(ctx, organizationID, namespace, id)
}

// GetGroups - queries groups
func (s *GroupServiceDB) GetGroups(
	ctx context.Context,
	organizationID string,
	namespace string,
	predicate map[string]string,
	offset string,
	limit int64) (res []*types.Group, nextOffset string, err error) {
	defer s.metricsRegistry.Elapsed("groups_svc_query", "org", organizationID)()
	if predicate["id"] != "" {
		group, err := s.GetGroup(ctx, organizationID, namespace, predicate["id"])
		if err != nil {
			return nil, "", err
		}
		return []*types.Group{group}, "", nil
	}
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, namespace); err != nil {
		return nil, "", err
	}
	return s.groupRepository.Query(
		ctx,
		organizationID,
		namespace,
		predicate,
		offset,
		limit)
}

// AddRolesToGroup helper
func (s *GroupServiceDB) AddRolesToGroup(
	ctx context.Context,
	organizationID string,
	namespace string,
	groupID string,
	roleIDs ...string,
) error {
	defer s.metricsRegistry.Elapsed("groups_svc_add_roles", "org", organizationID)()
	group, err := s.groupRepository.GetByID(
		ctx,
		organizationID,
		namespace,
		groupID)
	if err != nil {
		return err
	}
	group.RoleIds = utils.AddSlice(group.RoleIds, roleIDs...)
	version := group.Version
	group.Version++
	// update group
	return s.updateGroup(ctx, organizationID, version, domain.NewGroupExt(group))
}

// DeleteRolesToGroup helper
func (s *GroupServiceDB) DeleteRolesToGroup(
	ctx context.Context,
	organizationID string,
	namespace string,
	groupID string,
	roleIDs ...string,
) error {
	defer s.metricsRegistry.Elapsed("groups_svc_delete_roles", "org", organizationID)()
	group, err := s.groupRepository.GetByID(
		ctx,
		organizationID,
		namespace,
		groupID)
	if err != nil {
		return err
	}
	group.RoleIds = utils.RemoveSlice(group.RoleIds, roleIDs...)
	version := group.Version
	group.Version++

	// update group
	return s.updateGroup(ctx, organizationID, version, domain.NewGroupExt(group))
}

func (s *GroupServiceDB) updateGroup(
	ctx context.Context,
	organizationID string,
	version int64,
	xGroup *domain.GroupExt) (err error) {
	if _, err := s.orgService.verifyOrganizationNamespace(
		ctx, organizationID, xGroup.Delegate.Namespace); err != nil {
		return err
	}
	if version == 0 {
		err = s.groupRepository.Create(
			ctx,
			organizationID,
			xGroup.Delegate.Namespace,
			xGroup.Delegate.Id,
			xGroup.Delegate,
			time.Duration(0))
	} else {
		err = s.groupRepository.Update(
			ctx,
			organizationID,
			xGroup.Delegate.Namespace,
			xGroup.Delegate.Id,
			version,
			xGroup.Delegate,
			time.Duration(0))
	}
	if err != nil {
		return err
	}
	hash := xGroup.Hash()
	return s.hashRepository.Update(
		ctx,
		organizationID,
		xGroup.Delegate.Namespace,
		hash,
		-1,
		domain.NewHashIndex(hash, []string{xGroup.Delegate.Id}),
		time.Duration(0),
	)
}
