package server

import (
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/metrics"
	"github.com/bhatti/PlexAuthZ/internal/service/db"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"testing"
)

// SetupGrpcServerForTesting helper
func SetupGrpcServerForTesting(
	t *testing.T,
	cfg *domain.Config,
	clientType domain.ClientType,
	fn func(config *domain.Config)) (
	clients Clients,
	teardown func()) {
	t.Helper()

	adapter := &GrpcAdapter{}
	err := adapter.listen("127.0.0.1:0")
	require.NoError(t, err)
	cfg.GrpcListenPort = adapter.Addr().String()

	var tlsC domain.TLSConfig
	if clientType == domain.RootClientType {
		tlsC, err = cfg.TLSRootClient()
		require.NoError(t, err)
	} else if clientType == domain.NobodyClientType {
		tlsC, err = cfg.TLSNobodyClient()
		require.NoError(t, err)
	} else {
		tlsC, err = cfg.TLSClient()
		require.NoError(t, err)
	}
	cc, clients, err := NewClients(
		tlsC.CAFile,
		tlsC.CertFile,
		tlsC.KeyFile,
		adapter.Addr().String(),
	)
	require.NoError(t, err)
	clients.ClientType = clientType

	if fn != nil {
		fn(cfg)
	}

	authService, _, err := db.CreateDatabaseAuthService(cfg, metrics.New())
	require.NoError(t, err)

	opts := make([]grpc.ServerOption, 0)
	err = adapter.startServer(cfg, authService, opts)
	require.NoError(t, err)

	go func() {
		_ = adapter.Serve()
	}()

	return clients, func() {
		_ = adapter.Close()
		_ = cc.Close()
	}
}
