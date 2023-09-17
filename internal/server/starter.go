package server

import (
	"fmt"
	api "github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/internal/authz"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/service"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	log "github.com/sirupsen/logrus"
	"github.com/twinj/uuid"
	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"os"
	"strings"
	"time"
)

const (
	objectWildcard = "*"
	updateAction   = "update"
	queryAction    = "query"
	deleteAction   = "delete"
	authAction     = "auth"
)

// GrpcAdapter for managing gRPC server.
type GrpcAdapter struct {
	id         string
	grpcServer *grpc.Server
	listener   net.Listener
}

// Addr for managing address of gRPC server.
func (a *GrpcAdapter) Addr() net.Addr {
	return a.listener.Addr()
}

// Close stops server.
func (a *GrpcAdapter) Close() (err error) {
	log.WithField("gRPCListen", a.Addr().String()).
		Infof("##################### stopping gRPC server %s #####################", a.id)
	if a.grpcServer != nil {
		a.grpcServer.Stop()
	}
	if a.listener != nil {
		err = a.listener.Close()
	}
	a.grpcServer = nil
	a.listener = nil
	return
}

// Serve starts serving requests.
func (a *GrpcAdapter) Serve() (err error) {
	if a.grpcServer == nil {
		return fmt.Errorf("grpcServer not initialized")
	}
	if a.listener == nil {
		return fmt.Errorf("listner not initialized")
	}
	log.WithField("gRPCListen", a.Addr().String()).
		Infof("##################### serving gRPC server %s #####################", a.id)
	return a.grpcServer.Serve(a.listener)
}

// Serve starts listening to port.
func (a *GrpcAdapter) listen(listenPort string) (err error) {
	a.id = uuid.NewV4().String()
	a.listener, err = net.Listen("tcp", listenPort)
	return err
}

func (a *GrpcAdapter) registerServers(
	authorizer authz.Authorizer,
	authService service.AuthAdminService) error {
	if srv, err := NewAuthServer(
		authService,
		authorizer,
	); err == nil {
		api.RegisterAuthZServiceServer(a.grpcServer, srv)
	} else {
		return err
	}

	if srv, err := NewGroupsServer(
		authService,
		authorizer,
	); err == nil {
		api.RegisterGroupsServiceServer(a.grpcServer, srv)
	} else {
		return err
	}

	if srv, err := NewOrganizationsServer(
		authService,
		authorizer,
	); err == nil {
		api.RegisterOrganizationsServiceServer(a.grpcServer, srv)
	} else {
		return err
	}

	if srv, err := NewPermissionsServer(
		authService,
		authorizer,
	); err == nil {
		api.RegisterPermissionsServiceServer(a.grpcServer, srv)
	} else {
		return err
	}

	if srv, err := NewPrincipalsServer(
		authService,
		authorizer,
	); err == nil {
		api.RegisterPrincipalsServiceServer(a.grpcServer, srv)
	} else {
		return err
	}

	if srv, err := NewRelationshipsServer(
		authService,
		authorizer,
	); err == nil {
		api.RegisterRelationshipsServiceServer(a.grpcServer, srv)
	} else {
		return err
	}

	if srv, err := NewResourcesServer(
		authService,
		authorizer,
	); err == nil {
		api.RegisterResourcesServiceServer(a.grpcServer, srv)
	} else {
		return err
	}

	if srv, err := NewRolesServer(
		authService,
		authorizer,
	); err == nil {
		api.RegisterRolesServiceServer(a.grpcServer, srv)
	} else {
		return err
	}
	return nil
}

// StartServers starts gRPC server.
func StartServers(
	config *domain.Config,
	authService service.AuthAdminService,
	grpcOpts ...grpc.ServerOption) (adapter *GrpcAdapter, err error) {
	adapter = &GrpcAdapter{}
	err = adapter.listen(config.GrpcListenPort)
	if err != nil {
		return nil, err
	}
	err = adapter.startServer(config, authService, grpcOpts)
	if err != nil {
		return nil, err
	}
	return adapter, nil
}

func (a *GrpcAdapter) startServer(
	config *domain.Config,
	authService service.AuthAdminService,
	grpcOpts []grpc.ServerOption,
) (err error) {
	var authorizer authz.Authorizer
	if config.GrpcSasl {
		serverTLSConfig, err := config.SetupTLSServer(a.Addr().String())
		if err != nil {
			log.WithField("Error", err).Fatalf("Could setup TLS for gRPC")
			return err
		}
		serverCreds := credentials.NewTLS(serverTLSConfig)
		grpcOpts = append(grpcOpts, grpc.Creds(serverCreds))

		authorizer, err = authz.CreateAuthorizer(authz.CasbinAuthorizerKind, config, authService)
		if err != nil {
			return err
		}
	} else {
		authorizer, err = authz.CreateAuthorizer(authz.NullAuthorizerKind, config, authService)
		if err != nil {
			return err
		}
	}

	logger := zap.L().Named("PlexAuthZ")
	zapOpts := []grpc_zap.Option{
		grpc_zap.WithDurationField(
			func(duration time.Duration) zapcore.Field {
				return zap.Int64("grpc.time_ns", duration.Nanoseconds())
			},
		),
	}

	halfSampler := trace.ProbabilitySampler(0.5)
	//trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	trace.ApplyConfig(trace.Config{DefaultSampler: func(p trace.SamplingParameters) trace.SamplingDecision {
		if strings.Contains(p.Name, "query") {
			return trace.SamplingDecision{Sample: true}
		}
		return halfSampler(p)
	}})

	err = view.Register(ocgrpc.DefaultServerViews...)
	if err != nil {
		return err
	}

	if config.GrpcSasl {
		grpcOpts = append(grpcOpts,
			grpc.StreamInterceptor(
				grpc_middleware.ChainStreamServer(
					grpc_ctxtags.StreamServerInterceptor(),
					grpc_zap.StreamServerInterceptor(
						logger, zapOpts...,
					),
					grpc_auth.StreamServerInterceptor(
						authz.Authenticate,
					),
				)), grpc.UnaryInterceptor(
				grpc_middleware.ChainUnaryServer(
					grpc_ctxtags.UnaryServerInterceptor(),
					grpc_zap.UnaryServerInterceptor(
						logger, zapOpts...,
					),
					grpc_auth.UnaryServerInterceptor(
						authz.Authenticate,
					),
				)),
			grpc.StatsHandler(&ocgrpc.ServerHandler{}),
		)
	}

	if config.Debug {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return err
		}
		zap.ReplaceGlobals(logger)

		metricsLogFile, err := os.CreateTemp("", "metrics-*.log")
		if err != nil {
			return err
		}
		tracesLogFile, err := os.CreateTemp("", "traces-*.log")
		if err != nil {
			return err
		}

		telemeterExporter, err := exporter.NewLogExporter(exporter.Options{
			MetricsLogFile:    metricsLogFile.Name(),
			TracesLogFile:     tracesLogFile.Name(),
			ReportingInterval: time.Second,
		})
		if err != nil {
			return err
		}

		err = telemeterExporter.Start()
		if err != nil {
			return err
		}
	}
	a.grpcServer = grpc.NewServer(grpcOpts...)

	if err = a.registerServers(authorizer, authService); err != nil {
		return err
	}

	return nil
}
