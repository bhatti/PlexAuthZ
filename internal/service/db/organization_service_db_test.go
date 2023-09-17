package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Should_CRUD_Organizations(t *testing.T) {
	// GIVEN auth-service
	ctx := context.TODO()

	// Creating an organization
	store, org, err := newAuthServiceAndOrg()
	require.NoError(t, err)

	// WHEN finding an organization
	saved, err := store.GetOrganization(ctx, org.Id)
	// THEN it should not fail
	require.NoError(t, err)

	require.Equal(t, int64(1), saved.Version)
	require.Equal(t, org.Name, saved.Name)
	require.Equal(t, 3, len(saved.Namespaces))
	require.Equal(t, org.Namespaces, saved.Namespaces)
	require.Equal(t, org.Url, saved.Url)

	// WHEN updating an organization
	saved.Name = "new-name"
	err = store.UpdateOrganization(ctx, saved)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN finding an organization
	saved, err = store.GetOrganization(ctx, saved.Id)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, int64(2), saved.Version)
	require.Equal(t, org.Name, saved.Name)
	require.Equal(t, org.Namespaces, saved.Namespaces)
	require.Equal(t, org.Url, saved.Url)
	require.Equal(t, 0, len(saved.ParentIds))

	// WHEN querying organizations
	res, _, err := store.GetOrganizations(ctx, nil, "", 0)
	require.NoError(t, err)
	// THEN it should not fail
	require.True(t, len(res) > 0)

	// WHEN querying organizations by id
	res, _, err = store.GetOrganizations(ctx, map[string]string{"id": saved.Id}, "", 0)
	require.NoError(t, err)
	// THEN it should not fail
	require.Equal(t, 1, len(res))

	// WHEN Deleting an organization
	err = store.DeleteOrganization(ctx, org.Id)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN Deleting an organization without id
	err = store.DeleteOrganization(ctx, "")
	// THEN it should fail
	require.Error(t, err)
}
