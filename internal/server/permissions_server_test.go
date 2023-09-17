package server

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_ShouldCreateAndGetAndDeletePermission(t *testing.T) {
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
		"Should Create/Update/Get/Delete Permissions": testShouldCRUDPermissions,
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

func testShouldCRUDPermissions(t *testing.T, clients Clients) {
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
	resource, err := clients.ResourcesClient.Create(ctx, &services.CreateResourceRequest{
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
		Name:           "res-1",
		AllowedActions: []string{"read", "write"},
	})
	require.NoError(t, err)

	permissionRes, err := clients.PermissionsClient.Create(ctx, &services.CreatePermissionRequest{
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
		ResourceId:     resource.Id,
		Actions:        []string{"read"},
	})
	require.NoError(t, err)

	queryRes, err := clients.PermissionsClient.Query(ctx, &services.QueryPermissionRequest{
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
		Predicates:     map[string]string{"id": permissionRes.Id},
	})
	require.NoError(t, err)
	permission, err := queryRes.Recv()
	require.NoError(t, err)
	require.Equal(t, int64(1), permission.Version)
	require.Equal(t, resource.Id, permission.ResourceId)
	require.Equal(t, 1, len(permission.Actions))
	require.Equal(t, "read", permission.Actions[0])

	_, err = clients.PermissionsClient.Update(ctx, &services.UpdatePermissionRequest{
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
		ResourceId:     resource.Id,
		Id:             permission.Id,
		Actions:        []string{"write"},
	})
	require.NoError(t, err)

	query, err := clients.PermissionsClient.Query(context.Background(), &services.QueryPermissionRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
	})
	require.NoError(t, err)
	permission, err = query.Recv()
	require.NoError(t, err)
	require.Equal(t, int64(2), permission.Version)
	require.Equal(t, resource.Id, permission.ResourceId)
	require.Equal(t, 1, len(permission.Actions))
	require.Equal(t, "write", permission.Actions[0])

	_, err = clients.PermissionsClient.Delete(context.Background(), &services.DeletePermissionRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
		Id:             permission.Id,
	})
	require.NoError(t, err)

	// should not find it after deleting it
	query, err = clients.PermissionsClient.Query(context.Background(), &services.QueryPermissionRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
	})
	require.NoError(t, err)
	_, err = query.Recv()
	require.Error(t, err)
}
