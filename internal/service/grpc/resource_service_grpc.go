package grpc

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/server"
	"google.golang.org/protobuf/types/known/durationpb"
	"time"
)

// ResourceServiceGrpc - manages persistence of resources and resource instances.
type ResourceServiceGrpc struct {
	clients server.Clients
}

// NewResourceServiceGrpc manages persistence of resources.
func NewResourceServiceGrpc(
	clients server.Clients,
) *ResourceServiceGrpc {
	return &ResourceServiceGrpc{
		clients: clients,
	}
}

// CreateResource - creates a new instance of resource.
func (s *ResourceServiceGrpc) CreateResource(
	ctx context.Context,
	organizationID string,
	resource *types.Resource) (*types.Resource, error) {
	if organizationID == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	res, err := s.clients.ResourcesClient.Create(
		ctx,
		&services.CreateResourceRequest{
			OrganizationId: organizationID,
			Namespace:      resource.Namespace,
			Name:           resource.Name,
			Capacity:       resource.Capacity,
			Attributes:     resource.Attributes,
			AllowedActions: resource.AllowedActions,
		})
	if err != nil {
		return nil, err
	}
	resource.Id = res.Id
	return resource, nil
}

// UpdateResource - updates an existing resource.
func (s *ResourceServiceGrpc) UpdateResource(
	ctx context.Context,
	organizationID string,
	resource *types.Resource) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	_, err := s.clients.ResourcesClient.Update(
		ctx,
		&services.UpdateResourceRequest{
			Id:             resource.Id,
			OrganizationId: organizationID,
			Namespace:      resource.Namespace,
			Name:           resource.Name,
			Capacity:       resource.Capacity,
			Attributes:     resource.Attributes,
			AllowedActions: resource.AllowedActions,
		})
	return err
}

// DeleteResource removes resource.
func (s *ResourceServiceGrpc) DeleteResource(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	if namespace == "" {
		return domain.NewValidationError(
			fmt.Sprintf("namespace is not defined"))
	}
	if id == "" {
		return domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	_, err := s.clients.ResourcesClient.Delete(
		ctx,
		&services.DeleteResourceRequest{
			OrganizationId: organizationID,
			Namespace:      namespace,
			Id:             id,
		})
	return err
}

// GetResource - finds resource by id.
func (s *ResourceServiceGrpc) GetResource(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string,
) (*types.Resource, error) {
	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	resources, _, err := s.QueryResources(
		ctx,
		organizationID,
		namespace,
		map[string]string{"id": id},
		"",
		1,
	)
	if err != nil {
		return nil, err
	}
	if len(resources) == 0 {
		return nil, domain.NewNotFoundError(fmt.Sprintf("resource %s is not found", id))
	}
	return resources[0], nil
}

// QueryResources - queries resources by predicates.
func (s *ResourceServiceGrpc) QueryResources(
	ctx context.Context,
	organizationID string,
	namespace string,
	predicates map[string]string,
	offset string,
	limit int64) (arr []*types.Resource, nextOffset string, err error) {
	if organizationID == "" {
		return nil, "", domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	if namespace == "" {
		return nil, "", domain.NewValidationError(
			fmt.Sprintf("namespace is not defined"))
	}
	res, err := s.clients.ResourcesClient.Query(
		ctx,
		&services.QueryResourceRequest{
			OrganizationId: organizationID,
			Namespace:      namespace,
			Predicates:     predicates,
			Offset:         offset,
			Limit:          limit,
		})
	if err != nil {
		return nil, "", err
	}
	for {
		resourceRes, err := res.Recv()
		if err != nil {
			break
		}
		nextOffset = resourceRes.NextOffset
		arr = append(arr, &types.Resource{
			Id:             resourceRes.Id,
			Version:        resourceRes.Version,
			Namespace:      resourceRes.Namespace,
			Created:        resourceRes.Created,
			Updated:        resourceRes.Updated,
			Name:           resourceRes.Name,
			Capacity:       resourceRes.Capacity,
			Attributes:     resourceRes.Attributes,
			AllowedActions: resourceRes.AllowedActions,
		})
	}
	return
}

// QueryResourceInstances - queries resource-instances.
func (s *ResourceServiceGrpc) QueryResourceInstances(
	ctx context.Context,
	organizationID string,
	namespace string,
	resourceID string,
	predicates map[string]string,
	offset string,
	limit int64) (arr []*types.ResourceInstance, nextOffset string, err error) {
	if organizationID == "" {
		return nil, "", domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	if namespace == "" {
		return nil, "", domain.NewValidationError(
			fmt.Sprintf("namespace is not defined"))
	}
	res, err := s.clients.ResourcesClient.QueryResourceInstances(
		ctx,
		&services.QueryResourceInstanceRequest{
			OrganizationId: organizationID,
			Namespace:      namespace,
			ResourceId:     resourceID,
			Predicates:     predicates,
			Offset:         offset,
			Limit:          limit,
		})
	if err != nil {
		return nil, "", err
	}
	for {
		resourceRes, err := res.Recv()
		if err != nil {
			break
		}
		nextOffset = resourceRes.NextOffset
		arr = append(arr, &types.ResourceInstance{
			Id:          resourceRes.Id,
			Version:     resourceRes.Version,
			Namespace:   resourceRes.Namespace,
			ResourceId:  resourceRes.ResourceId,
			PrincipalId: resourceRes.PrincipalId,
			State:       resourceRes.State,
			Created:     resourceRes.Created,
			Updated:     resourceRes.Updated,
		})
	}
	return
}

// CountResourceInstances - size of total and allocated resource-instances.
func (s *ResourceServiceGrpc) CountResourceInstances(
	ctx context.Context,
	organizationID string,
	namespace string,
	resourceID string,
) (capacity int32, allocated int32, err error) {
	if organizationID == "" {
		return 0, 0, domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	if namespace == "" {
		return 0, 0, domain.NewValidationError(
			fmt.Sprintf("namespace is not defined"))
	}
	if resourceID == "" {
		return 0, 0, domain.NewValidationError(
			fmt.Sprintf("resource-id is not defined"))
	}
	res, err := s.clients.ResourcesClient.CountResourceInstances(
		ctx,
		&services.CountResourceInstancesRequest{
			OrganizationId: organizationID,
			Namespace:      namespace,
			ResourceId:     resourceID,
		})
	if err != nil {
		return 0, 0, err
	}
	return res.Capacity, res.Allocated, nil
}

// AllocateResourceInstance - allocates resource-instance.
func (s *ResourceServiceGrpc) AllocateResourceInstance(
	ctx context.Context,
	organizationID string,
	namespace string,
	resourceID string,
	principalID string,
	constraints string,
	expiry time.Duration,
	context map[string]string,
) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	if namespace == "" {
		return domain.NewValidationError(
			fmt.Sprintf("namespace is not defined"))
	}
	if resourceID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("resource-id is not defined"))
	}
	if principalID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("principal-id is not defined"))
	}
	_, err := s.clients.AuthClient.Allocate(
		ctx,
		&services.AllocateResourceRequest{
			OrganizationId: organizationID,
			Namespace:      namespace,
			ResourceId:     resourceID,
			PrincipalId:    principalID,
			Constraints:    constraints,
			Expiry:         durationpb.New(expiry),
			Context:        context,
		})
	return err
}

// DeallocateResourceInstance - deallocates resource-instance.
func (s *ResourceServiceGrpc) DeallocateResourceInstance(
	ctx context.Context,
	organizationID string,
	namespace string,
	resourceID string,
	principalID string,
) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	if namespace == "" {
		return domain.NewValidationError(
			fmt.Sprintf("namespace is not defined"))
	}
	if resourceID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("resource-id is not defined"))
	}
	if principalID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("principal-id is not defined"))
	}
	_, err := s.clients.AuthClient.Deallocate(
		ctx,
		&services.DeallocateResourceRequest{
			OrganizationId: organizationID,
			Namespace:      namespace,
			ResourceId:     resourceID,
			PrincipalId:    principalID,
		})
	return err
}
