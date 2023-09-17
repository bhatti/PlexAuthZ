package ddb

import (
	"fmt"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ShouldSaveAndGetData(t *testing.T) {
	// GIVEN config and ddb-service
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := NewDDBStore(cfg)
	require.NoError(t, err)
	tenant := "test-save-get"
	namespace := "ns1"
	err = store.CreateTable("table1", "")
	require.NoError(t, err)
	err = store.Create("table1", "", tenant, namespace, "id1", []byte("data1"), 0)
	require.NoError(t, err)

	saved, err := store.Get("table1", "", tenant, namespace, "id1")
	require.NoError(t, err)
	require.Equal(t, "data1", string(saved["id1"]))

	err = store.Update("table1", "", tenant, namespace, "id1", 1, []byte("data2"), 0)
	require.NoError(t, err)

	saved, err = store.Get("table1", "", tenant, namespace, "id1")
	require.NoError(t, err)
	require.Equal(t, "data2", string(saved["id1"]))

	count, err := store.Size("table1", "", tenant, namespace)
	require.NoError(t, err)
	require.Equal(t, int64(1), count)

	err = store.ClearTable("table1", "", tenant, namespace)
	require.NoError(t, err)
}

func Test_ShouldSaveAndDeleteData(t *testing.T) {
	// GIVEN config and ddb-service
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := NewDDBStore(cfg)
	require.NoError(t, err)

	tenant := "test-save-del"
	namespace := "ns1"
	err = store.CreateTable("table1", "")
	require.NoError(t, err)
	data := []byte("data1")
	err = store.Create("table1", "", tenant, namespace, "id1", data, 0)
	require.NoError(t, err)

	_, err = store.Get("table1", "", tenant, namespace, "id1")
	require.NoError(t, err)

	err = store.Delete("table1", "", tenant, namespace, "id1")
	require.NoError(t, err)

	_, err = store.Get("table1", "", tenant, namespace, "id1")
	require.Error(t, err)

	err = store.ClearTable("table1", "", tenant, namespace)
	require.NoError(t, err)
}

func Test_ShouldSaveAndQueryData(t *testing.T) {
	// GIVEN config and ddb-service
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	store, err := NewDDBStore(cfg)
	require.NoError(t, err)
	tenant := "test-query-data"
	namespace := "ns1"
	err = store.CreateTable("table1", "")
	require.NoError(t, err)
	for i := 0; i < 200; i++ {
		id := fmt.Sprintf("id_%d", i)
		data := []byte(fmt.Sprintf("data_%d", i))
		err = store.Update("table1", "", tenant, namespace, id, -1, data, 0)
		require.NoError(t, err)
	}
	res, _, err := store.Query("table1", "", tenant, namespace, nil, "", 0)
	require.NoError(t, err)
	require.Equal(t, 200, len(res))
	res, _, err = store.Query("table1", "", tenant, namespace, map[string]string{}, "", 0)
	require.NoError(t, err)
	require.Equal(t, "data_0", string(res["id_0"]))

	err = store.ClearTable("table1", "", tenant, namespace)
	require.NoError(t, err)

	res, _, err = store.Query("table1", "", tenant, namespace, nil, "", 0)
	require.NoError(t, err)
	require.Equal(t, 0, len(res))
}
