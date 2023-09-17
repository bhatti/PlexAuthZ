package authz

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/internal/domain"
)

// NullAuthorizer for implementing no authorization.
type NullAuthorizer struct {
}

// Check null implementation.
func (n NullAuthorizer) Check(
	_ context.Context,
	_ *services.CheckConstraintsRequest) (*services.CheckConstraintsResponse, error) {
	return &services.CheckConstraintsResponse{}, nil
}

// Authorize returns empty response.
func (n NullAuthorizer) Authorize(
	_ context.Context,
	_ *services.AuthRequest,
) (*services.AuthResponse, error) {
	return &services.AuthResponse{}, nil
}

// NoAuthorizer rejects all authorization requests.
type NoAuthorizer struct {
}

// Check without implementation.
func (n NoAuthorizer) Check(
	_ context.Context,
	_ *services.CheckConstraintsRequest) (*services.CheckConstraintsResponse, error) {
	return nil, domain.NewInternalError("check method is not implemented for NullAuthorizer", domain.InternalCode)
}

// Authorize without any enforcement.
func (n NoAuthorizer) Authorize(
	_ context.Context,
	_ *services.AuthRequest,
) (*services.AuthResponse, error) {
	return nil, domain.NewAuthError(fmt.Sprintf("authz error"))
}
