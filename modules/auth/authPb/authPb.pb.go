//Version

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.23.3
// source: modules/auth/authPb/authPb.proto

package Go_MicroService

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Structures
type AccessToKenSearchReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccessToken string `protobuf:"bytes,1,opt,name=accessToken,proto3" json:"accessToken,omitempty"`
}

func (x *AccessToKenSearchReq) Reset() {
	*x = AccessToKenSearchReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_modules_auth_authPb_authPb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessToKenSearchReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessToKenSearchReq) ProtoMessage() {}

func (x *AccessToKenSearchReq) ProtoReflect() protoreflect.Message {
	mi := &file_modules_auth_authPb_authPb_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessToKenSearchReq.ProtoReflect.Descriptor instead.
func (*AccessToKenSearchReq) Descriptor() ([]byte, []int) {
	return file_modules_auth_authPb_authPb_proto_rawDescGZIP(), []int{0}
}

func (x *AccessToKenSearchReq) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

type AccessToKenSearchRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsValid bool `protobuf:"varint,1,opt,name=isValid,proto3" json:"isValid,omitempty"`
}

func (x *AccessToKenSearchRes) Reset() {
	*x = AccessToKenSearchRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_modules_auth_authPb_authPb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessToKenSearchRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessToKenSearchRes) ProtoMessage() {}

func (x *AccessToKenSearchRes) ProtoReflect() protoreflect.Message {
	mi := &file_modules_auth_authPb_authPb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessToKenSearchRes.ProtoReflect.Descriptor instead.
func (*AccessToKenSearchRes) Descriptor() ([]byte, []int) {
	return file_modules_auth_authPb_authPb_proto_rawDescGZIP(), []int{1}
}

func (x *AccessToKenSearchRes) GetIsValid() bool {
	if x != nil {
		return x.IsValid
	}
	return false
}

type RoleCountReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RoleCountReq) Reset() {
	*x = RoleCountReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_modules_auth_authPb_authPb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoleCountReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoleCountReq) ProtoMessage() {}

func (x *RoleCountReq) ProtoReflect() protoreflect.Message {
	mi := &file_modules_auth_authPb_authPb_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoleCountReq.ProtoReflect.Descriptor instead.
func (*RoleCountReq) Descriptor() ([]byte, []int) {
	return file_modules_auth_authPb_authPb_proto_rawDescGZIP(), []int{2}
}

type RoleCountRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count int64 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *RoleCountRes) Reset() {
	*x = RoleCountRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_modules_auth_authPb_authPb_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoleCountRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoleCountRes) ProtoMessage() {}

func (x *RoleCountRes) ProtoReflect() protoreflect.Message {
	mi := &file_modules_auth_authPb_authPb_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoleCountRes.ProtoReflect.Descriptor instead.
func (*RoleCountRes) Descriptor() ([]byte, []int) {
	return file_modules_auth_authPb_authPb_proto_rawDescGZIP(), []int{3}
}

func (x *RoleCountRes) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

var File_modules_auth_authPb_authPb_proto protoreflect.FileDescriptor

var file_modules_auth_authPb_authPb_proto_rawDesc = []byte{
	0x0a, 0x20, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x61,
	0x75, 0x74, 0x68, 0x50, 0x62, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x50, 0x62, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x38, 0x0a, 0x14, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x4b, 0x65,
	0x6e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x71, 0x12, 0x20, 0x0a, 0x0b, 0x61, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x30, 0x0a, 0x14,
	0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x4b, 0x65, 0x6e, 0x53, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x52, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x69, 0x73, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x22, 0x0e,
	0x0a, 0x0c, 0x52, 0x6f, 0x6c, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x22, 0x24,
	0x0a, 0x0c, 0x52, 0x6f, 0x6c, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x12, 0x14,
	0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x32, 0x7f, 0x0a, 0x0f, 0x41, 0x75, 0x74, 0x68, 0x47, 0x72, 0x70, 0x63,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x41, 0x0a, 0x11, 0x41, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x54, 0x6f, 0x4b, 0x65, 0x6e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x12, 0x15, 0x2e, 0x41,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x4b, 0x65, 0x6e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x52, 0x65, 0x71, 0x1a, 0x15, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x4b, 0x65,
	0x6e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x73, 0x12, 0x29, 0x0a, 0x09, 0x52, 0x6f,
	0x6c, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x0d, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x0d, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x52, 0x65, 0x73, 0x42, 0x26, 0x5a, 0x24, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x75, 0x61, 0x74, 0x6f, 0x6d, 0x39, 0x39, 0x39, 0x2f, 0x47, 0x6f,
	0x2d, 0x4d, 0x69, 0x63, 0x72, 0x6f, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_modules_auth_authPb_authPb_proto_rawDescOnce sync.Once
	file_modules_auth_authPb_authPb_proto_rawDescData = file_modules_auth_authPb_authPb_proto_rawDesc
)

func file_modules_auth_authPb_authPb_proto_rawDescGZIP() []byte {
	file_modules_auth_authPb_authPb_proto_rawDescOnce.Do(func() {
		file_modules_auth_authPb_authPb_proto_rawDescData = protoimpl.X.CompressGZIP(file_modules_auth_authPb_authPb_proto_rawDescData)
	})
	return file_modules_auth_authPb_authPb_proto_rawDescData
}

var file_modules_auth_authPb_authPb_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_modules_auth_authPb_authPb_proto_goTypes = []interface{}{
	(*AccessToKenSearchReq)(nil), // 0: AccessToKenSearchReq
	(*AccessToKenSearchRes)(nil), // 1: AccessToKenSearchRes
	(*RoleCountReq)(nil),         // 2: RoleCountReq
	(*RoleCountRes)(nil),         // 3: RoleCountRes
}
var file_modules_auth_authPb_authPb_proto_depIdxs = []int32{
	0, // 0: AuthGrpcService.AccessToKenSearch:input_type -> AccessToKenSearchReq
	2, // 1: AuthGrpcService.RoleCount:input_type -> RoleCountReq
	1, // 2: AuthGrpcService.AccessToKenSearch:output_type -> AccessToKenSearchRes
	3, // 3: AuthGrpcService.RoleCount:output_type -> RoleCountRes
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_modules_auth_authPb_authPb_proto_init() }
func file_modules_auth_authPb_authPb_proto_init() {
	if File_modules_auth_authPb_authPb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_modules_auth_authPb_authPb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessToKenSearchReq); i {
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
		file_modules_auth_authPb_authPb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessToKenSearchRes); i {
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
		file_modules_auth_authPb_authPb_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoleCountReq); i {
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
		file_modules_auth_authPb_authPb_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoleCountRes); i {
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
			RawDescriptor: file_modules_auth_authPb_authPb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_modules_auth_authPb_authPb_proto_goTypes,
		DependencyIndexes: file_modules_auth_authPb_authPb_proto_depIdxs,
		MessageInfos:      file_modules_auth_authPb_authPb_proto_msgTypes,
	}.Build()
	File_modules_auth_authPb_authPb_proto = out.File
	file_modules_auth_authPb_authPb_proto_rawDesc = nil
	file_modules_auth_authPb_authPb_proto_goTypes = nil
	file_modules_auth_authPb_authPb_proto_depIdxs = nil
}
