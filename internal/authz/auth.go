package authz

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
)

// AuthorizerKind defines enum for authorization implementations.
type AuthorizerKind string

const (
	// DefaultAuthorizerKind for authorization implementation.
	DefaultAuthorizerKind AuthorizerKind = "DEFAULT"

	// CasbinAuthorizerKind based on Casbin implementation.
	CasbinAuthorizerKind AuthorizerKind = "CASBIN"

	// NullAuthorizerKind based on NULL implementation.
	NullAuthorizerKind AuthorizerKind = "NULL"

	// NoneAuthorizerKind based on None implementation.
	NoneAuthorizerKind AuthorizerKind = "NONE"
)

// Authorizer interface for authorizing access requests and checking constraints.
type Authorizer interface {
	// Authorize checks permissions for access.
	Authorize(
		ctx context.Context,
		req *services.AuthRequest,
	) (*services.AuthResponse, error)

	// Check inspects constraints for access.
	Check(
		ctx context.Context,
		req *services.CheckConstraintsRequest,
	) (*services.CheckConstraintsResponse, error)
}
