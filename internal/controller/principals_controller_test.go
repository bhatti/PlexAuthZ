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

func Test_ShouldSucceedWithPrincipalsGetAndQuery(t *testing.T) {
	to, ctrl, err := newTestPrincipalsController()
	require.NoError(t, err)

	principle := &types.Principal{
		Namespaces: to.org.Namespaces,
		Username:   uuid.NewV4().String(),
	}
	{
		reqB, err := json.Marshal(principle)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/principles")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId

		// WHEN creating principle
		err = ctrl.create(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		createRes := ctx.Result.(*services.CreatePrincipalResponse)
		require.NotEqual(t, "", createRes.Id)
		principle.Id = createRes.Id
	}

	// Now querying...
	{
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/principles/")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId

		// WHEN querying principle
		err = ctrl.query(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		queryRes := ctx.Result.([]*types.Principal)
		require.True(t, len(queryRes) > 0)
	}

	// Now fetching...
	{
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/principles/" + principle.Id)
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{URL: u})
		ctx.Params["id"] = principle.Id
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN getting principle
		err = ctrl.get(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		getRes := ctx.Result.(*services.GetPrincipalResponse)
		require.Equal(t, principle.Id, getRes.Id)
	}
}

func Test_ShouldSucceedWithPrincipalsCreateAndDelete(t *testing.T) {
	to, ctrl, err := newTestPrincipalsController()
	require.NoError(t, err)

	principle := &types.Principal{
		Namespaces: to.org.Namespaces,
		Username:   uuid.NewV4().String(),
	}

	{
		reqB, err := json.Marshal(principle)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/principles/")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId

		// WHEN creating principle
		err = ctrl.create(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		createRes := ctx.Result.(*services.CreatePrincipalResponse)
		require.NotEqual(t, "", createRes.Id)
		principle.Id = createRes.Id
	}

	// Now deleting ...
	{
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/principles/" + principle.Id)
		require.NoError(t, err)
		reqB, err := json.Marshal(principle)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["id"] = principle.Id
		ctx.Params["organization_id"] = to.principal.OrganizationId

		// WHEN deleting principle
		err = ctrl.delete(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func Test_ShouldSucceedWithPrincipalsCreateAndUpdate(t *testing.T) {
	to, ctrl, err := newTestPrincipalsController()
	require.NoError(t, err)

	principle := &types.Principal{
		Namespaces: to.org.Namespaces,
		Username:   uuid.NewV4().String(),
	}
	{
		reqB, err := json.Marshal(principle)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/principles/")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["id"] = to.principal.Id
		ctx.Params["namespace"] = to.group.Namespace

		// WHEN creating principle
		err = ctrl.create(ctx)
		// THEN it should not fail
		require.NoError(t, err)
		createRes := ctx.Result.(*services.CreatePrincipalResponse)
		require.NotEqual(t, "", createRes.Id)
		principle.Id = createRes.Id
	}

	// Now updating...
	{
		reqB, err := json.Marshal(principle)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/principles" + principle.Id)
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["id"] = principle.Id
		ctx.Params["organization_id"] = to.principal.OrganizationId

		// WHEN updating principle
		err = ctrl.update(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func Test_ShouldSucceedWithPrincipalsAddDeleteRoles(t *testing.T) {
	to, ctrl, err := newTestPrincipalsController()
	require.NoError(t, err)
	{
		req := &services.AddRolesToPrincipalRequest{
			OrganizationId: to.org.Id,
			Namespace:      to.role.Namespace,
			PrincipalId:    to.principal.Id,
			RoleIds:        []string{to.role.Id},
		}
		reqB, err := json.Marshal(req)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/principles/" + to.principal.Id + "/roles/add")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["namespace"] = to.role.Namespace
		ctx.Params["id"] = to.principal.Id

		// WHEN adding roles
		err = ctrl.addRoles(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
	{
		req := &services.DeleteRolesToPrincipalRequest{
			OrganizationId: to.org.Id,
			Namespace:      to.role.Namespace,
			PrincipalId:    to.principal.Id,
			RoleIds:        []string{to.role.Id},
		}
		reqB, err := json.Marshal(req)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/principles/" + to.principal.Id + "/roles/delete")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["namespace"] = to.role.Namespace
		ctx.Params["id"] = to.principal.Id

		// WHEN deleting roles
		err = ctrl.deleteRoles(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func Test_ShouldSucceedWithPrincipalsAddDeletePermissions(t *testing.T) {
	to, ctrl, err := newTestPrincipalsController()
	require.NoError(t, err)
	{
		req := &services.AddPermissionsToPrincipalRequest{
			OrganizationId: to.org.Id,
			Namespace:      to.permission.Namespace,
			PrincipalId:    to.principal.Id,
			PermissionIds:  []string{to.permission.Id},
		}
		reqB, err := json.Marshal(req)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/principles/" + to.principal.Id + "/permissions/add")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["namespace"] = to.permission.Namespace
		ctx.Params["id"] = to.principal.Id

		// WHEN adding permissions
		err = ctrl.addPermissions(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
	{
		req := &services.DeletePermissionsToPrincipalRequest{
			OrganizationId: to.org.Id,
			Namespace:      to.permission.Namespace,
			PrincipalId:    to.principal.Id,
			PermissionIds:  []string{to.permission.Id},
		}
		reqB, err := json.Marshal(req)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/principles/" + to.principal.Id + "/permissions/delete")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["namespace"] = to.permission.Namespace
		ctx.Params["id"] = to.principal.Id

		// WHEN deleting permissions
		err = ctrl.deletePermissions(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func Test_ShouldSucceedWithPrincipalsAddDeleteGroups(t *testing.T) {
	to, ctrl, err := newTestPrincipalsController()
	require.NoError(t, err)
	{
		req := &services.AddGroupsToPrincipalRequest{
			OrganizationId: to.org.Id,
			Namespace:      to.group.Namespace,
			PrincipalId:    to.principal.Id,
			GroupIds:       []string{to.group.Id},
		}
		reqB, err := json.Marshal(req)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/principles/" + to.principal.Id + "/groups/add")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["namespace"] = to.group.Namespace
		ctx.Params["id"] = to.principal.Id

		// WHEN adding groups
		err = ctrl.addGroups(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
	{
		req := &services.DeleteGroupsToPrincipalRequest{
			OrganizationId: to.org.Id,
			Namespace:      to.group.Namespace,
			PrincipalId:    to.principal.Id,
			GroupIds:       []string{to.group.Id},
		}
		reqB, err := json.Marshal(req)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/principles/" + to.principal.Id + "/groups/delete")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["namespace"] = to.group.Namespace
		ctx.Params["id"] = to.principal.Id

		// WHEN deleting groups
		err = ctrl.deleteGroups(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func Test_ShouldSucceedWithPrincipalsAddDeleteRelationships(t *testing.T) {
	to, ctrl, err := newTestPrincipalsController()
	require.NoError(t, err)
	{
		req := &services.AddRelationshipsToPrincipalRequest{
			OrganizationId:  to.org.Id,
			Namespace:       to.relation.Namespace,
			PrincipalId:     to.principal.Id,
			RelationshipIds: []string{to.relation.Id},
		}
		reqB, err := json.Marshal(req)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/principles/" + to.principal.Id + "/relations/add")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["namespace"] = to.relation.Namespace
		ctx.Params["id"] = to.principal.Id

		// WHEN adding relations
		err = ctrl.addRelationships(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
	{
		req := &services.DeleteRelationshipsToPrincipalRequest{
			OrganizationId:  to.org.Id,
			Namespace:       to.relation.Namespace,
			PrincipalId:     to.principal.Id,
			RelationshipIds: []string{to.relation.Id},
		}
		reqB, err := json.Marshal(req)
		require.NoError(t, err)

		reader := io.NopCloser(bytes.NewReader(reqB))
		u, err := url.Parse("https://localhost:8080/api/v1/" + to.org.Id + "/principles/" + to.principal.Id + "/relations/delete")
		require.NoError(t, err)

		ctx := web.NewStubContext(&http.Request{Body: reader, URL: u})
		ctx.Params["organization_id"] = to.principal.OrganizationId
		ctx.Params["namespace"] = to.relation.Namespace
		ctx.Params["id"] = to.principal.Id

		// WHEN deleting relations
		err = ctrl.deleteRelationships(ctx)
		// THEN it should not fail
		require.NoError(t, err)
	}
}

func newTestPrincipalsController() (to *testObjects, ctrl *PrincipalsController, err error) {
	webServer := web.NewStubWebServer()
	if to, err = newTestObjects(); err != nil {
		return
	}
	ctrl = NewPrincipalsController(to.config, to.authService, webServer)
	return
}
