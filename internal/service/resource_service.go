package service

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"time"
)

// ResourceService - admin APIs for resources data
type ResourceService interface {
	// CreateResource - create resource
	CreateResource(
		ctx context.Context,
		organizationID string,
		resource *types.Resource) (*types.Resource, error)

	// UpdateResource - updates resource
	UpdateResource(
		ctx context.Context,
		organizationID string,
		resource *types.Resource) error

	// DeleteResource removes resource
	DeleteResource(
		ctx context.Context,
		organizationID string,
		namespace string,
		id string) error

	// GetResource - finds resource
	GetResource(
		ctx context.Context,
		organizationID string,
		namespace string,
		id string,
	) (*types.Resource, error)

	// QueryResources - queries resources
	QueryResources(
		ctx context.Context,
		organizationID string,
		namespace string,
		predicate map[string]string,
		offset string,
		limit int64) (res []*types.Resource, nextOffset string, err error)

	// AllocateResourceInstance - allocates resource-instance
	AllocateResourceInstance(
		ctx context.Context,
		organizationID string,
		namespace string,
		resourceID string,
		principalID string,
		constraints string,
		expiry time.Duration,
		context map[string]string,
	) error

	// DeallocateResourceInstance - deallocates resource-instance
	DeallocateResourceInstance(
		ctx context.Context,
		organizationID string,
		namespace string,
		resourceID string,
		principalID string,
	) error

	// CountResourceInstances - size of total and allocated resource-instances
	CountResourceInstances(
		ctx context.Context,
		organizationID string,
		namespace string,
		resourceID string,
	) (capacity int32, allocated int32, err error)

	// QueryResourceInstances - queries resource-instances
	QueryResourceInstances(
		ctx context.Context,
		organizationID string,
		namespace string,
		resourceID string,
		predicate map[string]string,
		offset string,
		limit int64) (res []*types.ResourceInstance, nextOffset string, err error)
}
