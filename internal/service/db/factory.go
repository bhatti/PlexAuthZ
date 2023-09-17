package db

import (
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/metrics"
	"github.com/bhatti/PlexAuthZ/internal/repository"
	"github.com/bhatti/PlexAuthZ/internal/repository/ddb"
	"github.com/bhatti/PlexAuthZ/internal/repository/redis"
	"github.com/bhatti/PlexAuthZ/internal/service"
	"time"
)

// CreateDatabaseAuthService factory method
func CreateDatabaseAuthService(
	cfg *domain.Config,
	metricsRegistry *metrics.Registry,
) (service.AuthAdminService, domain.Closeable, error) {
	store, err := CreateDataStore(cfg)
	if err != nil {
		return nil, nil, err
	}
	orgRepository, err := repository.NewOrganizationRepository(store)
	if err != nil {
		return nil, nil, err
	}
	principalRepository, err := repository.NewPrincipalRepository(store)
	if err != nil {
		return nil, nil, err
	}
	groupRepository, err := repository.NewGroupRepository(store)
	if err != nil {
		return nil, nil, err
	}
	permissionRepository, err := repository.NewPermissionRepository(store)
	if err != nil {
		return nil, nil, err
	}
	relationshipRepository, err := repository.NewRelationshipRepository(store)
	if err != nil {
		return nil, nil, err
	}
	resourceRepository, err := repository.NewResourceRepository(store)
	if err != nil {
		return nil, nil, err
	}
	resourceInstanceRepository, err := repository.NewResourceInstanceRepository(store, time.Minute)
	if err != nil {
		return nil, nil, err
	}
	roleRepository, err := repository.NewRoleRepository(store)
	if err != nil {
		return nil, nil, err
	}
	hashRepository, err := repository.NewHashIndexRepository(store)
	if err != nil {
		return nil, nil, err
	}
	authService := NewAuthAdminServiceDB(
		cfg,
		metricsRegistry,
		orgRepository,
		principalRepository,
		groupRepository,
		permissionRepository,
		relationshipRepository,
		resourceRepository,
		resourceInstanceRepository,
		roleRepository,
		hashRepository,
		cfg.MaxCacheSize,
		cfg.CacheExpirationMillis,
	)
	return authService, authService, nil
}

// CreateDataStore factory
func CreateDataStore(cfg *domain.Config) (repository.DataStore, error) {
	if cfg.PersistenceProvider == domain.DynamoDBPersistenceProvider {
		return ddb.NewDDBStore(cfg)
	}
	return redis.NewRedisStore(cfg)
}
