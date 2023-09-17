package repository

import (
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"time"
)

// NewPrincipalRepository creates repository for persisting principals
func NewPrincipalRepository(
	store DataStore,
) (Repository[types.Principal], error) {
	return NewBaseRepository[types.Principal](store,
		"Principal",
		"",
		time.Duration(0),
		func() *types.Principal {
			return &types.Principal{}
		})
}
