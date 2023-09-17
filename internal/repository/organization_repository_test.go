package repository

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/repository/redis"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_ShouldSaveAndGetOrganization(t *testing.T) {
	// GIVEN config, redis-service and org repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	repository, err := NewOrganizationRepository(store)
	require.NoError(t, err)
	org := buildTestOrg(1)
	err = repository.Create(ctx, org.Id, "", org.Id, &org, time.Duration(0))
	require.NoError(t, err)

	saved, err := repository.GetByID(ctx, org.Id, "", org.Id)
	require.NoError(t, err)
	require.Equal(t, org.Name, saved.Name)
	require.Equal(t, org.Url, saved.Url)
	require.Equal(t, 2, len(saved.Namespaces))

	saved.Name = "new-name"
	err = repository.Update(ctx, org.Id, "", org.Id, saved.Version, saved, time.Duration(0))
	require.NoError(t, err)

	saved, err = repository.GetByID(ctx, org.Id, "", org.Id)
	require.NoError(t, err)
	require.Equal(t, "new-name", saved.Name)
	require.Equal(t, org.Url, saved.Url)
	require.Equal(t, 2, len(saved.Namespaces))
}

func Test_ShouldSaveAndDeleteOrganization(t *testing.T) {
	// GIVEN config, redis-service and org repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	repository, err := NewOrganizationRepository(store)
	require.NoError(t, err)
	org := buildTestOrg(1)

	namespace := "org-del-namespace"

	err = repository.Create(ctx, org.Id, namespace, org.Id, &org, time.Duration(0))
	require.NoError(t, err)

	_, err = repository.GetByIDs(ctx, org.Id, namespace, org.Id)
	require.NoError(t, err)

	err = repository.Delete(ctx, org.Id, namespace, org.Id)
	require.NoError(t, err)

	_, err = repository.GetByIDs(ctx, org.Id, namespace, org.Id)
	require.Error(t, err)
}

func Test_ShouldSaveAndQueryOrganization(t *testing.T) {
	// GIVEN config, redis-service and org repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	repository, err := NewOrganizationRepository(store)
	require.NoError(t, err)
	for i := 0; i < 200; i++ {
		org := buildTestOrg(i)
		err = repository.Create(ctx, org.Id, "", org.Id, &org, time.Duration(0))
		require.NoError(t, err)
	}
	res, _, err := repository.Query(ctx, "", "", nil, "", 200)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = repository.Query(ctx, "", "", map[string]string{"id": "id_0"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = repository.Query(ctx, "", "", map[string]string{"url:>=": "url_"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = repository.Query(ctx, "", "", map[string]string{}, "100", 10)
	require.NoError(t, err)
	require.Equal(t, 10, len(res))
	res, _, err = repository.Query(ctx, "", "", map[string]string{}, "2000000", 10)
	require.NoError(t, err)
	require.Equal(t, 0, len(res))
	res, _, err = repository.Query(ctx, "", "", map[string]string{"url:==": "url_0"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = repository.Query(ctx, "", "", map[string]string{"url:!=": "1111"}, "", 200)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = repository.Query(ctx, "", "", map[string]string{"url:>=": "url_199"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 89, len(res))
	res, _, err = repository.Query(ctx, "", "", map[string]string{"url:<=": "url_1"}, "", 2)
	require.NoError(t, err)
	require.Equal(t, 2, len(res))
	res, _, err = repository.Query(ctx, "", "", map[string]string{"url:<": "url_1"}, "", 1)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = repository.Query(ctx, "", "", map[string]string{"url:>": "url_198"}, "", 50)
	require.NoError(t, err)
	require.Equal(t, 50, len(res))
}

func buildTestOrg(i int) types.Organization {
	return types.Organization{
		Id:         fmt.Sprintf("id_%d", i),
		Url:        fmt.Sprintf("url_%d", i),
		Namespaces: []string{"1", "2"},
	}
}
