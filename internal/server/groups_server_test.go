package server

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_ShouldCreateAndGetAndDeleteGroup(t *testing.T) {
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
		"Should Create/Update/Get/Delete Groups": testShouldCRUDGroups,
	} {
		t.Run(scenario, func(t *testing.T) {
			for _, clientType := range clientTypes {
				clients, teardown := SetupGrpcServerForTesting(t, cfg, clientType, nil)
				fn(t, clients)
				teardown()
			}
		})
	}
}

func testShouldCRUDGroups(t *testing.T, clients Clients) {
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
	groupRes, err := clients.GroupsClient.Create(ctx, &services.CreateGroupRequest{
		Name:           "group-name",
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
	})
	require.NoError(t, err)

	queryRes, err := clients.GroupsClient.Query(ctx, &services.QueryGroupRequest{
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
		Predicates:     map[string]string{"id": groupRes.Id},
	})
	require.NoError(t, err)
	group, err := queryRes.Recv()
	require.NoError(t, err)
	require.Equal(t, int64(1), group.Version)
	require.Equal(t, "group-name", group.Name)

	_, err = clients.GroupsClient.Update(ctx, &services.UpdateGroupRequest{
		Id:             group.Id,
		Name:           "new-name",
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
	})
	require.NoError(t, err)

	role1, err := clients.RolesClient.Create(ctx, &services.CreateRoleRequest{
		Name:           "role1",
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
	})
	require.NoError(t, err)
	role2, err := clients.RolesClient.Create(ctx, &services.CreateRoleRequest{
		Name:           "role2",
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
	})
	require.NoError(t, err)
	_, err = clients.GroupsClient.AddRoles(ctx, &services.AddRolesToGroupRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
		RoleIds:        []string{role1.Id, role2.Id},
		GroupId:        group.Id,
	})
	require.NoError(t, err)

	query, err := clients.GroupsClient.Query(context.Background(), &services.QueryGroupRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
	})
	require.NoError(t, err)

	group, err = query.Recv()
	require.NoError(t, err)
	require.Equal(t, 2, len(group.RoleIds))
	require.Equal(t, "new-name", group.Name)

	_, err = clients.GroupsClient.DeleteRoles(ctx, &services.DeleteRolesToGroupRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
		RoleIds:        []string{role1.Id, role2.Id},
		GroupId:        group.Id,
	})
	require.NoError(t, err)

	_, err = clients.GroupsClient.Delete(context.Background(), &services.DeleteGroupRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
		Id:             group.Id,
	})
	require.NoError(t, err)

	// should not find it after deleting it
	query, err = clients.GroupsClient.Query(context.Background(), &services.QueryGroupRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
	})
	require.NoError(t, err)
	_, err = query.Recv()
	require.Error(t, err)
}
