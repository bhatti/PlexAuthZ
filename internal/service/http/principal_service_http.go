package http

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/web"
)

// PrincipalServiceHTTP - manages persistence of principal objects
type PrincipalServiceHTTP struct {
	*baseHTTPClient
}

// NewPrincipalServiceHTTP manages persistence of principal data
func NewPrincipalServiceHTTP(
	client web.HTTPClient,
	baseURL string,
) *PrincipalServiceHTTP {
	return &PrincipalServiceHTTP{
		baseHTTPClient: &baseHTTPClient{
			client:  client,
			baseURL: baseURL,
		},
	}
}

// CreatePrincipal - creates new instance of principal
// Note - this method won't be used to update any role-ids, group-ids, relations, and permission-ids
func (h *PrincipalServiceHTTP) CreatePrincipal(
	ctx context.Context,
	principal *types.Principal) (*types.Principal, error) {
	req := &services.CreatePrincipalRequest{
		OrganizationId: principal.OrganizationId,
		Namespaces:     principal.Namespaces,
		Username:       principal.Username,
		Name:           principal.Name,
		Attributes:     principal.Attributes,
	}
	res := &services.CreatePrincipalResponse{}

	_, _, err := h.post(ctx,
		fmt.Sprintf("/api/v1/%s/principals", principal.OrganizationId),
		req,
		res,
	)
	if err != nil {
		return nil, err
	}
	principal.Id = res.Id
	return principal, nil
}

// UpdatePrincipal - updates existing instance of principal
// Note - this method won't be used to update any role-ids, group-ids, relations, and permission-ids
func (h *PrincipalServiceHTTP) UpdatePrincipal(
	ctx context.Context,
	principal *types.Principal) error {
	req := &services.UpdatePrincipalRequest{
		Id:             principal.Id,
		OrganizationId: principal.OrganizationId,
		Namespaces:     principal.Namespaces,
		Username:       principal.Username,
		Name:           principal.Name,
		Attributes:     principal.Attributes,
	}
	res := &services.UpdatePrincipalResponse{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/%s/principals/%s", principal.OrganizationId, principal.Id),
		req,
		res,
	)
	return err
}

// DeletePrincipal removes principal
func (h *PrincipalServiceHTTP) DeletePrincipal(
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
	_, _, err := h.del(ctx,
		fmt.Sprintf("/api/v1/%s/principals/%s", organizationID, id),
	)
	return err
}

// AddGroupsToPrincipal helper
func (h *PrincipalServiceHTTP) AddGroupsToPrincipal(
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
			fmt.Sprintf("principal-id is not defined for adding group"))
	}
	if len(groupIDs) == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("group-ids is not defined"))
	}
	req := &services.AddGroupsToPrincipalRequest{
		OrganizationId: organizationID,
		Namespace:      namespace,
		PrincipalId:    principalID,
		GroupIds:       groupIDs,
	}
	res := &services.AddGroupsToPrincipalResponse{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/%s/%s/principals/%s/groups/add", organizationID, namespace, principalID),
		req,
		res,
	)
	return err
}

// DeleteGroupsToPrincipal helper
func (h *PrincipalServiceHTTP) DeleteGroupsToPrincipal(
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
			fmt.Sprintf("principal-id is not defined for deleting group"))
	}
	if len(groupIDs) == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("group-ids is not defined"))
	}
	req := &services.DeleteGroupsToPrincipalRequest{
		OrganizationId: organizationID,
		Namespace:      namespace,
		PrincipalId:    principalID,
		GroupIds:       groupIDs,
	}
	res := &services.DeleteGroupsToPrincipalResponse{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/%s/%s/principals/%s/groups/delete", organizationID, namespace, principalID),
		req,
		res,
	)
	return err
}

// AddRolesToPrincipal helper
func (h *PrincipalServiceHTTP) AddRolesToPrincipal(
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
			fmt.Sprintf("principal-id is not defined for adding role"))
	}
	if len(roleIDs) == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("role-ids is not defined"))
	}
	req := &services.AddRolesToPrincipalRequest{
		OrganizationId: organizationID,
		Namespace:      namespace,
		PrincipalId:    principalID,
		RoleIds:        roleIDs,
	}
	res := &services.AddRolesToPrincipalResponse{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/%s/%s/principals/%s/roles/add", organizationID, namespace, principalID),
		req,
		res,
	)
	return err
}

// DeleteRolesToPrincipal helper
func (h *PrincipalServiceHTTP) DeleteRolesToPrincipal(
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
			fmt.Sprintf("principal-id is not defined for deleting role"))
	}
	if len(roleIDs) == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("role-ids is not defined"))
	}
	req := &services.DeleteRolesToPrincipalRequest{
		OrganizationId: organizationID,
		Namespace:      namespace,
		PrincipalId:    principalID,
		RoleIds:        roleIDs,
	}
	res := &services.DeleteRolesToPrincipalResponse{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/%s/%s/principals/%s/roles/delete", organizationID, namespace, principalID),
		req,
		res,
	)
	return err
}

// AddPermissionsToPrincipal helper
func (h *PrincipalServiceHTTP) AddPermissionsToPrincipal(
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
			fmt.Sprintf("principal-id is not defined for adding permission"))
	}
	if len(permissionIds) == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("permission-ids is not defined"))
	}
	req := &services.AddPermissionsToPrincipalRequest{
		OrganizationId: organizationID,
		Namespace:      namespace,
		PrincipalId:    principalID,
		PermissionIds:  permissionIds,
	}
	res := &services.AddPermissionsToPrincipalResponse{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/%s/%s/principals/%s/permissions/add", organizationID, namespace, principalID),
		req,
		res,
	)
	return err
}

// DeletePermissionsToPrincipal helper
func (h *PrincipalServiceHTTP) DeletePermissionsToPrincipal(
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
			fmt.Sprintf("principal-id is not defined for deleting permission"))
	}
	if len(permissionIds) == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("permission-ids is not defined"))
	}
	req := &services.DeletePermissionsToPrincipalRequest{
		OrganizationId: organizationID,
		Namespace:      namespace,
		PrincipalId:    principalID,
		PermissionIds:  permissionIds,
	}
	res := &services.DeletePermissionsToPrincipalResponse{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/%s/%s/principals/%s/permissions/delete", organizationID, namespace, principalID),
		req,
		res,
	)
	return err
}

// AddRelationshipsToPrincipal helper
func (h *PrincipalServiceHTTP) AddRelationshipsToPrincipal(
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
			fmt.Sprintf("principal-id is not defined for adding relation"))
	}
	if len(relationshipIds) == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("relationship-ids is not defined"))
	}
	req := &services.AddRelationshipsToPrincipalRequest{
		OrganizationId:  organizationID,
		Namespace:       namespace,
		PrincipalId:     principalID,
		RelationshipIds: relationshipIds,
	}
	res := &services.AddRelationshipsToPrincipalResponse{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/%s/%s/principals/%s/relations/add", organizationID, namespace, principalID),
		req,
		res,
	)
	return err
}

// DeleteRelationshipsToPrincipal helper
func (h *PrincipalServiceHTTP) DeleteRelationshipsToPrincipal(
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
			fmt.Sprintf("principal-id is not defined for deleting relation"))
	}
	if len(relationshipIds) == 0 {
		return domain.NewValidationError(
			fmt.Sprintf("relationship-ids is not defined"))
	}
	req := &services.DeleteRelationshipsToPrincipalRequest{
		OrganizationId:  organizationID,
		Namespace:       namespace,
		PrincipalId:     principalID,
		RelationshipIds: relationshipIds,
	}
	res := &services.DeleteRelationshipsToPrincipalResponse{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/%s/%s/principals/%s/relations/delete", organizationID, namespace, principalID),
		req,
		res,
	)
	return err
}

// GetPrincipal - retrieves principal
func (h *PrincipalServiceHTTP) GetPrincipal(
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
	principals, _, err := h.GetPrincipals(
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
func (h *PrincipalServiceHTTP) GetPrincipals(
	ctx context.Context,
	organizationID string,
	predicates map[string]string,
	offset string,
	limit int64) (arr []*types.Principal, nextOffset string, err error) {
	if predicates == nil {
		predicates = make(map[string]string)
	}
	res := &[]services.QueryPrincipalResponse{}
	predicates["offset"] = offset
	predicates["limit"] = fmt.Sprintf("%d", limit)
	_, resHeaders, err := h.get(
		ctx,
		fmt.Sprintf("/api/v1/%s/principals", organizationID),
		predicates,
		res,
	)
	if err != nil {
		return nil, "", err
	}
	for _, next := range *res {
		arr = append(arr, &types.Principal{
			Id:             next.Id,
			Version:        next.Version,
			Namespaces:     next.Namespaces,
			OrganizationId: next.OrganizationId,
			Username:       next.Username,
			Name:           next.Name,
			Email:          next.Email,
			Attributes:     next.Attributes,
			GroupIds:       next.GroupIds,
			RoleIds:        next.RoleIds,
			PermissionIds:  next.PermissionIds,
			RelationIds:    next.RelationIds,
			Created:        next.Created,
			Updated:        next.Updated,
		})
	}
	nextOffset = resHeaders[domain.NextOffsetHeader]
	return
}

// GetPrincipalExt - retrieves full principal
func (h *PrincipalServiceHTTP) GetPrincipalExt(
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
	res := &services.GetPrincipalResponse{}

	_, _, err = h.get(
		ctx,
		fmt.Sprintf("/api/v1/%s/%s/principals/%s", organizationID, namespace, id),
		nil,
		res,
	)
	if err != nil {
		return nil, err
	}
	principalExt := domain.NewPrincipalExtFromResponse(res)
	if err = principalExt.Validate(); err != nil {
		return nil, err
	}
	return principalExt, nil
}
