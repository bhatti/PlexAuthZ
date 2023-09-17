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

// PermissionsController - provides persistence for telemetry
type PermissionsController struct {
	config    *domain.Config
	authAdminService service.AuthAdminService
}

// NewPermissionsController instantiates controller for managing permissions
func NewPermissionsController(
	config *domain.Config,
	authAdminService service.AuthAdminService,
	webserver web.Server) *PermissionsController {
	ctrl := &PermissionsController{
		config:    config,
		authAdminService: authAdminService,
	}

	webserver.POST("/api/v1/:organization_id/:namespace/permissions", ctrl.create)
	webserver.PUT("/api/v1/:organization_id/:namespace/permissions/:id", ctrl.update)
	webserver.GET("/api/v1/:organization_id/:namespace/permissions", ctrl.query)
	webserver.DELETE("/api/v1/:organization_id/:namespace/permissions/:id", ctrl.delete)
	return ctrl
}

// create handler
func (ctr *PermissionsController) create(c web.APIContext) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	permission := &types.Permission{}
	if err = json.Unmarshal(b, permission); err != nil {
		return err
	}
	permission.Namespace = c.Param("namespace")
	if permission, err = ctr.authAdminService.CreatePermission(
		context.Background(),
		c.Param("organization_id"),
		permission); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.CreatePermissionResponse{
		Id: permission.Id,
	})
}

// update handler
func (ctr *PermissionsController) update(c web.APIContext) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	permission := &types.Permission{}
	if err = json.Unmarshal(b, permission); err != nil {
		return err
	}
	permission.Id = c.Param("id")
	permission.Namespace = c.Param("namespace")
	if err = ctr.authAdminService.UpdatePermission(
		context.Background(),
		c.Param("organization_id"),
		permission); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.UpdatePermissionResponse{})
}

// query handler
func (ctr *PermissionsController) query(c web.APIContext) (err error) {
	predicates, offset, limit := toPredicates(c, "id", "namespace", "scope", "resource_id")
	res, nextOffset, err := ctr.authAdminService.GetPermissions(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
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
func (ctr *PermissionsController) delete(c web.APIContext) (err error) {
	if err = ctr.authAdminService.DeletePermission(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
	); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.DeletePermissionResponse{})
}
