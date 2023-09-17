package authz

import (
	"fmt"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/service"
)

// CreateAuthorizer factory
func CreateAuthorizer(
	kind AuthorizerKind,
	config *domain.Config,
	authService service.AuthAdminService,
) (Authorizer, error) {
	if kind == DefaultAuthorizerKind {
		return NewDefaultAuthorizer(authService), nil
	} else if kind == CasbinAuthorizerKind {
		return NewGrpcAuth(config)
	} else if kind == NullAuthorizerKind {
		return NullAuthorizer{}, nil
	} else if kind == NoneAuthorizerKind {
		return NoAuthorizer{}, nil
	} else {
		return nil, fmt.Errorf("unknown kind %s", kind)
	}
}
