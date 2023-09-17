package grpc

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/server"
)

// GroupServiceGrpc - manages persistence of groups data
type GroupServiceGrpc struct {
	clients server.Clients
}

// NewGroupServiceGrpc manages persistence of groups data
func NewGroupServiceGrpc(
	clients server.Clients,
) *GroupServiceGrpc {
	return &GroupServiceGrpc{
		clients: clients,
	}
}

// CreateGroup - creates a new group
func (s *GroupServiceGrpc) CreateGroup(
	ctx context.Context,
	organizationID string,
	group *types.Group) (*types.Group, error) {
	if organizationID == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	res, err := s.clients.GroupsClient.Create(
		ctx,
		&services.CreateGroupRequest{
			OrganizationId: organizationID,
			Namespace:      group.Namespace,
			Name:           group.Name,
			ParentIds:      group.ParentIds,
		})
	if err != nil {
		return nil, err
	}
	group.Id = res.Id
	return group, nil
}

// UpdateGroup - updates an existing group
func (s *GroupServiceGrpc) UpdateGroup(
	ctx context.Context,
	organizationID string,
	group *types.Group) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	_, err := s.clients.GroupsClient.Update(
		ctx,
		&services.UpdateGroupRequest{
			Id:             group.Id,
			OrganizationId: organizationID,
			Namespace:      group.Namespace,
			Name:           group.Name,
			ParentIds:      group.ParentIds,
		})
	return err
}

// DeleteGroup removes group
func (s *GroupServiceGrpc) DeleteGroup(
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
	_, err := s.clients.GroupsClient.Delete(
		ctx,
		&services.DeleteGroupRequest{
			OrganizationId: organizationID,
			Namespace:      namespace,
			Id:             id,
		})
	return err
}

// GetGroup - finds group
func (s *GroupServiceGrpc) GetGroup(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string,
) (*types.Group, error) {
	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	groups, _, err := s.GetGroups(
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
func (s *GroupServiceGrpc) GetGroups(
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
	res, err := s.clients.GroupsClient.Query(
		ctx,
		&services.QueryGroupRequest{
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
		group, err := res.Recv()
		if err != nil {
			break
		}
		nextOffset = group.NextOffset
		arr = append(arr, &types.Group{
			Id:        group.Id,
			Version:   group.Version,
			Namespace: group.Namespace,
			Name:      group.Name,
			RoleIds:   group.RoleIds,
			Updated:   group.Updated,
		})
	}
	return
}

// AddRolesToGroup helper
func (s *GroupServiceGrpc) AddRolesToGroup(
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
	_, err := s.clients.GroupsClient.AddRoles(ctx, &services.AddRolesToGroupRequest{
		OrganizationId: organizationID,
		Namespace:      namespace,
		GroupId:        groupID,
		RoleIds:        roleIDs,
	})
	return err
}

// DeleteRolesToGroup helper
func (s *GroupServiceGrpc) DeleteRolesToGroup(
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
	_, err := s.clients.GroupsClient.DeleteRoles(ctx, &services.DeleteRolesToGroupRequest{
		OrganizationId: organizationID,
		Namespace:      namespace,
		GroupId:        groupID,
		RoleIds:        roleIDs,
	})
	return err
}
