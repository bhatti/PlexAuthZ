package server

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_ShouldCreateAndGetAndDeleteOrganization(t *testing.T) {
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
		"Should Create/Update/Get/Delete Organizations": testShouldCRUDOrganizations,
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

func testShouldCRUDOrganizations(t *testing.T, clients Clients) {
	ctx := context.Background()
	orgRes, err := clients.OrganizationsClient.Create(ctx, &services.CreateOrganizationRequest{
		Name:       "old-name",
		Namespaces: []string{"admin", "finance", "engineering"},
	})
	if clients.ClientType == domain.RootClientType {
		require.NoError(t, err)
	} else {
		require.Error(t, err)
		return
	}
	org, err := clients.OrganizationsClient.Get(ctx, &services.GetOrganizationRequest{
		Id: orgRes.Id,
	})
	require.NoError(t, err)
	require.Equal(t, int64(1), org.Version)
	require.Equal(t, "old-name", org.Name, org.Name)
	require.Equal(t, 3, len(org.Namespaces))

	_, err = clients.OrganizationsClient.Update(ctx, &services.UpdateOrganizationRequest{
		Id:         orgRes.Id,
		Name:       "new-name",
		Namespaces: []string{"admin", "engineering"},
	})
	require.NoError(t, err)

	org, err = clients.OrganizationsClient.Get(ctx, &services.GetOrganizationRequest{
		Id: orgRes.Id,
	})
	require.NoError(t, err)
	require.Equal(t, int64(2), org.Version)
	require.Equal(t, "new-name", org.Name, org.Name)
	require.Equal(t, 2, len(org.Namespaces))

	query, err := clients.OrganizationsClient.Query(ctx, &services.QueryOrganizationRequest{})
	require.NoError(t, err)
	_, err = query.Recv()
	require.NoError(t, err)

	_, err = clients.OrganizationsClient.Delete(ctx, &services.DeleteOrganizationRequest{
		Id: orgRes.Id,
	})
	require.NoError(t, err)

	// should not find it after deleting it
	org, err = clients.OrganizationsClient.Get(ctx, &services.GetOrganizationRequest{
		Id: orgRes.Id,
	})
	require.Error(t, err)
}
