package grpc

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/server"
)

// PrincipalServiceGrpc - manages persistence of principal objects
type PrincipalServiceGrpc struct {
	clients server.Clients
}

// NewPrincipalServiceGrpc manages persistence of principal data
func NewPrincipalServiceGrpc(
	clients server.Clients,
) *PrincipalServiceGrpc {
	return &PrincipalServiceGrpc{
		clients: clients,
	}
}

// CreatePrincipal - creates new instance of principal
// Note - this method won't be used to update any role-ids, group-ids, relations, and permission-ids
func (s *PrincipalServiceGrpc) CreatePrincipal(
	ctx context.Context,
	principal *types.Principal) (*types.Principal, error) {
	res, err := s.clients.PrincipalsClient.Create(
		ctx,
		&services.CreatePrincipalRequest{
			OrganizationId: principal.OrganizationId,
			Namespaces:     principal.Namespaces,
			Username:       principal.Username,
			Name:           principal.Name,
			Attributes:     principal.Attributes,
		})
	if err != nil {
		return nil, err
	}
	principal.Id = res.Id
	return principal, nil
}

// UpdatePrincipal - updates existing instance of principal
// Note - this method won't be used to update any role-ids, group-ids, relations, and permission-ids
func (s *PrincipalServiceGrpc) UpdatePrincipal(
	ctx context.Context,
	principal *types.Principal) error {
	_, err := s.clients.PrincipalsClient.Update(
		ctx,
		&services.UpdatePrincipalRequest{
			Id:             principal.Id,
			OrganizationId: principal.OrganizationId,
			Namespaces:     principal.Namespaces,
			Username:       principal.Username,
			Name:           principal.Name,
			Attributes:     principal.Attributes,
		})
	return err
}

// DeletePrincipal removes principal
func (s *PrincipalServiceGrpc) DeletePrincipal(
	ctx context.Context,
	organizationID string,
	id string) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	if id == "" {
		return domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	_, err := s.clients.PrincipalsClient.Delete(
		ctx,
		&services.DeletePrincipalRequest{
			OrganizationId: organizationID,
			Id:             id,
		})
	return err
}

// AddGroupsToPrincipal helper
func (s *PrincipalServiceGrpc) AddGroupsToPrincipal(
	ctx context.Context,
	organizationID string,
	namespace string,
	principalID string,
	groupIDs ...string,
) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	if namespace == "" {
		return domain.NewValidationError(
			fmt.Sprintf("namespace is not defined"))
	}
	if principalID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("principal-id is not defined"))
	}
	if len(groupIDs) == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("group-ids is not defined"))
	}
	_, err := s.clients.PrincipalsClient.AddGroups(ctx, &services.AddGroupsToPrincipalRequest{
		OrganizationId: organizationID,
		Namespace:      namespace,
		PrincipalId:    principalID,
		GroupIds:       groupIDs,
	})
	return err
}

// DeleteGroupsToPrincipal helper
func (s *PrincipalServiceGrpc) DeleteGroupsToPrincipal(
	ctx context.Context,
	organizationID string,
	namespace string,
	principalID string,
	groupIDs ...string,
) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	if namespace == "" {
		return domain.NewValidationError(
			fmt.Sprintf("namespace is not defined"))
	}
	if principalID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("principal-id is not defined"))
	}
	if len(groupIDs) == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("group-ids is not defined"))
	}
	_, err := s.clients.PrincipalsClient.DeleteGroups(ctx, &services.DeleteGroupsToPrincipalRequest{
		OrganizationId: organizationID,
		Namespace:      namespace,
		PrincipalId:    principalID,
		GroupIds:       groupIDs,
	})
	return err
}

// AddRolesToPrincipal helper
func (s *PrincipalServiceGrpc) AddRolesToPrincipal(
	ctx context.Context,
	organizationID string,
	namespace string,
	principalID string,
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
	if principalID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("principal-id is not defined"))
	}
	if len(roleIDs) == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("role-ids is not defined"))
	}
	_, err := s.clients.PrincipalsClient.AddRoles(ctx, &services.AddRolesToPrincipalRequest{
		OrganizationId: organizationID,
		Namespace:      namespace,
		PrincipalId:    principalID,
		RoleIds:        roleIDs,
	})
	return err
}

// DeleteRolesToPrincipal helper
func (s *PrincipalServiceGrpc) DeleteRolesToPrincipal(
	ctx context.Context,
	organizationID string,
	namespace string,
	principalID string,
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
	if principalID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("principal-id is not defined"))
	}
	if len(roleIDs) == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("role-ids is not defined"))
	}
	_, err := s.clients.PrincipalsClient.DeleteRoles(ctx, &services.DeleteRolesToPrincipalRequest{
		OrganizationId: organizationID,
		Namespace:      namespace,
		PrincipalId:    principalID,
		RoleIds:        roleIDs,
	})
	return err
}

// AddPermissionsToPrincipal helper
func (s *PrincipalServiceGrpc) AddPermissionsToPrincipal(
	ctx context.Context,
	organizationID string,
	namespace string,
	principalID string,
	permissionIds ...string,
) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	if namespace == "" {
		return domain.NewValidationError(
			fmt.Sprintf("namespace is not defined"))
	}
	if principalID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("principal-id is not defined"))
	}
	if len(permissionIds) == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("permission-ids is not defined"))
	}
	_, err := s.clients.PrincipalsClient.AddPermissions(ctx, &services.AddPermissionsToPrincipalRequest{
		OrganizationId: organizationID,
		Namespace:      namespace,
		PrincipalId:    principalID,
		PermissionIds:  permissionIds,
	})
	return err
}

// DeletePermissionsToPrincipal helper
func (s *PrincipalServiceGrpc) DeletePermissionsToPrincipal(
	ctx context.Context,
	organizationID string,
	namespace string,
	principalID string,
	permissionIds ...string,
) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	if namespace == "" {
		return domain.NewValidationError(
			fmt.Sprintf("namespace is not defined"))
	}
	if principalID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("principal-id is not defined"))
	}
	if len(permissionIds) == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("permission-ids is not defined"))
	}
	_, err := s.clients.PrincipalsClient.DeletePermissions(ctx, &services.DeletePermissionsToPrincipalRequest{
		OrganizationId: organizationID,
		Namespace:      namespace,
		PrincipalId:    principalID,
		PermissionIds:  permissionIds,
	})
	return err
}

// AddRelationshipsToPrincipal helper
func (s *PrincipalServiceGrpc) AddRelationshipsToPrincipal(
	ctx context.Context,
	organizationID string,
	namespace string,
	principalID string,
	relationshipIds ...string,
) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	if namespace == "" {
		return domain.NewValidationError(
			fmt.Sprintf("namespace is not defined"))
	}
	if principalID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("principal-id is not defined"))
	}
	if len(relationshipIds) == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("relationship-ids is not defined"))
	}
	_, err := s.clients.PrincipalsClient.AddRelationships(ctx, &services.AddRelationshipsToPrincipalRequest{
		OrganizationId:  organizationID,
		Namespace:       namespace,
		PrincipalId:     principalID,
		RelationshipIds: relationshipIds,
	})
	return err
}

// DeleteRelationshipsToPrincipal helper
func (s *PrincipalServiceGrpc) DeleteRelationshipsToPrincipal(
	ctx context.Context,
	organizationID string,
	namespace string,
	principalID string,
	relationshipIds ...string,
) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	if namespace == "" {
		return domain.NewValidationError(
			fmt.Sprintf("namespace is not defined"))
	}
	if principalID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("principal-id is not defined"))
	}
	if len(relationshipIds) == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("relationship-ids is not defined"))
	}
	_, err := s.clients.PrincipalsClient.DeleteRelationships(ctx, &services.DeleteRelationshipsToPrincipalRequest{
		OrganizationId:  organizationID,
		Namespace:       namespace,
		PrincipalId:     principalID,
		RelationshipIds: relationshipIds,
	})
	return err
}

// GetPrincipal - retrieves principal
func (s *PrincipalServiceGrpc) GetPrincipal(
	ctx context.Context,
	organizationID string,
	id string,
) (*types.Principal, error) {
	if organizationID == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("organization_id is not defined"))
	}

	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	principals, _, err := s.GetPrincipals(
		ctx,
		organizationID,
		map[string]string{"id": id},
		"",
		1,
	)
	if err != nil {
		return nil, err
	}
	if len(principals) == 0 {
		return nil, domain.NewNotFoundError(fmt.Sprintf("principal %s is not found", id))
	}
	return principals[0], nil
}

// GetPrincipals - queries principals
func (s *PrincipalServiceGrpc) GetPrincipals(
	ctx context.Context,
	organizationID string,
	predicates map[string]string,
	offset string,
	limit int64) (arr []*types.Principal, nextToken string, err error) {
	res, err := s.clients.PrincipalsClient.Query(
		ctx,
		&services.QueryPrincipalRequest{
			OrganizationId: organizationID,
			Predicates:     predicates,
			Offset:         offset,
			Limit:          limit,
		})
	if err != nil {
		return nil, "", err
	}
	for {
		principalRes, err := res.Recv()
		if err != nil {
			break
		}
		principal := &types.Principal{
			Id:             principalRes.Id,
			Version:        principalRes.Version,
			Namespaces:     principalRes.Namespaces,
			OrganizationId: principalRes.OrganizationId,
			Username:       principalRes.Username,
			Name:           principalRes.Name,
			Email:          principalRes.Email,
			Attributes:     principalRes.Attributes,
			GroupIds:       principalRes.GroupIds,
			RoleIds:        principalRes.RoleIds,
			PermissionIds:  principalRes.PermissionIds,
			RelationIds:    principalRes.RelationIds,
			Created:        principalRes.Created,
			Updated:        principalRes.Updated,
		}
		nextToken = principalRes.NextOffset
		arr = append(arr, principal)
	}
	return
}

// GetPrincipalExt - retrieves full principal
func (s *PrincipalServiceGrpc) GetPrincipalExt(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string,
) (xPrincipal *domain.PrincipalExt, err error) {
	if organizationID == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("organization_id is not defined"))
	}

	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	res, err := s.clients.PrincipalsClient.Get(
		ctx,
		&services.GetPrincipalRequest{
			OrganizationId: organizationID,
			Namespace:      namespace,
			Id:             id,
		},
	)
	if err != nil {
		return nil, err
	}
	principalExt := domain.NewPrincipalExtFromResponse(res)
	return principalExt, nil
}
