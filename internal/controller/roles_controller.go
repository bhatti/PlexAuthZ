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

// RolesController - provides persistence for telemetry
type RolesController struct {
	config    *domain.Config
	authAdminService service.AuthAdminService
}

// NewRolesController instantiates controller for managing roles
func NewRolesController(
	config *domain.Config,
	authAdminService service.AuthAdminService,
	webserver web.Server) *RolesController {
	ctrl := &RolesController{
		config:    config,
		authAdminService: authAdminService,
	}

	webserver.POST("/api/v1/:organization_id/:namespace/roles", ctrl.create)
	webserver.PUT("/api/v1/:organization_id/:namespace/roles/:id", ctrl.update)
	webserver.GET("/api/v1/:organization_id/:namespace/roles", ctrl.query)
	webserver.DELETE("/api/v1/:organization_id/:namespace/roles/:id", ctrl.delete)
	webserver.PUT("/api/v1/:organization_id/:namespace/roles/:id/permissions/add", ctrl.addPermissions)
	webserver.PUT("/api/v1/:organization_id/:namespace/roles/:id/permissions/delete", ctrl.deletePermissions)
	return ctrl
}

// create handler
func (ctr *RolesController) create(c web.APIContext) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	role := &types.Role{}
	err = json.Unmarshal(b, role)
	if err != nil {
		return err
	}
	role.Namespace = c.Param("namespace")
	role, err = ctr.authAdminService.CreateRole(
		context.Background(),
		c.Param("organization_id"),
		role)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.CreateRoleResponse{
		Id: role.Id,
	})
}

// update handler
func (ctr *RolesController) update(c web.APIContext) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	role := &types.Role{}
	err = json.Unmarshal(b, role)
	if err != nil {
		return err
	}
	role.Id = c.Param("id")
	role.Namespace = c.Param("namespace")
	if err = ctr.authAdminService.UpdateRole(
		context.Background(),
		c.Param("organization_id"),
		role); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.UpdateRoleResponse{})
}

// query handler
func (ctr *RolesController) query(c web.APIContext) (err error) {
	predicates, offset, limit := toPredicates(c, "id", "name")
	res, nextOffset, err := ctr.authAdminService.GetRoles(
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
func (ctr *RolesController) delete(c web.APIContext) (err error) {
	err = ctr.authAdminService.DeleteRole(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.DeleteRoleResponse{})
}

// addPermissions handler
func (ctr *RolesController) addPermissions(c web.APIContext) (err error) {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	req := services.AddPermissionsToRoleRequest{}
	if err = json.Unmarshal(b, &req); err != nil {
		return err
	}
	if err = ctr.authAdminService.AddPermissionsToRole(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
		req.PermissionIds...,
	); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.AddPermissionsToRoleResponse{})
}

// deletePermissions handler
func (ctr *RolesController) deletePermissions(c web.APIContext) (err error) {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	req := services.DeletePermissionsToRoleRequest{}
	if err = json.Unmarshal(b, &req); err != nil {
		return err
	}
	if err = ctr.authAdminService.DeletePermissionsToRole(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
		req.PermissionIds...,
	); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.DeletePermissionsToRoleResponse{})
}
