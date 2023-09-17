package authz

import (
	"context"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_NullAuthorizerAuthorize(t *testing.T) {
	ctx := context.TODO()
	authorizer := NullAuthorizer{}
	req := &services.AuthRequest{
		PrincipalId: "root",
		Action:      "query",
		Resource:    "*",
	}
	_, err := authorizer.Authorize(ctx, req)
	require.Nil(t, err)
	_, err = authorizer.Check(ctx, nil)
	require.Nil(t, err) // not implemented
}

func Test_NoAuthorizerAuthorize(t *testing.T) {
	ctx := context.TODO()
	authorizer := NoAuthorizer{}
	req := &services.AuthRequest{
		PrincipalId: "root",
		Action:      "query",
		Resource:    "*",
	}
	_, err := authorizer.Authorize(ctx, req)
	require.Error(t, err) // not accessible
}
