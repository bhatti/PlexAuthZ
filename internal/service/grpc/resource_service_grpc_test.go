package grpc

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/service"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func testCRUDResources(
	ctx context.Context,
	t *testing.T,
	authService service.AuthAdminService,
	org *types.Organization,
) {
	principal, err := domain.NewPrincipalBuilder().
		WithNamespaces(org.Namespaces...).
		WithUsername("user1").
		WithOrganizationId(org.Id).Build()
	require.NoError(t, err)
	// WHEN Creating principal
	principal, err = authService.CreatePrincipal(ctx, principal)
	// THEN it should not fail
	require.NoError(t, err)

	// AND resource
	resource, err := domain.NewResourceBuilder().
		WithNamespace(org.Namespaces[0]).
		WithName("res-1").
		WithAllowedActions("read").
		WithAttribute("k1", "v1").
		WithAttribute("k2", "v2").
		WithCapacity(2).Build()
	require.NoError(t, err)

	// WHEN creating a resource
	savedResource, err := authService.CreateResource(ctx, org.Id, resource)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN creating a resource
	err = authService.UpdateResource(ctx, org.Id, resource)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN getting a resource
	_, err = authService.GetResource(ctx, org.Id, resource.Namespace, resource.Id)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN querying an organization by id
	res, _, err := authService.QueryResources(ctx, org.Id, org.Namespaces[0], map[string]string{"id": savedResource.Id}, "", 0)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, 1, len(res))

	// WHEN querying an organization by scope
	res, _, err = authService.QueryResources(ctx, org.Id, org.Namespaces[0], map[string]string{"name": "res-1"}, "", 0)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, 1, len(res))

	// Iterating resources
	for _, next := range res {
		require.Equal(t, resource.Namespace, next.Namespace)
		require.Equal(t, resource.Name, next.Name)
		require.Equal(t, 2, len(resource.Attributes))
		err := authService.AllocateResourceInstance(ctx, org.Id, org.Namespaces[0], next.Id, principal.Id, "", time.Hour, nil)
		require.NoError(t, err)
		capacity, allocated, err := authService.CountResourceInstances(ctx, org.Id, org.Namespaces[0], next.Id)
		require.Equal(t, int32(2), capacity)
		require.Equal(t, int32(1), allocated)
		instances, _, err := authService.QueryResourceInstances(ctx, org.Id, org.Namespaces[0], next.Id, nil, "", 0)
		require.Equal(t, 1, len(instances))
		err = authService.DeallocateResourceInstance(ctx, org.Id, org.Namespaces[0], next.Id, principal.Id)
		require.NoError(t, err)
		capacity, allocated, err = authService.CountResourceInstances(ctx, org.Id, org.Namespaces[0], next.Id)
		require.Equal(t, int32(2), capacity)
		require.Equal(t, int32(0), allocated)
		instances, _, err = authService.QueryResourceInstances(ctx, org.Id, org.Namespaces[0], next.Id, nil, "", 0)
		require.Equal(t, 0, len(instances))
		// WHEN deleting a resource
		err = authService.DeleteResource(ctx, org.Id, org.Namespaces[0], next.Id)
		// THEN it should not fail
		require.NoError(t, err)
		// WHEN deleting a resource without org-id, namespace, id - then should fail
		require.Error(t, authService.DeleteResource(ctx, "", org.Namespaces[0], next.Id))
		require.Error(t, authService.DeleteResource(ctx, org.Id, "", next.Id))
		require.Error(t, authService.DeleteResource(ctx, org.Id, org.Namespaces[0], ""))
	}
}

func testCRUDResourcesWithInstances(
	ctx context.Context,
	t *testing.T,
	authService service.AuthAdminService,
	org *types.Organization,
) {
	principal, err := authService.CreatePrincipal(ctx, &types.Principal{
		Username:       "user1",
		OrganizationId: org.Id,
		Namespaces:     org.Namespaces,
	})
	require.NoError(t, err)
	resource, err := domain.NewResourceBuilder().
		WithName(fmt.Sprintf("/file/%d", 1)).
		WithNamespace(org.Namespaces[0]).
		WithCapacity(20).
		WithAllowedActions("read", "write").Build()
	_, err = authService.CreateResource(ctx, org.Id, resource)
	require.NoError(t, err)
	res, _, err := authService.QueryResources(ctx, org.Id, org.Namespaces[0], map[string]string{"name": "/file/1"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	for _, next := range res {
		for i := 0; i < 20; i++ {
			err = authService.AllocateResourceInstance(
				ctx,
				org.Id,
				org.Namespaces[0],
				next.Id,
				principal.Id, "eq true true", time.Hour, nil)
			require.NoError(t, err)
		}
		err = authService.AllocateResourceInstance(
			ctx,
			org.Id,
			org.Namespaces[0],
			next.Id,
			principal.Id,
			"", time.Hour, nil)
		require.NoError(t, err)
		capacity, allocated, err := authService.CountResourceInstances(ctx, org.Id, org.Namespaces[0], next.Id)
		require.Equal(t, int32(20), capacity)
		require.Equal(t, int32(1), allocated)
		for i := 0; i < 20; i++ {
			err = authService.DeallocateResourceInstance(ctx, org.Id, org.Namespaces[0], next.Id, principal.Id)
			require.NoError(t, err)
		}
		capacity, allocated, err = authService.CountResourceInstances(ctx, org.Id, org.Namespaces[0], next.Id)
		require.Equal(t, int32(20), capacity)
		require.Equal(t, int32(0), allocated)
		err = authService.DeleteResource(ctx, org.Id, org.Namespaces[0], next.Id)
		require.NoError(t, err)
	}
}
