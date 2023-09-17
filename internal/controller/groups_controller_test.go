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

func Test_ShouldSucceedWithGroupCreateAndUpdate(t *testing.T) {
	to, ctrl, err := newTestGroupController()
	require.NoError(t, err)

	group := &types.Group{
		Namespace: to.permission.Namespace,
		Name:      uuid.NewV4().String(),
	}
	{
		reqB, err := json.Marshal(group)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" +
			to.principal.OrganizationId + "/" + to.permission.Namespace + "/groups")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace
		// WHEN creating group
		err = ctrl.create(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		createRes := ctx.Result.(*services.CreateGroupResponse)
		require.NotEqual(t, "", createRes.Id)
		group.Id = createRes.Id
	}
	// Now updating...
	{
		reqB, err := json.Marshal(group)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" +
			to.principal.OrganizationId + "/" + to.permission.Namespace + "/groups")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = group.Namespace
		ctx.Params["id"] = group.Id

		// WHEN updating group
		err = ctrl.update(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func Test_ShouldSucceedWithGroupQuery(t *testing.T) {
	to, ctrl, err := newTestGroupController()
	require.NoError(t, err)
	u, err := url.Parse("https://localhost:8080/api/v1/" +
		to.principal.OrganizationId + "/" + to.group.Namespace + "/groups?offset=0&limit=10")
	require.NoError(t, err)

	ctx := web.NewStubContext(&http.Request{Body: nil, URL: u})
	ctx.Params["organization_id"] = to.principal.OrganizationId
	ctx.Params["principal_id"] = to.principal.Id
	ctx.Params["namespace"] = to.group.Namespace

	// WHEN querying groups
	err = ctrl.query(ctx)
	// THEN it should not fail
	require.NoError(t, err)
	res := ctx.Result.([]*types.Group)
	require.True(t, len(res) > 0)
}

func Test_ShouldSucceedWithGroupDelete(t *testing.T) {
	to, ctrl, err := newTestGroupController()
	require.NoError(t, err)
	u, err := url.Parse("https://localhost:8080/api/v1/" +
		to.principal.OrganizationId + "/" + to.permission.Namespace + "/groups/")
	require.NoError(t, err)

	ctx := web.NewStubContext(&http.Request{Body: nil, URL: u})
	ctx.Params["organization_id"] = to.principal.OrganizationId
	ctx.Params["principal_id"] = to.principal.Id
	ctx.Params["namespace"] = to.group.Namespace
	ctx.Params["id"] = to.group.Id
	// WHEN deleting group
	err = ctrl.delete(ctx)
	// THEN it should not fail
	require.NoError(t, err)
}

func Test_ShouldSucceedWithAddAndDeleteRolesToGroup(t *testing.T) {
	to, ctrl, err := newTestGroupController()
	require.NoError(t, err)
	{
		req := &services.AddRolesToGroupRequest{
			Namespace: to.group.Namespace,
			RoleIds:   []string{to.role.Id},
		}
		reqB, err := json.Marshal(req)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" +
			to.principal.OrganizationId + "/" + to.permission.Namespace + "/groups/" + to.group.Id)
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})

		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace
		ctx.Params["id"] = to.group.Id

		// WHEN adding role to group
		err = ctrl.addRoles(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}

	// Now deleting...
	{
		req := &services.DeleteRolesToGroupRequest{
			Namespace: to.group.Namespace,
			RoleIds:   []string{to.role.Id},
		}
		reqB, err := json.Marshal(req)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" +
			to.principal.OrganizationId + "/" + to.permission.Namespace + "/groups/" + to.group.Id)
		require.NoError(t, err)
		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})

		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["principal_id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace
		ctx.Params["id"] = to.group.Id

		// WHEN deleting role to group
		err = ctrl.deleteRoles(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func newTestGroupController() (to *testObjects, ctrl *GroupsController, err error) {
	webServer := web.NewStubWebServer()
	if to, err = newTestObjects(); err != nil {
		return
	}
	ctrl = NewGroupsController(to.config, to.authService, webServer)
	return
}
