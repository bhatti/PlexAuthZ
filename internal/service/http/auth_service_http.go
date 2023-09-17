package http

import (
	"github.com/bhatti/PlexAuthZ/internal/service"
	"github.com/bhatti/PlexAuthZ/internal/web"
)

// authAdminServiceHTTP - manages persistence of AuthZ data.
type authAdminServiceHTTP struct {
	*OrganizationServiceHTTP // implementation for organization admin service
	*PrincipalServiceHTTP    // implementation for principal admin service
	*ResourceServiceHTTP     // implementation for resource admin service
	*PermissionServiceHTTP   // implementation for permission admin service
	*RoleServiceHTTP         // implementation for role admin service
	*GroupServiceHTTP        // implementation for group admin service
	*RelationshipServiceHTTP // implementation for relationships admin service
}

// NewAuthAdminServiceHTTP manages persistence of AuthZ data
func NewAuthAdminServiceHTTP(
	client web.HTTPClient,
	baseURL string,
) service.AuthAdminService {
	return &authAdminServiceHTTP{
		OrganizationServiceHTTP: NewOrganizationServiceHTTP(client, baseURL),
		PrincipalServiceHTTP:    NewPrincipalServiceHTTP(client, baseURL),
		ResourceServiceHTTP:     NewResourceServiceHTTP(client, baseURL),
		PermissionServiceHTTP:   NewPermissionServiceHTTP(client, baseURL),
		RoleServiceHTTP:         NewRoleServiceHTTP(client, baseURL),
		GroupServiceHTTP:        NewGroupServiceHTTP(client, baseURL),
		RelationshipServiceHTTP: NewRelationshipServiceHTTP(client, baseURL),
	}
}
