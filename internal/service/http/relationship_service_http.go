package http

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/web"
)

// RelationshipServiceHTTP - manages persistence of relationship data
type RelationshipServiceHTTP struct {
	*baseHTTPClient
}

// NewRelationshipServiceHTTP manages persistence of relationship data
func NewRelationshipServiceHTTP(
	client web.HTTPClient,
	baseURL string,
) *RelationshipServiceHTTP {
	return &RelationshipServiceHTTP{
		baseHTTPClient: &baseHTTPClient{
			client:  client,
			baseURL: baseURL,
		},
	}
}

// CreateRelationship - creates a new relationship
func (h *RelationshipServiceHTTP) CreateRelationship(
	ctx context.Context,
	organizationID string,
	relationship *types.Relationship) (*types.Relationship, error) {
	if organizationID == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	req := &services.CreateRelationshipRequest{
		OrganizationId: organizationID,
		Namespace:      relationship.Namespace,
		Relation:       relationship.Relation,
		PrincipalId:    relationship.PrincipalId,
		ResourceId:     relationship.ResourceId,
		Attributes:     relationship.Attributes,
	}
	res := &services.CreateRelationshipResponse{}
	_, _, err := h.post(ctx,
		fmt.Sprintf("/api/v1/%s/%s/relations", organizationID, relationship.Namespace),
		req,
		res,
	)
	if err != nil {
		return nil, err
	}
	relationship.Id = res.Id
	return relationship, nil
}

// UpdateRelationship - updates an existing relationship
func (h *RelationshipServiceHTTP) UpdateRelationship(
	ctx context.Context,
	organizationID string,
	relationship *types.Relationship) error {
	if organizationID == "" {
		return domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	req := &services.UpdateRelationshipRequest{
		OrganizationId: organizationID,
		Namespace:      relationship.Namespace,
		Id:             relationship.Id,
		Relation:       relationship.Relation,
		PrincipalId:    relationship.PrincipalId,
		ResourceId:     relationship.ResourceId,
		Attributes:     relationship.Attributes,
	}
	res := &services.UpdateRelationshipResponse{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/%s/%s/relations/%s", organizationID, relationship.Namespace, relationship.Id),
		req,
		res,
	)
	return err
}

// DeleteRelationship removes relationship
func (h *RelationshipServiceHTTP) DeleteRelationship(
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
		fmt.Sprintf("/api/v1/%s/%s/relations/%s", organizationID, namespace, id),
	)
	return err
}

// GetRelationship - finds relationship
func (h *RelationshipServiceHTTP) GetRelationship(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string,
) (*types.Relationship, error) {
	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	relations, _, err := h.GetRelationships(
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
	if len(relations) == 0 {
		return nil, domain.NewNotFoundError(fmt.Sprintf("relationship %s is not found", id))
	}
	return relations[0], nil
}

// GetRelationships - queries relationships
func (h *RelationshipServiceHTTP) GetRelationships(
	ctx context.Context,
	organizationID string,
	namespace string,
	predicates map[string]string,
	offset string,
	limit int64) (arr []*types.Relationship, nextOffset string, err error) {
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
	res := &[]services.QueryRelationshipResponse{}
	predicates["offset"] = offset
	predicates["limit"] = fmt.Sprintf("%d", limit)
	_, resHeaders, err := h.get(
		ctx,
		fmt.Sprintf("/api/v1/%s/%s/relations", organizationID, namespace),
		predicates,
		res,
	)
	if err != nil {
		return nil, "", err
	}
	for _, relationRes := range *res {
		arr = append(arr, &types.Relationship{
			Id:          relationRes.Id,
			Version:     relationRes.Version,
			Namespace:   relationRes.Namespace,
			Relation:    relationRes.Relation,
			PrincipalId: relationRes.PrincipalId,
			ResourceId:  relationRes.ResourceId,
			Attributes:  relationRes.Attributes,
			Created:     relationRes.Created,
			Updated:     relationRes.Updated,
		})
	}
	nextOffset = resHeaders[domain.NextOffsetHeader]
	return
}
