package controller

import (
	"context"
	"encoding/json"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/internal/authz"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/service"
	"github.com/bhatti/PlexAuthZ/internal/web"
	"io"
	"net/http"
	"time"
)

// AuthController - authorization controller
type AuthController struct {
	config      *domain.Config
	authService service.AuthAdminService
	authorizer  authz.Authorizer
}

// NewAuthController instantiates controller for managing groups
func NewAuthController(
	config *domain.Config,
	authService service.AuthAdminService,
	webserver web.Server) (*AuthController, error) {
	authorizer, err := authz.CreateAuthorizer(authz.DefaultAuthorizerKind, config, authService)
	if err != nil {
		return nil, err
	}
	ctrl := &AuthController{
		config:      config,
		authService: authService,
		authorizer:  authorizer,
	}
	webserver.POST("/api/v1/:organization_id/:namespace/:principal_id/auth", ctrl.auth)
	webserver.POST("/api/v1/:organization_id/:namespace/:principal_id/auth/constraints", ctrl.check)
	webserver.PUT("/api/v1/:organization_id/:namespace/resources/:id/allocate/:principal_id", ctrl.allocate)
	webserver.PUT("/api/v1/:organization_id/:namespace/resources/:id/deallocate/:principal_id", ctrl.deallocate)
	return ctrl, nil
}

// auth handler
func (ctr *AuthController) auth(c web.APIContext) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	req := &services.AuthRequest{}
	if err = json.Unmarshal(b, req); err != nil {
		return err
	}
	req.OrganizationId = c.Param("organization_id")
	req.Namespace = c.Param("namespace")
	req.PrincipalId = c.Param("principal_id")

	res, err := ctr.authorizer.Authorize(
		context.Background(),
		req)

	if err != nil {
		return c.String(domain.ErrorToHTTPStatus(err), err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

// auth handler
func (ctr *AuthController) check(c web.APIContext) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	req := &services.CheckConstraintsRequest{}
	if err = json.Unmarshal(b, req); err != nil {
		return err
	}
	req.OrganizationId = c.Param("organization_id")
	req.Namespace = c.Param("namespace")
	req.PrincipalId = c.Param("principal_id")

	res, err := ctr.authorizer.Check(
		context.Background(),
		req)

	if err != nil {
		return c.String(domain.ErrorToHTTPStatus(err), err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

// allocate handler
func (ctr *AuthController) allocate(c web.APIContext) (err error) {
	req := &services.AllocateResourceRequest{}
	if b, err := io.ReadAll(c.Request().Body); err == nil {
		_ = json.Unmarshal(b, req)
	}
	var expiry time.Duration
	if req.Expiry != nil {
		expiry = req.Expiry.AsDuration()
	}

	err = ctr.authService.AllocateResourceInstance(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
		c.Param("principal_id"),
		req.Constraints,
		expiry,
		req.Context,
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.AllocateResourceResponse{})
}

// deallocate handler
func (ctr *AuthController) deallocate(c web.APIContext) (err error) {
	err = ctr.authService.DeallocateResourceInstance(
		context.Background(),
		c.Param("organization_id"),
		c.Param("namespace"),
		c.Param("id"),
		c.Param("principal_id"),
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &services.DeallocateResourceResponse{})
}
