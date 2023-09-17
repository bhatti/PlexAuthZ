package service

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
)

// GroupService - admin APIs for groups
type GroupService interface {
	// CreateGroup - creates a new group
	CreateGroup(
		ctx context.Context,
		organizationID string,
		group *types.Group) (*types.Group, error)

	// UpdateGroup - updates an exising group
	UpdateGroup(
		ctx context.Context,
		organizationID string,
		group *types.Group) error

	// DeleteGroup removes group
	DeleteGroup(
		ctx context.Context,
		organizationID string,
		namespace string,
		id string) error

	// GetGroup - finds group
	GetGroup(
		ctx context.Context,
		organizationID string,
		namespace string,
		id string,
	) (*types.Group, error)

	// GetGroups - queries groups
	GetGroups(
		ctx context.Context,
		organizationID string,
		namespace string,
		predicate map[string]string,
		offset string,
		limit int64) (res []*types.Group, nextOffset string, err error)

	// AddRolesToGroup helper
	AddRolesToGroup(
		ctx context.Context,
		organizationID string,
		namespace string,
		groupID string,
		roleIDs ...string,
	) error

	// DeleteRolesToGroup helper
	DeleteRolesToGroup(
		ctx context.Context,
		organizationID string,
		namespace string,
		groupID string,
		roleIDs ...string,
	) error
}
