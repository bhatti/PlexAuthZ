package grpc

import (
	"github.com/bhatti/PlexAuthZ/internal/server"
	"github.com/bhatti/PlexAuthZ/internal/service"
)

// authAdminServiceGrpc - manages persistence of AuthZ data.
type authAdminServiceGrpc struct {
	*OrganizationServiceGrpc // implementation for organization admin service
	*PrincipalServiceGrpc    // implementation for principal admin service
	*ResourceServiceGrpc     // implementation for resource admin service
	*PermissionServiceGrpc   // implementation for permission admin service
	*RoleServiceGrpc         // implementation for role admin service
	*GroupServiceGrpc        // implementation for group admin service
	*RelationshipServiceGrpc // implementation for relationships admin service
}

// NewAuthAdminServiceGrpc manages persistence of AuthZ data
func NewAuthAdminServiceGrpc(
	clients server.Clients,
) service.AuthAdminService {
	return &authAdminServiceGrpc{
		OrganizationServiceGrpc: NewOrganizationServiceGrpc(clients),
		PrincipalServiceGrpc:    NewPrincipalServiceGrpc(clients),
		ResourceServiceGrpc:     NewResourceServiceGrpc(clients),
		PermissionServiceGrpc:   NewPermissionServiceGrpc(clients),
		RoleServiceGrpc:         NewRoleServiceGrpc(clients),
		GroupServiceGrpc:        NewGroupServiceGrpc(clients),
		RelationshipServiceGrpc: NewRelationshipServiceGrpc(clients),
	}
}
