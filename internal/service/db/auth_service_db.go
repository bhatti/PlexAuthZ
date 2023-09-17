package db

import (
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/metrics"
	"github.com/bhatti/PlexAuthZ/internal/repository"
)

// authAdminServiceDB - manages persistence of AuthZ data.
type authAdminServiceDB struct {
	*OrganizationServiceDB // implementation for organization admin service
	*PrincipalServiceDB    // implementation for principal admin service
	*ResourceServiceDB     // implementation for resource admin service
	*PermissionServiceDB   // implementation for permission admin service
	*RoleServiceDB         // implementation for role admin service
	*GroupServiceDB        // implementation for group admin service
	*RelationshipServiceDB // implementation for relationships admin service
}

// NewAuthAdminServiceDB manages persistence of AuthZ data
func NewAuthAdminServiceDB(
	config *domain.Config,
	metricsRegistry *metrics.Registry,
	orgRepository repository.Repository[types.Organization],
	principalRepository repository.Repository[types.Principal],
	groupsRepository repository.Repository[types.Group],
	permissionRepository repository.Repository[types.Permission],
	relationshipRepository repository.Repository[types.Relationship],
	resourceRepository repository.Repository[types.Resource],
	resourceInstanceRepositoryFactory repository.ResourceInstanceRepositoryFactory,
	roleRepository repository.Repository[types.Role],
	hashRepository repository.Repository[domain.HashIndex],
	maxCacheSize int,
	cacheExpirationMillis int,
) *authAdminServiceDB {
	orgService := NewOrganizationServiceDB(metricsRegistry, orgRepository, maxCacheSize, cacheExpirationMillis)
	principalService := NewPrincipalServiceDB(
		config,
		metricsRegistry,
		orgService,
		principalRepository,
		groupsRepository,
		permissionRepository,
		relationshipRepository,
		resourceRepository,
		roleRepository,
		hashRepository,
		maxCacheSize,
		cacheExpirationMillis)
	resourceService := NewResourceServiceDB(
		metricsRegistry,
		orgService,
		principalService,
		resourceRepository,
		resourceInstanceRepositoryFactory,
		hashRepository)
	permissionService := NewPermissionServiceDB(
		metricsRegistry,
		orgService,
		resourceRepository,
		permissionRepository,
		hashRepository)
	roleService := NewRoleServiceDB(
		metricsRegistry,
		orgService,
		roleRepository,
		hashRepository)
	groupService := NewGroupServiceDB(
		metricsRegistry,
		orgService,
		groupsRepository,
		hashRepository)
	relationshipService := NewRelationshipServiceDB(
		metricsRegistry,
		orgService,
		relationshipRepository,
		hashRepository)
	return &authAdminServiceDB{
		OrganizationServiceDB: orgService,
		PrincipalServiceDB:    principalService,
		ResourceServiceDB:     resourceService,
		PermissionServiceDB:   permissionService,
		RoleServiceDB:         roleService,
		GroupServiceDB:        groupService,
		RelationshipServiceDB: relationshipService,
	}
}

// Close no-op
func (s *authAdminServiceDB) Close() error {
	return nil
}
