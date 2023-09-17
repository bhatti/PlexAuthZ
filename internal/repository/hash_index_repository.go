package repository

import (
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"time"
)

// NewHashIndexRepository creates repository for indexes
func NewHashIndexRepository(
	store DataStore,
) (Repository[domain.HashIndex], error) {
	return NewBaseRepository[domain.HashIndex](store,
		"HashIndex",
		"",
		time.Duration(0),
		func() *domain.HashIndex {
			return &domain.HashIndex{}
		})
}
