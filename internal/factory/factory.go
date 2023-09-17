package factory

import (
	"bytes"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/metrics"
	"github.com/bhatti/PlexAuthZ/internal/server"
	"github.com/bhatti/PlexAuthZ/internal/service"
	"github.com/bhatti/PlexAuthZ/internal/service/db"
	"github.com/bhatti/PlexAuthZ/internal/service/grpc"
	"github.com/bhatti/PlexAuthZ/internal/service/http"
	"github.com/bhatti/PlexAuthZ/internal/web"
	"io"
)

// CreateAuthAdminService factory
func CreateAuthAdminService(
	cfg *domain.Config,
	registry *metrics.Registry,
	clientType domain.ClientType,
	addr string,
) (authService service.AuthAdminService, cc domain.Closeable, err error) {
	if cfg.AuthServiceProvider == domain.HttpAuthServiceProvider {
		client := web.NewHTTPClient(cfg)
		return http.NewAuthAdminServiceHTTP(client, addr), io.NopCloser(bytes.NewReader([]byte{})), nil
	} else if cfg.AuthServiceProvider == domain.GrpcAuthServiceProvider {
		var tlsC domain.TLSConfig
		if clientType == domain.RootClientType {
			tlsC, err = cfg.TLSRootClient()
			if err != nil {
				return nil, nil, err
			}
		} else if clientType == domain.NobodyClientType {
			tlsC, err = cfg.TLSNobodyClient()
			if err != nil {
				return nil, nil, err
			}
		} else {
			tlsC, err = cfg.TLSClient()
			if err != nil {
				return nil, nil, err
			}
		}
		cc, clients, err := server.NewClients(
			tlsC.CAFile,
			tlsC.CertFile,
			tlsC.KeyFile,
			addr,
		)
		if err != nil {
			return nil, nil, err
		}
		authService := grpc.NewAuthAdminServiceGrpc(clients)
		return authService, cc, nil
	} else {
		return db.CreateDatabaseAuthService(cfg, registry)
	}
}
