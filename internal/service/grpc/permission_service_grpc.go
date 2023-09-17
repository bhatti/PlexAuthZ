package grpc

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/server"
)

// PermissionServiceGrpc - manages persistence of permission data
type PermissionServiceGrpc struct {
	clients server.Clients
}

// NewPermissionServiceGrpc manages persistence of permission data
func NewPermissionServiceGrpc(
	clients server.Clients,
) *PermissionServiceGrpc {
	return &PermissionServiceGrpc{
		clients: clients,
	}
}

// CreatePermission - creates a new permission
func (s *PermissionServiceGrpc) CreatePermission(
	ctx context.Context,
	organizationID string,
	permission *types.Permission) (*types.Permission, error) {
	if organizationID == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	res, err := s.clients.PermissionsClient.Create(
		ctx,
		&services.CreatePermissionRequest{
			OrganizationId: organizationID,
			Namespace:      permission.Namespace,
			Scope:          permission.Scope,
			Actions:        permission.Actions,
			ResourceId:     permission.ResourceId,
			Effect:         permission.Effect,
			Constraints:    permission.Constraints,
		})
	if err != nil {
		return nil, err
	}
	permission.Id = res.Id
	return permission, nil
}

// UpdatePermission - updates an existing permission
func (s *PermissionServiceGrpc) UpdatePermission(
	ctx context.Context,
	organizationID string,
	permission *types.Permission) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	_, err := s.clients.PermissionsClient.Update(
		ctx,
		&services.UpdatePermissionRequest{
			Id:             permission.Id,
			OrganizationId: organizationID,
			Namespace:      permission.Namespace,
			Scope:          permission.Scope,
			Actions:        permission.Actions,
			ResourceId:     permission.ResourceId,
			Effect:         permission.Effect,
			Constraints:    permission.Constraints,
		})
	return err
}

// DeletePermission removes permission
func (s *PermissionServiceGrpc) DeletePermission(
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
	_, err := s.clients.PermissionsClient.Delete(
		ctx,
		&services.DeletePermissionRequest{
			OrganizationId: organizationID,
			Namespace:      namespace,
			Id:             id,
		})
	return err
}

// GetPermission - finds permission
func (s *PermissionServiceGrpc) GetPermission(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string,
) (*types.Permission, error) {
	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	permissions, _, err := s.GetPermissions(
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
func (s *PermissionServiceGrpc) GetPermissions(
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
	res, err := s.clients.PermissionsClient.Query(
		ctx,
		&services.QueryPermissionRequest{
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
		permission, err := res.Recv()
		if err != nil {
			break
		}
		nextOffset = permission.NextOffset
		arr = append(arr, &types.Permission{
			Id:          permission.Id,
			Version:     permission.Version,
			Namespace:   permission.Namespace,
			Scope:       permission.Scope,
			Actions:     permission.Actions,
			ResourceId:  permission.ResourceId,
			Effect:      permission.Effect,
			Constraints: permission.Constraints,
			Created:     permission.Created,
			Updated:     permission.Updated,
		})
	}
	return
}
