package authz

import (
	"context"
	"fmt"
	"github.com/bhatti/PlexAuthZ/api/v1/services"
	"github.com/bhatti/PlexAuthZ/api/v1/types"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/casbin/casbin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type authorizer struct {
	enforcer *casbin.Enforcer
}

type subjectContextKey struct{}

// NewGrpcAuth constructor
func NewGrpcAuth(config *domain.Config) (Authorizer, error) {
	aclFile, err := config.ACLModelFile()
	if err != nil {
		return nil, err
	}
	policyFile, err := config.ACLPolicyFile()
	if err != nil {
		return nil, err
	}
	return &authorizer{
		enforcer: casbin.NewEnforcer(aclFile, policyFile),
	}, nil
}

// Subject returns subject from context
func Subject(ctx context.Context) string {
	v := ctx.Value(subjectContextKey{})
	if v == nil {
		return ""
	}
	return v.(string)
}

// Authenticate checks access
func Authenticate(ctx context.Context) (context.Context, error) {
	peer, ok := peer.FromContext(ctx)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"ctx": ctx,
		}).Warn("authenticator couldn't find peer info")
		return ctx, status.New(
			codes.Unknown,
			"authenticator couldn't find peer info",
		).Err()
	}

	if peer.AuthInfo == nil {
		logrus.WithFields(logrus.Fields{
			"ctx":  ctx,
			"peer": peer,
		}).Warn("authenticator couldn't find peer authz info")
		return context.WithValue(ctx, subjectContextKey{}, ""), nil
	}

	tlsInfo := peer.AuthInfo.(credentials.TLSInfo)
	subject := tlsInfo.State.VerifiedChains[0][0].Subject.CommonName
	return context.WithValue(ctx, subjectContextKey{}, subject), nil
}

// Authorize checks authorization permission
func (a *authorizer) Authorize(
	_ context.Context,
	req *services.AuthRequest,
) (*services.AuthResponse, error) {
	valid := a.enforcer.Enforce(req.PrincipalId, req.Resource, req.Action)
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		logrus.WithFields(logrus.Fields{
			"subject": req.PrincipalId,
			"object":  req.Resource,
			"action":  req.Action,
			"valid":   valid,
		}).Debugf("authenticator authorizing")
	}
	if !valid {
		return nil, status.New(
			codes.PermissionDenied,
			fmt.Sprintf("%s not permitted to %s to %s",
				req.PrincipalId, req.Action, req.Resource)).Err()
	}
	return &services.AuthResponse{
		Effect: types.Effect_PERMITTED,
	}, nil
}

// Check enforces constraints but not implemented
func (a *authorizer) Check(
	_ context.Context,
	_ *services.CheckConstraintsRequest,
) (*services.CheckConstraintsResponse, error) {
	return nil, domain.NewInternalError("check method is not implemented for GRPC", domain.InternalCode)
}
