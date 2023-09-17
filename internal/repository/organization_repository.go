package repository

import (
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"time"
)

// NewOrganizationRepository creates repository for persisting organization
func NewOrganizationRepository(
	store DataStore,
) (Repository[types.Organization], error) {
	return NewBaseRepository[types.Organization](store,
		"Organization",
		"",
		time.Duration(0),
		func() *types.Organization {
			return &types.Organization{}
		})
}
