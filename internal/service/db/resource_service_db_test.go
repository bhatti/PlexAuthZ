package db

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_Should_CRUD_Resources(t *testing.T) {
	// GIVEN auth-service
	ctx := context.TODO()
	store, org, err := newAuthServiceAndOrg()
	require.NoError(t, err)

	// AND principal
	principal, err := domain.NewPrincipalBuilder().
		WithNamespaces(org.Namespaces...).
		WithUsername("user1").
		WithName("john").
		WithOrganizationId(org.Id).Build()
	require.NoError(t, err)
	// WHEN Creating principal
	principal, err = store.CreatePrincipal(ctx, principal)
	// THEN it should not fail
	require.NoError(t, err)

	// AND resource
	resource, err := domain.NewResourceBuilder().
		WithNamespace(org.Namespaces[0]).
		WithName("res-1").
		WithAllowedActions("read").
		WithCapacity(2).Build()
	require.NoError(t, err)

	// WHEN creating a resource
	savedResource, err := store.CreateResource(ctx, org.Id, resource)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN updating a resource
	err = store.UpdateResource(ctx, org.Id, resource)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN querying an organization by id
	res, _, err := store.QueryResources(ctx, org.Id, org.Namespaces[0], map[string]string{"id": savedResource.Id}, "", 0)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, 1, len(res))

	// WHEN querying an organization by scope
	res, _, err = store.QueryResources(ctx, org.Id, org.Namespaces[0], map[string]string{"name": "res-1"}, "", 0)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, 1, len(res))

	// Iterating resources
	for _, next := range res {
		require.Equal(t, resource.Namespace, next.Namespace)
		require.Equal(t, resource.Name, next.Name)
		err := store.AllocateResourceInstance(ctx,
			org.Id,
			org.Namespaces[0],
			next.Id,
			principal.Id,
			"eq true true",
			time.Hour, nil)
		require.NoError(t, err)
		instances, _, err := store.QueryResourceInstances(ctx, org.Id, org.Namespaces[0], next.Id, nil, "", 10)
		require.NoError(t, err)
		require.True(t, len(instances) > 0)
		err = store.DeallocateResourceInstance(ctx, org.Id, org.Namespaces[0], next.Id, principal.Id)
		require.NoError(t, err)
		// WHEN deleting a resource
		err = store.DeleteResource(ctx, org.Id, org.Namespaces[0], next.Id)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func Test_ShouldCreateAndGetAndDeleteResourceInstance(t *testing.T) {
	// GIVEN auth-service
	ctx := context.TODO()
	store, org, err := newAuthServiceAndOrg()
	require.NoError(t, err)
	principal, err := store.CreatePrincipal(ctx, &types.Principal{
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
	_, err = store.CreateResource(ctx, org.Id, resource)
	require.NoError(t, err)
	res, _, err := store.QueryResources(ctx, org.Id, org.Namespaces[0], map[string]string{"name": "/file/1"}, "", 0)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	for _, next := range res {
		for i := 0; i < 20; i++ {
			err = store.AllocateResourceInstance(
				ctx,
				org.Id,
				org.Namespaces[0],
				next.Id,
				principal.Id, "eq 1 1", time.Hour, nil)
			require.NoError(t, err)
		}
		err = store.AllocateResourceInstance(
			ctx,
			org.Id,
			org.Namespaces[0],
			next.Id,
			principal.Id, "", time.Hour, nil)
		require.NoError(t, err)
		capacity, allocated, err := store.CountResourceInstances(ctx, org.Id, org.Namespaces[0], next.Id)
		require.Equal(t, int32(20), capacity)
		require.Equal(t, int32(1), allocated) // allocated to same principal
		for i := 0; i < 20; i++ {
			err = store.DeallocateResourceInstance(
				ctx,
				org.Id,
				org.Namespaces[0],
				next.Id,
				principal.Id)
			require.NoError(t, err)
		}
		capacity, allocated, err = store.CountResourceInstances(ctx, org.Id, org.Namespaces[0], next.Id)
		require.Equal(t, int32(20), capacity)
		require.Equal(t, int32(0), allocated)
		err = store.DeleteResource(ctx, org.Id, org.Namespaces[0], next.Id)
		require.NoError(t, err)
	}
}
