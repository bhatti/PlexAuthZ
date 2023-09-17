package redis

import (
	"fmt"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ShouldSaveAndGetData(t *testing.T) {
	// GIVEN config and redis-service
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := NewRedisStore(cfg)
	require.NoError(t, err)
	baseTable := "table1"
	namespace := "test-save-get"
	tenant := "123"
	err = store.Create(baseTable, "", tenant, namespace, "id1", []byte("data1"), 0)
	require.NoError(t, err)

	saved, err := store.Get(baseTable, "", tenant, namespace, "id1")
	require.NoError(t, err)
	require.Equal(t, "data1", string(saved["id1"]))

	err = store.Update(baseTable, "", tenant, namespace, "id1", 1, []byte("data2"), 0)
	require.NoError(t, err)

	saved, err = store.Get(baseTable, "", tenant, namespace, "id1")
	require.NoError(t, err)
	require.Equal(t, "data2", string(saved["id1"]))

	count, err := store.Size(baseTable, "", tenant, namespace)
	require.NoError(t, err)
	require.Equal(t, int64(1), count)

	err = store.ClearTable(baseTable, "", tenant, namespace)
	require.NoError(t, err)
}

func Test_ShouldSaveAndDeleteData(t *testing.T) {
	// GIVEN config and redis-service
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := NewRedisStore(cfg)
	require.NoError(t, err)

	data := []byte("data1")
	namespace := "test-save-del"
	baseTable := "table1"
	tenant := "123"
	err = store.Create(baseTable, "", tenant, namespace, "id1", data, 0)
	require.NoError(t, err)

	_, err = store.Get(baseTable, "", tenant, namespace, "id1")
	require.NoError(t, err)

	err = store.Delete(baseTable, "", tenant, namespace, "id1")
	require.NoError(t, err)

	_, err = store.Get(baseTable, "", tenant, namespace, "id1")
	require.Error(t, err)

	err = store.ClearTable(baseTable, "", tenant, namespace)
	require.NoError(t, err)
}

func Test_ShouldSaveAndQueryData(t *testing.T) {
	// GIVEN config and redis-service
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := NewRedisStore(cfg)
	require.NoError(t, err)
	namespace := "test-query-data"
	baseTable := "table1"
	tenant := "123"
	for i := 0; i < 200; i++ {
		id := fmt.Sprintf("id_%d", i)
		data := []byte(fmt.Sprintf("data_%d", i))
		err = store.Create(baseTable, "", tenant, namespace, id, data, 0)
		require.NoError(t, err)
	}
	res, _, err := store.Query(baseTable, "", tenant, namespace, nil, "", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = store.Query(baseTable, "", tenant, namespace, map[string]string{}, "", 0)
	require.NoError(t, err)
	require.Equal(t, "data_0", string(res["id_0"]))

	err = store.ClearTable(baseTable, "", tenant, namespace)
	require.NoError(t, err)

	res, _, err = store.Query(baseTable, "", tenant, namespace, nil, "", 0)
	require.NoError(t, err)
	require.Equal(t, 0, len(res))
}
