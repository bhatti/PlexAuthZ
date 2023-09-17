package service

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
)

// OrganizationService - admin APIs for organization data
type OrganizationService interface {
	// GetOrganization finds organization
	GetOrganization(
		ctx context.Context,
		id string) (*types.Organization, error)

	// GetOrganizations - queries organizations
	GetOrganizations(
		ctx context.Context,
		predicate map[string]string,
		offset string,
		limit int64) (res []*types.Organization, nextOffset string, err error)

	// CreateOrganization - adds an organization
	CreateOrganization(
		ctx context.Context,
		org *types.Organization) (*types.Organization, error)

	// UpdateOrganization - updates organization
	UpdateOrganization(
		ctx context.Context,
		org *types.Organization) error

	// DeleteOrganization deletes organization
	DeleteOrganization(
		ctx context.Context,
		id string) error
}
