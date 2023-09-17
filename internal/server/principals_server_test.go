package server

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_ShouldCreateAndGetAndDeletePrincipal(t *testing.T) {
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
		"Should Create/Update/Get/Delete Principals": testShouldCRUDPrincipals,
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

func testShouldCRUDPrincipals(t *testing.T, clients Clients) {
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

	principalRes, err := clients.PrincipalsClient.Create(ctx, &services.CreatePrincipalRequest{
		Username:       "principal1",
		Name:           "john",
		OrganizationId: orgRes.Id,
		Namespaces:     []string{"admin", "finance", "engineering"},
	})
	require.NoError(t, err)

	relationRes, err := clients.RelationshipsClient.Create(ctx, &services.CreateRelationshipRequest{
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
		Relation:       "relation",
		ResourceId:     resource.Id,
		PrincipalId:    principalRes.Id,
	})
	require.NoError(t, err)

	queryRes, err := clients.PrincipalsClient.Query(ctx, &services.QueryPrincipalRequest{
		OrganizationId: orgRes.Id,
		Predicates:     map[string]string{"id": principalRes.Id},
	})
	require.NoError(t, err)
	principal, err := queryRes.Recv()
	require.NoError(t, err)
	require.Equal(t, int64(1), principal.Version)
	require.Equal(t, "principal1", principal.Username)
	require.Equal(t, "john", principal.Name)

	_, err = clients.PrincipalsClient.Update(ctx, &services.UpdatePrincipalRequest{
		Id:             principal.Id,
		Username:       "user1",
		Name:           "jane",
		OrganizationId: orgRes.Id,
		Namespaces:     []string{"admin", "finance"},
	})
	require.NoError(t, err)

	saved, err := clients.PrincipalsClient.Get(ctx, &services.GetPrincipalRequest{
		OrganizationId: orgRes.Id,
		Id:             principal.Id,
		Namespace:      "admin",
	})
	require.NoError(t, err)
	require.Equal(t, "user1", saved.Username)
	require.Equal(t, "jane", saved.Name)

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
	_, err = clients.PrincipalsClient.AddRoles(ctx, &services.AddRolesToPrincipalRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
		RoleIds:        []string{role1.Id, role2.Id},
		PrincipalId:    principal.Id,
	})
	require.NoError(t, err)

	_, err = clients.PrincipalsClient.AddGroups(ctx, &services.AddGroupsToPrincipalRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
		GroupIds:       []string{groupRes.Id},
		PrincipalId:    principal.Id,
	})
	require.NoError(t, err)

	_, err = clients.PrincipalsClient.AddPermissions(ctx, &services.AddPermissionsToPrincipalRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
		PermissionIds:  []string{permissionRes.Id},
		PrincipalId:    principal.Id,
	})
	require.NoError(t, err)

	_, err = clients.PrincipalsClient.AddRelationships(ctx, &services.AddRelationshipsToPrincipalRequest{
		Namespace:       "admin",
		OrganizationId:  orgRes.Id,
		RelationshipIds: []string{relationRes.Id},
		PrincipalId:     principal.Id,
	})
	require.NoError(t, err)

	query, err := clients.PrincipalsClient.Query(context.Background(), &services.QueryPrincipalRequest{
		OrganizationId: orgRes.Id,
	})
	require.NoError(t, err)

	principal, err = query.Recv()
	require.NoError(t, err)
	require.Equal(t, 2, len(principal.RoleIds))
	require.Equal(t, "user1", principal.Username)

	require.NoError(t, err)
	_, err = clients.PrincipalsClient.DeleteRoles(ctx, &services.DeleteRolesToPrincipalRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
		RoleIds:        []string{role1.Id, role2.Id},
		PrincipalId:    principal.Id,
	})
	require.NoError(t, err)

	_, err = clients.PrincipalsClient.DeleteGroups(ctx, &services.DeleteGroupsToPrincipalRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
		GroupIds:       []string{groupRes.Id},
		PrincipalId:    principal.Id,
	})
	require.NoError(t, err)

	_, err = clients.PrincipalsClient.DeletePermissions(ctx, &services.DeletePermissionsToPrincipalRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
		PermissionIds:  []string{permissionRes.Id},
		PrincipalId:    principal.Id,
	})
	require.NoError(t, err)

	_, err = clients.PrincipalsClient.DeleteRelationships(ctx, &services.DeleteRelationshipsToPrincipalRequest{
		Namespace:       "admin",
		OrganizationId:  orgRes.Id,
		RelationshipIds: []string{relationRes.Id},
		PrincipalId:     principal.Id,
	})
	require.NoError(t, err)

	_, err = clients.PrincipalsClient.Delete(context.Background(), &services.DeletePrincipalRequest{
		OrganizationId: orgRes.Id,
		Id:             principal.Id,
	})
	require.NoError(t, err)

	// should not find it after deleting it
	query, err = clients.PrincipalsClient.Query(context.Background(), &services.QueryPrincipalRequest{
		OrganizationId: orgRes.Id,
	})
	require.NoError(t, err)
	_, err = query.Recv()
	require.Error(t, err)
}
