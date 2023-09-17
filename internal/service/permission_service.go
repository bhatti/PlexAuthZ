package service

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
)

// PermissionService - admin APIs for permissions data
type PermissionService interface {
	// CreatePermission - creates a new permission
	CreatePermission(
		ctx context.Context,
		organizationID string,
		permission *types.Permission) (*types.Permission, error)

	// UpdatePermission - updates an existing permission
	UpdatePermission(
		ctx context.Context,
		organizationID string,
		permission *types.Permission) error

	// DeletePermission removes permission
	DeletePermission(
		ctx context.Context,
		organizationID string,
		namespace string,
		id string) error

	// GetPermission - finds permission
	GetPermission(
		ctx context.Context,
		organizationID string,
		namespace string,
		id string,
	) (*types.Permission, error)

	// GetPermissions - queries permissions
	GetPermissions(
		ctx context.Context,
		organizationID string,
		namespace string,
		predicate map[string]string,
		offset string,
		limit int64) (res []*types.Permission, nextOffset string, err error)
}
