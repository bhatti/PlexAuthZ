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
	"strconv"
)

// GroupsController - provides persistence for telemetry
type GroupsController struct {
	config    *domain.Config
	authAdminService service.AuthAdminService
}

// NewGroupsController instantiates controller for managing groups
func NewGroupsController(
	config *domain.Config,
	authAdminService service.AuthAdminService,
	webserver web.Server) *GroupsController {
	ctrl := &GroupsController{
		config:    config,
		authAdminService: authAdminService,
	}

	webserver.POST("/api/v1/:organization_id/:namespace/groups", ctrl.create)
	webserver.PUT("/api/v1/:organization_id/:namespace/groups/:id", ctrl.update)
	webserver.GET("/api/v1/:organization_id/:namespace/groups", ctrl.query)
	webserver.DELETE("/api/v1/:organization_id/:namespace/groups/:id", ctrl.delete)
	webserver.PUT("/api/v1/:organization_id/:namespace/groups/:id/roles/add", ctrl.addRoles)
	webserver.PUT("/api/v1/:organization_id/:namespace/groups/:id/roles/delete", ctrl.deleteRoles)

	return ctrl
}

// create handler
func (ctr *GroupsController) create(c web.APIContext) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	group := &types.Group{}
	if err = json.Unmarshal(b, group); err != nil {
		return err
	}
	group.Namespace = c.Param("namespace")
	if group, err = ctr.authAdminService.CreateGroup(
		context.Background(),
		c.Param("organization_id"),
		group); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.CreateGroupResponse{
		Id: group.Id,
	})
}

// update handler
func (ctr *GroupsController) update(c web.APIContext) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	group := &types.Group{}
	if err = json.Unmarshal(b, group); err != nil {
		return err
	}
	group.Id = c.Param("id")
	group.Namespace = c.Param("namespace")
	if err = ctr.authAdminService.UpdateGroup(
		context.Background(),
		c.Param("organization_id"),
		group); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.UpdateGroupResponse{})
}

// query handler
func (ctr *GroupsController) query(c web.APIContext) (err error) {
	predicates, offset, limit := toPredicates(c, "id", "name")
	res, nextOffset, err := ctr.authAdminService.GetGroups(
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
func (ctr *GroupsController) delete(c web.APIContext) (err error) {
	if err = ctr.authAdminService.DeleteGroup(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
	); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.DeleteGroupResponse{})
}

// addRoles handler
func (ctr *GroupsController) addRoles(c web.APIContext) (err error) {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	req := &services.AddRolesToGroupRequest{}
	if err = json.Unmarshal(b, req); err != nil {
		return err
	}
	if err = ctr.authAdminService.AddRolesToGroup(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
		req.RoleIds...,
	); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.AddRolesToGroupResponse{})
}

// deleteRoles handler
func (ctr *GroupsController) deleteRoles(c web.APIContext) (err error) {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	req := &services.DeleteRolesToGroupRequest{}
	if err = json.Unmarshal(b, req); err != nil {
		return err
	}
	if err = ctr.authAdminService.DeleteRolesToGroup(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
		req.RoleIds...,
	); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.DeleteRolesToGroupResponse{})
}

func toPredicates(c web.APIContext, keys ...string) (predicates map[string]string, offset string, limit int64) {
	predicates = make(map[string]string)
	for _, key := range keys {
		if c.QueryParams().Get(key) != "" {
			predicates[key] = c.QueryParams().Get(key)
		}
	}
	offset = c.QueryParams().Get("offset")
	if c.QueryParams().Get("limit") != "" {
		limit, _ = strconv.ParseInt(c.QueryParams().Get("limit"), 10, 64)
	}
	return
}
