// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.2
// source: goim_client.proto

package pb

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

type HandshakeStatus int32

const (
	HandshakeStatus_Success      HandshakeStatus = 0
	HandshakeStatus_AccessDenied HandshakeStatus = 1
)

// Enum value maps for HandshakeStatus.
var (
	HandshakeStatus_name = map[int32]string{
		0: "Success",
		1: "AccessDenied",
	}
	HandshakeStatus_value = map[string]int32{
		"Success":      0,
		"AccessDenied": 1,
	}
)

func (x HandshakeStatus) Enum() *HandshakeStatus {
	p := new(HandshakeStatus)
	*p = x
	return p
}

func (x HandshakeStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (HandshakeStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_goim_client_proto_enumTypes[0].Descriptor()
}

func (HandshakeStatus) Type() protoreflect.EnumType {
	return &file_goim_client_proto_enumTypes[0]
}

func (x HandshakeStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use HandshakeStatus.Descriptor instead.
func (HandshakeStatus) EnumDescriptor() ([]byte, []int) {
	return file_goim_client_proto_rawDescGZIP(), []int{0}
}

type HandshakeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *HandshakeRequest) Reset() {
	*x = HandshakeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goim_client_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandshakeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandshakeRequest) ProtoMessage() {}

func (x *HandshakeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_goim_client_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HandshakeRequest.ProtoReflect.Descriptor instead.
func (*HandshakeRequest) Descriptor() ([]byte, []int) {
	return file_goim_client_proto_rawDescGZIP(), []int{0}
}

func (x *HandshakeRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type HandshakeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status HandshakeStatus `protobuf:"varint,1,opt,name=status,proto3,enum=goim_client.HandshakeStatus" json:"status,omitempty"`
}

func (x *HandshakeResponse) Reset() {
	*x = HandshakeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goim_client_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandshakeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandshakeResponse) ProtoMessage() {}

func (x *HandshakeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_goim_client_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HandshakeResponse.ProtoReflect.Descriptor instead.
func (*HandshakeResponse) Descriptor() ([]byte, []int) {
	return file_goim_client_proto_rawDescGZIP(), []int{1}
}

func (x *HandshakeResponse) GetStatus() HandshakeStatus {
	if x != nil {
		return x.Status
	}
	return HandshakeStatus_Success
}

var File_goim_client_proto protoreflect.FileDescriptor

var file_goim_client_proto_rawDesc = []byte{
	0x0a, 0x11, 0x67, 0x6f, 0x69, 0x6d, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x67, 0x6f, 0x69, 0x6d, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x22, 0x28, 0x0a, 0x10, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x49, 0x0a, 0x11, 0x48, 0x61,
	0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x34, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x1c, 0x2e, 0x67, 0x6f, 0x69, 0x6d, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x48, 0x61,
	0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x2a, 0x30, 0x0a, 0x0f, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61,
	0x6b, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x44,
	0x65, 0x6e, 0x69, 0x65, 0x64, 0x10, 0x01, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_goim_client_proto_rawDescOnce sync.Once
	file_goim_client_proto_rawDescData = file_goim_client_proto_rawDesc
)

func file_goim_client_proto_rawDescGZIP() []byte {
	file_goim_client_proto_rawDescOnce.Do(func() {
		file_goim_client_proto_rawDescData = protoimpl.X.CompressGZIP(file_goim_client_proto_rawDescData)
	})
	return file_goim_client_proto_rawDescData
}

var file_goim_client_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_goim_client_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_goim_client_proto_goTypes = []interface{}{
	(HandshakeStatus)(0),      // 0: goim_client.HandshakeStatus
	(*HandshakeRequest)(nil),  // 1: goim_client.HandshakeRequest
	(*HandshakeResponse)(nil), // 2: goim_client.HandshakeResponse
}
var file_goim_client_proto_depIdxs = []int32{
	0, // 0: goim_client.HandshakeResponse.status:type_name -> goim_client.HandshakeStatus
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_goim_client_proto_init() }
func file_goim_client_proto_init() {
	if File_goim_client_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_goim_client_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HandshakeRequest); i {
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
		file_goim_client_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HandshakeResponse); i {
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
			RawDescriptor: file_goim_client_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_goim_client_proto_goTypes,
		DependencyIndexes: file_goim_client_proto_depIdxs,
		EnumInfos:         file_goim_client_proto_enumTypes,
		MessageInfos:      file_goim_client_proto_msgTypes,
	}.Build()
	File_goim_client_proto = out.File
	file_goim_client_proto_rawDesc = nil
	file_goim_client_proto_goTypes = nil
	file_goim_client_proto_depIdxs = nil
}
