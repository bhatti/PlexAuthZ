package server

import (
	"crypto/tls"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Clients for GRPC server
type Clients struct {
	AuthClient          services.AuthZServiceClient
	GroupsClient        services.GroupsServiceClient
	OrganizationsClient services.OrganizationsServiceClient
	PermissionsClient   services.PermissionsServiceClient
	PrincipalsClient    services.PrincipalsServiceClient
	RelationshipsClient services.RelationshipsServiceClient
	ResourcesClient     services.ResourcesServiceClient
	RolesClient         services.RolesServiceClient
	ClientType          domain.ClientType
}

// NewClients constructor
func NewClients(caFile string, certFile string, keyFile string, addr string) (
	conn *grpc.ClientConn,
	clients Clients,
	err error) {
	var tlsConfig *tls.Config
	tlsConfig, err = domain.TLSConfig{
		CAFile:   caFile,
		CertFile: certFile,
		KeyFile:  keyFile,
		Server:   false,
	}.SetupTLS()
	if err != nil {
		return
	}
	tlsCreds := credentials.NewTLS(tlsConfig)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(tlsCreds)}
	conn, err = grpc.Dial(addr, opts...)
	if err != nil {
		return
	}

	clients.AuthClient = services.NewAuthZServiceClient(conn)
	clients.GroupsClient = services.NewGroupsServiceClient(conn)
	clients.OrganizationsClient = services.NewOrganizationsServiceClient(conn)
	clients.PermissionsClient = services.NewPermissionsServiceClient(conn)
	clients.PrincipalsClient = services.NewPrincipalsServiceClient(conn)
	clients.RelationshipsClient = services.NewRelationshipsServiceClient(conn)
	clients.ResourcesClient = services.NewResourcesServiceClient(conn)
	clients.RolesClient = services.NewRolesServiceClient(conn)
	return
}

// ClientTypesMap builds map with different client types
func ClientTypesMap(cfg *domain.Config) (map[domain.ClientType]domain.TLSConfig, error) {
	tlsC, err := cfg.TLSClient()
	if err != nil {
		return nil, err
	}
	tlsNobody, err := cfg.TLSNobodyClient()
	if err != nil {
		return nil, err
	}
	tlsRoot, err := cfg.TLSRootClient()
	if err != nil {
		return nil, err
	}

	return map[domain.ClientType]domain.TLSConfig{
		domain.DefaultClientType: tlsC,
		domain.NobodyClientType:  tlsNobody,
		domain.RootClientType:    tlsRoot,
	}, nil
}
