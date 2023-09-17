package service

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
)

// PrincipalService - admin APIs for principals data
type PrincipalService interface {
	// CreatePrincipal - creates new principal object
	CreatePrincipal(
		ctx context.Context,
		principal *types.Principal) (*types.Principal, error)

	// UpdatePrincipal - updates principal
	UpdatePrincipal(
		ctx context.Context,
		principal *types.Principal) error

	// DeletePrincipal removes principal
	DeletePrincipal(
		ctx context.Context,
		organizationID string,
		id string) error

	// GetPrincipalExt - retrieves full principal
	GetPrincipalExt(
		ctx context.Context,
		organizationID string,
		namespace string,
		id string,
	) (ext *domain.PrincipalExt, err error)

	// GetPrincipal - retrieves principal
	GetPrincipal(
		ctx context.Context,
		organizationID string,
		id string,
	) (*types.Principal, error)

	// GetPrincipals - queries principals
	GetPrincipals(
		ctx context.Context,
		organizationID string,
		predicate map[string]string,
		offset string,
		limit int64) (res []*types.Principal, nextOffset string, err error)
	// AddGroupsToPrincipal helper
	AddGroupsToPrincipal(
		ctx context.Context,
		organizationID string,
		namespace string,
		principalID string,
		groupIDs ...string,
	) error

	// DeleteGroupsToPrincipal helper
	DeleteGroupsToPrincipal(
		ctx context.Context,
		organizationID string,
		namespace string,
		principalID string,
		groupIDs ...string,
	) error

	// AddRolesToPrincipal helper
	AddRolesToPrincipal(
		ctx context.Context,
		organizationID string,
		namespace string,
		principalID string,
		roleIDs ...string,
	) error

	// DeleteRolesToPrincipal helper
	DeleteRolesToPrincipal(
		ctx context.Context,
		organizationID string,
		namespace string,
		principalID string,
		roleIDs ...string,
	) error

	// AddPermissionsToPrincipal helper
	AddPermissionsToPrincipal(
		ctx context.Context,
		organizationID string,
		namespace string,
		principalID string,
		permissionIds ...string,
	) error

	// DeletePermissionsToPrincipal helper
	DeletePermissionsToPrincipal(
		ctx context.Context,
		organizationID string,
		namespace string,
		principalID string,
		permissionIds ...string,
	) error

	// AddRelationshipsToPrincipal helper
	AddRelationshipsToPrincipal(
		ctx context.Context,
		organizationID string,
		namespace string,
		principalID string,
		relationshipIds ...string,
	) error

	// DeleteRelationshipsToPrincipal helper
	DeleteRelationshipsToPrincipal(
		ctx context.Context,
		organizationID string,
		namespace string,
		principalID string,
		relationshipIds ...string,
	) error
}
