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

// PrincipalsController - provides persistence for telemetry
type PrincipalsController struct {
	config    *domain.Config
	authAdminService service.AuthAdminService
}

// NewPrincipalsController instantiates controller for managing principals
func NewPrincipalsController(
	config *domain.Config,
	authAdminService service.AuthAdminService,
	webserver web.Server) *PrincipalsController {
	ctrl := &PrincipalsController{
		config:    config,
		authAdminService: authAdminService,
	}

	webserver.POST("/api/v1/:organization_id/principals", ctrl.create)
	webserver.PUT("/api/v1/:organization_id/principals/:id", ctrl.update)
	webserver.GET("/api/v1/:organization_id/:namespace/principals/:id", ctrl.get)
	webserver.GET("/api/v1/:organization_id/principals", ctrl.query)
	webserver.DELETE("/api/v1/:organization_id/principals/:id", ctrl.delete)
	webserver.PUT("/api/v1/:organization_id/:namespace/principals/:id/groups/add", ctrl.addGroups)
	webserver.PUT("/api/v1/:organization_id/:namespace/principals/:id/groups/delete", ctrl.deleteGroups)
	webserver.PUT("/api/v1/:organization_id/:namespace/principals/:id/roles/add", ctrl.addRoles)
	webserver.PUT("/api/v1/:organization_id/:namespace/principals/:id/roles/delete", ctrl.deleteRoles)
	webserver.PUT("/api/v1/:organization_id/:namespace/principals/:id/permissions/add", ctrl.addPermissions)
	webserver.PUT("/api/v1/:organization_id/:namespace/principals/:id/permissions/delete", ctrl.deletePermissions)
	webserver.PUT("/api/v1/:organization_id/:namespace/principals/:id/relations/add", ctrl.addRelationships)
	webserver.PUT("/api/v1/:organization_id/:namespace/principals/:id/relations/delete", ctrl.deleteRelationships)
	return ctrl
}

// create handler
func (ctr *PrincipalsController) create(c web.APIContext) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	principal := &types.Principal{}
	if err = json.Unmarshal(b, principal); err != nil {
		return err
	}
	principal.OrganizationId = c.Param("organization_id")
	if principal, err = ctr.authAdminService.CreatePrincipal(
		context.Background(),
		principal); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.CreatePrincipalResponse{
		Id: principal.Id,
	})
}

// update handler
func (ctr *PrincipalsController) update(c web.APIContext) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	principal := &types.Principal{}
	if err = json.Unmarshal(b, principal); err != nil {
		return err
	}
	principal.OrganizationId = c.Param("organization_id")
	principal.Id = c.Param("id")
	if err = ctr.authAdminService.UpdatePrincipal(
		context.Background(),
		principal); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.UpdatePrincipalResponse{})
}

// get handler
func (ctr *PrincipalsController) get(c web.APIContext) (err error) {
	principal, err := ctr.authAdminService.GetPrincipalExt(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
	)
	if err != nil {
		return err
	}
	res := principal.ToGetPrincipalResponse()
	return c.JSON(http.StatusOK, res)
}

// query handler
func (ctr *PrincipalsController) query(c web.APIContext) (err error) {
	predicates, offset, limit := toPredicates(c, "id", "name")
	res, nextOffset, err := ctr.authAdminService.GetPrincipals(
		context.Background(),
		c.Param("organization_id"),
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
func (ctr *PrincipalsController) delete(c web.APIContext) (err error) {
	if err = ctr.authAdminService.DeletePrincipal(
		context.Background(),
		c.Param("organization_id"),
		c.Param("id"),
	); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.DeletePrincipalResponse{})
}

// addGroups handler
func (ctr *PrincipalsController) addGroups(c web.APIContext) (err error) {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	req := &services.AddGroupsToPrincipalRequest{}
	if err = json.Unmarshal(b, req); err != nil {
		return err
	}
	err = ctr.authAdminService.AddGroupsToPrincipal(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
		req.GroupIds...,
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.AddGroupsToPrincipalResponse{})
}

// deleteGroups handler
func (ctr *PrincipalsController) deleteGroups(c web.APIContext) (err error) {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	req := &services.DeleteGroupsToPrincipalRequest{}
	if err = json.Unmarshal(b, req); err != nil {
		return err
	}
	if err = ctr.authAdminService.DeleteGroupsToPrincipal(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
		req.GroupIds...,
	); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.DeleteGroupsToPrincipalResponse{})
}

// addRoles handler
func (ctr *PrincipalsController) addRoles(c web.APIContext) (err error) {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	req := &services.AddRolesToPrincipalRequest{}
	if err = json.Unmarshal(b, req); err != nil {
		return err
	}
	if err = ctr.authAdminService.AddRolesToPrincipal(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
		req.RoleIds...,
	); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.AddRolesToPrincipalResponse{})
}

// deleteRoles handler
func (ctr *PrincipalsController) deleteRoles(c web.APIContext) (err error) {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	req := &services.DeleteRolesToPrincipalRequest{}
	if err = json.Unmarshal(b, req); err != nil {
		return err
	}
	if err = ctr.authAdminService.DeleteRolesToPrincipal(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
		req.RoleIds...,
	); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.DeleteRolesToPrincipalResponse{})
}

// addPermissions handler
func (ctr *PrincipalsController) addPermissions(c web.APIContext) (err error) {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	req := &services.AddPermissionsToPrincipalRequest{}
	if err = json.Unmarshal(b, req); err != nil {
		return err
	}
	if err = ctr.authAdminService.AddPermissionsToPrincipal(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
		req.PermissionIds...,
	); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.AddPermissionsToPrincipalResponse{})
}

// deletePermissions handler
func (ctr *PrincipalsController) deletePermissions(c web.APIContext) (err error) {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	req := &services.DeletePermissionsToPrincipalRequest{}
	if err = json.Unmarshal(b, req); err != nil {
		return err
	}
	if err = ctr.authAdminService.DeletePermissionsToPrincipal(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
		req.PermissionIds...,
	); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.DeletePermissionsToPrincipalResponse{})
}

// addRelationships handler
func (ctr *PrincipalsController) addRelationships(c web.APIContext) (err error) {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	req := &services.AddRelationshipsToPrincipalRequest{}
	if err = json.Unmarshal(b, req); err != nil {
		return err
	}
	if err = ctr.authAdminService.AddRelationshipsToPrincipal(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
		req.RelationshipIds...,
	); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.AddRelationshipsToPrincipalResponse{})
}

// deleteRelationships handler
func (ctr *PrincipalsController) deleteRelationships(c web.APIContext) (err error) {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	req := &services.DeleteRelationshipsToPrincipalRequest{}
	if err = json.Unmarshal(b, req); err != nil {
		return err
	}
	if err = ctr.authAdminService.DeleteRelationshipsToPrincipal(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
		req.RelationshipIds...,
	); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.DeleteRelationshipsToPrincipalResponse{})
}
