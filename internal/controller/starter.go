package controller

import (
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/service"
	"github.com/bhatti/PlexAuthZ/internal/web"
)

// StartControllers registers controllers for REST APIs.
func StartControllers(
	config *domain.Config,
	authService service.AuthAdminService,
	webServer web.Server,
) error {
	// Start controllers
	if _, err := NewAuthController(
		config,
		authService,
		webServer); err != nil {
		return err
	}

	_ = NewOrganizationsController(
		config,
		authService,
		webServer)

	_ = NewPrincipalsController(
		config,
		authService,
		webServer)

	_ = NewGroupsController(
		config,
		authService,
		webServer)

	_ = NewPermissionsController(
		config,
		authService,
		webServer)

	_ = NewRelationshipsController(
		config,
		authService,
		webServer)

	_ = NewResourcesController(
		config,
		authService,
		webServer)

	_ = NewRolesController(
		config,
		authService,
		webServer)
	return nil
}
