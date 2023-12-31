// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: api/v1/services/relationship_service.proto

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

// RelationshipsServiceClient is the client API for RelationshipsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RelationshipsServiceClient interface {
	// Create Relationships swagger:route POST /api/v1/{organization_id}/{namespace}/relations relationships createRelationshipRequest
	//
	// Responses:
	// 200: createRelationshipResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Create(ctx context.Context, in *CreateRelationshipRequest, opts ...grpc.CallOption) (*CreateRelationshipResponse, error)
	// Update Relationships swagger:route PUT /api/v1/{organization_id}/{namespace}/relations/{id} relationships updateRelationshipRequest
	//
	// Responses:
	// 200: updateRelationshipResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Update(ctx context.Context, in *UpdateRelationshipRequest, opts ...grpc.CallOption) (*UpdateRelationshipResponse, error)
	// Query Relationship swagger:route GET /api/v1/{organization_id}/{namespace}/relations relationships queryRelationshipRequest
	//
	// Responses:
	// 200: queryRelationshipResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Query(ctx context.Context, in *QueryRelationshipRequest, opts ...grpc.CallOption) (RelationshipsService_QueryClient, error)
	// Delete Relationship swagger:route DELETE /api/v1/{organization_id}/{namespace}/relations/{id} relationships deleteRelationshipRequest
	//
	// Responses:
	// 200: deleteRelationshipResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Delete(ctx context.Context, in *DeleteRelationshipRequest, opts ...grpc.CallOption) (*DeleteRelationshipResponse, error)
}

type relationshipsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRelationshipsServiceClient(cc grpc.ClientConnInterface) RelationshipsServiceClient {
	return &relationshipsServiceClient{cc}
}

func (c *relationshipsServiceClient) Create(ctx context.Context, in *CreateRelationshipRequest, opts ...grpc.CallOption) (*CreateRelationshipResponse, error) {
	out := new(CreateRelationshipResponse)
	err := c.cc.Invoke(ctx, "/api.authz.services.RelationshipsService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationshipsServiceClient) Update(ctx context.Context, in *UpdateRelationshipRequest, opts ...grpc.CallOption) (*UpdateRelationshipResponse, error) {
	out := new(UpdateRelationshipResponse)
	err := c.cc.Invoke(ctx, "/api.authz.services.RelationshipsService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationshipsServiceClient) Query(ctx context.Context, in *QueryRelationshipRequest, opts ...grpc.CallOption) (RelationshipsService_QueryClient, error) {
	stream, err := c.cc.NewStream(ctx, &RelationshipsService_ServiceDesc.Streams[0], "/api.authz.services.RelationshipsService/Query", opts...)
	if err != nil {
		return nil, err
	}
	x := &relationshipsServiceQueryClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RelationshipsService_QueryClient interface {
	Recv() (*QueryRelationshipResponse, error)
	grpc.ClientStream
}

type relationshipsServiceQueryClient struct {
	grpc.ClientStream
}

func (x *relationshipsServiceQueryClient) Recv() (*QueryRelationshipResponse, error) {
	m := new(QueryRelationshipResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *relationshipsServiceClient) Delete(ctx context.Context, in *DeleteRelationshipRequest, opts ...grpc.CallOption) (*DeleteRelationshipResponse, error) {
	out := new(DeleteRelationshipResponse)
	err := c.cc.Invoke(ctx, "/api.authz.services.RelationshipsService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RelationshipsServiceServer is the server API for RelationshipsService service.
// All implementations must embed UnimplementedRelationshipsServiceServer
// for forward compatibility
type RelationshipsServiceServer interface {
	// Create Relationships swagger:route POST /api/v1/{organization_id}/{namespace}/relations relationships createRelationshipRequest
	//
	// Responses:
	// 200: createRelationshipResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Create(context.Context, *CreateRelationshipRequest) (*CreateRelationshipResponse, error)
	// Update Relationships swagger:route PUT /api/v1/{organization_id}/{namespace}/relations/{id} relationships updateRelationshipRequest
	//
	// Responses:
	// 200: updateRelationshipResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Update(context.Context, *UpdateRelationshipRequest) (*UpdateRelationshipResponse, error)
	// Query Relationship swagger:route GET /api/v1/{organization_id}/{namespace}/relations relationships queryRelationshipRequest
	//
	// Responses:
	// 200: queryRelationshipResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Query(*QueryRelationshipRequest, RelationshipsService_QueryServer) error
	// Delete Relationship swagger:route DELETE /api/v1/{organization_id}/{namespace}/relations/{id} relationships deleteRelationshipRequest
	//
	// Responses:
	// 200: deleteRelationshipResponse
	// 400	Bad Request
	// 401	Not Authorized
	// 500	Internal Error
	Delete(context.Context, *DeleteRelationshipRequest) (*DeleteRelationshipResponse, error)
	mustEmbedUnimplementedRelationshipsServiceServer()
}

// UnimplementedRelationshipsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRelationshipsServiceServer struct {
}

func (UnimplementedRelationshipsServiceServer) Create(context.Context, *CreateRelationshipRequest) (*CreateRelationshipResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedRelationshipsServiceServer) Update(context.Context, *UpdateRelationshipRequest) (*UpdateRelationshipResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedRelationshipsServiceServer) Query(*QueryRelationshipRequest, RelationshipsService_QueryServer) error {
	return status.Errorf(codes.Unimplemented, "method Query not implemented")
}
func (UnimplementedRelationshipsServiceServer) Delete(context.Context, *DeleteRelationshipRequest) (*DeleteRelationshipResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedRelationshipsServiceServer) mustEmbedUnimplementedRelationshipsServiceServer() {}

// UnsafeRelationshipsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RelationshipsServiceServer will
// result in compilation errors.
type UnsafeRelationshipsServiceServer interface {
	mustEmbedUnimplementedRelationshipsServiceServer()
}

func RegisterRelationshipsServiceServer(s grpc.ServiceRegistrar, srv RelationshipsServiceServer) {
	s.RegisterService(&RelationshipsService_ServiceDesc, srv)
}

func _RelationshipsService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRelationshipRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationshipsServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.authz.services.RelationshipsService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationshipsServiceServer).Create(ctx, req.(*CreateRelationshipRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelationshipsService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRelationshipRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationshipsServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.authz.services.RelationshipsService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationshipsServiceServer).Update(ctx, req.(*UpdateRelationshipRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelationshipsService_Query_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(QueryRelationshipRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RelationshipsServiceServer).Query(m, &relationshipsServiceQueryServer{stream})
}

type RelationshipsService_QueryServer interface {
	Send(*QueryRelationshipResponse) error
	grpc.ServerStream
}

type relationshipsServiceQueryServer struct {
	grpc.ServerStream
}

func (x *relationshipsServiceQueryServer) Send(m *QueryRelationshipResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _RelationshipsService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRelationshipRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationshipsServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.authz.services.RelationshipsService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationshipsServiceServer).Delete(ctx, req.(*DeleteRelationshipRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RelationshipsService_ServiceDesc is the grpc.ServiceDesc for RelationshipsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RelationshipsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.authz.services.RelationshipsService",
	HandlerType: (*RelationshipsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _RelationshipsService_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _RelationshipsService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _RelationshipsService_Delete_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Query",
			Handler:       _RelationshipsService_Query_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/v1/services/relationship_service.proto",
}
