package service

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
)

// RoleService - admin APIs for roles data
type RoleService interface {
	// CreateRole - creates a new role
	CreateRole(
		ctx context.Context,
		organizationID string,
		role *types.Role) (*types.Role, error)

	// UpdateRole - updates an existing role
	UpdateRole(
		ctx context.Context,
		organizationID string,
		role *types.Role) error

	// DeleteRole removes role
	DeleteRole(
		ctx context.Context,
		organizationID string,
		namespace string,
		id string) error

	// GetRole - finds role
	GetRole(
		ctx context.Context,
		organizationID string,
		namespace string,
		id string,
	) (*types.Role, error)

	// GetRoles - queries roles
	GetRoles(
		ctx context.Context,
		organizationID string,
		namespace string,
		predicate map[string]string,
		offset string,
		limit int64) (res []*types.Role, nextOffset string, err error)

	// AddPermissionsToRole helper
	AddPermissionsToRole(
		ctx context.Context,
		organizationID string,
		namespace string,
		roleID string,
		permissionIds ...string,
	) error

	// DeletePermissionsToRole helper
	DeletePermissionsToRole(
		ctx context.Context,
		organizationID string,
		namespace string,
		roleID string,
		permissionIds ...string,
	) error
}
