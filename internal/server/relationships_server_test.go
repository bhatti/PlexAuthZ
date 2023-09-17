package server

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_ShouldCreateAndGetAndDeleteRelationship(t *testing.T) {
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
		"Should Create/Update/Get/Delete Relationships": testShouldCRUDRelationships,
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

func testShouldCRUDRelationships(t *testing.T, clients Clients) {
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
		OrganizationId: orgRes.Id,
		Username:       "principal1",
		Namespaces:     []string{"admin", "finance", "engineering"},
	})
	require.NoError(t, err)

	resource, err := clients.ResourcesClient.Create(ctx, &services.CreateResourceRequest{
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
		Name:           "res-1",
		AllowedActions: []string{"read", "write"},
	})
	require.NoError(t, err)

	relationshipRes, err := clients.RelationshipsClient.Create(ctx, &services.CreateRelationshipRequest{
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
		ResourceId:     resource.Id,
		PrincipalId:    principalRes.Id,
		Relation:       "rel-1",
	})
	require.NoError(t, err)

	queryRes, err := clients.RelationshipsClient.Query(ctx, &services.QueryRelationshipRequest{
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
		Predicates:     map[string]string{"id": relationshipRes.Id},
	})
	require.NoError(t, err)
	relationship, err := queryRes.Recv()
	require.NoError(t, err)
	require.Equal(t, int64(1), relationship.Version)
	require.Equal(t, resource.Id, relationship.ResourceId)

	_, err = clients.RelationshipsClient.Update(ctx, &services.UpdateRelationshipRequest{
		OrganizationId: orgRes.Id,
		Namespace:      "admin",
		ResourceId:     resource.Id,
		Id:             relationship.Id,
		PrincipalId:    principalRes.Id,
		Relation:       "rel-2",
	})
	require.NoError(t, err)

	query, err := clients.RelationshipsClient.Query(context.Background(), &services.QueryRelationshipRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
	})
	require.NoError(t, err)
	relationship, err = query.Recv()
	require.NoError(t, err)
	require.Equal(t, int64(2), relationship.Version)
	require.Equal(t, resource.Id, relationship.ResourceId)
	require.Equal(t, "rel-2", relationship.Relation)

	_, err = clients.RelationshipsClient.Delete(context.Background(), &services.DeleteRelationshipRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
		Id:             relationship.Id,
	})
	require.NoError(t, err)

	// should not find it after deleting it
	query, err = clients.RelationshipsClient.Query(context.Background(), &services.QueryRelationshipRequest{
		Namespace:      "admin",
		OrganizationId: orgRes.Id,
	})
	require.NoError(t, err)
	_, err = query.Recv()
	require.Error(t, err)
}
