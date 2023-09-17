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

// ResourcesController - provides persistence for telemetry
type ResourcesController struct {
	config    *domain.Config
	authAdminService service.AuthAdminService
}

// NewResourcesController instantiates controller for managing resources
func NewResourcesController(
	config *domain.Config,
	authAdminService service.AuthAdminService,
	webserver web.Server) *ResourcesController {
	ctrl := &ResourcesController{
		config:    config,
		authAdminService: authAdminService,
	}

	webserver.POST("/api/v1/:organization_id/:namespace/resources", ctrl.create)
	webserver.PUT("/api/v1/:organization_id/:namespace/resources/:id", ctrl.update)
	webserver.GET("/api/v1/:organization_id/:namespace/resources", ctrl.query)
	webserver.DELETE("/api/v1/:organization_id/:namespace/resources/:id", ctrl.delete)
	webserver.GET("/api/v1/:organization_id/:namespace/resources/:id/instances", ctrl.queryAllocatedInstances)
	webserver.GET("/api/v1/:organization_id/:namespace/resources/:id/instance_count", ctrl.allocatedInstancesCount)
	return ctrl
}

// create handler
func (ctr *ResourcesController) create(c web.APIContext) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	resource := &types.Resource{}
	err = json.Unmarshal(b, resource)
	if err != nil {
		return err
	}
	resource.Namespace = c.Param("namespace")
	resource, err = ctr.authAdminService.CreateResource(
		context.Background(),
		c.Param("organization_id"),
		resource)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.CreateResourceResponse{
		Id: resource.Id,
	})
}

// update handler
func (ctr *ResourcesController) update(c web.APIContext) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	resource := &types.Resource{}
	err = json.Unmarshal(b, resource)
	if err != nil {
		return err
	}
	resource.Id = c.Param("id")
	resource.Namespace = c.Param("namespace")
	if err = ctr.authAdminService.UpdateResource(
		context.Background(),
		c.Param("organization_id"),
		resource); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.UpdateResourceResponse{})
}

// query handler
func (ctr *ResourcesController) query(c web.APIContext) (err error) {
	predicates, offset, limit := toPredicates(c, "id", "name")
	res, nextOffset, err := ctr.authAdminService.QueryResources(
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
func (ctr *ResourcesController) delete(c web.APIContext) (err error) {
	err = ctr.authAdminService.DeleteResource(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.DeleteResourceResponse{})
}

// queryAllocatedInstances handler
func (ctr *ResourcesController) queryAllocatedInstances(c web.APIContext) (err error) {
	predicates, offset, limit := toPredicates(c, "name") // no id as it will search for instance-id
	instances, nextOffset, err := ctr.authAdminService.QueryResourceInstances(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
		predicates,
		offset,
		limit,
	)
	if err != nil {
		return err
	}
	c.Response().Header().Set(domain.NextOffsetHeader, nextOffset)
	return c.JSON(http.StatusOK, instances)
}

// allocatedInstancesCount handler
func (ctr *ResourcesController) allocatedInstancesCount(c web.APIContext) (err error) {
	capacity, allocated, err := ctr.authAdminService.CountResourceInstances(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.CountResourceInstancesResponse{
		Capacity:  capacity,
		Allocated: allocated,
	})
}
