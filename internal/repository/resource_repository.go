package repository

import (
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"time"
)

// NewResourceRepository creates repository for persisting resources
func NewResourceRepository(
	store DataStore,
) (Repository[types.Resource], error) {
	return NewBaseRepository[types.Resource](store,
		"Resource",
		"",
		time.Duration(0),
		func() *types.Resource {
			return &types.Resource{}
		})
}
