package http

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/web"
)

// OrganizationServiceHTTP - manages persistence of AuthZ data.
type OrganizationServiceHTTP struct {
	*baseHTTPClient
}

// NewOrganizationServiceHTTP manages persistence of organization.
func NewOrganizationServiceHTTP(
	client web.HTTPClient,
	baseURL string,
) *OrganizationServiceHTTP {
	return &OrganizationServiceHTTP{
		baseHTTPClient: &baseHTTPClient{
			client:  client,
			baseURL: baseURL,
		},
	}
}

// CreateOrganization - adds an organization.
func (h *OrganizationServiceHTTP) CreateOrganization(
	ctx context.Context,
	org *types.Organization) (*types.Organization, error) {
	req := &services.CreateOrganizationRequest{
		Namespaces: org.Namespaces,
		Name:       org.Name,
		Url:        org.Url,
		ParentIds:  org.ParentIds,
	}
	res := &services.CreateOrganizationResponse{}
	_, _, err := h.post(ctx,
		"/api/v1/organizations",
		req,
		res,
	)
	if err != nil {
		return nil, err
	}
	org.Id = res.Id
	return org, nil
}

// UpdateOrganization - updates organization in the database.
func (h *OrganizationServiceHTTP) UpdateOrganization(
	ctx context.Context,
	org *types.Organization) error {
	req := &services.UpdateOrganizationRequest{
		Id:         org.Id,
		Version:    org.Version,
		Namespaces: org.Namespaces,
		Name:       org.Name,
		Url:        org.Url,
		ParentIds:  org.ParentIds,
	}
	res := &services.UpdateOrganizationRequest{}
	_, _, err := h.put(ctx,
		fmt.Sprintf("/api/v1/organizations/%s", org.Id),
		req,
		res,
	)
	return err
}

// DeleteOrganization removes organization.
func (h *OrganizationServiceHTTP) DeleteOrganization(
	ctx context.Context,
	id string) error {
	if id == "" {
		return domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	_, _, err := h.del(ctx,
		fmt.Sprintf("/api/v1/organizations/%s", id),
	)
	return err
}

// GetOrganization finds organization
func (h *OrganizationServiceHTTP) GetOrganization(
	ctx context.Context,
	id string) (org *types.Organization, err error) {
	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	res := &services.GetOrganizationResponse{}
	_, _, err = h.get(
		ctx,
		fmt.Sprintf("/api/v1/organizations/%s", id),
		nil,
		res,
	)
	if err != nil {
		return nil, err
	}
	return &types.Organization{
		Id:         res.Id,
		Version:    res.Version,
		Name:       res.Name,
		Namespaces: res.Namespaces,
		Url:        res.Url,
		ParentIds:  res.ParentIds,
		Created:    res.Created,
		Updated:    res.Updated,
	}, nil
}

// GetOrganizations - queries organizations
func (h *OrganizationServiceHTTP) GetOrganizations(
	ctx context.Context,
	predicates map[string]string,
	offset string,
	limit int64) (arr []*types.Organization, nextOffset string, err error) {
	if predicates == nil {
		predicates = make(map[string]string)
	}
	res := &[]services.QueryOrganizationResponse{}
	predicates["offset"] = offset
	predicates["limit"] = fmt.Sprintf("%d", limit)
	_, resHeaders, err := h.get(
		ctx,
		"/api/v1/organizations",
		predicates,
		res,
	)
	if err != nil {
		return nil, "", err
	}
	for _, next := range *res {
		arr = append(arr, &types.Organization{
			Id:         next.Id,
			Version:    next.Version,
			Name:       next.Name,
			Namespaces: next.Namespaces,
			Url:        next.Url,
			ParentIds:  next.ParentIds,
			Created:    next.Created,
			Updated:    next.Updated,
		})
	}
	nextOffset = resHeaders[domain.NextOffsetHeader]
	return
}
