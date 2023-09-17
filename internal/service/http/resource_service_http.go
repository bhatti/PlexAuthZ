package http

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/web"
	"google.golang.org/protobuf/types/known/durationpb"
	"time"
)

// ResourceServiceHTTP - manages persistence of resources and resource instances
type ResourceServiceHTTP struct {
	*baseHTTPClient
}

// NewResourceServiceHTTP manages persistence of resources
func NewResourceServiceHTTP(
	client web.HTTPClient,
	baseURL string,
) *ResourceServiceHTTP {
	return &ResourceServiceHTTP{
		baseHTTPClient: &baseHTTPClient{
			client:  client,
			baseURL: baseURL,
		},
	}
}

// CreateResource - creates a new instance of resource
func (h *ResourceServiceHTTP) CreateResource(
	ctx context.Context,
	organizationID string,
	resource *types.Resource) (*types.Resource, error) {
	if organizationID == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	req := &services.CreateResourceRequest{
		OrganizationId: organizationID,
		Namespace:      resource.Namespace,
		Name:           resource.Name,
		Capacity:       resource.Capacity,
		Attributes:     resource.Attributes,
		AllowedActions: resource.AllowedActions,
	}
	res := &services.CreateResourceResponse{}
	_, _, err := h.post(ctx,
		fmt.Sprintf("/api/v1/%s/%s/resources", organizationID, resource.Namespace),
		req,
		res,
	)
	if err != nil {
		return nil, err
	}
	resource.Id = res.Id
	return resource, nil
}

// UpdateResource - updates an existing resource
func (h *ResourceServiceHTTP) UpdateResource(
	ctx context.Context,
	organizationID string,
	resource *types.Resource) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	req := &services.UpdateResourceRequest{
		Id:             resource.Id,
		OrganizationId: organizationID,
		Namespace:      resource.Namespace,
		Name:           resource.Name,
		Capacity:       resource.Capacity,
		Attributes:     resource.Attributes,
		AllowedActions: resource.AllowedActions,
	}
	res := &services.UpdateResourceResponse{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/%s/%s/resources/%s", organizationID, resource.Namespace, resource.Id),
		req,
		res,
	)
	return err
}

// DeleteResource removes resource
func (h *ResourceServiceHTTP) DeleteResource(
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
	_, _, err := h.del(ctx,
		fmt.Sprintf("/api/v1/%s/%s/resources/%s", organizationID, namespace, id),
	)
	return err
}

// GetResource - finds resource
func (h *ResourceServiceHTTP) GetResource(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string,
) (*types.Resource, error) {
	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	resources, _, err := h.QueryResources(
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

// QueryResources - queries resources.
func (h *ResourceServiceHTTP) QueryResources(
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
	if predicates == nil {
		predicates = make(map[string]string)
	}
	res := &[]services.QueryResourceResponse{}
	predicates["offset"] = offset
	predicates["limit"] = fmt.Sprintf("%d", limit)
	_, resHeaders, err := h.get(
		ctx,
		fmt.Sprintf("/api/v1/%s/%s/resources", organizationID, namespace),
		predicates,
		res,
	)
	if err != nil {
		return nil, "", err
	}
	for _, resourceRes := range *res {
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
	nextOffset = resHeaders[domain.NextOffsetHeader]
	return
}

// QueryResourceInstances - queries resource-instances.
func (h *ResourceServiceHTTP) QueryResourceInstances(
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
	if predicates == nil {
		predicates = make(map[string]string)
	}
	predicates["offset"] = offset
	predicates["limit"] = fmt.Sprintf("%d", limit)
	res := &[]services.QueryResourceInstanceResponse{}
	_, resHeaders, err := h.get(
		ctx,
		fmt.Sprintf("/api/v1/%s/%s/resources/%s/instances", organizationID, namespace, resourceID),
		predicates,
		res,
	)

	if err != nil {
		return nil, "", err
	}
	for _, resourceRes := range *res {
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
	nextOffset = resHeaders[domain.NextOffsetHeader]
	return
}

// CountResourceInstances - size of total and allocated resource-instances
func (h *ResourceServiceHTTP) CountResourceInstances(
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
	res := &services.CountResourceInstancesResponse{}
	_, _, err = h.get(
		ctx,
		fmt.Sprintf("/api/v1/%s/%s/resources/%s/instance_count", organizationID, namespace, resourceID),
		nil,
		res,
	)
	if err != nil {
		return 0, 0, err
	}
	return res.Capacity, res.Allocated, nil
}

// AllocateResourceInstance - allocates resource-instance
func (h *ResourceServiceHTTP) AllocateResourceInstance(
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
			fmt.Sprintf("principal-id is not defined for allocating resource"))
	}
	req := &services.AllocateResourceRequest{
		ResourceId:     resourceID,
		OrganizationId: organizationID,
		Namespace:      namespace,
		PrincipalId:    principalID,
		Constraints:    constraints,
		Expiry:         durationpb.New(expiry),
		Context:        context,
	}
	res := &services.AllocateResourceResponse{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/%s/%s/resources/%s/allocate/%s", organizationID, namespace, resourceID, principalID),
		req,
		res,
	)
	return err
}

// DeallocateResourceInstance - deallocates resource-instance
func (h *ResourceServiceHTTP) DeallocateResourceInstance(
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
			fmt.Sprintf("principal-id is not defined for deallocating resource"))
	}
	req := &services.DeallocateResourceRequest{
		ResourceId:     resourceID,
		OrganizationId: organizationID,
		Namespace:      namespace,
		PrincipalId:    principalID,
	}
	res := &services.DeallocateResourceResponse{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/%s/%s/resources/%s/deallocate/%s", organizationID, namespace, resourceID, principalID),
		req,
		res,
	)
	return err
}
