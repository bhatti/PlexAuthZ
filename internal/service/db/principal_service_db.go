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

// PrincipalServiceDB - manages persistence of principal objects
type PrincipalServiceDB struct {
	config                 *domain.Config
	metricsRegistry        *metrics.Registry
	orgService             *OrganizationServiceDB
	principalRepository    repository.Repository[types.Principal]
	groupRepository        repository.Repository[types.Group]
	permissionRepository   repository.Repository[types.Permission]
	relationshipRepository repository.Repository[types.Relationship]
	resourceRepository     repository.Repository[types.Resource]
	roleRepository         repository.Repository[types.Role]
	hashRepository         repository.Repository[domain.HashIndex]
	principalCache         *expirable.LRU[string, *domain.PrincipalExt]
}

// NewPrincipalServiceDB manages persistence of principal data
func NewPrincipalServiceDB(
	config *domain.Config,
	metricsRegistry *metrics.Registry,
	orgService *OrganizationServiceDB,
	principalRepository repository.Repository[types.Principal],
	groupsRepository repository.Repository[types.Group],
	permissionRepository repository.Repository[types.Permission],
	relationshipRepository repository.Repository[types.Relationship],
	resourceRepository repository.Repository[types.Resource],
	roleRepository repository.Repository[types.Role],
	hashRepository repository.Repository[domain.HashIndex],
	maxCacheSize int,
	cacheExpirationMillis int,
) *PrincipalServiceDB {
	return &PrincipalServiceDB{
		config:                 config,
		metricsRegistry:        metricsRegistry,
		orgService:             orgService,
		principalRepository:    principalRepository,
		groupRepository:        groupsRepository,
		permissionRepository:   permissionRepository,
		relationshipRepository: relationshipRepository,
		resourceRepository:     resourceRepository,
		roleRepository:         roleRepository,
		hashRepository:         hashRepository,
		principalCache: expirable.NewLRU[string, *domain.PrincipalExt](
			maxCacheSize,
			nil,
			time.Millisecond*time.Duration(cacheExpirationMillis)),
	}
}

// CreatePrincipal - creates new instance of principal
// Note - this method won't be used to update any role-ids, group-ids, relations, and permission-ids
func (s *PrincipalServiceDB) CreatePrincipal(
	ctx context.Context,
	principal *types.Principal) (*types.Principal, error) {
	defer s.metricsRegistry.Elapsed("principals_svc_create", "org", principal.OrganizationId)()
	xPrincipal := domain.NewPrincipalExt(principal)
	if err := xPrincipal.Validate(); err != nil {
		return nil, err
	}

	principalHash := xPrincipal.Hash()
	hashIndex, _ := s.hashRepository.GetByID(ctx, principal.OrganizationId, "", principalHash)
	if hashIndex != nil {
		return nil, domain.NewDuplicateError(
			fmt.Sprintf("principal with username %s already exists with id %v",
				principal.Username, hashIndex.Ids))
	}

	// initialize
	principal.Id = uuid.NewV4().String()
	principal.Created = timestamppb.Now()
	principal.Version = 1
	principal.GroupIds = make([]string, 0)
	principal.RoleIds = make([]string, 0)
	principal.PermissionIds = make([]string, 0)
	principal.RelationIds = make([]string, 0)

	// update principal in database
	err := s.updatePrincipal(
		ctx,
		0, // first version
		xPrincipal)
	if err != nil {
		return nil, err
	}

	return principal, nil
}

// UpdatePrincipal - updates existing instance of principal
// Note - this method won't be used to update any role-ids, group-ids, relations, and permission-ids
func (s *PrincipalServiceDB) UpdatePrincipal(
	ctx context.Context,
	principal *types.Principal) error {
	defer s.metricsRegistry.Elapsed("principals_svc_update", "org", principal.OrganizationId)()
	xPrincipal := domain.NewPrincipalExt(principal)
	if err := xPrincipal.Validate(); err != nil {
		return err
	}
	if principal.Id == "" {
		return domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}

	existing, err := s.principalRepository.GetByID(
		ctx,
		principal.OrganizationId,
		"", // no namespace
		principal.Id)
	if err != nil {
		return err
	}
	version := principal.Version
	if version == 0 {
		version = existing.Version
	}
	principal.Created = existing.Created
	principal.Version = version + 1
	principal.GroupIds = existing.GroupIds
	principal.RoleIds = existing.RoleIds
	principal.PermissionIds = existing.RoleIds
	principal.RelationIds = existing.RelationIds

	// update principal in database
	return s.updatePrincipal(ctx, version, xPrincipal)
}

// DeletePrincipal removes principal
func (s *PrincipalServiceDB) DeletePrincipal(
	ctx context.Context,
	organizationID string,
	id string) error {
	defer s.metricsRegistry.Elapsed("principals_svc_delete", "org", organizationID)()
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organizationID is not defined"))
	}
	if id == "" {
		return domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}

	principal, err := s.principalRepository.GetByID(ctx, organizationID, "", id) // no namespace
	if err != nil {
		return err
	}
	xPrincipal := domain.NewPrincipalExt(principal)
	principalHash := xPrincipal.Hash()

	if err = s.principalRepository.Delete(
		ctx,
		organizationID,
		"", // no namespace
		id); err != nil {
		return err
	}

	// clear cache
	key := toKey(organizationID, "", id)
	_ = s.principalCache.Remove(key)

	for _, namespace := range principal.Namespaces {
		// ignore errors
		err = s.deletePrincipalAllGroupAndRoleIds(
			ctx,
			organizationID,
			namespace,
			domain.NewPrincipalExt(principal))
		if err != nil {
			log.WithFields(log.Fields{
				"Component":    "PrincipalServiceDB",
				"Organization": organizationID,
				"Id":           id,
				"Error":        err,
			}).
				Warnf("failed to delete groups and role ids")
		}
	}

	// remove mapping between username and principal-id and ignore errors
	err = s.hashRepository.Delete(ctx, xPrincipal.Delegate.OrganizationId, "", principalHash)
	if err != nil {
		log.WithFields(log.Fields{
			"Component":    "PrincipalServiceDB",
			"Organization": organizationID,
			"Id":           id,
			"Error":        err,
		}).
			Warnf("failed to delete mapping between username and principal-id")
	}
	log.WithFields(log.Fields{
		"Component":   "PrincipalServiceDB",
		"PrincipalId": id,
	}).
		Infof("deleted principal")
	return nil
}

// AddGroupsToPrincipal helper
func (s *PrincipalServiceDB) AddGroupsToPrincipal(
	ctx context.Context,
	organizationID string,
	namespace string,
	principalID string,
	groupIDs ...string,
) error {
	defer s.metricsRegistry.Elapsed("principals_svc_add_groups", "org", organizationID)()
	if principalID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("principal-id is not defined for add-group"))
	}
	principal, err := s.principalRepository.GetByID(
		ctx,
		organizationID,
		"", // no namespace
		principalID)
	if err != nil {
		return err
	}
	if !utils.Includes(principal.Namespaces, namespace) {
		return domain.NewValidationError(fmt.Sprintf("namespace %s is not allowed", namespace))
	}
	version := principal.Version
	principal.GroupIds = utils.AddSlice(principal.GroupIds, groupIDs...)
	principal.Version++

	xPrincipal := domain.NewPrincipalExt(principal)
	// update principal
	err = s.updatePrincipal(ctx, version, xPrincipal)
	if err != nil {
		return err
	}

	//update role/group-ids and clear cache
	return s.updatePrincipalAllGroupAndRoleIds(ctx, namespace, xPrincipal)
}

// DeleteGroupsToPrincipal helper
func (s *PrincipalServiceDB) DeleteGroupsToPrincipal(
	ctx context.Context,
	organizationID string,
	namespace string,
	principalID string,
	groupIDs ...string,
) error {
	defer s.metricsRegistry.Elapsed("principals_svc_delete_groups", "org", organizationID)()
	if principalID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("principal-id is not defined for delete-group"))
	}
	principal, err := s.principalRepository.GetByID(
		ctx,
		organizationID,
		"", // no namespace
		principalID)
	if err != nil {
		return err
	}

	if !utils.Includes(principal.Namespaces, namespace) {
		return domain.NewValidationError(fmt.Sprintf("namespace %s is not allowed", namespace))
	}
	version := principal.Version
	principal.GroupIds = utils.RemoveSlice(principal.GroupIds, groupIDs...)
	principal.Version++

	xPrincipal := domain.NewPrincipalExt(principal)
	// update principal
	err = s.updatePrincipal(ctx, version, xPrincipal)
	if err != nil {
		return err
	}

	//update role/group-ids and clear cache
	return s.updatePrincipalAllGroupAndRoleIds(ctx, namespace, xPrincipal)
}

// AddRolesToPrincipal helper
func (s *PrincipalServiceDB) AddRolesToPrincipal(
	ctx context.Context,
	organizationID string,
	namespace string,
	principalID string,
	roleIDs ...string,
) error {
	defer s.metricsRegistry.Elapsed("principals_svc_add_roles", "org", organizationID)()
	if principalID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("principal-id is not defined for add-role"))
	}
	principal, err := s.principalRepository.GetByID(
		ctx,
		organizationID,
		"", // no namespace
		principalID)
	if err != nil {
		return err
	}
	if !utils.Includes(principal.Namespaces, namespace) {
		return domain.NewValidationError(fmt.Sprintf("namespace %s is not allowed", namespace))
	}
	version := principal.Version
	principal.RoleIds = utils.AddSlice(principal.RoleIds, roleIDs...)
	principal.Version++

	xPrincipal := domain.NewPrincipalExt(principal)
	// update principal
	err = s.updatePrincipal(ctx, version, xPrincipal)
	if err != nil {
		return err
	}
	//update role/group-ids and clear cache
	return s.updatePrincipalAllGroupAndRoleIds(ctx, namespace, xPrincipal)
}

// DeleteRolesToPrincipal helper
func (s *PrincipalServiceDB) DeleteRolesToPrincipal(
	ctx context.Context,
	organizationID string,
	namespace string,
	principalID string,
	roleIDs ...string,
) error {
	defer s.metricsRegistry.Elapsed("principals_svc_delete_roles", "org", organizationID)()
	if principalID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("principal-id is not defined for delete-role"))
	}
	principal, err := s.principalRepository.GetByID(
		ctx,
		organizationID,
		"", // no namespace
		principalID)
	if err != nil {
		return err
	}
	if !utils.Includes(principal.Namespaces, namespace) {
		return domain.NewValidationError(fmt.Sprintf("namespace %s is not allowed", namespace))
	}
	version := principal.Version
	principal.RoleIds = utils.RemoveSlice(principal.RoleIds, roleIDs...)
	principal.Version++

	xPrincipal := domain.NewPrincipalExt(principal)
	// update principal
	err = s.updatePrincipal(ctx, version, xPrincipal)
	if err != nil {
		return err
	}
	//update role/group-ids and clear cache
	return s.updatePrincipalAllGroupAndRoleIds(ctx, namespace, xPrincipal)
}

// AddPermissionsToPrincipal helper
func (s *PrincipalServiceDB) AddPermissionsToPrincipal(
	ctx context.Context,
	organizationID string,
	namespace string,
	principalID string,
	permissionIds ...string,
) error {
	defer s.metricsRegistry.Elapsed("principals_svc_add_permissions", "org", organizationID)()
	if principalID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("principal-id is not defined for add-permission"))
	}
	principal, err := s.principalRepository.GetByID(
		ctx,
		organizationID,
		"", // no namespace
		principalID)
	if err != nil {
		return err
	}
	if !utils.Includes(principal.Namespaces, namespace) {
		return domain.NewValidationError(fmt.Sprintf("namespace %s is not allowed", namespace))
	}
	version := principal.Version
	principal.PermissionIds = utils.AddSlice(principal.PermissionIds, permissionIds...)
	principal.Version++

	xPrincipal := domain.NewPrincipalExt(principal)
	// update principal
	err = s.updatePrincipal(ctx, version, xPrincipal)
	if err != nil {
		return err
	}
	//update role/group-ids and clear cache
	return s.updatePrincipalAllGroupAndRoleIds(ctx, namespace, xPrincipal)
}

// DeletePermissionsToPrincipal helper
func (s *PrincipalServiceDB) DeletePermissionsToPrincipal(
	ctx context.Context,
	organizationID string,
	namespace string,
	principalID string,
	permissionIds ...string,
) error {
	defer s.metricsRegistry.Elapsed("principals_svc_delete_permissions", "org", organizationID)()
	if principalID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("principal-id is not defined for delete-permission"))
	}
	principal, err := s.principalRepository.GetByID(
		ctx,
		organizationID,
		"", // no namespace
		principalID)
	if err != nil {
		return err
	}
	if !utils.Includes(principal.Namespaces, namespace) {
		return domain.NewValidationError(fmt.Sprintf("namespace %s is not allowed", namespace))
	}
	version := principal.Version
	principal.PermissionIds = utils.RemoveSlice(principal.PermissionIds, permissionIds...)
	principal.Version++

	xPrincipal := domain.NewPrincipalExt(principal)
	// update principal
	err = s.updatePrincipal(ctx, version, xPrincipal)
	if err != nil {
		return err
	}
	//update role/group-ids and clear cache
	return s.updatePrincipalAllGroupAndRoleIds(ctx, namespace, xPrincipal)
}

// AddRelationshipsToPrincipal helper
func (s *PrincipalServiceDB) AddRelationshipsToPrincipal(
	ctx context.Context,
	organizationID string,
	namespace string,
	principalID string,
	relationshipIds ...string,
) error {
	defer s.metricsRegistry.Elapsed("principals_svc_add_relations", "org", organizationID)()
	if principalID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("principal-id is not defined for add-relation"))
	}
	principal, err := s.principalRepository.GetByID(
		ctx,
		organizationID,
		"", // no namespace
		principalID)
	if err != nil {
		return err
	}
	if !utils.Includes(principal.Namespaces, namespace) {
		return domain.NewValidationError(fmt.Sprintf("namespace %s is not allowed", namespace))
	}
	version := principal.Version
	principal.RelationIds = utils.AddSlice(principal.RelationIds, relationshipIds...)
	principal.Version++

	xPrincipal := domain.NewPrincipalExt(principal)
	// update principal
	err = s.updatePrincipal(ctx, version, xPrincipal)
	if err != nil {
		return err
	}
	//update role/group-ids and clear cache
	return s.updatePrincipalAllGroupAndRoleIds(ctx, namespace, xPrincipal)
}

// DeleteRelationshipsToPrincipal helper
func (s *PrincipalServiceDB) DeleteRelationshipsToPrincipal(
	ctx context.Context,
	organizationID string,
	namespace string,
	principalID string,
	relationshipIds ...string,
) error {
	defer s.metricsRegistry.Elapsed("principals_svc_delete_relations", "org", organizationID)()
	if principalID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("principal-id is not defined for delete-relation"))
	}
	principal, err := s.principalRepository.GetByID(
		ctx,
		organizationID,
		"", // no namespace
		principalID)
	if err != nil {
		return err
	}
	if !utils.Includes(principal.Namespaces, namespace) {
		return domain.NewValidationError(fmt.Sprintf("namespace %s is not allowed", namespace))
	}
	version := principal.Version
	principal.RelationIds = utils.RemoveSlice(principal.RelationIds, relationshipIds...)
	principal.Version++

	xPrincipal := domain.NewPrincipalExt(principal)
	// update principal
	err = s.updatePrincipal(ctx, version, xPrincipal)
	if err != nil {
		return err
	}
	//update role/group-ids and clear cache
	return s.updatePrincipalAllGroupAndRoleIds(ctx, namespace, xPrincipal)
}

// GetPrincipal - retrieves principal
func (s *PrincipalServiceDB) GetPrincipal(
	ctx context.Context,
	organizationID string,
	id string,
) (*types.Principal, error) {
	defer s.metricsRegistry.Elapsed("principals_svc_get", "org", organizationID)()
	if organizationID == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("organization_id is not defined"))
	}

	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}

	// check cache
	key := toKey(organizationID, "", id) // no namespace for key
	xPrincipal, _ := s.principalCache.Get(key)
	if xPrincipal != nil {
		return xPrincipal.Delegate, nil
	}

	// load from database
	return s.principalRepository.GetByID(ctx, organizationID, "", id) // no namespace for principal
}

// GetPrincipals - queries principals
func (s *PrincipalServiceDB) GetPrincipals(
	ctx context.Context,
	organizationID string,
	predicates map[string]string,
	offset string,
	limit int64) (res []*types.Principal, nextToken string, err error) {
	defer s.metricsRegistry.Elapsed("principals_svc_query", "org", organizationID)()
	if predicates["id"] != "" {
		principal, err := s.GetPrincipal(ctx, organizationID, predicates["id"])
		if err != nil {
			return nil, "", err
		}
		return []*types.Principal{principal}, "", nil
	}
	if organizationID == "" {
		return nil, "", domain.NewValidationError(
			fmt.Sprintf("organizationID is not defined"))
	}

	return s.principalRepository.Query(
		ctx,
		organizationID,
		"", // no namespace
		predicates,
		offset, limit)
}

// GetPrincipalExt - retrieves full principal
func (s *PrincipalServiceDB) GetPrincipalExt(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string,
) (xPrincipal *domain.PrincipalExt, err error) {
	defer s.metricsRegistry.Elapsed("principals_svc_getx", "org", organizationID)()
	org, err := s.orgService.verifyOrganizationNamespace(ctx, organizationID, namespace)
	if err != nil {
		return nil, err
	}
	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}

	// check cache
	key := toKey(organizationID, "", id) // no namespace for key
	xPrincipal, _ = s.principalCache.Get(key)
	if xPrincipal != nil {
		return xPrincipal, nil
	}

	// load from database
	principal, err := s.principalRepository.GetByID(ctx, organizationID, "", id) // no namespace for principal
	if err != nil {
		return nil, err
	}
	if !utils.Includes(principal.Namespaces, namespace) {
		return nil, domain.NewValidationError(fmt.Sprintf("namespace %s is not allowed", namespace))
	}
	xPrincipal = domain.NewPrincipalExt(principal)
	xPrincipal.Organization = org

	// get cache of all group and role ids
	groupIDs, roleIDs, updatedGroupRoleIds := s.getPrincipalAllGroupAndRoleIds(
		ctx,
		organizationID,
		namespace,
		xPrincipal)

	// populate groups
	if len(groupIDs) > 0 {
		groups, err := s.groupRepository.GetByIDs(ctx, organizationID, namespace, groupIDs...)
		if err != nil {
			return nil, err
		}
		for _, group := range groups {
			xPrincipal.GroupsByName[group.Name] = group
			for _, nextRoleId := range group.RoleIds {
				roleIDs = utils.AddSlice(roleIDs, nextRoleId)
			}
		}
	}

	// populate roles
	if len(roleIDs) > 0 {
		roles, err := s.roleRepository.GetByIDs(ctx, organizationID, namespace, roleIDs...)
		if err != nil {
			return nil, err
		}
		for _, role := range roles {
			xPrincipal.RolesByName[role.Name] = role
		}
	}

	// populate permissions
	if len(principal.PermissionIds) > 0 {
		err := s.populatePermissions(ctx, organizationID, namespace, xPrincipal, principal.PermissionIds...)
		if err != nil {
			return nil, err
		}
	}

	for _, role := range xPrincipal.RolesByName {
		err = s.populatePermissions(ctx, organizationID, namespace, xPrincipal, role.PermissionIds...)
		if err != nil {
			return nil, err
		}
	}

	// populate relationships
	if len(principal.RelationIds) > 0 {
		relations, err := s.relationshipRepository.GetByIDs(ctx, organizationID, namespace, principal.RelationIds...)
		if err != nil {
			return nil, err
		}
		for _, rel := range relations {
			xPrincipal.RelationsById[rel.Id] = rel
		}
	}

	latestGroupRoleDate := xPrincipal.LatestGroupRoleDate()
	if updatedGroupRoleIds == nil || (latestGroupRoleDate != nil && updatedGroupRoleIds != nil &&
		latestGroupRoleDate.Seconds > updatedGroupRoleIds.Seconds) {
		// role or group was updated, so we need to populate cache again
		err = s.updatePrincipalAllGroupAndRoleIds(
			ctx,
			namespace,
			xPrincipal)
	}
	s.principalCache.Add(key, xPrincipal)
	return
}

func (s *PrincipalServiceDB) updatePrincipal(
	ctx context.Context,
	version int64,
	xPrincipal *domain.PrincipalExt,
) (err error) {
	// Verify organization
	org, err := s.orgService.GetOrganization(
		ctx,
		xPrincipal.Delegate.OrganizationId)
	if err != nil {
		return err
	}

	for _, namespace := range xPrincipal.Delegate.Namespaces {
		if !utils.Includes(org.Namespaces, namespace) {
			return domain.NewValidationError(fmt.Sprintf("namespace %s is not allowed", namespace))
		}
	}

	// save principal
	xPrincipal.Delegate.Updated = timestamppb.Now()

	if version == 0 {
		err = s.principalRepository.Create(
			ctx,
			xPrincipal.Delegate.OrganizationId,
			"", // namespace
			xPrincipal.Delegate.Id,
			xPrincipal.Delegate,
			time.Duration(0),
		)
	} else {
		err = s.principalRepository.Update(
			ctx,
			xPrincipal.Delegate.OrganizationId,
			"", // namespace
			xPrincipal.Delegate.Id,
			version,
			xPrincipal.Delegate,
			time.Duration(0),
		)
	}
	if err != nil {
		return err
	}

	// clear cache
	key := toKey(xPrincipal.Delegate.OrganizationId, "", xPrincipal.Delegate.Id)
	_ = s.principalCache.Remove(key)

	principalHash := xPrincipal.Hash()
	// update mapping between username and principal-id
	if log.IsLevelEnabled(log.DebugLevel) {
		log.WithFields(log.Fields{
			"Component":    "PrincipalServiceDB",
			"Organization": xPrincipal.Delegate.OrganizationId,
			"Principal":    xPrincipal.Delegate.Id,
		}).
			Debugf("updated principal")
	}
	return s.hashRepository.Update(
		ctx,
		xPrincipal.Delegate.OrganizationId,
		"", // no namespace,
		principalHash,
		-1, // no version
		domain.NewHashIndex(principalHash, []string{xPrincipal.Delegate.Id}),
		time.Duration(0),
	)
}

// updatePrincipalAllGroupAndRoleIds
func (s *PrincipalServiceDB) updatePrincipalAllGroupAndRoleIds(
	ctx context.Context,
	namespace string,
	xPrincipal *domain.PrincipalExt,
) error {
	allGroupIds := make(map[string]int32)
	err := s.populateAllGroups(
		ctx,
		xPrincipal.Delegate.OrganizationId,
		namespace,
		allGroupIds,
		0,
		xPrincipal.Delegate.GroupIds...,
	)
	if err != nil {
		return err
	}
	allRoleIds := make(map[string]int32)
	err = s.populateAllRoles(
		ctx,
		xPrincipal.Delegate.OrganizationId,
		namespace,
		allRoleIds,
		0,
		xPrincipal.Delegate.RoleIds...,
	)
	if err != nil {
		return err
	}

	err = s.hashRepository.Update(
		ctx,
		xPrincipal.Delegate.OrganizationId,
		namespace,
		xPrincipal.GroupHashIndex(),
		-1, // no version
		domain.NewHashIndex(xPrincipal.GroupHashIndex(), utils.MapIntToArray(allGroupIds)),
		time.Duration(0),
	)
	if err != nil {
		return err
	}
	return s.hashRepository.Update(
		ctx,
		xPrincipal.Delegate.OrganizationId,
		namespace,
		xPrincipal.RoleHashIndex(),
		-1, // no version
		domain.NewHashIndex(xPrincipal.RoleHashIndex(), utils.MapIntToArray(allRoleIds)),
		time.Duration(0),
	)
}

func (s *PrincipalServiceDB) getPrincipalAllGroupAndRoleIds(
	ctx context.Context,
	organizationID string,
	namespace string,
	xPrincipal *domain.PrincipalExt,
) (groupIndexIds []string, roleIndexIds []string, updated *timestamppb.Timestamp) {
	groupIndex, _ := s.hashRepository.GetByID(ctx, organizationID, namespace, xPrincipal.GroupHashIndex())
	roleIndex, _ := s.hashRepository.GetByID(ctx, organizationID, namespace, xPrincipal.RoleHashIndex())
	if groupIndex != nil {
		groupIndexIds = groupIndex.Ids
		updated = groupIndex.Updated
	}
	if roleIndex != nil {
		roleIndexIds = roleIndex.Ids
		if updated == nil || updated.Seconds > roleIndex.Updated.Seconds {
			updated = roleIndex.Updated
		}
	}
	return
}

func (s *PrincipalServiceDB) deletePrincipalAllGroupAndRoleIds(
	ctx context.Context,
	organizationID string,
	namespace string,
	xPrincipal *domain.PrincipalExt,
) error {
	err1 := s.hashRepository.Delete(ctx, organizationID, namespace, xPrincipal.GroupHashIndex())
	err2 := s.hashRepository.Delete(ctx, organizationID, namespace, xPrincipal.RoleHashIndex())
	if err1 != nil {
		return err1
	}
	return err2
}

// populatePermissions
func (s *PrincipalServiceDB) populatePermissions(
	ctx context.Context,
	organizationID string,
	namespace string,
	xPrincipal *domain.PrincipalExt,
	permissionIds ...string,
) error {
	if len(permissionIds) == 0 {
		return nil
	}
	permissions, err := s.permissionRepository.GetByIDs(ctx, organizationID, namespace, permissionIds...)
	if err != nil {
		return err
	}
	var resourceIDs []string
	for _, perm := range permissions {
		resourceIDs = utils.AddSlice(resourceIDs, perm.ResourceId)
	}
	resources, err := s.resourceRepository.GetByIDs(ctx, organizationID, namespace, resourceIDs...)
	if err != nil {
		return err
	}
	for id, resource := range resources {
		xPrincipal.ResourcesById[id] = resource
	}

	for _, perm := range permissions {
		if err = xPrincipal.AddPermission(perm); err != nil {
			return err
		}
	}
	return nil
}

// populateAllGroups
func (s *PrincipalServiceDB) populateAllGroups(
	ctx context.Context,
	organizationID string,
	namespace string,
	allGroupIds map[string]int32,
	level int,
	groupIDs ...string,
) error {
	if len(groupIDs) == 0 {
		return nil
	}
	if level > s.config.MaxGroupRoleLevels {
		return nil
	}
	for _, nextId := range groupIDs {
		group, err := s.groupRepository.GetByID(ctx, organizationID, namespace, nextId)
		if err != nil {
			return err
		}
		allGroupIds[nextId] = allGroupIds[nextId] + 1
		for _, parentId := range group.ParentIds {
			err = s.populateAllGroups(ctx, organizationID, namespace, allGroupIds, level+1, parentId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// populateAllRoles
func (s *PrincipalServiceDB) populateAllRoles(
	ctx context.Context,
	organizationID string,
	namespace string,
	allRoleIds map[string]int32,
	level int,
	roleIDs ...string,
) error {
	if len(roleIDs) == 0 {
		return nil
	}
	if level > s.config.MaxGroupRoleLevels {
		return nil
	}
	for _, nextId := range roleIDs {
		role, err := s.roleRepository.GetByID(ctx, organizationID, namespace, nextId)
		if err != nil {
			return err
		}
		allRoleIds[nextId] = allRoleIds[nextId] + 1
		for _, parentId := range role.ParentIds {
			err = s.populateAllRoles(ctx, organizationID, namespace, allRoleIds, level+1, parentId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
