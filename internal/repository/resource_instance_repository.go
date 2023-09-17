package repository

import (
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"time"
)

type resourceInstanceRepository struct {
	store      DataStore
	expiration time.Duration
}

// NewResourceInstanceRepository creates repository for persisting resources instance
func NewResourceInstanceRepository(
	store DataStore,
	expiration time.Duration,
) (ResourceInstanceRepositoryFactory, error) {
	return &resourceInstanceRepository{
		store:      store,
		expiration: expiration,
	}, nil
}

func (r *resourceInstanceRepository) CreateResourceInstanceRepository(
	resourceID string,
) (Repository[types.ResourceInstance], error) {
	return NewBaseRepository[types.ResourceInstance](
		r.store,
		"ResourceInstance",
		resourceID,
		r.expiration,
		func() *types.ResourceInstance {
			return &types.ResourceInstance{}
		})

}
