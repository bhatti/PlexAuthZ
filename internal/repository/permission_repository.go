package repository

import (
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"time"
)

// NewPermissionRepository creates repository for persisting permission
func NewPermissionRepository(
	store DataStore,
) (Repository[types.Permission], error) {
	return NewBaseRepository[types.Permission](
		store,
		"Permission",
		"",
		time.Duration(0),
		func() *types.Permission {
			return &types.Permission{}
		})
}
