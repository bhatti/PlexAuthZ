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

func Test_ShouldSaveAndGetRole(t *testing.T) {
	// GIVEN config, redis-service and role repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	repository, err := NewRoleRepository(store)
	require.NoError(t, err)
	role := buildTestRole(1)
	testOrgId := uuid.NewV4().String()
	namespace := "role-save-namespace"
	err = repository.Create(ctx, testOrgId, namespace, role.Id, &role, time.Duration(0))
	require.NoError(t, err)

	saved, err := repository.GetByID(ctx, testOrgId, namespace, role.Id)
	require.NoError(t, err)
	require.Equal(t, role.Name, saved.Name)
	require.Equal(t, 2, len(role.PermissionIds))
	require.Equal(t, 2, len(role.ParentIds))

	saved.Name = "new-role"
	err = repository.Update(ctx, testOrgId, namespace, role.Id, saved.Version, saved, time.Duration(0))
	require.NoError(t, err)

	saved, err = repository.GetByID(ctx, testOrgId, namespace, role.Id)
	require.NoError(t, err)
	require.Equal(t, "new-role", saved.Name)

	err = store.ClearTable("Role", "", testOrgId, namespace)
	require.NoError(t, err)
}

func Test_ShouldSaveAndDeleteRole(t *testing.T) {
	// GIVEN config, redis-service and role repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	repository, err := NewRoleRepository(store)
	require.NoError(t, err)
	role := buildTestRole(1)
	testOrgId := uuid.NewV4().String()
	namespace := "role-del-namespace"

	err = repository.Create(ctx, testOrgId, namespace, role.Id, &role, time.Duration(0))
	require.NoError(t, err)

	_, err = repository.GetByIDs(ctx, testOrgId, namespace, role.Id)
	require.NoError(t, err)

	err = repository.Delete(ctx, testOrgId, namespace, role.Id)
	require.NoError(t, err)

	_, err = repository.GetByIDs(ctx, testOrgId, namespace, role.Id)
	require.Error(t, err)

	err = store.ClearTable("Role", "", testOrgId, namespace)
	require.NoError(t, err)
}

func Test_ShouldSaveAndQueryRole(t *testing.T) {
	// GIVEN config, redis-service and role repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	testOrgId := uuid.NewV4().String()
	namespace := "role-query-namespace"
	repository, err := NewRoleRepository(store)
	require.NoError(t, err)
	for i := 0; i < 200; i++ {
		role := buildTestRole(i)
		err = repository.Create(ctx, testOrgId, namespace, role.Id, &role, time.Duration(0))
		require.NoError(t, err)
	}
	res, _, err := repository.Query(ctx, testOrgId, namespace, nil, "", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"id": "id_0"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"name": "/file/1"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"name": "/file/000"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 0, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{}, "100", 10)
	require.NoError(t, err)
	require.Equal(t, 10, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{}, "200", 10)
	require.NoError(t, err)
	require.Equal(t, 0, len(res))

	err = store.ClearTable("Role", "", testOrgId, namespace)
	require.NoError(t, err)
}

func buildTestRole(i int) types.Role {
	return types.Role{
		Id:            fmt.Sprintf("id_%d", i),
		Name:          fmt.Sprintf("/file/%d", i),
		PermissionIds: []string{"1", "2"},
		ParentIds:     []string{"1", "2"},
	}
}
