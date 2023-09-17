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

func Test_ShouldSaveAndGetResource(t *testing.T) {
	// GIVEN config, redis-service and resource repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	repository, err := NewResourceRepository(store)
	require.NoError(t, err)
	resource := buildTestResource(1)
	testOrgId := uuid.NewV4().String()
	namespace := "resource-save-namespace"
	err = repository.Create(ctx, testOrgId, namespace, resource.Id, &resource, time.Duration(0))
	require.NoError(t, err)

	saved, err := repository.GetByID(ctx, testOrgId, namespace, resource.Id)
	require.NoError(t, err)
	require.Equal(t, resource.Name, saved.Name)
	require.Equal(t, resource.Capacity, saved.Capacity)
	require.Equal(t, resource.AllowedActions, saved.AllowedActions)

	saved.Name = "new-name"
	err = repository.Update(ctx, testOrgId, namespace, resource.Id, saved.Version, saved, time.Duration(0))
	require.NoError(t, err)

	saved, err = repository.GetByID(ctx, testOrgId, namespace, resource.Id)
	require.NoError(t, err)
	require.Equal(t, "new-name", saved.Name)

	err = store.ClearTable("Resource", "", testOrgId, namespace)
	require.NoError(t, err)
}

func Test_ShouldSaveAndDeleteResource(t *testing.T) {
	// GIVEN config, redis-service and resource repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	repository, err := NewResourceRepository(store)
	require.NoError(t, err)
	resource := buildTestResource(1)
	testOrgId := uuid.NewV4().String()
	namespace := "resource-del-namespace"

	err = repository.Create(ctx, testOrgId, namespace, resource.Id, &resource, time.Duration(0))
	require.NoError(t, err)

	_, err = repository.GetByIDs(ctx, testOrgId, namespace, resource.Id)
	require.NoError(t, err)

	err = repository.Delete(ctx, testOrgId, namespace, resource.Id)
	require.NoError(t, err)

	_, err = repository.GetByIDs(ctx, testOrgId, namespace, resource.Id)
	require.Error(t, err)

	err = store.ClearTable("Resource", "", testOrgId, namespace)
	require.NoError(t, err)
}

func Test_ShouldSaveAndQueryResource(t *testing.T) {
	// GIVEN config, redis-service and resource repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	testOrgId := uuid.NewV4().String()
	namespace := "resource-query-namespace"
	repository, err := NewResourceRepository(store)
	require.NoError(t, err)
	for i := 0; i < 200; i++ {
		resource := buildTestResource(i)
		err = repository.Create(ctx, testOrgId, namespace, resource.Id, &resource, time.Duration(0))
		require.NoError(t, err)
	}
	res, _, err := repository.Query(ctx, testOrgId, namespace, nil, "", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"id": "id_0"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	require.Equal(t, int32(1), res[0].Capacity)
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{}, "100", 10)
	require.NoError(t, err)
	require.Equal(t, 10, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{}, "200", 10)
	require.NoError(t, err)
	require.Equal(t, 0, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"capacity:==": "1"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"capacity:!=": "1111"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"capacity:>=": "100"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 101, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"capacity:<=": "100"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 100, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"capacity:<": "100"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 99, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"capacity:>": "100"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 100, len(res))
	err = store.ClearTable("Resource", "", testOrgId, namespace)
	require.NoError(t, err)
}

func buildTestResource(i int) types.Resource {
	return types.Resource{
		Id:             fmt.Sprintf("id_%d", i),
		Name:           fmt.Sprintf("/file/%d", i),
		Capacity:       int32(i + 1),
		Attributes:     make(map[string]string),
		AllowedActions: []string{"read", "write"},
	}
}
