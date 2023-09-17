package db

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/metrics"
	"github.com/bhatti/PlexAuthZ/internal/service"
)

func newAuthServiceAndOrg() (service.AuthAdminService, *types.Organization, error) {
	ctx := context.TODO()
	cfg, err := domain.NewConfig("")
	if err != nil {
		return nil, nil, err
	}
	authService, _, err := CreateDatabaseAuthService(cfg, metrics.New())
	if err != nil {
		return nil, nil, err
	}
	org, err := domain.NewOrganizationBuilder().
		WithName("test-org1").
		WithNamespaces("admin", "finance", "engineering").Build()
	if err != nil {
		return nil, nil, err
	}

	savedOrg, err := authService.CreateOrganization(ctx, org)
	return authService, savedOrg, err
}
