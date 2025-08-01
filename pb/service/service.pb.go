// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: service/service.proto

package service

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	common "github.com/fathurzoy/go-grpc-ecommerce-be/pb/common"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type HelloWolrdRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HelloWolrdRequest) Reset() {
	*x = HelloWolrdRequest{}
	mi := &file_service_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HelloWolrdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloWolrdRequest) ProtoMessage() {}

func (x *HelloWolrdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloWolrdRequest.ProtoReflect.Descriptor instead.
func (*HelloWolrdRequest) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{0}
}

func (x *HelloWolrdRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type HelloWorldResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Base          *common.BaseResponse   `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HelloWorldResponse) Reset() {
	*x = HelloWorldResponse{}
	mi := &file_service_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HelloWorldResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloWorldResponse) ProtoMessage() {}

func (x *HelloWorldResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloWorldResponse.ProtoReflect.Descriptor instead.
func (*HelloWorldResponse) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{1}
}

func (x *HelloWorldResponse) GetBase() *common.BaseResponse {
	if x != nil {
		return x.Base
	}
	return nil
}

func (x *HelloWorldResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_service_service_proto protoreflect.FileDescriptor

const file_service_service_proto_rawDesc = "" +
	"\n" +
	"\x15service/service.proto\x12\aservice\x1a\x1acommon/base_response.proto\x1a\x1bbuf/validate/validate.proto\"3\n" +
	"\x11HelloWolrdRequest\x12\x1e\n" +
	"\x04name\x18\x01 \x01(\tB\n" +
	"\xbaH\ar\x05\x10\x01\x18\xff\x01R\x04name\"X\n" +
	"\x12HelloWorldResponse\x12(\n" +
	"\x04base\x18\x01 \x01(\v2\x14.common.BaseResponseR\x04base\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage2Z\n" +
	"\x11HelloWorldService\x12E\n" +
	"\n" +
	"HelloWorld\x12\x1a.service.HelloWolrdRequest\x1a\x1b.service.HelloWorldResponseB6Z4github.com/fathurzoy/go-grpc-ecommerce-be/pb/serviceb\x06proto3"

var (
	file_service_service_proto_rawDescOnce sync.Once
	file_service_service_proto_rawDescData []byte
)

func file_service_service_proto_rawDescGZIP() []byte {
	file_service_service_proto_rawDescOnce.Do(func() {
		file_service_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_service_service_proto_rawDesc), len(file_service_service_proto_rawDesc)))
	})
	return file_service_service_proto_rawDescData
}

var file_service_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_service_service_proto_goTypes = []any{
	(*HelloWolrdRequest)(nil),   // 0: service.HelloWolrdRequest
	(*HelloWorldResponse)(nil),  // 1: service.HelloWorldResponse
	(*common.BaseResponse)(nil), // 2: common.BaseResponse
}
var file_service_service_proto_depIdxs = []int32{
	2, // 0: service.HelloWorldResponse.base:type_name -> common.BaseResponse
	0, // 1: service.HelloWorldService.HelloWorld:input_type -> service.HelloWolrdRequest
	1, // 2: service.HelloWorldService.HelloWorld:output_type -> service.HelloWorldResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_service_service_proto_init() }
func file_service_service_proto_init() {
	if File_service_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_service_service_proto_rawDesc), len(file_service_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_service_proto_goTypes,
		DependencyIndexes: file_service_service_proto_depIdxs,
		MessageInfos:      file_service_service_proto_msgTypes,
	}.Build()
	File_service_service_proto = out.File
	file_service_service_proto_goTypes = nil
	file_service_service_proto_depIdxs = nil
}
