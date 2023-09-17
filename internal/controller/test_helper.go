package controller

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/metrics"
	"github.com/bhatti/PlexAuthZ/internal/service"
	"github.com/bhatti/PlexAuthZ/internal/service/db"
	"github.com/bhatti/PlexAuthZ/internal/web"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// SetupWebServerForTesting helper
func SetupWebServerForTesting(
	t *testing.T,
	cfg *domain.Config,
	fn func(config *domain.Config)) (
	client web.HTTPClient,
	teardown func()) {
	t.Helper()

	webServer := web.NewDefaultWebServer(cfg)

	client = web.NewHTTPClient(cfg)
	serverAuthService, _, err := db.CreateDatabaseAuthService(cfg, metrics.New())
	require.NoError(t, err)

	err = StartControllers(cfg, serverAuthService, webServer)
	require.NoError(t, err)

	if fn != nil {
		fn(cfg)
	}

	go func() {
		_ = webServer.Start(cfg.HttpListenPort)
	}()
	time.Sleep(time.Millisecond * 500)

	return client, func() {
		_ = webServer.Stop()
	}
}

type testObjects struct {
	ctx         context.Context
	config      *domain.Config
	authService service.AuthAdminService
	principal   *types.Principal
	org         *types.Organization
	resource    *types.Resource
	permission  *types.Permission
	role        *types.Role
	group       *types.Group
	relation    *types.Relationship
}

func newTestObjects() (to *testObjects, err error) {
	to = &testObjects{ctx: context.TODO()}
	to.config, err = domain.NewConfig("")
	if to.authService, _, err = db.CreateDatabaseAuthService(to.config, metrics.New()); err != nil {
		return
	}
	if err = to.createTestOrg(); err != nil {
		return
	}
	if err = to.createTestPrincipal(); err != nil {
		return
	}
	if err = to.addTestRole(); err != nil {
		return
	}
	if err = to.addTestGroup(); err != nil {
		return
	}
	if err = to.createTestResource(); err != nil {
		return
	}
	if err = to.addTestPermission(); err != nil {
		return
	}
	if err = to.createTestRelationship(); err != nil {
		return
	}
	return
}

func (to *testObjects) createTestOrg() (err error) {
	to.org, err = to.authService.CreateOrganization(to.ctx, &types.Organization{
		Name:       "test-org",
		Namespaces: []string{"admin", "finance"},
	})
	return
}

func (to *testObjects) createTestPrincipal() (err error) {
	to.principal, err = to.authService.CreatePrincipal(to.ctx, &types.Principal{
		OrganizationId: to.org.Id,
		Namespaces:     to.org.Namespaces,
		Username:       "john",
	})
	return
}

func (to *testObjects) addTestPermission() (err error) {
	if err = to.createTestPermission(); err != nil {
		return err
	}
	return to.authService.AddPermissionsToPrincipal(
		to.ctx, to.org.Id, to.org.Namespaces[0], to.principal.Id, to.permission.Id)
}

func (to *testObjects) addTestRole() (err error) {
	if err = to.createTestRole(); err != nil {
		return err
	}
	return to.authService.AddRolesToPrincipal(
		to.ctx, to.org.Id, to.org.Namespaces[0], to.principal.Id, to.role.Id)
}

func (to *testObjects) createTestRole() (err error) {
	to.role, err = to.authService.CreateRole(to.ctx, to.org.Id, &types.Role{
		Namespace: to.org.Namespaces[0],
		Name:      "role",
	})
	return
}

func (to *testObjects) addTestGroup() (err error) {
	if err = to.createTestGroup(); err != nil {
		return err
	}
	return to.authService.AddGroupsToPrincipal(
		to.ctx, to.org.Id, to.org.Namespaces[0], to.principal.Id, to.group.Id)
}

func (to *testObjects) createTestGroup() (err error) {
	to.group, err = to.authService.CreateGroup(to.ctx, to.org.Id, &types.Group{
		Namespace: to.org.Namespaces[0],
		Name:      "group",
	})
	return
}

func (to *testObjects) createTestPermission() (err error) {
	to.permission, err = to.authService.CreatePermission(to.ctx, to.org.Id, &types.Permission{
		Namespace:  to.org.Namespaces[0],
		Actions:    []string{"read", "write"},
		ResourceId: to.resource.Id,
	})
	return
}

func (to *testObjects) createTestResource() (err error) {
	to.resource, err = to.authService.CreateResource(to.ctx, to.org.Id, &types.Resource{
		Namespace:      to.org.Namespaces[0],
		Name:           "paper",
		AllowedActions: []string{"read", "write"},
		Capacity:       10,
	})
	return
}

func (to *testObjects) createTestRelationship() (err error) {
	to.relation, err = to.authService.CreateRelationship(to.ctx, to.org.Id, &types.Relationship{
		Namespace:   to.org.Namespaces[0],
		Relation:    "relation",
		PrincipalId: to.principal.Id,
		ResourceId:  to.resource.Id,
	})
	return
}
