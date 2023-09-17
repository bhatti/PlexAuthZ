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

func Test_ShouldSucceedWithOrgGetAndQuery(t *testing.T) {
	_, ctrl, err := newTestOrgController()
	require.NoError(t, err)

	org := &types.Organization{
		Namespaces: []string{"finance"},
		Name:       uuid.NewV4().String(),
	}
	{
		reqB, err := json.Marshal(org)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/organizations")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		// WHEN creating org
		err = ctrl.create(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		createRes := ctx.Result.(*services.CreateOrganizationResponse)
		require.NotEqual(t, "", createRes.Id)
		org.Id = createRes.Id
	}

	// Now getting ...
	{
		u, err := url.Parse("https://localhost:8080/api/v1/organizations/" + org.Id + "?id=" + org.Id)
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{URL: u})
		ctx.Params["id"] = org.Id

		// WHEN getting org
		err = ctrl.get(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		getRes := ctx.Result.(*types.Organization)
		require.Equal(t, org.Id, getRes.Id)
	}

	// Now querying...
	{
		u, err := url.Parse("https://localhost:8080/api/v1/organizations/" + org.Id + "?id=" + org.Id)
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{URL: u})
		ctx.Params["id"] = org.Id

		// WHEN querying org
		err = ctrl.query(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		queryRes := ctx.Result.([]*types.Organization)
		require.True(t, len(queryRes) > 0)
	}
}

func Test_ShouldSucceedWithOrgCreateAndDelete(t *testing.T) {
	_, ctrl, err := newTestOrgController()
	org := &types.Organization{
		Namespaces: []string{"finance"},
		Name:       uuid.NewV4().String(),
	}
	require.NoError(t, err)
	{
		reqB, err := json.Marshal(org)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/organizations")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		// WHEN creating org
		err = ctrl.create(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		createRes := ctx.Result.(*services.CreateOrganizationResponse)
		require.NotEqual(t, "", createRes.Id)
		org.Id = createRes.Id
	}

	// Now deleting ...
	{
		reqB, err := json.Marshal(org)
		require.NoError(t, err)
		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/organizations/" + org.Id)
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["id"] = org.Id

		// WHEN deleting org
		err = ctrl.delete(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func Test_ShouldSucceedWithOrgCreateAndUpdate(t *testing.T) {
	_, ctrl, err := newTestOrgController()
	require.NoError(t, err)

	org := &types.Organization{
		Namespaces: []string{"finance"},
		Name:       uuid.NewV4().String(),
	}
	{
		reqB, err := json.Marshal(org)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/organizations")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		// WHEN creating org
		err = ctrl.create(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		createRes := ctx.Result.(*services.CreateOrganizationResponse)
		require.NotEqual(t, "", createRes.Id)
		org.Id = createRes.Id
	}

	// Now updating...
	{
		reqB, err := json.Marshal(org)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/organizations/" + org.Id)
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["id"] = org.Id

		// WHEN updating org
		err = ctrl.update(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func newTestOrgController() (to *testObjects, ctrl *OrganizationsController, err error) {
	webServer := web.NewStubWebServer()
	if to, err = newTestObjects(); err != nil {
		return
	}
	ctrl = NewOrganizationsController(to.config, to.authService, webServer)
	return
}
