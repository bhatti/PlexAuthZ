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

func Test_ShouldSucceedWithRolesGetAndQuery(t *testing.T) {
	to, ctrl, err := newTestRolesController()
	require.NoError(t, err)

	role := &types.Role{
		Namespace: to.role.Namespace,
		Name:      uuid.NewV4().String(),
	}
	{
		reqB, err := json.Marshal(role)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.role.Namespace + "/roles")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN creating role
		err = ctrl.create(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		createRes := ctx.Result.(*services.CreateRoleResponse)
		require.NotEqual(t, "", createRes.Id)
		role.Id = createRes.Id
	}

	// Now querying...
	{
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.role.Namespace + "/roles/")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN querying role
		err = ctrl.query(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		queryRes := ctx.Result.([]*types.Role)
		require.True(t, len(queryRes) > 0)
	}
}

func Test_ShouldSucceedWithRolesCreateAndDelete(t *testing.T) {
	to, ctrl, err := newTestRolesController()
	require.NoError(t, err)

	role := &types.Role{
		Namespace: to.role.Namespace,
		Name:      uuid.NewV4().String(),
	}
	{
		reqB, err := json.Marshal(role)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.role.Namespace + "/roles")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN creating role
		err = ctrl.create(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		createRes := ctx.Result.(*services.CreateRoleResponse)
		require.NotEqual(t, "", createRes.Id)
		role.Id = createRes.Id
	}

	// Now deleting ...
	{
		reqB, err := json.Marshal(role)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.role.Namespace + "/roles/" + role.Id)
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["id"] = role.Id
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN deleting role
		err = ctrl.delete(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func Test_ShouldSucceedWithRolesCreateAndUpdate(t *testing.T) {
	to, ctrl, err := newTestRolesController()
	require.NoError(t, err)

	role := &types.Role{
		Namespace: to.role.Namespace,
		Name:      uuid.NewV4().String(),
	}
	{
		reqB, err := json.Marshal(role)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.role.Namespace + "/roles")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN creating role
		err = ctrl.create(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		createRes := ctx.Result.(*services.CreateRoleResponse)
		require.NotEqual(t, "", createRes.Id)
		role.Id = createRes.Id
	}

	// Now updating...
	{
		reqB, err := json.Marshal(role)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.role.Namespace + "/roles" + role.Id)
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["id"] = role.Id
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN updating role
		err = ctrl.update(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func Test_ShouldSucceedWithRolePermission(t *testing.T) {
	to, ctrl, err := newTestRolesController()
	require.NoError(t, err)
	{
		req := &services.AddPermissionsToRoleRequest{
			Namespace:      to.role.Namespace,
			OrganizationId: to.org.Id,
			RoleId:         to.role.Id,
			PermissionIds:  []string{to.permission.Id},
		}
		reqB, err := json.Marshal(req)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.role.Namespace +
			"/roles/" + to.role.Id + "/permissions/" + to.permission.Id + "/add")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.role.Namespace
		ctx.Params["id"] = to.role.Id

		// WHEN adding permission to role
		err = ctrl.addPermissions(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
	{
		req := &services.DeletePermissionsToRoleRequest{
			Namespace:      to.role.Namespace,
			OrganizationId: to.org.Id,
			RoleId:         to.role.Id,
			PermissionIds:  []string{to.permission.Id},
		}
		reqB, err := json.Marshal(req)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/" + to.role.Namespace +
			"/roles/" + to.role.Id + "/permissions/" + to.permission.Id + "/delete")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.role.Namespace
		ctx.Params["id"] = to.role.Id

		// WHEN deleting permission to role
		err = ctrl.deletePermissions(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func newTestRolesController() (to *testObjects, ctrl *RolesController, err error) {
	webServer := web.NewStubWebServer()
	if to, err = newTestObjects(); err != nil {
		return
	}
	ctrl = NewRolesController(to.config, to.authService, webServer)
	return
}
