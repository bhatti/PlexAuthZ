package client

import (
	"github.com/bhatti/PlexAuthZ/internal/authz"
	"github.com/bhatti/PlexAuthZ/internal/benchmark"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/factory"
	"github.com/bhatti/PlexAuthZ/internal/metrics"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func Test_BenchmarkDB(t *testing.T) {
	_ = os.Setenv("CONFIG_DIR", "../../config")
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	cfg.AuthServiceProvider = domain.DatabaseAuthServiceProvider
	providers := []domain.PersistenceProvider{domain.RedisPersistenceProvider, domain.DynamoDBPersistenceProvider}
	for _, provider := range providers {
		registry := metrics.New()
		cfg.PersistenceProvider = provider
		runDBTests(t,
			cfg,
			10,
			time.Second*1,
			registry,
			testPermissionsForDepositAccount,
			testPermissionsForProjects,
			testPermissionsWithScopeAndWithinGeoLocation,
			testForDetectingAmbiguousPermission,
			testForPermissionsForFeatureFlagsWithGeoFencing,
			testForPermissionsWithWildcardActions,
			testForPermissionsWithQuota,
			testReBACPermissions,
			testAttributeBasedPermissions,
			testForPermissionsWithAttributes,
			testRBAC,
			testCRUD,
			testForPermissionsWithIPAddresses,
			testPermissionsWithOwners,
		)
		for k, v := range registry.Summary() {
			t.Logf("benchmark summary %s -- %s = %v", providers, k, v)
		}
	}
}

func runDBTests(
	t *testing.T,
	cfg *domain.Config,
	tps int,
	duration time.Duration,
	registry *metrics.Registry,
	fns ...func(t *testing.T, authAdapter *AuthAdapter),
) {

	cfg.PersistenceProvider = domain.RedisPersistenceProvider
	cfg.AuthServiceProvider = domain.DatabaseAuthServiceProvider
	authSvc, cc, err := factory.CreateAuthAdminService(cfg, registry, domain.RootClientType, "")
	require.NoError(t, err)

	authorizer, err := authz.CreateAuthorizer(authz.DefaultAuthorizerKind, cfg, authSvc)
	require.NoError(t, err)
	authAdapter := New(authorizer, authSvc)

	fn := func() error {
		for _, f := range fns {
			f(t, authAdapter)
		}
		return nil
	}

	res := benchmark.Benchmark(fn, benchmark.Request{
		TPS:         tps,
		Duration:    duration,
		Percentiles: []float64{95.0, 99.0},
	})
	logrus.WithFields(
		logrus.Fields{
			"Component": "ClientTests",
			"Results":   res.String(),
		}).Debugf("benchmark results")
	_ = cc.Close()
}
