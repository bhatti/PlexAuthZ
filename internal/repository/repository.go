package repository

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"time"
)

type Repository[T any] interface {
	// GetByIDs - finds objects matching ids.
	GetByIDs(
		ctx context.Context,
		organizationID string,
		namespace string,
		ids ...string,
	) (map[string]*T, error)

	// GetByID - finds an object by id.
	GetByID(
		ctx context.Context,
		organizationID string,
		namespace string,
		id string,
	) (*T, error)

	// Create - creates a new object.
	Create(
		ctx context.Context,
		organizationID string,
		namespace string,
		id string,
		obj *T,
		expiration time.Duration,
	) error

	// Update - saves existing object.
	Update(
		ctx context.Context,
		organizationID string,
		namespace string,
		id string,
		version int64,
		obj *T,
		expiration time.Duration,
	) error

	// Query - queries objects
	Query(
		ctx context.Context,
		organizationID string,
		namespace string,
		predicate map[string]string,
		offset string,
		limit int64,
	) (res []*T, nextOffset string, err error)

	// Delete - remove the object
	Delete(
		ctx context.Context,
		organizationID string,
		namespace string,
		id string,
	) error

	// Size of table
	Size(
		ctx context.Context,
		organizationID string,
		namespace string,
	) (int64, error)
}

// ResourceInstanceRepositoryFactory helper
type ResourceInstanceRepositoryFactory interface {
	CreateResourceInstanceRepository(resourceID string) (Repository[types.ResourceInstance], error)
}
