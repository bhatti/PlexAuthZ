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

func Test_ShouldSucceedWithResourcesGetAndQuery(t *testing.T) {
	to, ctrl, err := newTestResourcesController()
	require.NoError(t, err)

	resource := &types.Resource{
		Namespace:      to.resource.Namespace,
		Name:           uuid.NewV4().String(),
		Capacity:       10,
		AllowedActions: []string{"read", "write"},
	}
	{
		reqB, err := json.Marshal(resource)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.resource.Namespace + "/resources")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN creating resource
		err = ctrl.create(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		createRes := ctx.Result.(*services.CreateResourceResponse)
		require.NotEqual(t, "", createRes.Id)
		resource.Id = createRes.Id
	}

	// Now querying...
	{
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.resource.Namespace + "/resources/")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN querying resource
		err = ctrl.query(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		queryRes := ctx.Result.([]*types.Resource)
		require.True(t, len(queryRes) > 0)
	}
}

func Test_ShouldSucceedWithResourcesCreateAndDelete(t *testing.T) {
	to, ctrl, err := newTestResourcesController()
	require.NoError(t, err)

	resource := &types.Resource{
		Namespace:      to.resource.Namespace,
		Name:           uuid.NewV4().String(),
		Capacity:       10,
		AllowedActions: []string{"read", "write"},
	}
	{
		reqB, err := json.Marshal(resource)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.resource.Namespace + "/resources")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN creating resource
		err = ctrl.create(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		createRes := ctx.Result.(*services.CreateResourceResponse)
		require.NotEqual(t, "", createRes.Id)
		resource.Id = createRes.Id
	}

	// Now deleting ...
	{
		reqB, err := json.Marshal(resource)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.resource.Namespace + "/resources/" + resource.Id)
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["id"] = resource.Id
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN deleting resource
		err = ctrl.delete(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func Test_ShouldSucceedWithResourcesCreateAndUpdate(t *testing.T) {
	to, ctrl, err := newTestResourcesController()
	require.NoError(t, err)

	resource := &types.Resource{
		Namespace:      to.resource.Namespace,
		Name:           uuid.NewV4().String(),
		Capacity:       10,
		AllowedActions: []string{"read", "write"},
	}
	{
		reqB, err := json.Marshal(resource)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.resource.Namespace + "/resources")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN creating resource
		err = ctrl.create(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		createRes := ctx.Result.(*services.CreateResourceResponse)
		require.NotEqual(t, "", createRes.Id)
		resource.Id = createRes.Id
	}

	// Now updating...
	{
		reqB, err := json.Marshal(resource)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.resource.Namespace + "/resources" + resource.Id)
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["id"] = resource.Id
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN updating resource
		err = ctrl.update(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func newTestResourcesController() (to *testObjects, ctrl *ResourcesController, err error) {
	webServer := web.NewStubWebServer()
	if to, err = newTestObjects(); err != nil {
		return
	}
	ctrl = NewResourcesController(to.config, to.authService, webServer)
	return
}
