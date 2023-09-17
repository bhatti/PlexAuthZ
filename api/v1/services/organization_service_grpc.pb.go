// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: api/v1/services/organization_service.proto

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

// OrganizationsServiceClient is the client API for OrganizationsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrganizationsServiceClient interface {
	// Create Organizations swagger:route POST /api/v1/organizations organizations createOrganizationRequest
	//
	// Responses:
	// 200: createOrganizationResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Create(ctx context.Context, in *CreateOrganizationRequest, opts ...grpc.CallOption) (*CreateOrganizationResponse, error)
	// Update Organizations swagger:route PUT /api/v1/organizations/{id} organizations updateOrganizationRequest
	//
	// Responses:
	// 200: updateOrganizationResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Update(ctx context.Context, in *UpdateOrganizationRequest, opts ...grpc.CallOption) (*UpdateOrganizationResponse, error)
	// Get Organization swagger:route GET /api/v1/organizations/{id} organizations getOrganizationRequest
	//
	// Responses:
	// 200: getOrganizationResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Get(ctx context.Context, in *GetOrganizationRequest, opts ...grpc.CallOption) (*GetOrganizationResponse, error)
	// Query Organization swagger:route GET /api/v1/organizations organizations queryOrganizationRequest
	//
	// Responses:
	// 200: queryOrganizationResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Query(ctx context.Context, in *QueryOrganizationRequest, opts ...grpc.CallOption) (OrganizationsService_QueryClient, error)
	// Delete Organization swagger:route DELETE /api/v1/organizations/{id} organizations deleteOrganizationRequest
	//
	// Responses:
	// 200: deleteOrganizationResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Delete(ctx context.Context, in *DeleteOrganizationRequest, opts ...grpc.CallOption) (*DeleteOrganizationResponse, error)
}

type organizationsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOrganizationsServiceClient(cc grpc.ClientConnInterface) OrganizationsServiceClient {
	return &organizationsServiceClient{cc}
}

func (c *organizationsServiceClient) Create(ctx context.Context, in *CreateOrganizationRequest, opts ...grpc.CallOption) (*CreateOrganizationResponse, error) {
	out := new(CreateOrganizationResponse)
	err := c.cc.Invoke(ctx, "/api.authz.services.OrganizationsService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *organizationsServiceClient) Update(ctx context.Context, in *UpdateOrganizationRequest, opts ...grpc.CallOption) (*UpdateOrganizationResponse, error) {
	out := new(UpdateOrganizationResponse)
	err := c.cc.Invoke(ctx, "/api.authz.services.OrganizationsService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *organizationsServiceClient) Get(ctx context.Context, in *GetOrganizationRequest, opts ...grpc.CallOption) (*GetOrganizationResponse, error) {
	out := new(GetOrganizationResponse)
	err := c.cc.Invoke(ctx, "/api.authz.services.OrganizationsService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *organizationsServiceClient) Query(ctx context.Context, in *QueryOrganizationRequest, opts ...grpc.CallOption) (OrganizationsService_QueryClient, error) {
	stream, err := c.cc.NewStream(ctx, &OrganizationsService_ServiceDesc.Streams[0], "/api.authz.services.OrganizationsService/Query", opts...)
	if err != nil {
		return nil, err
	}
	x := &organizationsServiceQueryClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type OrganizationsService_QueryClient interface {
	Recv() (*QueryOrganizationResponse, error)
	grpc.ClientStream
}

type organizationsServiceQueryClient struct {
	grpc.ClientStream
}

func (x *organizationsServiceQueryClient) Recv() (*QueryOrganizationResponse, error) {
	m := new(QueryOrganizationResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *organizationsServiceClient) Delete(ctx context.Context, in *DeleteOrganizationRequest, opts ...grpc.CallOption) (*DeleteOrganizationResponse, error) {
	out := new(DeleteOrganizationResponse)
	err := c.cc.Invoke(ctx, "/api.authz.services.OrganizationsService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrganizationsServiceServer is the server API for OrganizationsService service.
// All implementations must embed UnimplementedOrganizationsServiceServer
// for forward compatibility
type OrganizationsServiceServer interface {
	// Create Organizations swagger:route POST /api/v1/organizations organizations createOrganizationRequest
	//
	// Responses:
	// 200: createOrganizationResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Create(context.Context, *CreateOrganizationRequest) (*CreateOrganizationResponse, error)
	// Update Organizations swagger:route PUT /api/v1/organizations/{id} organizations updateOrganizationRequest
	//
	// Responses:
	// 200: updateOrganizationResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Update(context.Context, *UpdateOrganizationRequest) (*UpdateOrganizationResponse, error)
	// Get Organization swagger:route GET /api/v1/organizations/{id} organizations getOrganizationRequest
	//
	// Responses:
	// 200: getOrganizationResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Get(context.Context, *GetOrganizationRequest) (*GetOrganizationResponse, error)
	// Query Organization swagger:route GET /api/v1/organizations organizations queryOrganizationRequest
	//
	// Responses:
	// 200: queryOrganizationResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Query(*QueryOrganizationRequest, OrganizationsService_QueryServer) error
	// Delete Organization swagger:route DELETE /api/v1/organizations/{id} organizations deleteOrganizationRequest
	//
	// Responses:
	// 200: deleteOrganizationResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Delete(context.Context, *DeleteOrganizationRequest) (*DeleteOrganizationResponse, error)
	mustEmbedUnimplementedOrganizationsServiceServer()
}

// UnimplementedOrganizationsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedOrganizationsServiceServer struct {
}

func (UnimplementedOrganizationsServiceServer) Create(context.Context, *CreateOrganizationRequest) (*CreateOrganizationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedOrganizationsServiceServer) Update(context.Context, *UpdateOrganizationRequest) (*UpdateOrganizationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedOrganizationsServiceServer) Get(context.Context, *GetOrganizationRequest) (*GetOrganizationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedOrganizationsServiceServer) Query(*QueryOrganizationRequest, OrganizationsService_QueryServer) error {
	return status.Errorf(codes.Unimplemented, "method Query not implemented")
}
func (UnimplementedOrganizationsServiceServer) Delete(context.Context, *DeleteOrganizationRequest) (*DeleteOrganizationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedOrganizationsServiceServer) mustEmbedUnimplementedOrganizationsServiceServer() {}

// UnsafeOrganizationsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrganizationsServiceServer will
// result in compilation errors.
type UnsafeOrganizationsServiceServer interface {
	mustEmbedUnimplementedOrganizationsServiceServer()
}

func RegisterOrganizationsServiceServer(s grpc.ServiceRegistrar, srv OrganizationsServiceServer) {
	s.RegisterService(&OrganizationsService_ServiceDesc, srv)
}

func _OrganizationsService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrganizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganizationsServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.authz.services.OrganizationsService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganizationsServiceServer).Create(ctx, req.(*CreateOrganizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrganizationsService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateOrganizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganizationsServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.authz.services.OrganizationsService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganizationsServiceServer).Update(ctx, req.(*UpdateOrganizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrganizationsService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOrganizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganizationsServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.authz.services.OrganizationsService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganizationsServiceServer).Get(ctx, req.(*GetOrganizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrganizationsService_Query_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(QueryOrganizationRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(OrganizationsServiceServer).Query(m, &organizationsServiceQueryServer{stream})
}

type OrganizationsService_QueryServer interface {
	Send(*QueryOrganizationResponse) error
	grpc.ServerStream
}

type organizationsServiceQueryServer struct {
	grpc.ServerStream
}

func (x *organizationsServiceQueryServer) Send(m *QueryOrganizationResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _OrganizationsService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteOrganizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganizationsServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.authz.services.OrganizationsService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganizationsServiceServer).Delete(ctx, req.(*DeleteOrganizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OrganizationsService_ServiceDesc is the grpc.ServiceDesc for OrganizationsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OrganizationsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.authz.services.OrganizationsService",
	HandlerType: (*OrganizationsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _OrganizationsService_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _OrganizationsService_Update_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _OrganizationsService_Get_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _OrganizationsService_Delete_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Query",
			Handler:       _OrganizationsService_Query_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/v1/services/organization_service.proto",
}