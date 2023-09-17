package repository

import (
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"time"
)

// NewRelationshipRepository creates repository for persisting relationships
func NewRelationshipRepository(
	store DataStore,
) (Repository[types.Relationship], error) {
	return NewBaseRepository[types.Relationship](store,
		"Relationship",
		"",
		time.Duration(0),
		func() *types.Relationship {
			return &types.Relationship{}
		})
}
