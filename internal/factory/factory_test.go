package factory

import (
	"github.com/bhatti/PlexAuthZ/internal/authz"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/metrics"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ShouldCreateAuthorizerWithHTTP(t *testing.T) {
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	cfg.PersistenceProvider = domain.RedisPersistenceProvider
	cfg.AuthServiceProvider = domain.HttpAuthServiceProvider
	authSvc, cc, err := CreateAuthAdminService(cfg, metrics.New(), domain.RootClientType, cfg.HttpListenPort)
	require.NoError(t, err)
	defer func() {
		_ = cc.Close()
	}()
	_, err = authz.CreateAuthorizer(authz.DefaultAuthorizerKind, cfg, authSvc)
	require.NoError(t, err)
}

func Test_ShouldCreateAuthorizerWithGRPC(t *testing.T) {
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	cfg.PersistenceProvider = domain.RedisPersistenceProvider
	cfg.AuthServiceProvider = domain.GrpcAuthServiceProvider
	authSvc, cc, err := CreateAuthAdminService(cfg, metrics.New(), domain.RootClientType, cfg.GrpcListenPort)
	require.NoError(t, err)
	defer func() {
		_ = cc.Close()
	}()
	_, err = authz.CreateAuthorizer(authz.CasbinAuthorizerKind, cfg, authSvc)
	require.Error(t, err)
}

func Test_ShouldCreateAuthorizerWithGRPCNull(t *testing.T) {
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	cfg.PersistenceProvider = domain.RedisPersistenceProvider
	cfg.AuthServiceProvider = domain.GrpcAuthServiceProvider
	authSvc, cc, err := CreateAuthAdminService(cfg, metrics.New(), domain.RootClientType, cfg.GrpcListenPort)
	require.NoError(t, err)
	defer func() {
		_ = cc.Close()
	}()
	_, err = authz.CreateAuthorizer(authz.NullAuthorizerKind, cfg, authSvc)
	require.NoError(t, err)
}

func Test_ShouldCreateAuthorizerWithGRPCNone(t *testing.T) {
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	cfg.PersistenceProvider = domain.RedisPersistenceProvider
	cfg.AuthServiceProvider = domain.GrpcAuthServiceProvider
	authSvc, cc, err := CreateAuthAdminService(cfg, metrics.New(), domain.RootClientType, cfg.GrpcListenPort)
	require.NoError(t, err)
	defer func() {
		_ = cc.Close()
	}()
	_, err = authz.CreateAuthorizer(authz.NoneAuthorizerKind, cfg, authSvc)
	require.NoError(t, err)
}

func Test_ShouldCreateAuthorizerWithGRPCUnknown(t *testing.T) {
	cfg, err := domain.NewConfig("")
	require.NoError(t, err)
	cfg.PersistenceProvider = domain.RedisPersistenceProvider
	cfg.AuthServiceProvider = domain.DatabaseAuthServiceProvider
	authSvc, cc, err := CreateAuthAdminService(cfg, metrics.New(), domain.RootClientType, cfg.GrpcListenPort)
	require.NoError(t, err)
	defer func() {
		_ = cc.Close()
	}()
	// should fail with unknown type
	_, err = authz.CreateAuthorizer(authz.AuthorizerKind("unknown"), cfg, authSvc)
	require.Error(t, err)
}
