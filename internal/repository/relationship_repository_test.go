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

func Test_ShouldSaveAndGetRelationship(t *testing.T) {
	// GIVEN config, redis-service and relation repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	repository, err := NewRelationshipRepository(store)
	require.NoError(t, err)
	relation := buildTestRelationship(1)
	testOrgId := uuid.NewV4().String()
	namespace := "relation-save-namespace"
	err = repository.Create(ctx, testOrgId, namespace, relation.Id, &relation, time.Duration(0))
	require.NoError(t, err)

	saved, err := repository.GetByID(ctx, testOrgId, namespace, relation.Id)
	require.NoError(t, err)
	require.Equal(t, relation.Relation, saved.Relation)
	require.Equal(t, relation.PrincipalId, saved.PrincipalId)
	require.Equal(t, relation.ResourceId, saved.ResourceId)

	saved.Relation = "new-rel"
	err = repository.Update(ctx, testOrgId, namespace, relation.Id, saved.Version, saved, time.Duration(0))
	require.NoError(t, err)

	saved, err = repository.GetByID(ctx, testOrgId, namespace, relation.Id)
	require.NoError(t, err)
	require.Equal(t, "new-rel", saved.Relation)

	err = store.ClearTable("Relationship", "", testOrgId, namespace)
	require.NoError(t, err)
}

func Test_ShouldSaveAndDeleteRelationship(t *testing.T) {
	// GIVEN config, redis-service and relation repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	repository, err := NewRelationshipRepository(store)
	require.NoError(t, err)
	relation := buildTestRelationship(1)

	testOrgId := uuid.NewV4().String()
	namespace := "relation-del-namespace"

	err = repository.Create(ctx, testOrgId, namespace, relation.Id, &relation, time.Duration(0))
	require.NoError(t, err)

	_, err = repository.GetByIDs(ctx, testOrgId, namespace, relation.Id)
	require.NoError(t, err)

	err = repository.Delete(ctx, testOrgId, namespace, relation.Id)
	require.NoError(t, err)

	_, err = repository.GetByIDs(ctx, testOrgId, namespace, relation.Id)
	require.Error(t, err)

	err = store.ClearTable("Relationship", "", testOrgId, namespace)
	require.NoError(t, err)
}

func Test_ShouldSaveAndQueryRelationship(t *testing.T) {
	// GIVEN config, redis-service and relation repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	testOrgId := uuid.NewV4().String()
	namespace := "relation-query-namespace"
	repository, err := NewRelationshipRepository(store)
	require.NoError(t, err)
	for i := 0; i < 200; i++ {
		relation := buildTestRelationship(i)
		err = repository.Create(ctx, testOrgId, namespace, relation.Id, &relation, time.Duration(0))
		require.NoError(t, err)
	}
	res, _, err := repository.Query(ctx, testOrgId, namespace, nil, "", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"id": "id_0"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"relation:>=": "rel_"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{}, "100", 10)
	require.NoError(t, err)
	require.Equal(t, 10, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{}, "200", 10)
	require.NoError(t, err)
	require.Equal(t, 0, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"relation:==": "rel_0"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"relation:!=": "1111"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"relation:>=": "rel_199"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 89, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"relation:<=": "rel_1"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 2, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"relation:<": "rel_1"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = repository.Query(ctx, testOrgId, namespace, map[string]string{"relation:>": "rel_198"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 89, len(res))
	err = store.ClearTable("Relationship", "", testOrgId, namespace)
	require.NoError(t, err)
}

func buildTestRelationship(i int) types.Relationship {
	return types.Relationship{
		Id:          fmt.Sprintf("id_%d", i),
		Relation:    fmt.Sprintf("rel_%d", i),
		PrincipalId: fmt.Sprintf("user_%d", i),
		ResourceId:  fmt.Sprintf("res_%d", i),
	}
}
