package authz

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/service"
)

// DefaultAuthorizer for defining authorization rules.
type DefaultAuthorizer struct {
	authAdminService service.AuthAdminService
}

// NewDefaultAuthorizer constructor
func NewDefaultAuthorizer(
	authAdminService service.AuthAdminService,
) Authorizer {
	return &DefaultAuthorizer{
		authAdminService: authAdminService,
	}
}

// Authorize checks access for principal, action and resource.
func (a *DefaultAuthorizer) Authorize(
	ctx context.Context,
	req *services.AuthRequest,
) (*services.AuthResponse, error) {
	principal, err := a.authAdminService.GetPrincipalExt(
		ctx,
		req.OrganizationId,
		req.Namespace,
		req.PrincipalId)
	if err != nil {
		return nil, err
	}
	res, err := principal.CheckPermission(
		req,
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Check ensures constraints matches for the principal.
func (a *DefaultAuthorizer) Check(
	ctx context.Context,
	req *services.CheckConstraintsRequest,
) (*services.CheckConstraintsResponse, error) {
	if req.Constraints == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("constraints is not defined"))
	}
	principal, err := a.authAdminService.GetPrincipalExt(
		ctx,
		req.OrganizationId,
		req.Namespace,
		req.PrincipalId)
	if err != nil {
		return nil, err
	}
	matched, output, err := principal.CheckConstraints(
		&services.AuthRequest{
			OrganizationId: req.OrganizationId,
			Namespace:      req.Namespace,
			PrincipalId:    req.PrincipalId,
			Context:        req.Context,
		},
		&types.Resource{},
		req.Constraints,
	)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, domain.NewAuthError(
			fmt.Sprintf("constraints '%s' not matched with context %v", req.Constraints, req.Context))

	}
	return &services.CheckConstraintsResponse{
		Matched: matched,
		Output:  output,
	}, nil
}
