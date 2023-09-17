package controller

import (
	"bytes"
	"encoding/json"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/web"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func Test_ShouldSucceedWithAuthorize(t *testing.T) {
	to, ctrl, err := newTestAuthController()
	require.NoError(t, err)
	req := &services.AuthRequest{
		Namespace:      to.permission.Namespace,
		OrganizationId: to.principal.OrganizationId,
		PrincipalId:    to.principal.Id,
		Action:         "read",
		Resource:       "paper",
	}
	reqB, err := json.Marshal(req)
	require.NoError(t, err)
	reader := io.NopCloser(bytes.NewReader(reqB))
	u, err := url.Parse("https://localhost:8080/api/v1/" +
		to.principal.OrganizationId + "/" + to.permission.Namespace + "/" + to.principal.Id + "/auth")
	require.NoError(t, err)

	ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
	ctx.Params["organization_id"] = to.principal.OrganizationId
	ctx.Params["principal_id"] = to.principal.Id
	ctx.Params["namespace"] = to.permission.Namespace
	// WHEN invoking auth with valid permission
	err = ctrl.auth(ctx)
	// THEN it should not fail
	require.NoError(t, err)
	res := ctx.Result.(*services.AuthResponse)
	require.Equal(t, types.Effect_PERMITTED, res.Effect)
}

func Test_ShouldSucceedWithCheck(t *testing.T) {
	to, ctrl, err := newTestAuthController()
	require.NoError(t, err)
	req := &services.CheckConstraintsRequest{
		Namespace:      to.permission.Namespace,
		OrganizationId: to.principal.OrganizationId,
		PrincipalId:    to.principal.Id,
		Constraints:    "eq 1 1",
	}
	reqB, err := json.Marshal(req)
	require.NoError(t, err)
	reader := io.NopCloser(bytes.NewReader(reqB))
	u, err := url.Parse("https://localhost:8080/api/v1/" +
		to.principal.OrganizationId + "/" + to.permission.Namespace + "/" + to.principal.Id + "/auth/constraints")
	require.NoError(t, err)

	ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
	ctx.Params["organization_id"] = to.principal.OrganizationId
	ctx.Params["principal_id"] = to.principal.Id
	ctx.Params["namespace"] = to.permission.Namespace

	// WHEN invoking auth with valid permission
	err = ctrl.check(ctx)
	// THEN it should not fail
	require.NoError(t, err)
	res := ctx.Result.(*services.CheckConstraintsResponse)
	require.True(t, res.Matched)
}

func Test_ShouldSucceedWithResourceAllocation(t *testing.T) {
	to, ctrl, err := newTestAuthController()
	require.NoError(t, err)
	{
		req := &services.AllocateResourceRequest{
			Namespace:      to.resource.Namespace,
			OrganizationId: to.org.Id,
			ResourceId:     to.resource.Id,
			PrincipalId:    to.principal.Id,
		}
		reqB, err := json.Marshal(req)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.resource.Namespace +
			"/resources/" + to.resource.Id + "/allocate/" + to.principal.Id)
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.resource.Namespace
		ctx.Params["id"] = to.resource.Id

		// WHEN allocating resource
		err = ctrl.allocate(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
	{
		req := &services.CountResourceInstancesRequest{
			Namespace:      to.resource.Namespace,
			OrganizationId: to.org.Id,
			ResourceId:     to.resource.Id,
		}
		reqB, err := json.Marshal(req)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.resource.Namespace +
			"/resources/" + to.resource.Id + "/instance_count")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.resource.Namespace
		ctx.Params["id"] = to.resource.Id

		// WHEN counting resource
		_, ctrl, err := newTestResourcesController()
		require.NoError(t, err)

		err = ctrl.allocatedInstancesCount(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		countRes := ctx.Result.(*services.CountResourceInstancesResponse)
		require.True(t, countRes.Allocated > 0)
	}
	{
		req := &services.QueryResourceInstanceRequest{
			Namespace:      to.resource.Namespace,
			OrganizationId: to.org.Id,
			ResourceId:     to.resource.Id,
		}
		reqB, err := json.Marshal(req)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.resource.Namespace +
			"/resources/" + to.resource.Id + "/instances")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.resource.Namespace
		ctx.Params["id"] = to.resource.Id

		// WHEN counting resource
		_, ctrl, err := newTestResourcesController()
		require.NoError(t, err)
		err = ctrl.queryAllocatedInstances(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		res := ctx.Result.([]*types.ResourceInstance)
		require.True(t, len(res) > 0)
	}
	{
		req := &services.DeallocateResourceRequest{
			Namespace:      to.resource.Namespace,
			OrganizationId: to.org.Id,
			ResourceId:     to.resource.Id,
			PrincipalId:    to.principal.Id,
		}
		reqB, err := json.Marshal(req)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.resource.Namespace +
			"/resources/" + to.resource.Id + "/deallocate/" + to.principal.Id)
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.resource.Namespace
		ctx.Params["id"] = to.resource.Id

		// WHEN allocating resource
		err = ctrl.deallocate(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func newTestAuthController() (to *testObjects, ctrl *AuthController, err error) {
	webServer := web.NewStubWebServer()
	if to, err = newTestObjects(); err != nil {
		return
	}
	ctrl, err = NewAuthController(to.config, to.authService, webServer)
	return
}
