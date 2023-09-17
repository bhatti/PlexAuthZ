package repository

import (
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"time"
)

// NewRoleRepository creates repository for persisting roles
func NewRoleRepository(
	store DataStore,
) (Repository[types.Role], error) {
	return NewBaseRepository[types.Role](store,
		"Role",
		"",
		time.Duration(0),
		func() *types.Role {
			return &types.Role{}
		})
}
