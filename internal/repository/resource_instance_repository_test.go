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

func Test_ShouldSaveAndGetResourceInstance(t *testing.T) {
	// GIVEN config, redis-service and instance instanceRepository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	resourceRepository, err := NewResourceInstanceRepository(store, time.Second*1)
	require.NoError(t, err)
	instanceRepository, err := resourceRepository.CreateResourceInstanceRepository("r1")
	require.NoError(t, err)
	instance := buildTestResourceInstance(1)
	testOrgId := uuid.NewV4().String()
	namespace := "instance-save-namespace"
	err = instanceRepository.Create(ctx, testOrgId, namespace, instance.Id, &instance, time.Duration(0))
	require.NoError(t, err)

	saved, err := instanceRepository.GetByID(ctx, testOrgId, namespace, instance.Id)
	require.NoError(t, err)
	require.Equal(t, instance.PrincipalId, saved.PrincipalId)
	require.Equal(t, instance.ResourceId, saved.ResourceId)
	require.Equal(t, instance.State, saved.State)

	saved.Version = 2
	err = instanceRepository.Update(ctx, testOrgId, namespace, instance.Id, saved.Version, saved, time.Duration(0))
	require.NoError(t, err)

	saved, err = instanceRepository.GetByID(ctx, testOrgId, namespace, instance.Id)
	require.NoError(t, err)
	require.Equal(t, int64(2), saved.Version)

	time.Sleep(time.Second * 2)
	_, err = instanceRepository.GetByIDs(ctx, testOrgId, namespace, instance.Id)
	require.Error(t, err)
}

func Test_ShouldSaveAndDeleteResourceInstance(t *testing.T) {
	// GIVEN config, redis-service and instance instanceRepository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	resourceRepository, err := NewResourceInstanceRepository(store, time.Second*1)
	require.NoError(t, err)
	instanceRepository, err := resourceRepository.CreateResourceInstanceRepository("r2")
	require.NoError(t, err)
	instance := buildTestResourceInstance(1)

	testOrgId := uuid.NewV4().String()
	namespace := "instance-del-namespace"

	err = instanceRepository.Create(ctx, testOrgId, namespace, instance.Id, &instance, time.Duration(0))
	require.NoError(t, err)

	_, err = instanceRepository.GetByIDs(ctx, testOrgId, namespace, instance.Id)
	require.NoError(t, err)

	err = instanceRepository.Delete(ctx, testOrgId, namespace, instance.Id)
	require.NoError(t, err)

	_, err = instanceRepository.GetByIDs(ctx, testOrgId, namespace, instance.Id)
	require.Error(t, err)
	err = store.ClearTable("ResourceInstance", "r2", testOrgId, namespace)
	require.NoError(t, err)
}

func Test_ShouldSaveAndQueryResourceInstance(t *testing.T) {
	// GIVEN config, redis-service and instance repository
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := redis.NewRedisStore(cfg)
	require.NoError(t, err)
	testOrgId := uuid.NewV4().String()
	namespace := "instance-query-namespace"
	resourceRepository, err := NewResourceInstanceRepository(store, time.Second*1)
	require.NoError(t, err)
	instanceRepository, err := resourceRepository.CreateResourceInstanceRepository("r3")
	require.NoError(t, err)
	for i := 0; i < 200; i++ {
		instance := buildTestResourceInstance(i)
		err = instanceRepository.Create(ctx, testOrgId, namespace, instance.Id, &instance, time.Duration(0))
		require.NoError(t, err)
	}
	res, _, err := instanceRepository.Query(ctx, testOrgId, namespace, nil, "", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = instanceRepository.Query(ctx, testOrgId, namespace, map[string]string{"id": "id_0"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = instanceRepository.Query(ctx, testOrgId, namespace, map[string]string{"principal_id:>=": "user_"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = instanceRepository.Query(ctx, testOrgId, namespace, map[string]string{}, "100", 10)
	require.NoError(t, err)
	require.Equal(t, 10, len(res))
	res, _, err = instanceRepository.Query(ctx, testOrgId, namespace, map[string]string{}, "200", 10)
	require.NoError(t, err)
	require.Equal(t, 0, len(res))
	res, _, err = instanceRepository.Query(ctx, testOrgId, namespace, map[string]string{"principal_id:==": "user_0"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = instanceRepository.Query(ctx, testOrgId, namespace, map[string]string{"principal_id:!=": "1111"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = instanceRepository.Query(ctx, testOrgId, namespace, map[string]string{"principal_id:>=": "user_199"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 89, len(res))
	res, _, err = instanceRepository.Query(ctx, testOrgId, namespace, map[string]string{"principal_id:<=": "user_1"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 2, len(res))
	res, _, err = instanceRepository.Query(ctx, testOrgId, namespace, map[string]string{"principal_id:<": "user_1"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	res, _, err = instanceRepository.Query(ctx, testOrgId, namespace, map[string]string{"principal_id:>": "user_198"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 89, len(res))
	err = store.ClearTable("ResourceInstance", "r3", testOrgId, namespace)
	require.NoError(t, err)
}

func buildTestResourceInstance(i int) types.ResourceInstance {
	return types.ResourceInstance{
		Id:          fmt.Sprintf("id_%d", i),
		PrincipalId: fmt.Sprintf("user_%d", i),
		State:       types.ResourceState_ALLOCATED,
	}
}
