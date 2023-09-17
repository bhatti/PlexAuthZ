package server

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_ShouldCreateAndGetAndDeleteResource(t *testing.T) {
	// GIVEN auth-client
	err := os.Setenv("CONFIG_DIR", "../../config")
	require.NoError(t, err)
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	cfg.GrpcSasl = true

	clientTypes := []domain.ClientType{domain.DefaultClientType, domain.NobodyClientType, domain.RootClientType}

	for scenario, fn := range map[string]func(
		t *testing.T,
		clients Clients,
	){
		"Should Create/Update/Get/Delete Resources": testShouldCRUDResources,
	} {
		t.Run(scenario, func(t *testing.T) {
			for _, clientType := range clientTypes {
				clients, teardown := SetupGrpcServerForTesting(t, cfg, clientType, nil)
				clients.ClientType = clientType
				fn(t, clients)
				teardown()
			}
		})
	}
}

func testShouldCRUDResources(t *testing.T, clients Clients) {
	ctx := context.Background()
	orgRes, err := clients.OrganizationsClient.Create(ctx, &services.CreateOrganizationRequest{
		Name:       "org-name",
		Namespaces: []string{"admin", "finance", "engineering"},
	})
	if clients.ClientType == domain.RootClientType {
		require.NoError(t, err)
	} else {
		require.Error(t, err)
		return
	}
	principalRes, err := clients.PrincipalsClient.Create(ctx, &services.CreatePrincipalRequest{
		Username:       "principal1",
		Name:           "john",
		OrganizationId: orgRes.Id,
		Namespaces:     []string{"admin", "finance", "engineering"},
	})
	require.NoError(t, err)

	resourceRes, err := clients.ResourcesClient.Create(ctx, &services.CreateResourceRequest{
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
		Name:           "res-1",
		Capacity:       10,
		AllowedActions: []string{"read", "write"},
	})
	require.NoError(t, err)

	_, err = clients.AuthClient.Allocate(ctx, &services.AllocateResourceRequest{
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
		ResourceId:     resourceRes.Id,
		PrincipalId:    principalRes.Id,
	})
	require.NoError(t, err)

	countRes, err := clients.ResourcesClient.CountResourceInstances(ctx, &services.CountResourceInstancesRequest{
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
		ResourceId:     resourceRes.Id,
	})
	require.NoError(t, err)
	require.True(t, countRes.Allocated > 0)

	_, err = clients.ResourcesClient.QueryResourceInstances(ctx, &services.QueryResourceInstanceRequest{
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
		ResourceId:     resourceRes.Id,
	})
	require.NoError(t, err)

	_, err = clients.AuthClient.Deallocate(ctx, &services.DeallocateResourceRequest{
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
		ResourceId:     resourceRes.Id,
		PrincipalId:    principalRes.Id,
	})
	require.NoError(t, err)

	queryRes, err := clients.ResourcesClient.Query(ctx, &services.QueryResourceRequest{
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
		Predicates:     map[string]string{"id": resourceRes.Id},
	})
	require.NoError(t, err)
	resource, err := queryRes.Recv()
	require.NoError(t, err)
	require.Equal(t, int64(1), resource.Version)
	require.Equal(t, "res-1", resource.Name)
	require.Equal(t, 2, len(resource.AllowedActions))

	_, err = clients.ResourcesClient.Update(ctx, &services.UpdateResourceRequest{
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
		Id:             resource.Id,
		Name:           "res-2",
		AllowedActions: []string{"write"},
	})
	require.NoError(t, err)

	query, err := clients.ResourcesClient.Query(context.Background(), &services.QueryResourceRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
	})
	require.NoError(t, err)
	resource, err = query.Recv()
	require.NoError(t, err)
	require.Equal(t, int64(2), resource.Version)
	require.Equal(t, "res-2", resource.Name)
	require.Equal(t, 1, len(resource.AllowedActions))

	_, err = clients.ResourcesClient.Delete(context.Background(), &services.DeleteResourceRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
		Id:             resource.Id,
	})
	require.NoError(t, err)

	// should not find it after deleting it
	query, err = clients.ResourcesClient.Query(context.Background(), &services.QueryResourceRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
	})
	require.NoError(t, err)
	_, err = query.Recv()
	require.Error(t, err)
}
