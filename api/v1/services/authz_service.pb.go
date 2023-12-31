// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.1
// source: api/v1/services/authz_service.proto

package services

import (
	types "github.com/bhatti/PlexAuthZ/api/v1/types"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// AuthRequest is request model for authorization access API.
//
// swagger:parameters authRequest
type AuthRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// in: path
	OrganizationId string `protobuf:"bytes,1,opt,name=organization_id,json=organizationId,proto3" json:"organization_id,omitempty"`
	// in: path
	Namespace string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// in: path
	PrincipalId string `protobuf:"bytes,3,opt,name=principal_id,json=principalId,proto3" json:"principal_id,omitempty"`
	// in: body
	Action string `protobuf:"bytes,4,opt,name=action,proto3" json:"action,omitempty"`
	// in: body
	Resource string `protobuf:"bytes,5,opt,name=resource,proto3" json:"resource,omitempty"`
	// in: body
	Scope string `protobuf:"bytes,6,opt,name=scope,proto3" json:"scope,omitempty"`
	// in: body
	Context map[string]string `protobuf:"bytes,7,rep,name=context,proto3" json:"context,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *AuthRequest) Reset() {
	*x = AuthRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_services_authz_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthRequest) ProtoMessage() {}

func (x *AuthRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_services_authz_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthRequest.ProtoReflect.Descriptor instead.
func (*AuthRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_services_authz_service_proto_rawDescGZIP(), []int{0}
}

func (x *AuthRequest) GetOrganizationId() string {
	if x != nil {
		return x.OrganizationId
	}
	return ""
}

func (x *AuthRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *AuthRequest) GetPrincipalId() string {
	if x != nil {
		return x.PrincipalId
	}
	return ""
}

func (x *AuthRequest) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *AuthRequest) GetResource() string {
	if x != nil {
		return x.Resource
	}
	return ""
}

func (x *AuthRequest) GetScope() string {
	if x != nil {
		return x.Scope
	}
	return ""
}

func (x *AuthRequest) GetContext() map[string]string {
	if x != nil {
		return x.Context
	}
	return nil
}

// AuthResponse is response model for authorization access API.
//
// swagger:parameters authResponse
type AuthResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// in: body
	Effect types.Effect `protobuf:"varint,1,opt,name=effect,proto3,enum=api.authz.types.Effect" json:"effect,omitempty"`
	// in: body
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *AuthResponse) Reset() {
	*x = AuthResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_services_authz_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthResponse) ProtoMessage() {}

func (x *AuthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_services_authz_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthResponse.ProtoReflect.Descriptor instead.
func (*AuthResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_services_authz_service_proto_rawDescGZIP(), []int{1}
}

func (x *AuthResponse) GetEffect() types.Effect {
	if x != nil {
		return x.Effect
	}
	return types.Effect(0)
}

func (x *AuthResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// CheckConstraintsRequest is request model for checking constraints and authorization access API.
//
// swagger:parameters checkConstraintsRequest
type CheckConstraintsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// in: path
	OrganizationId string `protobuf:"bytes,1,opt,name=organization_id,json=organizationId,proto3" json:"organization_id,omitempty"`
	// in: path
	Namespace string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// in: path
	PrincipalId string `protobuf:"bytes,3,opt,name=principal_id,json=principalId,proto3" json:"principal_id,omitempty"`
	// in: body
	Constraints string `protobuf:"bytes,4,opt,name=constraints,proto3" json:"constraints,omitempty"`
	// in: body
	Context map[string]string `protobuf:"bytes,5,rep,name=context,proto3" json:"context,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *CheckConstraintsRequest) Reset() {
	*x = CheckConstraintsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_services_authz_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckConstraintsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckConstraintsRequest) ProtoMessage() {}

func (x *CheckConstraintsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_services_authz_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckConstraintsRequest.ProtoReflect.Descriptor instead.
func (*CheckConstraintsRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_services_authz_service_proto_rawDescGZIP(), []int{2}
}

func (x *CheckConstraintsRequest) GetOrganizationId() string {
	if x != nil {
		return x.OrganizationId
	}
	return ""
}

func (x *CheckConstraintsRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *CheckConstraintsRequest) GetPrincipalId() string {
	if x != nil {
		return x.PrincipalId
	}
	return ""
}

func (x *CheckConstraintsRequest) GetConstraints() string {
	if x != nil {
		return x.Constraints
	}
	return ""
}

func (x *CheckConstraintsRequest) GetContext() map[string]string {
	if x != nil {
		return x.Context
	}
	return nil
}

// CheckConstraintsResponse is response model for checking constraints and authorization access API.
//
// swagger:parameters checkConstraintsResponse
type CheckConstraintsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// in: body
	Matched bool `protobuf:"varint,1,opt,name=matched,proto3" json:"matched,omitempty"`
	// in: body
	Output string `protobuf:"bytes,2,opt,name=output,proto3" json:"output,omitempty"`
}

func (x *CheckConstraintsResponse) Reset() {
	*x = CheckConstraintsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_services_authz_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckConstraintsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckConstraintsResponse) ProtoMessage() {}

func (x *CheckConstraintsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_services_authz_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckConstraintsResponse.ProtoReflect.Descriptor instead.
func (*CheckConstraintsResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_services_authz_service_proto_rawDescGZIP(), []int{3}
}

func (x *CheckConstraintsResponse) GetMatched() bool {
	if x != nil {
		return x.Matched
	}
	return false
}

func (x *CheckConstraintsResponse) GetOutput() string {
	if x != nil {
		return x.Output
	}
	return ""
}

// AllocateResourceRequest is request model for allocating resource.
//
// swagger:parameters allocateResourceRequest
type AllocateResourceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// in: path
	OrganizationId string `protobuf:"bytes,1,opt,name=organization_id,json=organizationId,proto3" json:"organization_id,omitempty"`
	// in: path
	Namespace string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// Resource Id.
	// in: path
	ResourceId string `protobuf:"bytes,3,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	// Principal Id.
	// in: path
	PrincipalId string `protobuf:"bytes,4,opt,name=principal_id,json=principalId,proto3" json:"principal_id,omitempty"`
	// Constraints of allocation.
	// in: body
	Constraints string `protobuf:"bytes,5,opt,name=constraints,proto3" json:"constraints,omitempty"`
	// Time duration in milliseconds after which instance will expire.
	// in: body
	Expiry *durationpb.Duration `protobuf:"bytes,6,opt,name=expiry,proto3" json:"expiry,omitempty"`
	// Context parameters.
	// in: body
	Context map[string]string `protobuf:"bytes,7,rep,name=context,proto3" json:"context,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *AllocateResourceRequest) Reset() {
	*x = AllocateResourceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_services_authz_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllocateResourceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllocateResourceRequest) ProtoMessage() {}

func (x *AllocateResourceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_services_authz_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllocateResourceRequest.ProtoReflect.Descriptor instead.
func (*AllocateResourceRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_services_authz_service_proto_rawDescGZIP(), []int{4}
}

func (x *AllocateResourceRequest) GetOrganizationId() string {
	if x != nil {
		return x.OrganizationId
	}
	return ""
}

func (x *AllocateResourceRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *AllocateResourceRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *AllocateResourceRequest) GetPrincipalId() string {
	if x != nil {
		return x.PrincipalId
	}
	return ""
}

func (x *AllocateResourceRequest) GetConstraints() string {
	if x != nil {
		return x.Constraints
	}
	return ""
}

func (x *AllocateResourceRequest) GetExpiry() *durationpb.Duration {
	if x != nil {
		return x.Expiry
	}
	return nil
}

func (x *AllocateResourceRequest) GetContext() map[string]string {
	if x != nil {
		return x.Context
	}
	return nil
}

// AllocateResourceResponse is response model for allocating resource.
//
// swagger:parameters allocateResourceResponse
type AllocateResourceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AllocateResourceResponse) Reset() {
	*x = AllocateResourceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_services_authz_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllocateResourceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllocateResourceResponse) ProtoMessage() {}

func (x *AllocateResourceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_services_authz_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllocateResourceResponse.ProtoReflect.Descriptor instead.
func (*AllocateResourceResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_services_authz_service_proto_rawDescGZIP(), []int{5}
}

// DeallocateResourceRequest is request model for deallocating resource.
//
// swagger:parameters deallocateResourceRequest
type DeallocateResourceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// in: path
	OrganizationId string `protobuf:"bytes,1,opt,name=organization_id,json=organizationId,proto3" json:"organization_id,omitempty"`
	// in: path
	Namespace string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// Resource Id
	// in: path
	ResourceId string `protobuf:"bytes,3,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	// Principal Id
	// in: path
	PrincipalId string `protobuf:"bytes,4,opt,name=principal_id,json=principalId,proto3" json:"principal_id,omitempty"`
}

func (x *DeallocateResourceRequest) Reset() {
	*x = DeallocateResourceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_services_authz_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeallocateResourceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeallocateResourceRequest) ProtoMessage() {}

func (x *DeallocateResourceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_services_authz_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeallocateResourceRequest.ProtoReflect.Descriptor instead.
func (*DeallocateResourceRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_services_authz_service_proto_rawDescGZIP(), []int{6}
}

func (x *DeallocateResourceRequest) GetOrganizationId() string {
	if x != nil {
		return x.OrganizationId
	}
	return ""
}

func (x *DeallocateResourceRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *DeallocateResourceRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *DeallocateResourceRequest) GetPrincipalId() string {
	if x != nil {
		return x.PrincipalId
	}
	return ""
}

// DeallocateResourceResponse is response model for deallocating resource.
//
// swagger:parameters deallocateResourceResponse
type DeallocateResourceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeallocateResourceResponse) Reset() {
	*x = DeallocateResourceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_services_authz_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeallocateResourceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeallocateResourceResponse) ProtoMessage() {}

func (x *DeallocateResourceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_services_authz_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeallocateResourceResponse.ProtoReflect.Descriptor instead.
func (*DeallocateResourceResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_services_authz_service_proto_rawDescGZIP(), []int{7}
}

var File_api_v1_services_authz_service_proto protoreflect.FileDescriptor

var file_api_v1_services_authz_service_proto_rawDesc = []byte{
	0x0a, 0x23, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x7a,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x1a, 0x18, 0x61, 0x70, 0x69, 0x2f, 0x76,
	0x31, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xc5, 0x02, 0x0a, 0x0b, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6f, 0x72,
	0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09,
	0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72,
	0x69, 0x6e, 0x63, 0x69, 0x70, 0x61, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x70, 0x72, 0x69, 0x6e, 0x63, 0x69, 0x70, 0x61, 0x6c, 0x49, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x12, 0x46, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x78, 0x74, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61,
	0x75, 0x74, 0x68, 0x7a, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x41, 0x75,
	0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78,
	0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x1a,
	0x3a, 0x0a, 0x0c, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x59, 0x0a, 0x0c, 0x41,
	0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f, 0x0a, 0x06, 0x65,
	0x66, 0x66, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x45, 0x66,
	0x66, 0x65, 0x63, 0x74, 0x52, 0x06, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0xb5, 0x02, 0x0a, 0x17, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x43, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6f, 0x72, 0x67,
	0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x6e,
	0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x69,
	0x6e, 0x63, 0x69, 0x70, 0x61, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x70, 0x72, 0x69, 0x6e, 0x63, 0x69, 0x70, 0x61, 0x6c, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b,
	0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x52,
	0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x38, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x43, 0x6f, 0x6e, 0x73, 0x74, 0x72,
	0x61, 0x69, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x43, 0x6f, 0x6e,
	0x74, 0x65, 0x78, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x78, 0x74, 0x1a, 0x3a, 0x0a, 0x0c, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x4c,
	0x0a, 0x18, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x43, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e,
	0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x61,
	0x74, 0x63, 0x68, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x6d, 0x61, 0x74,
	0x63, 0x68, 0x65, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x22, 0x89, 0x03, 0x0a,
	0x17, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x6f, 0x72, 0x67, 0x61,
	0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0e, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49,
	0x64, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12,
	0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64,
	0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x69, 0x6e, 0x63, 0x69, 0x70, 0x61, 0x6c, 0x5f, 0x69, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x69, 0x6e, 0x63, 0x69, 0x70, 0x61,
	0x6c, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e,
	0x74, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72,
	0x61, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x31, 0x0a, 0x06, 0x65, 0x78, 0x70, 0x69, 0x72, 0x79, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x06, 0x65, 0x78, 0x70, 0x69, 0x72, 0x79, 0x12, 0x52, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x78, 0x74, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x38, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x61, 0x75, 0x74, 0x68, 0x7a, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x41,
	0x6c, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x1a, 0x3a, 0x0a, 0x0c,
	0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x1a, 0x0a, 0x18, 0x41, 0x6c, 0x6c, 0x6f,
	0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0xa6, 0x01, 0x0a, 0x19, 0x44, 0x65, 0x61, 0x6c, 0x6c, 0x6f, 0x63,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6f, 0x72, 0x67,
	0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x6e,
	0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72,
	0x69, 0x6e, 0x63, 0x69, 0x70, 0x61, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x70, 0x72, 0x69, 0x6e, 0x63, 0x69, 0x70, 0x61, 0x6c, 0x49, 0x64, 0x22, 0x1c, 0x0a,
	0x1a, 0x44, 0x65, 0x61, 0x6c, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x96, 0x03, 0x0a, 0x0c,
	0x41, 0x75, 0x74, 0x68, 0x5a, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4e, 0x0a, 0x09,
	0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x12, 0x1f, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x61, 0x75, 0x74, 0x68, 0x7a, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x41,
	0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e,
	0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x62, 0x0a, 0x05,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x2b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74, 0x68,
	0x7a, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x43, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x43, 0x6f, 0x6e,
	0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x65, 0x0a, 0x08, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x12, 0x2b, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x2e, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x61, 0x75, 0x74, 0x68, 0x7a, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x41,
	0x6c, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6b, 0x0a, 0x0a, 0x44, 0x65, 0x61, 0x6c, 0x6c,
	0x6f, 0x63, 0x61, 0x74, 0x65, 0x12, 0x2d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74, 0x68,
	0x7a, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x44, 0x65, 0x61, 0x6c, 0x6c,
	0x6f, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x2e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x7a,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x44, 0x65, 0x61, 0x6c, 0x6c, 0x6f,
	0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x62, 0x68, 0x61, 0x74, 0x74, 0x69, 0x2f, 0x50, 0x6c, 0x65, 0x78, 0x41, 0x75,
	0x74, 0x68, 0x5a, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x2f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_services_authz_service_proto_rawDescOnce sync.Once
	file_api_v1_services_authz_service_proto_rawDescData = file_api_v1_services_authz_service_proto_rawDesc
)

func file_api_v1_services_authz_service_proto_rawDescGZIP() []byte {
	file_api_v1_services_authz_service_proto_rawDescOnce.Do(func() {
		file_api_v1_services_authz_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_services_authz_service_proto_rawDescData)
	})
	return file_api_v1_services_authz_service_proto_rawDescData
}

var file_api_v1_services_authz_service_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_api_v1_services_authz_service_proto_goTypes = []interface{}{
	(*AuthRequest)(nil),                // 0: api.authz.services.AuthRequest
	(*AuthResponse)(nil),               // 1: api.authz.services.AuthResponse
	(*CheckConstraintsRequest)(nil),    // 2: api.authz.services.CheckConstraintsRequest
	(*CheckConstraintsResponse)(nil),   // 3: api.authz.services.CheckConstraintsResponse
	(*AllocateResourceRequest)(nil),    // 4: api.authz.services.AllocateResourceRequest
	(*AllocateResourceResponse)(nil),   // 5: api.authz.services.AllocateResourceResponse
	(*DeallocateResourceRequest)(nil),  // 6: api.authz.services.DeallocateResourceRequest
	(*DeallocateResourceResponse)(nil), // 7: api.authz.services.DeallocateResourceResponse
	nil,                                // 8: api.authz.services.AuthRequest.ContextEntry
	nil,                                // 9: api.authz.services.CheckConstraintsRequest.ContextEntry
	nil,                                // 10: api.authz.services.AllocateResourceRequest.ContextEntry
	(types.Effect)(0),                  // 11: api.authz.types.Effect
	(*durationpb.Duration)(nil),        // 12: google.protobuf.Duration
}
var file_api_v1_services_authz_service_proto_depIdxs = []int32{
	8,  // 0: api.authz.services.AuthRequest.context:type_name -> api.authz.services.AuthRequest.ContextEntry
	11, // 1: api.authz.services.AuthResponse.effect:type_name -> api.authz.types.Effect
	9,  // 2: api.authz.services.CheckConstraintsRequest.context:type_name -> api.authz.services.CheckConstraintsRequest.ContextEntry
	12, // 3: api.authz.services.AllocateResourceRequest.expiry:type_name -> google.protobuf.Duration
	10, // 4: api.authz.services.AllocateResourceRequest.context:type_name -> api.authz.services.AllocateResourceRequest.ContextEntry
	0,  // 5: api.authz.services.AuthZService.Authorize:input_type -> api.authz.services.AuthRequest
	2,  // 6: api.authz.services.AuthZService.Check:input_type -> api.authz.services.CheckConstraintsRequest
	4,  // 7: api.authz.services.AuthZService.Allocate:input_type -> api.authz.services.AllocateResourceRequest
	6,  // 8: api.authz.services.AuthZService.Deallocate:input_type -> api.authz.services.DeallocateResourceRequest
	1,  // 9: api.authz.services.AuthZService.Authorize:output_type -> api.authz.services.AuthResponse
	3,  // 10: api.authz.services.AuthZService.Check:output_type -> api.authz.services.CheckConstraintsResponse
	5,  // 11: api.authz.services.AuthZService.Allocate:output_type -> api.authz.services.AllocateResourceResponse
	7,  // 12: api.authz.services.AuthZService.Deallocate:output_type -> api.authz.services.DeallocateResourceResponse
	9,  // [9:13] is the sub-list for method output_type
	5,  // [5:9] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_api_v1_services_authz_service_proto_init() }
func file_api_v1_services_authz_service_proto_init() {
	if File_api_v1_services_authz_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1_services_authz_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthRequest); i {
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
		file_api_v1_services_authz_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthResponse); i {
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
		file_api_v1_services_authz_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckConstraintsRequest); i {
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
		file_api_v1_services_authz_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckConstraintsResponse); i {
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
		file_api_v1_services_authz_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllocateResourceRequest); i {
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
		file_api_v1_services_authz_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllocateResourceResponse); i {
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
		file_api_v1_services_authz_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeallocateResourceRequest); i {
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
		file_api_v1_services_authz_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeallocateResourceResponse); i {
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
			RawDescriptor: file_api_v1_services_authz_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_v1_services_authz_service_proto_goTypes,
		DependencyIndexes: file_api_v1_services_authz_service_proto_depIdxs,
		MessageInfos:      file_api_v1_services_authz_service_proto_msgTypes,
	}.Build()
	File_api_v1_services_authz_service_proto = out.File
	file_api_v1_services_authz_service_proto_rawDesc = nil
	file_api_v1_services_authz_service_proto_goTypes = nil
	file_api_v1_services_authz_service_proto_depIdxs = nil
}
