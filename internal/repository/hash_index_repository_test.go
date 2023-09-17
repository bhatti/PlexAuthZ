package repository

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/repository/redis"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
	"testing"
	"time"
)

func Test_ShouldSaveAndGetHashIndex(t *testing.T) {
	// GIVEN config, redis-service and index repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	repository, err := NewHashIndexRepository(store)
	require.NoError(t, err)
	index := buildTestHashIndex(1)
	testOrgId := uuid.NewV4().String()
	namespace := "index-save-namespace"
	err = repository.Update(ctx, testOrgId, namespace, index.Hash, -1, &index, time.Duration(0))
	require.NoError(t, err)

	saved, err := repository.GetByIDs(ctx, testOrgId, namespace, index.Hash)
	require.NoError(t, err)
	require.Equal(t, 1, len(saved))
	require.Equal(t, index.Hash, saved[index.Hash].Hash)
	require.Equal(t, 2, len(saved[index.Hash].Ids))

	require.NoError(t, err)
}

func Test_ShouldSaveAndDeleteHashIndex(t *testing.T) {
	// GIVEN config, redis-service and index repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	repository, err := NewHashIndexRepository(store)
	require.NoError(t, err)
	index := buildTestHashIndex(1)

	testOrgId := uuid.NewV4().String()
	namespace := "index-del-namespace"

	err = repository.Update(ctx, testOrgId, namespace, index.Hash, -1, &index, time.Duration(0))
	require.NoError(t, err)

	_, err = repository.GetByIDs(ctx, testOrgId, namespace, index.Hash)
	require.NoError(t, err)

	err = repository.Delete(ctx, testOrgId, namespace, index.Hash)
	require.NoError(t, err)

	_, err = repository.GetByIDs(ctx, testOrgId, namespace, index.Hash)
	require.Error(t, err)
}

func Test_ShouldSaveAndQueryHashIndex(t *testing.T) {
	// GIVEN config, redis-service and index repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	testOrgId := uuid.NewV4().String()
	namespace := "index-query-namespace"
	repository, err := NewHashIndexRepository(store)
	require.NoError(t, err)
	for i := 0; i < 200; i++ {
		index := buildTestHashIndex(i)
		err = repository.Update(ctx, testOrgId, namespace, index.Hash, -1, &index, time.Duration(0))
		require.NoError(t, err)
	}
	res, _, err := repository.Query(ctx, testOrgId, namespace, nil, "", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
}

func buildTestHashIndex(i int) domain.HashIndex {
	return domain.HashIndex{
		Hash: fmt.Sprintf("hash_%d", i),
		Ids:  []string{"1", "2"},
	}
}
