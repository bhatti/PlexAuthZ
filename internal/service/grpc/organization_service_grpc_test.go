package grpc

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/service"
	"github.com/stretchr/testify/require"
	"testing"
)

func testCRUDOrganizations(
	ctx context.Context,
	t *testing.T,
	authService service.AuthAdminService,
	org *types.Organization,
) {
	// WHEN finding an organization
	saved, err := authService.GetOrganization(ctx, org.Id)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, int64(1), saved.Version)
	require.Equal(t, org.Name, saved.Name)
	require.Equal(t, 3, len(saved.Namespaces))
	require.Equal(t, org.Namespaces, saved.Namespaces)
	require.Equal(t, org.Url, saved.Url)

	// WHEN updating an organization
	saved.Name = "new-name"
	err = authService.UpdateOrganization(ctx, saved)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN finding an organization
	saved, err = authService.GetOrganization(ctx, saved.Id)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, int64(2), saved.Version)
	require.Equal(t, "new-name", saved.Name)
	require.Equal(t, org.Namespaces, saved.Namespaces)
	require.Equal(t, org.Url, saved.Url)
	require.Equal(t, 0, len(saved.ParentIds))

	// WHEN querying organizations
	res, _, err := authService.GetOrganizations(ctx, nil, "", 0)
	require.NoError(t, err)
	// THEN it should not fail
	require.True(t, len(res) > 0)

	// WHEN querying organizations by id
	res, _, err = authService.GetOrganizations(ctx, map[string]string{"id": saved.Id}, "", 0)
	require.NoError(t, err)
	// THEN it should not fail
	require.Equal(t, 1, len(res))

	// WHEN Deleting an organization
	err = authService.DeleteOrganization(ctx, org.Id)
	// THEN it should not fail
	require.NoError(t, err)
}
