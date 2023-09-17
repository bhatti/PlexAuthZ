package controller

import (
	"context"
	"encoding/json"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/service"
	"github.com/bhatti/PlexAuthZ/internal/web"
	"io"
	"net/http"
)

// OrganizationsController - provides persistence for telemetry
type OrganizationsController struct {
	config    *domain.Config
	authAdminService service.AuthAdminService
}

// NewOrganizationsController instantiates controller for managing organizations
func NewOrganizationsController(
	config *domain.Config,
	authAdminService service.AuthAdminService,
	webserver web.Server) *OrganizationsController {
	ctrl := &OrganizationsController{
		config:    config,
		authAdminService: authAdminService,
	}

	webserver.POST("/api/v1/organizations", ctrl.create)
	webserver.PUT("/api/v1/organizations/:id", ctrl.update)
	webserver.GET("/api/v1/organizations/:id", ctrl.get)
	webserver.GET("/api/v1/organizations", ctrl.query)
	webserver.DELETE("/api/v1/organizations/:id", ctrl.delete)
	return ctrl
}

// create handler
func (ctr *OrganizationsController) create(c web.APIContext) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	organization := &types.Organization{}
	if err = json.Unmarshal(b, organization); err != nil {
		return err
	}
	if organization, err = ctr.authAdminService.CreateOrganization(
		context.Background(),
		organization); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.CreateOrganizationResponse{
		Id: organization.Id,
	})
}

// update handler
func (ctr *OrganizationsController) update(c web.APIContext) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	organization := &types.Organization{}
	if err = json.Unmarshal(b, organization); err != nil {
		return err
	}
	organization.Id = c.Param("id")
	if err = ctr.authAdminService.UpdateOrganization(
		context.Background(),
		organization); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.UpdateOrganizationResponse{})
}

// get handler
func (ctr *OrganizationsController) get(c web.APIContext) (err error) {
	res, err := ctr.authAdminService.GetOrganization(
		context.Background(),
		c.Param("id"),
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

// query handler
func (ctr *OrganizationsController) query(c web.APIContext) (err error) {
	predicates, offset, limit := toPredicates(c, "id", "name")
	res, nextOffset, err := ctr.authAdminService.GetOrganizations(
		context.Background(),
		predicates,
		offset,
		limit,
	)
	if err != nil {
		return err
	}
	c.Response().Header().Set(domain.NextOffsetHeader, nextOffset)
	return c.JSON(http.StatusOK, res)
}

// delete handler
func (ctr *OrganizationsController) delete(c web.APIContext) (err error) {
	err = ctr.authAdminService.DeleteOrganization(
		context.Background(),
		c.Param("id"),
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.DeleteOrganizationResponse{})
}
