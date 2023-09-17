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

func Test_ShouldSaveAndGetGroup(t *testing.T) {
	// GIVEN config, redis-service and group repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	repository, err := NewGroupRepository(store)
	require.NoError(t, err)
	group := buildTestGroup(1)
	namespace := "group-save-namespace"
	testOrgId := uuid.NewV4().String()
	err = repository.Create(ctx, testOrgId, namespace, group.Id, &group, time.Duration(0))
	require.NoError(t, err)

	saved, err := repository.GetByID(ctx, testOrgId, namespace, group.Id)
	require.NoError(t, err)
	require.Equal(t, group.Name, saved.Name)
	require.Equal(t, 2, len(saved.RoleIds))

	saved.Name = "new-name"
	err = repository.Update(ctx, testOrgId, namespace, group.Id, group.Version, saved, time.Duration(0))
	require.NoError(t, err)

	saved, err = repository.GetByID(ctx, testOrgId, namespace, group.Id)
	require.NoError(t, err)
	require.Equal(t, "new-name", saved.Name)
	require.Equal(t, 2, len(saved.RoleIds))

	err = store.ClearTable("Group", "", testOrgId, namespace)
	require.NoError(t, err)
}

func Test_ShouldSaveAndDeleteGroup(t *testing.T) {
	// GIVEN config, redis-service and group repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	repository, err := NewGroupRepository(store)
	require.NoError(t, err)
	group := buildTestGroup(1)

	testOrgId := uuid.NewV4().String()
	namespace := "group-del-namespace"

	err = repository.Create(ctx, testOrgId, namespace, group.Id, &group, time.Duration(0))
	require.NoError(t, err)

	_, err = repository.GetByIDs(ctx, testOrgId, namespace, group.Id)
	require.NoError(t, err)

	err = repository.Delete(ctx, testOrgId, namespace, group.Id)
	require.NoError(t, err)

	_, err = repository.GetByIDs(ctx, testOrgId, namespace, group.Id)
	require.Error(t, err)

	err = store.ClearTable("Group", "", testOrgId, namespace)
	require.NoError(t, err)
}

func Test_ShouldSaveAndQueryGroup(t *testing.T) {
	// GIVEN config, redis-service and group repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	testOrgId := uuid.NewV4().String()
	namespace := "group-query-namespace"
	err = store.ClearTable("Group", "", testOrgId, namespace)
	require.NoError(t, err)
	repository, err := NewGroupRepository(store)
	require.NoError(t, err)
	for i := 0; i < 200; i++ {
		group := buildTestGroup(i)
		err = repository.Create(ctx, testOrgId, namespace, group.Id, &group, time.Duration(0))
		require.NoError(t, err)
	}
	res, _, err := repository.Query(ctx, testOrgId, namespace, nil, "0", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"id": "id_0"}, "0", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"name:>=": "name_"}, "0", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{}, "100", 10)
	require.NoError(t, err)
	require.Equal(t, 10, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{}, "200", 10)
	require.NoError(t, err)
	require.Equal(t, 0, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"name:==": "name_0"}, "0", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"name:!=": "1111"}, "0", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"name:>=": "name_199"}, "0", 0)
	require.NoError(t, err)
	require.Equal(t, 89, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"name:<=": "name_1"}, "0", 0)
	require.NoError(t, err)
	require.Equal(t, 2, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"name:<": "name_1"}, "0", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"name:>": "name_198"}, "0", 0)
	require.NoError(t, err)
	require.Equal(t, 89, len(res))

	count, err := repository.Size(ctx, testOrgId, namespace)
	require.NoError(t, err)
	require.True(t, count > 0)
}

func buildTestGroup(i int) types.Group {
	return types.Group{
		Id:      fmt.Sprintf("id_%d", i),
		Name:    fmt.Sprintf("name_%d", i),
		RoleIds: []string{"1", "2"},
	}
}
