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

// RelationshipsController - provides persistence for telemetry
type RelationshipsController struct {
	config    *domain.Config
	authAdminService service.AuthAdminService
}

// NewRelationshipsController instantiates controller for managing relations
func NewRelationshipsController(
	config *domain.Config,
	authAdminService service.AuthAdminService,
	webserver web.Server) *RelationshipsController {
	ctrl := &RelationshipsController{
		config:    config,
		authAdminService: authAdminService,
	}

	webserver.POST("/api/v1/:organization_id/:namespace/relations", ctrl.create)
	webserver.PUT("/api/v1/:organization_id/:namespace/relations/:id", ctrl.update)
	webserver.GET("/api/v1/:organization_id/:namespace/relations", ctrl.query)
	webserver.DELETE("/api/v1/:organization_id/:namespace/relations/:id", ctrl.delete)
	return ctrl
}

// create handler
func (ctr *RelationshipsController) create(c web.APIContext) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	relation := &types.Relationship{}
	err = json.Unmarshal(b, relation)
	if err != nil {
		return err
	}
	relation.Namespace = c.Param("namespace")
	relation, err = ctr.authAdminService.CreateRelationship(
		context.Background(),
		c.Param("organization_id"),
		relation)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.CreateRelationshipResponse{
		Id: relation.Id,
	})
}

// update handler
func (ctr *RelationshipsController) update(c web.APIContext) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	relation := &types.Relationship{}
	err = json.Unmarshal(b, relation)
	if err != nil {
		return err
	}
	relation.Id = c.Param("id")
	relation.Namespace = c.Param("namespace")
	if err = ctr.authAdminService.UpdateRelationship(
		context.Background(),
		c.Param("organization_id"),
		relation); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.UpdateRelationshipResponse{})
}

// query handler
func (ctr *RelationshipsController) query(c web.APIContext) (err error) {
	predicates, offset, limit := toPredicates(c, "id", "relation")
	res, nextOffset, err := ctr.authAdminService.GetRelationships(
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
func (ctr *RelationshipsController) delete(c web.APIContext) (err error) {
	err = ctr.authAdminService.DeleteRelationship(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.DeleteRelationshipResponse{})
}
