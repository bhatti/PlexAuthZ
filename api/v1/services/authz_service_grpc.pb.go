// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: api/v1/services/authz_service.proto

package services

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AuthZServiceClient is the client API for AuthZService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthZServiceClient interface {
	// Authorize swagger:route POST /api/v1/{organization_id}/{namespace}/{principal_id}/auth authz authRequest
	//
	// Responses:
	// 200: authResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Authorize(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (*AuthResponse, error)
	// Check swagger:route POST /api/v1/{organization_id}/{namespace}/{principal_id}/auth/constraints authz checkConstraintsRequest
	//
	// Responses:
	// 200: checkConstraintsResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Check(ctx context.Context, in *CheckConstraintsRequest, opts ...grpc.CallOption) (*CheckConstraintsResponse, error)
	// Allocate Resources swagger:route PUT /api/v1/{organization_id}/{namespace}/resources/{id}/allocate/{principal_id} resources allocateResourceRequest
	//
	// Responses:
	// 200: allocateResourceResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Allocate(ctx context.Context, in *AllocateResourceRequest, opts ...grpc.CallOption) (*AllocateResourceResponse, error)
	// Deallocate Resources swagger:route PUT /api/v1/{organization_id}/{namespace}/resources/{id}/deallocate/{principal_id} resources deallocateResourceRequest
	//
	// Responses:
	// 200: deallocateResourceResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Deallocate(ctx context.Context, in *DeallocateResourceRequest, opts ...grpc.CallOption) (*DeallocateResourceResponse, error)
}

type authZServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthZServiceClient(cc grpc.ClientConnInterface) AuthZServiceClient {
	return &authZServiceClient{cc}
}

func (c *authZServiceClient) Authorize(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (*AuthResponse, error) {
	out := new(AuthResponse)
	err := c.cc.Invoke(ctx, "/api.authz.services.AuthZService/Authorize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authZServiceClient) Check(ctx context.Context, in *CheckConstraintsRequest, opts ...grpc.CallOption) (*CheckConstraintsResponse, error) {
	out := new(CheckConstraintsResponse)
	err := c.cc.Invoke(ctx, "/api.authz.services.AuthZService/Check", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authZServiceClient) Allocate(ctx context.Context, in *AllocateResourceRequest, opts ...grpc.CallOption) (*AllocateResourceResponse, error) {
	out := new(AllocateResourceResponse)
	err := c.cc.Invoke(ctx, "/api.authz.services.AuthZService/Allocate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authZServiceClient) Deallocate(ctx context.Context, in *DeallocateResourceRequest, opts ...grpc.CallOption) (*DeallocateResourceResponse, error) {
	out := new(DeallocateResourceResponse)
	err := c.cc.Invoke(ctx, "/api.authz.services.AuthZService/Deallocate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthZServiceServer is the server API for AuthZService service.
// All implementations must embed UnimplementedAuthZServiceServer
// for forward compatibility
type AuthZServiceServer interface {
	// Authorize swagger:route POST /api/v1/{organization_id}/{namespace}/{principal_id}/auth authz authRequest
	//
	// Responses:
	// 200: authResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Authorize(context.Context, *AuthRequest) (*AuthResponse, error)
	// Check swagger:route POST /api/v1/{organization_id}/{namespace}/{principal_id}/auth/constraints authz checkConstraintsRequest
	//
	// Responses:
	// 200: checkConstraintsResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Check(context.Context, *CheckConstraintsRequest) (*CheckConstraintsResponse, error)
	// Allocate Resources swagger:route PUT /api/v1/{organization_id}/{namespace}/resources/{id}/allocate/{principal_id} resources allocateResourceRequest
	//
	// Responses:
	// 200: allocateResourceResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Allocate(context.Context, *AllocateResourceRequest) (*AllocateResourceResponse, error)
	// Deallocate Resources swagger:route PUT /api/v1/{organization_id}/{namespace}/resources/{id}/deallocate/{principal_id} resources deallocateResourceRequest
	//
	// Responses:
	// 200: deallocateResourceResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Deallocate(context.Context, *DeallocateResourceRequest) (*DeallocateResourceResponse, error)
	mustEmbedUnimplementedAuthZServiceServer()
}

// UnimplementedAuthZServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuthZServiceServer struct {
}

func (UnimplementedAuthZServiceServer) Authorize(context.Context, *AuthRequest) (*AuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authorize not implemented")
}
func (UnimplementedAuthZServiceServer) Check(context.Context, *CheckConstraintsRequest) (*CheckConstraintsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
func (UnimplementedAuthZServiceServer) Allocate(context.Context, *AllocateResourceRequest) (*AllocateResourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Allocate not implemented")
}
func (UnimplementedAuthZServiceServer) Deallocate(context.Context, *DeallocateResourceRequest) (*DeallocateResourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deallocate not implemented")
}
func (UnimplementedAuthZServiceServer) mustEmbedUnimplementedAuthZServiceServer() {}

// UnsafeAuthZServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthZServiceServer will
// result in compilation errors.
type UnsafeAuthZServiceServer interface {
	mustEmbedUnimplementedAuthZServiceServer()
}

func RegisterAuthZServiceServer(s grpc.ServiceRegistrar, srv AuthZServiceServer) {
	s.RegisterService(&AuthZService_ServiceDesc, srv)
}

func _AuthZService_Authorize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthZServiceServer).Authorize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.authz.services.AuthZService/Authorize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthZServiceServer).Authorize(ctx, req.(*AuthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthZService_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckConstraintsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthZServiceServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.authz.services.AuthZService/Check",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthZServiceServer).Check(ctx, req.(*CheckConstraintsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthZService_Allocate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AllocateResourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthZServiceServer).Allocate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.authz.services.AuthZService/Allocate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthZServiceServer).Allocate(ctx, req.(*AllocateResourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthZService_Deallocate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeallocateResourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthZServiceServer).Deallocate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.authz.services.AuthZService/Deallocate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthZServiceServer).Deallocate(ctx, req.(*DeallocateResourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthZService_ServiceDesc is the grpc.ServiceDesc for AuthZService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthZService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.authz.services.AuthZService",
	HandlerType: (*AuthZServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Authorize",
			Handler:    _AuthZService_Authorize_Handler,
		},
		{
			MethodName: "Check",
			Handler:    _AuthZService_Check_Handler,
		},
		{
			MethodName: "Allocate",
			Handler:    _AuthZService_Allocate_Handler,
		},
		{
			MethodName: "Deallocate",
			Handler:    _AuthZService_Deallocate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v1/services/authz_service.proto",
}
