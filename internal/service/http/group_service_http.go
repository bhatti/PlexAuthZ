package http

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/web"
)

// GroupServiceHTTP - manages persistence of groups data
type GroupServiceHTTP struct {
	*baseHTTPClient
}

// NewGroupServiceHTTP manages persistence of groups data
func NewGroupServiceHTTP(
	client web.HTTPClient,
	baseURL string,
) *GroupServiceHTTP {
	return &GroupServiceHTTP{
		baseHTTPClient: &baseHTTPClient{
			client:  client,
			baseURL: baseURL,
		},
	}
}

// CreateGroup - creates a new group
func (h *GroupServiceHTTP) CreateGroup(
	ctx context.Context,
	organizationID string,
	group *types.Group) (*types.Group, error) {
	if organizationID == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	req := &services.CreateGroupRequest{
		OrganizationId: organizationID,
		Namespace:      group.Namespace,
		Name:           group.Name,
		ParentIds:      group.ParentIds,
	}
	res := &services.CreateGroupResponse{}
	_, _, err := h.post(ctx,
		fmt.Sprintf("/api/v1/%s/%s/groups", organizationID, group.Namespace),
		req,
		res,
	)
	if err != nil {
		return nil, err
	}
	group.Id = res.Id
	return group, nil
}

// UpdateGroup - updates an existing group
func (h *GroupServiceHTTP) UpdateGroup(
	ctx context.Context,
	organizationID string,
	group *types.Group) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	req := &services.UpdateGroupRequest{
		Id:             group.Id,
		OrganizationId: organizationID,
		Namespace:      group.Namespace,
		Name:           group.Name,
		ParentIds:      group.ParentIds,
	}
	res := &services.UpdateGroupResponse{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/%s/%s/groups/%s", organizationID, group.Namespace, group.Id),
		req,
		res,
	)
	return err
}

// DeleteGroup removes group
func (h *GroupServiceHTTP) DeleteGroup(
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
		fmt.Sprintf("/api/v1/%s/%s/groups/%s", organizationID, namespace, id),
	)
	return err
}

// GetGroup - finds group
func (h *GroupServiceHTTP) GetGroup(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string,
) (*types.Group, error) {
	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	groups, _, err := h.GetGroups(
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
	if len(groups) == 0 {
		return nil, domain.NewNotFoundError(fmt.Sprintf("group %s is not found", id))
	}
	return groups[0], nil
}

// GetGroups - queries groups
func (h *GroupServiceHTTP) GetGroups(
	ctx context.Context,
	organizationID string,
	namespace string,
	predicates map[string]string,
	offset string,
	limit int64) (arr []*types.Group, nextOffset string, err error) {
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
	res := &[]services.QueryGroupResponse{}
	predicates["offset"] = offset
	predicates["limit"] = fmt.Sprintf("%d", limit)
	_, resHeaders, err := h.get(
		ctx,
		fmt.Sprintf("/api/v1/%s/%s/groups", organizationID, namespace),
		predicates,
		res,
	)
	if err != nil {
		return nil, "", err
	}
	for _, next := range *res {
		arr = append(arr, &types.Group{
			Id:        next.Id,
			Version:   next.Version,
			Namespace: next.Namespace,
			Name:      next.Name,
			RoleIds:   next.RoleIds,
			Updated:   next.Updated,
		})
	}
	nextOffset = resHeaders[domain.NextOffsetHeader]
	return
}

// AddRolesToGroup helper
func (h *GroupServiceHTTP) AddRolesToGroup(
	ctx context.Context,
	organizationID string,
	namespace string,
	groupID string,
	roleIDs ...string,
) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	if namespace == "" {
		return domain.NewValidationError(
			fmt.Sprintf("namespace is not defined"))
	}
	if groupID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("group-id is not defined"))
	}
	if len(roleIDs) == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("role-ids is not defined"))
	}
	req := &services.AddRolesToGroupRequest{
		OrganizationId: organizationID,
		Namespace:      namespace,
		GroupId:        groupID,
		RoleIds:        roleIDs,
	}
	res := &services.AddRolesToGroupResponse{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/%s/%s/groups/%s/roles/add", organizationID, namespace, groupID),
		req,
		res,
	)
	return err
}

// DeleteRolesToGroup helper
func (h *GroupServiceHTTP) DeleteRolesToGroup(
	ctx context.Context,
	organizationID string,
	namespace string,
	groupID string,
	roleIDs ...string,
) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	if namespace == "" {
		return domain.NewValidationError(
			fmt.Sprintf("namespace is not defined"))
	}
	if groupID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("group-id is not defined"))
	}
	if len(roleIDs) == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("role-ids is not defined"))
	}
	req := &services.DeleteRolesToGroupRequest{
		OrganizationId: organizationID,
		Namespace:      namespace,
		GroupId:        groupID,
		RoleIds:        roleIDs,
	}
	res := &services.DeleteRolesToGroupResponse{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/%s/%s/groups/%s/roles/delete", organizationID, namespace, groupID),
		req,
		res,
	)
	return err
}
