package http

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/web"
)

// PermissionServiceHTTP - manages persistence of permission data
type PermissionServiceHTTP struct {
	*baseHTTPClient
}

// NewPermissionServiceHTTP manages persistence of permission data
func NewPermissionServiceHTTP(
	client web.HTTPClient,
	baseURL string,
) *PermissionServiceHTTP {
	return &PermissionServiceHTTP{
		baseHTTPClient: &baseHTTPClient{
			client:  client,
			baseURL: baseURL,
		},
	}
}

// CreatePermission - creates a new permission
func (h *PermissionServiceHTTP) CreatePermission(
	ctx context.Context,
	organizationID string,
	permission *types.Permission) (*types.Permission, error) {
	if organizationID == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	req := &services.CreatePermissionRequest{
		OrganizationId: organizationID,
		Namespace:      permission.Namespace,
		Scope:          permission.Scope,
		Actions:        permission.Actions,
		ResourceId:     permission.ResourceId,
		Effect:         permission.Effect,
		Constraints:    permission.Constraints,
	}
	res := &services.CreatePermissionResponse{}
	_, _, err := h.post(ctx,
		fmt.Sprintf("/api/v1/%s/%s/permissions", organizationID, permission.Namespace),
		req,
		res,
	)
	if err != nil {
		return nil, err
	}
	permission.Id = res.Id
	return permission, nil
}

// UpdatePermission - updates an existing permission
func (h *PermissionServiceHTTP) UpdatePermission(
	ctx context.Context,
	organizationID string,
	permission *types.Permission) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	req := &services.UpdatePermissionRequest{
		Id:             permission.Id,
		OrganizationId: organizationID,
		Namespace:      permission.Namespace,
		Scope:          permission.Scope,
		Actions:        permission.Actions,
		ResourceId:     permission.ResourceId,
		Effect:         permission.Effect,
		Constraints:    permission.Constraints,
	}
	res := &services.UpdatePermissionResponse{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/%s/%s/permissions/%s", organizationID, permission.Namespace, permission.Id),
		req,
		res,
	)
	return err
}

// DeletePermission removes permission
func (h *PermissionServiceHTTP) DeletePermission(
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
		fmt.Sprintf("/api/v1/%s/%s/permissions/%s", organizationID, namespace, id),
	)
	return err
}

// GetPermission - finds permission
func (h *PermissionServiceHTTP) GetPermission(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string,
) (*types.Permission, error) {
	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	permissions, _, err := h.GetPermissions(
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
	if len(permissions) == 0 {
		return nil, domain.NewNotFoundError(fmt.Sprintf("permission %s is not found", id))
	}
	return permissions[0], nil
}

// GetPermissions - queries permissions
func (h *PermissionServiceHTTP) GetPermissions(
	ctx context.Context,
	organizationID string,
	namespace string,
	predicates map[string]string,
	offset string,
	limit int64) (arr []*types.Permission, nextOffset string, err error) {
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
	res := &[]services.QueryPermissionResponse{}
	predicates["offset"] = offset
	predicates["limit"] = fmt.Sprintf("%d", limit)
	_, resHeaders, err := h.get(
		ctx,
		fmt.Sprintf("/api/v1/%s/%s/permissions", organizationID, namespace),
		predicates,
		res,
	)
	if err != nil {
		return nil, "", err
	}
	for _, next := range *res {
		nextOffset = next.NextOffset
		arr = append(arr, &types.Permission{
			Id:          next.Id,
			Version:     next.Version,
			Namespace:   next.Namespace,
			Scope:       next.Scope,
			Actions:     next.Actions,
			ResourceId:  next.ResourceId,
			Effect:      next.Effect,
			Constraints: next.Constraints,
			Created:     next.Created,
			Updated:     next.Updated,
		})
	}
	nextOffset = resHeaders[domain.NextOffsetHeader]
	return
}
