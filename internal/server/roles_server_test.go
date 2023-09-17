package server

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_ShouldCreateAndGetAndDeleteRole(t *testing.T) {
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
		"Should Create/Update/Get/Delete Roles": testShouldCRUDRoles,
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

func testShouldCRUDRoles(t *testing.T, clients Clients) {
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

	roleRes, err := clients.RolesClient.Create(ctx, &services.CreateRoleRequest{
		Name:           "role-name",
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
	})
	require.NoError(t, err)

	_, err = clients.RolesClient.AddPermissions(ctx, &services.AddPermissionsToRoleRequest{
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
		RoleId:         roleRes.Id,
		PermissionIds:  []string{permissionRes.Id},
	})
	require.NoError(t, err)

	queryRes, err := clients.RolesClient.Query(ctx, &services.QueryRoleRequest{
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
		Predicates:     map[string]string{"id": roleRes.Id},
	})
	require.NoError(t, err)
	role, err := queryRes.Recv()
	require.NoError(t, err)
	require.Equal(t, int64(2), role.Version)
	require.Equal(t, "role-name", role.Name)
	require.Equal(t, 1, len(role.PermissionIds))

	_, err = clients.RolesClient.Update(ctx, &services.UpdateRoleRequest{
		Id:             role.Id,
		Name:           "new-name",
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
	})
	require.NoError(t, err)

	query, err := clients.RolesClient.Query(context.Background(), &services.QueryRoleRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
	})
	require.NoError(t, err)

	role, err = query.Recv()
	require.NoError(t, err)
	require.Equal(t, "new-name", role.Name)

	_, err = clients.RolesClient.DeletePermissions(ctx, &services.DeletePermissionsToRoleRequest{
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
		RoleId:         roleRes.Id,
		PermissionIds:  []string{permissionRes.Id},
	})
	require.NoError(t, err)

	_, err = clients.RolesClient.Delete(context.Background(), &services.DeleteRoleRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
		Id:             role.Id,
	})
	require.NoError(t, err)

	// should not find it after deleting it
	query, err = clients.RolesClient.Query(context.Background(), &services.QueryRoleRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
	})
	require.NoError(t, err)
	_, err = query.Recv()
	require.Error(t, err)
}
