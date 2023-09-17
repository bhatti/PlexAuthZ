package http

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/controller"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/service"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_HTTPBasedAuthService(t *testing.T) {
	runTests(t,
		testCRUDGroups,
		testCRUDOrganizations,
		testCRUDPermissions,
		testCRUDPrincipals,
		testCRUDPrincipalsWithPermissions,
		testCRUDRelationships,
		testCRUDResources,
		testCRUDResourcesWithInstances,
		testCRUDRoles,
		testCRUDRolesWithPermissions,
	)
}

func runTests(
	t *testing.T,
	fns ...func(ctx context.Context, t *testing.T, authService service.AuthAdminService, org *types.Organization),
) {
	_ = os.Setenv("CONFIG_DIR", "../../../config")
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	client, teardown := controller.SetupWebServerForTesting(t, cfg, nil)
	defer teardown()

	clientAuthService := NewAuthAdminServiceHTTP(client, "http://"+cfg.HttpListenPort)
	require.NoError(t, err)

	ctx := context.TODO()
	org, err := domain.NewOrganizationBuilder().
		WithName("test-org1").
		WithNamespaces("admin", "finance", "engineering").Build()
	require.NoError(t, err)

	for _, fn := range fns {
		org, err = clientAuthService.CreateOrganization(ctx, org)
		require.NoError(t, err)
		fn(ctx, t, clientAuthService, org)
	}
}
