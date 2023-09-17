package controller

import (
	"bytes"
	"encoding/json"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/web"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func Test_ShouldSucceedWithRelationshipsGetAndQuery(t *testing.T) {
	to, ctrl, err := newTestRelationshipsController()
	require.NoError(t, err)

	relation := &types.Relationship{
		Namespace:   to.relation.Namespace,
		Relation:    uuid.NewV4().String(),
		PrincipalId: to.principal.Id,
		ResourceId:  to.resource.Id,
	}
	{
		reqB, err := json.Marshal(relation)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.relation.Namespace + "/relations")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN creating relation
		err = ctrl.create(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		createRes := ctx.Result.(*services.CreateRelationshipResponse)
		require.NotEqual(t, "", createRes.Id)
		relation.Id = createRes.Id
	}

	// Now querying...
	{
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.relation.Namespace + "/relations/")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN querying relation
		err = ctrl.query(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		queryRes := ctx.Result.([]*types.Relationship)
		require.True(t, len(queryRes) > 0)
	}
}

func Test_ShouldSucceedWithRelationshipsCreateAndDelete(t *testing.T) {
	to, ctrl, err := newTestRelationshipsController()
	require.NoError(t, err)

	relation := &types.Relationship{
		Namespace:   to.relation.Namespace,
		Relation:    uuid.NewV4().String(),
		PrincipalId: to.principal.Id,
		ResourceId:  to.resource.Id,
	}
	{
		reqB, err := json.Marshal(relation)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.relation.Namespace + "/relations")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN creating relation
		err = ctrl.create(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		createRes := ctx.Result.(*services.CreateRelationshipResponse)
		require.NotEqual(t, "", createRes.Id)
		relation.Id = createRes.Id
	}

	// Now deleting ...
	{
		reqB, err := json.Marshal(relation)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.relation.Namespace + "/relations/" + relation.Id)
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["id"] = relation.Id
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN deleting relation
		err = ctrl.delete(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func Test_ShouldSucceedWithRelationshipsCreateAndUpdate(t *testing.T) {
	to, ctrl, err := newTestRelationshipsController()
	require.NoError(t, err)

	relation := &types.Relationship{
		Namespace:   to.relation.Namespace,
		Relation:    uuid.NewV4().String(),
		PrincipalId: to.principal.Id,
		ResourceId:  to.resource.Id,
	}
	{
		reqB, err := json.Marshal(relation)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.relation.Namespace + "/relations")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN creating relation
		err = ctrl.create(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		createRes := ctx.Result.(*services.CreateRelationshipResponse)
		require.NotEqual(t, "", createRes.Id)
		relation.Id = createRes.Id
	}

	// Now updating...
	{
		reqB, err := json.Marshal(relation)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.relation.Namespace + "/relations" + relation.Id)
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["id"] = relation.Id
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN updating relation
		err = ctrl.update(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func newTestRelationshipsController() (to *testObjects, ctrl *RelationshipsController, err error) {
	webServer := web.NewStubWebServer()
	if to, err = newTestObjects(); err != nil {
		return
	}
	ctrl = NewRelationshipsController(to.config, to.authService, webServer)
	return
}
