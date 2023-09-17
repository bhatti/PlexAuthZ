// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.1
// source: api/v1/services/permission_service.proto

package services

import (
	types "github.com/bhatti/PlexAuthZ/api/v1/types"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// CreatePermissionRequest is request model for creating permission.
//
// swagger:parameters createPermissionRequest
type CreatePermissionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// in: path
	OrganizationId string `protobuf:"bytes,1,opt,name=organization_id,json=organizationId,proto3" json:"organization_id,omitempty"`
	// in: path
	Namespace string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// Scope for permission.
	// in: body
	Scope string `protobuf:"bytes,3,opt,name=scope,proto3" json:"scope,omitempty"`
	// Actions that can be performed.
	// in: body
	Actions []string `protobuf:"bytes,4,rep,name=actions,proto3" json:"actions,omitempty"`
	// Resource for the action.
	// in: body
	ResourceId string `protobuf:"bytes,5,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	// Effect Permitted or Denied
	// in: body
	Effect types.Effect `protobuf:"varint,6,opt,name=effect,proto3,enum=api.authz.types.Effect" json:"effect,omitempty"`
	// Constraints expression with dynamic properties.
	// in: body
	Constraints string `protobuf:"bytes,7,opt,name=constraints,proto3" json:"constraints,omitempty"`
}

func (x *CreatePermissionRequest) Reset() {
	*x = CreatePermissionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_services_permission_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreatePermissionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePermissionRequest) ProtoMessage() {}

func (x *CreatePermissionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_services_permission_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePermissionRequest.ProtoReflect.Descriptor instead.
func (*CreatePermissionRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_services_permission_service_proto_rawDescGZIP(), []int{0}
}

func (x *CreatePermissionRequest) GetOrganizationId() string {
	if x != nil {
		return x.OrganizationId
	}
	return ""
}

func (x *CreatePermissionRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *CreatePermissionRequest) GetScope() string {
	if x != nil {
		return x.Scope
	}
	return ""
}

func (x *CreatePermissionRequest) GetActions() []string {
	if x != nil {
		return x.Actions
	}
	return nil
}

func (x *CreatePermissionRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *CreatePermissionRequest) GetEffect() types.Effect {
	if x != nil {
		return x.Effect
	}
	return types.Effect(0)
}

func (x *CreatePermissionRequest) GetConstraints() string {
	if x != nil {
		return x.Constraints
	}
	return ""
}

// CreatePermissionResponse is response model for creating permission.
//
// swagger:parameters createPermissionResponse
type CreatePermissionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ID unique identifier assigned to this permission.
	// in: body
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *CreatePermissionResponse) Reset() {
	*x = CreatePermissionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_services_permission_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreatePermissionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePermissionResponse) ProtoMessage() {}

func (x *CreatePermissionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_services_permission_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePermissionResponse.ProtoReflect.Descriptor instead.
func (*CreatePermissionResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_services_permission_service_proto_rawDescGZIP(), []int{1}
}

func (x *CreatePermissionResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// UpdatePermissionRequest is request model for updating permission.
//
// swagger:parameters updatePermissionRequest
type UpdatePermissionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// in: path
	OrganizationId string `protobuf:"bytes,1,opt,name=organization_id,json=organizationId,proto3" json:"organization_id,omitempty"`
	// in: path
	Namespace string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// in: path
	Id string `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	// Version
	// in: body
	Version int64 `protobuf:"varint,4,opt,name=version,proto3" json:"version,omitempty"`
	// Scope for permission.
	// in: body
	Scope string `protobuf:"bytes,5,opt,name=scope,proto3" json:"scope,omitempty"`
	// Actions that can be performed.
	// in: body
	Actions []string `protobuf:"bytes,6,rep,name=actions,proto3" json:"actions,omitempty"`
	// Resource for the action.
	// in: body
	ResourceId string `protobuf:"bytes,7,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	// Effect Permitted or Denied
	// in: body
	Effect types.Effect `protobuf:"varint,8,opt,name=effect,proto3,enum=api.authz.types.Effect" json:"effect,omitempty"`
	// Constraints expression with dynamic properties.
	// in: body
	Constraints string `protobuf:"bytes,9,opt,name=constraints,proto3" json:"constraints,omitempty"`
}

func (x *UpdatePermissionRequest) Reset() {
	*x = UpdatePermissionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_services_permission_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdatePermissionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePermissionRequest) ProtoMessage() {}

func (x *UpdatePermissionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_services_permission_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePermissionRequest.ProtoReflect.Descriptor instead.
func (*UpdatePermissionRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_services_permission_service_proto_rawDescGZIP(), []int{2}
}

func (x *UpdatePermissionRequest) GetOrganizationId() string {
	if x != nil {
		return x.OrganizationId
	}
	return ""
}

func (x *UpdatePermissionRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *UpdatePermissionRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdatePermissionRequest) GetVersion() int64 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *UpdatePermissionRequest) GetScope() string {
	if x != nil {
		return x.Scope
	}
	return ""
}

func (x *UpdatePermissionRequest) GetActions() []string {
	if x != nil {
		return x.Actions
	}
	return nil
}

func (x *UpdatePermissionRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *UpdatePermissionRequest) GetEffect() types.Effect {
	if x != nil {
		return x.Effect
	}
	return types.Effect(0)
}

func (x *UpdatePermissionRequest) GetConstraints() string {
	if x != nil {
		return x.Constraints
	}
	return ""
}

// UpdatePermissionResponse is response model for updating permission.
//
// swagger:parameters updatePermissionResponse
type UpdatePermissionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UpdatePermissionResponse) Reset() {
	*x = UpdatePermissionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_services_permission_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdatePermissionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePermissionResponse) ProtoMessage() {}

func (x *UpdatePermissionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_services_permission_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePermissionResponse.ProtoReflect.Descriptor instead.
func (*UpdatePermissionResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_services_permission_service_proto_rawDescGZIP(), []int{3}
}

// DeletePermissionRequest is request model for deleting permission.
//
// swagger:parameters deletePermissionRequest
type DeletePermissionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// in: path
	OrganizationId string `protobuf:"bytes,1,opt,name=organization_id,json=organizationId,proto3" json:"organization_id,omitempty"`
	// in: path
	Namespace string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// in: path
	Id string `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeletePermissionRequest) Reset() {
	*x = DeletePermissionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_services_permission_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeletePermissionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeletePermissionRequest) ProtoMessage() {}

func (x *DeletePermissionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_services_permission_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeletePermissionRequest.ProtoReflect.Descriptor instead.
func (*DeletePermissionRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_services_permission_service_proto_rawDescGZIP(), []int{4}
}

func (x *DeletePermissionRequest) GetOrganizationId() string {
	if x != nil {
		return x.OrganizationId
	}
	return ""
}

func (x *DeletePermissionRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *DeletePermissionRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// DeletePermissionResponse is response model for deleting permission by id.
//
// swagger:parameters deletePermissionResponse
type DeletePermissionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeletePermissionResponse) Reset() {
	*x = DeletePermissionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_services_permission_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeletePermissionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeletePermissionResponse) ProtoMessage() {}

func (x *DeletePermissionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_services_permission_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeletePermissionResponse.ProtoReflect.Descriptor instead.
func (*DeletePermissionResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_services_permission_service_proto_rawDescGZIP(), []int{5}
}

// QueryPermissionRequest is request model for querying permissions.
//
// swagger:parameters queryPermissionRequest
type QueryPermissionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// in: path
	OrganizationId string `protobuf:"bytes,1,opt,name=organization_id,json=organizationId,proto3" json:"organization_id,omitempty"`
	// in: path
	Namespace string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// Name of the permission.
	// in:query
	Predicates map[string]string `protobuf:"bytes,3,rep,name=predicates,proto3" json:"predicates,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// in: query
	Offset string `protobuf:"bytes,4,opt,name=offset,proto3" json:"offset,omitempty"`
	// in: query
	Limit int64 `protobuf:"varint,5,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *QueryPermissionRequest) Reset() {
	*x = QueryPermissionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_services_permission_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryPermissionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryPermissionRequest) ProtoMessage() {}

func (x *QueryPermissionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_services_permission_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryPermissionRequest.ProtoReflect.Descriptor instead.
func (*QueryPermissionRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_services_permission_service_proto_rawDescGZIP(), []int{6}
}

func (x *QueryPermissionRequest) GetOrganizationId() string {
	if x != nil {
		return x.OrganizationId
	}
	return ""
}

func (x *QueryPermissionRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *QueryPermissionRequest) GetPredicates() map[string]string {
	if x != nil {
		return x.Predicates
	}
	return nil
}

func (x *QueryPermissionRequest) GetOffset() string {
	if x != nil {
		return x.Offset
	}
	return ""
}

func (x *QueryPermissionRequest) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

// QueryPermissionResponse is response model for querying permission.
//
// swagger:parameters queryPermissionResponse
type QueryPermissionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ID unique identifier assigned to this permission.
	// in: body
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Version
	// in: body
	Version int64 `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
	// Namespace of Permission.
	// in: Body
	Namespace string `protobuf:"bytes,3,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// Scope for permission.
	// in: body
	Scope string `protobuf:"bytes,4,opt,name=scope,proto3" json:"scope,omitempty"`
	// Actions that can be performed.
	// in: body
	Actions []string `protobuf:"bytes,5,rep,name=actions,proto3" json:"actions,omitempty"`
	// Resource for the action.
	// in: body
	ResourceId string `protobuf:"bytes,6,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	// Effect Permitted or Denied
	// in: body
	Effect types.Effect `protobuf:"varint,7,opt,name=effect,proto3,enum=api.authz.types.Effect" json:"effect,omitempty"`
	// Constraints expression with dynamic properties.
	// in: body
	Constraints string `protobuf:"bytes,8,opt,name=constraints,proto3" json:"constraints,omitempty"`
	// in: body
	NextOffset string `protobuf:"bytes,9,opt,name=next_offset,json=nextOffset,proto3" json:"next_offset,omitempty"`
	// Created date
	// in: body
	Created *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=created,proto3" json:"created,omitempty"`
	// Updated date
	// in: body
	Updated *timestamppb.Timestamp `protobuf:"bytes,11,opt,name=updated,proto3" json:"updated,omitempty"`
}

func (x *QueryPermissionResponse) Reset() {
	*x = QueryPermissionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_services_permission_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryPermissionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryPermissionResponse) ProtoMessage() {}

func (x *QueryPermissionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_services_permission_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryPermissionResponse.ProtoReflect.Descriptor instead.
func (*QueryPermissionResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_services_permission_service_proto_rawDescGZIP(), []int{7}
}

func (x *QueryPermissionResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *QueryPermissionResponse) GetVersion() int64 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *QueryPermissionResponse) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *QueryPermissionResponse) GetScope() string {
	if x != nil {
		return x.Scope
	}
	return ""
}

func (x *QueryPermissionResponse) GetActions() []string {
	if x != nil {
		return x.Actions
	}
	return nil
}

func (x *QueryPermissionResponse) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *QueryPermissionResponse) GetEffect() types.Effect {
	if x != nil {
		return x.Effect
	}
	return types.Effect(0)
}

func (x *QueryPermissionResponse) GetConstraints() string {
	if x != nil {
		return x.Constraints
	}
	return ""
}

func (x *QueryPermissionResponse) GetNextOffset() string {
	if x != nil {
		return x.NextOffset
	}
	return ""
}

func (x *QueryPermissionResponse) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

func (x *QueryPermissionResponse) GetUpdated() *timestamppb.Timestamp {
	if x != nil {
		return x.Updated
	}
	return nil
}

var File_api_v1_services_permission_service_proto protoreflect.FileDescriptor

var file_api_v1_services_permission_service_proto_rawDesc = []byte{
	0x0a, 0x28, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x2f, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x61, 0x70, 0x69, 0x2e,
	0x61, 0x75, 0x74, 0x68, 0x7a, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x1a, 0x18,
	0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x61, 0x75, 0x74,
	0x68, 0x7a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x84, 0x02, 0x0a, 0x17, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e,
	0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1c,
	0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x73, 0x63, 0x6f, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x63, 0x6f,
	0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x07, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1f, 0x0a, 0x0b,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x2f, 0x0a,
	0x06, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e,
	0x45, 0x66, 0x66, 0x65, 0x63, 0x74, 0x52, 0x06, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x12, 0x20,
	0x0a, 0x0b, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x73,
	0x22, 0x2a, 0x0a, 0x18, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0xae, 0x02, 0x0a,
	0x17, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x6f, 0x72, 0x67, 0x61,
	0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0e, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49,
	0x64, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f,
	0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x07, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x2f, 0x0a, 0x06, 0x65, 0x66,
	0x66, 0x65, 0x63, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x45, 0x66, 0x66,
	0x65, 0x63, 0x74, 0x52, 0x06, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x63,
	0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x73, 0x22, 0x1a, 0x0a,
	0x18, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x70, 0x0a, 0x17, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6f,
	0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1c, 0x0a,
	0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x1a, 0x0a, 0x18, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xa8, 0x02, 0x0a, 0x16, 0x51, 0x75, 0x65, 0x72,
	0x79, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6f, 0x72, 0x67,
	0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x6e,
	0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x5a, 0x0a, 0x0a, 0x70, 0x72, 0x65,
	0x64, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x3a, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63,
	0x61, 0x74, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x70, 0x72, 0x65, 0x64, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x69,
	0x6d, 0x69, 0x74, 0x1a, 0x3d, 0x0a, 0x0f, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x92, 0x03, 0x0a, 0x17, 0x51, 0x75, 0x65, 0x72, 0x79, 0x50, 0x65, 0x72, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18,
	0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d,
	0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x2f, 0x0a, 0x06, 0x65, 0x66, 0x66, 0x65, 0x63,
	0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x75,
	0x74, 0x68, 0x7a, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x45, 0x66, 0x66, 0x65, 0x63, 0x74,
	0x52, 0x06, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6f, 0x6e, 0x73,
	0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63,
	0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x6e, 0x65,
	0x78, 0x74, 0x5f, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x6e, 0x65, 0x78, 0x74, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x34, 0x0a, 0x07, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x12, 0x34, 0x0a, 0x07, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x32, 0xa7, 0x03, 0x0a, 0x12, 0x50, 0x65, 0x72, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x63,
	0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x2b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61,
	0x75, 0x74, 0x68, 0x7a, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74, 0x68,
	0x7a, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x63, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x2b, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x62, 0x0a, 0x05, 0x51, 0x75, 0x65, 0x72,
	0x79, 0x12, 0x2a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x50, 0x65, 0x72, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x63, 0x0a, 0x06,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x2b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74,
	0x68, 0x7a, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x50,
	0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x62, 0x68, 0x61, 0x74, 0x74, 0x69, 0x2f, 0x50, 0x6c, 0x65, 0x78, 0x41, 0x75, 0x74, 0x68, 0x5a,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_services_permission_service_proto_rawDescOnce sync.Once
	file_api_v1_services_permission_service_proto_rawDescData = file_api_v1_services_permission_service_proto_rawDesc
)

func file_api_v1_services_permission_service_proto_rawDescGZIP() []byte {
	file_api_v1_services_permission_service_proto_rawDescOnce.Do(func() {
		file_api_v1_services_permission_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_services_permission_service_proto_rawDescData)
	})
	return file_api_v1_services_permission_service_proto_rawDescData
}

var file_api_v1_services_permission_service_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_api_v1_services_permission_service_proto_goTypes = []interface{}{
	(*CreatePermissionRequest)(nil),  // 0: api.authz.services.CreatePermissionRequest
	(*CreatePermissionResponse)(nil), // 1: api.authz.services.CreatePermissionResponse
	(*UpdatePermissionRequest)(nil),  // 2: api.authz.services.UpdatePermissionRequest
	(*UpdatePermissionResponse)(nil), // 3: api.authz.services.UpdatePermissionResponse
	(*DeletePermissionRequest)(nil),  // 4: api.authz.services.DeletePermissionRequest
	(*DeletePermissionResponse)(nil), // 5: api.authz.services.DeletePermissionResponse
	(*QueryPermissionRequest)(nil),   // 6: api.authz.services.QueryPermissionRequest
	(*QueryPermissionResponse)(nil),  // 7: api.authz.services.QueryPermissionResponse
	nil,                              // 8: api.authz.services.QueryPermissionRequest.PredicatesEntry
	(types.Effect)(0),                // 9: api.authz.types.Effect
	(*timestamppb.Timestamp)(nil),    // 10: google.protobuf.Timestamp
}
var file_api_v1_services_permission_service_proto_depIdxs = []int32{
	9,  // 0: api.authz.services.CreatePermissionRequest.effect:type_name -> api.authz.types.Effect
	9,  // 1: api.authz.services.UpdatePermissionRequest.effect:type_name -> api.authz.types.Effect
	8,  // 2: api.authz.services.QueryPermissionRequest.predicates:type_name -> api.authz.services.QueryPermissionRequest.PredicatesEntry
	9,  // 3: api.authz.services.QueryPermissionResponse.effect:type_name -> api.authz.types.Effect
	10, // 4: api.authz.services.QueryPermissionResponse.created:type_name -> google.protobuf.Timestamp
	10, // 5: api.authz.services.QueryPermissionResponse.updated:type_name -> google.protobuf.Timestamp
	0,  // 6: api.authz.services.PermissionsService.Create:input_type -> api.authz.services.CreatePermissionRequest
	2,  // 7: api.authz.services.PermissionsService.Update:input_type -> api.authz.services.UpdatePermissionRequest
	6,  // 8: api.authz.services.PermissionsService.Query:input_type -> api.authz.services.QueryPermissionRequest
	4,  // 9: api.authz.services.PermissionsService.Delete:input_type -> api.authz.services.DeletePermissionRequest
	1,  // 10: api.authz.services.PermissionsService.Create:output_type -> api.authz.services.CreatePermissionResponse
	3,  // 11: api.authz.services.PermissionsService.Update:output_type -> api.authz.services.UpdatePermissionResponse
	7,  // 12: api.authz.services.PermissionsService.Query:output_type -> api.authz.services.QueryPermissionResponse
	5,  // 13: api.authz.services.PermissionsService.Delete:output_type -> api.authz.services.DeletePermissionResponse
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_api_v1_services_permission_service_proto_init() }
func file_api_v1_services_permission_service_proto_init() {
	if File_api_v1_services_permission_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1_services_permission_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreatePermissionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_services_permission_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreatePermissionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_services_permission_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdatePermissionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_services_permission_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdatePermissionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_services_permission_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeletePermissionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_services_permission_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeletePermissionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_services_permission_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryPermissionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_services_permission_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryPermissionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_v1_services_permission_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_v1_services_permission_service_proto_goTypes,
		DependencyIndexes: file_api_v1_services_permission_service_proto_depIdxs,
		MessageInfos:      file_api_v1_services_permission_service_proto_msgTypes,
	}.Build()
	File_api_v1_services_permission_service_proto = out.File
	file_api_v1_services_permission_service_proto_rawDesc = nil
	file_api_v1_services_permission_service_proto_goTypes = nil
	file_api_v1_services_permission_service_proto_depIdxs = nil
}
