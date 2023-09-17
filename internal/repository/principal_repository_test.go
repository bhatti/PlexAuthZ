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

func Test_ShouldSaveAndGetPrincipal(t *testing.T) {
	// GIVEN config, redis-service and principal repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	repository, err := NewPrincipalRepository(store)
	require.NoError(t, err)
	principal := buildTestPrincipal(1)
	testOrgId := uuid.NewV4().String()
	namespace := "principal-save-namespace"
	err = repository.Create(ctx, testOrgId, namespace, principal.Id, &principal, time.Duration(0))
	require.NoError(t, err)

	saved, err := repository.GetByID(ctx, testOrgId, namespace, principal.Id)
	require.NoError(t, err)
	require.Equal(t, principal.Username, saved.Username)
	require.Equal(t, principal.OrganizationId, saved.OrganizationId)
	require.Equal(t, 2, len(saved.RoleIds))
	require.Equal(t, 2, len(saved.GroupIds))
	require.Equal(t, 2, len(saved.PermissionIds))
	require.Equal(t, 2, len(saved.RelationIds))

	saved.Username = "user2"
	saved.Name = "jane"
	err = repository.Update(ctx, testOrgId, namespace, principal.Id, saved.Version, saved, time.Duration(0))
	require.NoError(t, err)

	saved, err = repository.GetByID(ctx, testOrgId, namespace, principal.Id)
	require.NoError(t, err)
	require.Equal(t, "user2", saved.Username)
	require.Equal(t, "jane", saved.Name)

	err = store.ClearTable("Principal", "", testOrgId, namespace)
	require.NoError(t, err)
}

func Test_ShouldSaveAndDeletePrincipal(t *testing.T) {
	// GIVEN config, redis-service and principal repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	repository, err := NewPrincipalRepository(store)
	require.NoError(t, err)
	principal := buildTestPrincipal(1)

	testOrgId := uuid.NewV4().String()
	namespace := "principal-del-namespace"

	err = repository.Create(ctx, testOrgId, namespace, principal.Id, &principal, time.Duration(0))
	require.NoError(t, err)

	_, err = repository.GetByIDs(ctx, testOrgId, namespace, principal.Id)
	require.NoError(t, err)

	err = repository.Delete(ctx, testOrgId, namespace, principal.Id)
	require.NoError(t, err)

	_, err = repository.GetByIDs(ctx, testOrgId, namespace, principal.Id)
	require.Error(t, err)

	err = store.ClearTable("Principal", "", testOrgId, namespace)
	require.NoError(t, err)
}

func Test_ShouldSaveAndQueryPrincipal(t *testing.T) {
	// GIVEN config, redis-service and principal repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	testOrgId := uuid.NewV4().String()
	namespace := "principal-query-namespace"
	err = store.ClearTable("Principal", "", testOrgId, namespace)
	require.NoError(t, err)
	repository, err := NewPrincipalRepository(store)
	require.NoError(t, err)
	for i := 0; i < 200; i++ {
		principal := buildTestPrincipal(i)
		err = repository.Create(ctx, testOrgId, namespace, principal.Id, &principal, time.Duration(0))
		require.NoError(t, err)
	}
	res, _, err := repository.Query(ctx, testOrgId, namespace, nil, "", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"id": "id_0"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"username:>=": "name_"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{}, "100", 10)
	require.NoError(t, err)
	require.Equal(t, 10, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{}, "200", 10)
	require.NoError(t, err)
	require.Equal(t, 0, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"username:==": "name_0"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"username:!=": "1111"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"username:>=": "name_199"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 89, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"username:<=": "name_1"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 2, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"username:<": "name_1"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"username:>": "name_198"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 89, len(res))
}

func buildTestPrincipal(i int) types.Principal {
	return types.Principal{
		Id:             fmt.Sprintf("id_%d", i),
		Username:       fmt.Sprintf("name_%d", i),
		OrganizationId: fmt.Sprintf("org_%d", i),
		Name:           fmt.Sprintf("john doe_%d", i),
		RoleIds:        []string{"1", "2"},
		GroupIds:       []string{"1", "2"},
		PermissionIds:  []string{"1", "2"},
		RelationIds:    []string{"1", "2"},
	}
}
