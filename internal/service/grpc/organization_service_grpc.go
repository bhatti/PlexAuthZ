package grpc

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/server"
)

// OrganizationServiceGrpc - manages persistence of AuthZ data.
type OrganizationServiceGrpc struct {
	clients server.Clients
}

// NewOrganizationServiceGrpc manages persistence of organization.
func NewOrganizationServiceGrpc(
	clients server.Clients,
) *OrganizationServiceGrpc {
	return &OrganizationServiceGrpc{
		clients: clients,
	}
}

// GetOrganization finds organization.
func (s *OrganizationServiceGrpc) GetOrganization(
	ctx context.Context,
	id string) (org *types.Organization, err error) {
	if id == "" {
		return nil, domain.NewValidationError(
			fmt.Sprintf("organization-id is not defined"))
	}
	res, err := s.clients.OrganizationsClient.Get(
		ctx,
		&services.GetOrganizationRequest{
			Id: id,
		},
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

// GetOrganizations - queries organizations.
func (s *OrganizationServiceGrpc) GetOrganizations(
	ctx context.Context,
	predicates map[string]string,
	offset string,
	limit int64) (arr []*types.Organization, nextToken string, err error) {
	res, err := s.clients.OrganizationsClient.Query(
		ctx,
		&services.QueryOrganizationRequest{
			Predicates: predicates,
			Offset:     offset,
			Limit:      limit,
		})
	if err != nil {
		return nil, "", err
	}
	for {
		orgRes, err := res.Recv()
		if err != nil {
			break
		}
		org := &types.Organization{
			Id:         orgRes.Id,
			Version:    orgRes.Version,
			Name:       orgRes.Name,
			Namespaces: orgRes.Namespaces,
			Url:        orgRes.Url,
			ParentIds:  orgRes.ParentIds,
			Created:    orgRes.Created,
			Updated:    orgRes.Updated,
		}
		nextToken = orgRes.NextOffset
		arr = append(arr, org)
	}
	return
}

// CreateOrganization - adds an organization.
func (s *OrganizationServiceGrpc) CreateOrganization(
	ctx context.Context,
	org *types.Organization) (*types.Organization, error) {
	res, err := s.clients.OrganizationsClient.Create(
		ctx,
		&services.CreateOrganizationRequest{
			Namespaces: org.Namespaces,
			Name:       org.Name,
			Url:        org.Url,
			ParentIds:  org.ParentIds,
		})
	if err != nil {
		return nil, err
	}
	org.Id = res.Id
	return org, nil
}

// UpdateOrganization - updates organization.
func (s *OrganizationServiceGrpc) UpdateOrganization(
	ctx context.Context,
	org *types.Organization) error {
	_, err := s.clients.OrganizationsClient.Update(
		ctx,
		&services.UpdateOrganizationRequest{
			Id:         org.Id,
			Version:    org.Version,
			Namespaces: org.Namespaces,
			Name:       org.Name,
			Url:        org.Url,
			ParentIds:  org.ParentIds,
		})
	return err
}

// DeleteOrganization removes organization.
func (s *OrganizationServiceGrpc) DeleteOrganization(
	ctx context.Context,
	id string) error {
	if id == "" {
		return domain.NewValidationError(
			fmt.Sprintf("id is not defined"))
	}
	_, err := s.clients.OrganizationsClient.Delete(
		ctx,
		&services.DeleteOrganizationRequest{
			Id: id,
		})
	return err
}
