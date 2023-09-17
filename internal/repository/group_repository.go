package repository

import (
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"time"
)

// NewGroupRepository creates repository for persisting groups
func NewGroupRepository(
	store DataStore,
) (Repository[types.Group], error) {
	return NewBaseRepository[types.Group](
		store,
		"Group",
		"",
		time.Duration(0),
		func() *types.Group {
			return &types.Group{}
		})
}
