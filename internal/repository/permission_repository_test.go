package repository

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/repository/redis"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
	"testing"
	"time"
)

func Test_ShouldSaveAndGetPermission(t *testing.T) {
	// GIVEN config, redis-service and permission repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	repository, err := NewPermissionRepository(store)
	require.NoError(t, err)
	permission := buildTestPermission(1)
	testOrgId := uuid.NewV4().String()
	namespace := "permission-save-namespace"
	err = repository.Create(ctx, testOrgId, namespace, permission.Id, &permission, time.Duration(0))
	require.NoError(t, err)

	saved, err := repository.GetByID(ctx, testOrgId, namespace, permission.Id)
	require.NoError(t, err)
	require.Equal(t, permission.Namespace, saved.Namespace)
	require.Equal(t, permission.Scope, saved.Scope)
	require.Equal(t, permission.ResourceId, saved.ResourceId)
	require.Equal(t, permission.Effect, saved.Effect)
	require.Equal(t, permission.Constraints, saved.Constraints)
	require.Equal(t, 2, len(saved.Actions))

	saved.Namespace = "ns2"
	err = repository.Update(ctx, testOrgId, namespace, permission.Id, saved.Version, saved, time.Duration(0))
	require.NoError(t, err)

	saved, err = repository.GetByID(ctx, testOrgId, namespace, permission.Id)
	require.NoError(t, err)
	require.Equal(t, "ns2", saved.Namespace)

	err = store.ClearTable("Permission", "", testOrgId, namespace)
	require.NoError(t, err)
}

func Test_ShouldSaveAndDeletePermission(t *testing.T) {
	// GIVEN config, redis-service and permission repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	repository, err := NewPermissionRepository(store)
	require.NoError(t, err)
	permission := buildTestPermission(1)

	testOrgId := uuid.NewV4().String()
	namespace := "permission-del-namespace"

	err = repository.Create(ctx, testOrgId, namespace, permission.Id, &permission, time.Duration(0))
	require.NoError(t, err)

	_, err = repository.GetByIDs(ctx, testOrgId, namespace, permission.Id)
	require.NoError(t, err)

	err = repository.Delete(ctx, testOrgId, namespace, permission.Id)
	require.NoError(t, err)

	_, err = repository.GetByIDs(ctx, testOrgId, namespace, permission.Id)
	require.Error(t, err)

	err = store.ClearTable("Permission", "", testOrgId, namespace)
	require.NoError(t, err)
}

func Test_ShouldSaveAndQueryPermission(t *testing.T) {
	// GIVEN config, redis-service and permission repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	testOrgId := uuid.NewV4().String()
	namespace := "permission-query-namespace"
	err = store.ClearTable("Permission", "", testOrgId, namespace)
	require.NoError(t, err)
	repository, err := NewPermissionRepository(store)
	require.NoError(t, err)
	for i := 0; i < 200; i++ {
		permission := buildTestPermission(i)
		err = repository.Create(ctx, testOrgId, namespace, permission.Id, &permission, time.Duration(0))
		require.NoError(t, err)
	}
	res, _, err := repository.Query(ctx, testOrgId, namespace, nil, "", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"id": "id_0"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"namespace:>=": "ns_"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{}, "100", 10)
	require.NoError(t, err)
	require.Equal(t, 10, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{}, "200", 10)
	require.NoError(t, err)
	require.Equal(t, 0, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"namespace:==": "ns_0"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"namespace:!=": "1111"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"namespace:>=": "ns_199"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 89, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"namespace:<=": "ns_1"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 2, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"namespace:<": "ns_1"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"namespace:>": "ns_198"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 89, len(res))
}

func buildTestPermission(i int) types.Permission {
	return types.Permission{
		Id:          fmt.Sprintf("id_%d", i),
		Namespace:   fmt.Sprintf("ns_%d", i),
		Scope:       fmt.Sprintf("scope_%d", i),
		Actions:     []string{"read", "write"},
		ResourceId:  "1",
		Effect:      types.Effect_PERMITTED,
		Constraints: "time > 10",
	}
}
