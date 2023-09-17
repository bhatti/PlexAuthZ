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

func Test_ShouldSucceedWithPermissionsGetAndQuery(t *testing.T) {
	to, ctrl, err := newTestPermissionsController()
	require.NoError(t, err)

	perm := &types.Permission{
		Namespace:  to.permission.Namespace,
		Scope:      uuid.NewV4().String(),
		Actions:    []string{"*"},
		ResourceId: to.resource.Id,
	}
	{

		reqB, err := json.Marshal(perm)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.permission.Namespace + "/permissions")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN creating permission
		err = ctrl.create(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		createRes := ctx.Result.(*services.CreatePermissionResponse)
		require.NotEqual(t, "", createRes.Id)
		perm.Id = createRes.Id
	}

	// Now querying...
	{
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.permission.Namespace + "/permissions/")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN querying permission
		err = ctrl.query(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		queryRes := ctx.Result.([]*types.Permission)
		require.True(t, len(queryRes) > 0)
	}
}

func Test_ShouldSucceedWithPermissionsCreateAndDelete(t *testing.T) {
	to, ctrl, err := newTestPermissionsController()
	require.NoError(t, err)

	perm := &types.Permission{
		Namespace:  to.permission.Namespace,
		Scope:      uuid.NewV4().String(),
		Actions:    []string{"*"},
		ResourceId: to.resource.Id,
	}
	{
		reqB, err := json.Marshal(perm)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.permission.Namespace + "/permissions")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN creating perm
		err = ctrl.create(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		createRes := ctx.Result.(*services.CreatePermissionResponse)
		require.NotEqual(t, "", createRes.Id)
		perm.Id = createRes.Id
	}

	// Now deleting ...
	{
		reqB, err := json.Marshal(perm)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.permission.Namespace + "/permissions/" + perm.Id)
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["id"] = perm.Id
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN deleting perm
		err = ctrl.delete(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func Test_ShouldSucceedWithPermissionsCreateAndUpdate(t *testing.T) {
	to, ctrl, err := newTestPermissionsController()
	require.NoError(t, err)

	perm := &types.Permission{
		Namespace:  to.permission.Namespace,
		Scope:      uuid.NewV4().String(),
		Actions:    []string{"*"},
		ResourceId: to.resource.Id,
	}
	{
		reqB, err := json.Marshal(perm)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.permission.Namespace + "/permissions")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN creating perm
		err = ctrl.create(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		createRes := ctx.Result.(*services.CreatePermissionResponse)
		require.NotEqual(t, "", createRes.Id)
		perm.Id = createRes.Id
	}

	// Now updating...
	{
		reqB, err := json.Marshal(perm)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.permission.Namespace + "/permissions" + perm.Id)
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["id"] = perm.Id
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN updating perm
		err = ctrl.update(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func newTestPermissionsController() (to *testObjects, ctrl *PermissionsController, err error) {
	webServer := web.NewStubWebServer()
	if to, err = newTestObjects(); err != nil {
		return
	}
	ctrl = NewPermissionsController(to.config, to.authService, webServer)
	return
}
