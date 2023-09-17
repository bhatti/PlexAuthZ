package db

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Should_CRUD_Relationships(t *testing.T) {
	// GIVEN auth-service
	ctx := context.TODO()
	store, org, err := newAuthServiceAndOrg()
	require.NoError(t, err)

	resource, err := store.CreateResource(ctx, org.Id, &types.Resource{
		Namespace:      org.Namespaces[0],
		Name:           "resource1",
		AllowedActions: []string{"read", "write"},
	})
	require.NoError(t, err)
	// AND principal
	principal, err := domain.NewPrincipalBuilder().
		WithNamespaces(org.Namespaces...).
		WithUsername("user1").
		WithOrganizationId(org.Id).Build()
	require.NoError(t, err)
	// WHEN Creating principal
	principal, err = store.CreatePrincipal(ctx, principal)
	// THEN it should not fail
	require.NoError(t, err)

	// AND relationship
	relationship, err := domain.NewRelationshipBuilder().
		WithNamespace(org.Namespaces[0]).
		WithRelation("rel-1").
		WithResourceId(resource.Id).
		WithAttribute("k1", "v1").
		WithAttribute("k2", "v2").
		WithPrincipalId(principal.Id).Build()
	require.NoError(t, err)

	// WHEN creating a relationship
	savedRelationship, err := store.CreateRelationship(ctx, org.Id, relationship)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN updating a relationship
	err = store.UpdateRelationship(ctx, org.Id, relationship)
	// THEN it should not fail
	require.NoError(t, err)

	// WHEN querying an organization by id
	res, _, err := store.GetRelationships(ctx, org.Id, org.Namespaces[0], map[string]string{"id": savedRelationship.Id}, "", 0)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, 1, len(res))

	// WHEN querying an organization by scope
	res, _, err = store.GetRelationships(ctx, org.Id, org.Namespaces[0], map[string]string{"relation": "rel-1"}, "", 0)
	// THEN it should not fail
	require.NoError(t, err)
	require.Equal(t, 1, len(res))

	// Iterating relationships
	for _, next := range res {
		require.Equal(t, relationship.Namespace, next.Namespace)
		require.Equal(t, relationship.Relation, next.Relation)
		require.Equal(t, resource.Id, next.ResourceId)
		require.Equal(t, 2, len(next.Attributes))
		// WHEN deleting a relationship
		err := store.DeleteRelationship(ctx, org.Id, org.Namespaces[0], next.Id)
		// THEN it should not fail
		require.NoError(t, err)
	}
}
