// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023 Datadog, Inc.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: fixtures_test.proto

package grpc

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

// The request message containing the user's name.
type FixtureRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *FixtureRequest) Reset() {
	*x = FixtureRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fixtures_test_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FixtureRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FixtureRequest) ProtoMessage() {}

func (x *FixtureRequest) ProtoReflect() protoreflect.Message {
	mi := &file_fixtures_test_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FixtureRequest.ProtoReflect.Descriptor instead.
func (*FixtureRequest) Descriptor() ([]byte, []int) {
	return file_fixtures_test_proto_rawDescGZIP(), []int{0}
}

func (x *FixtureRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// The response message containing the greetings
type FixtureReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *FixtureReply) Reset() {
	*x = FixtureReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fixtures_test_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FixtureReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FixtureReply) ProtoMessage() {}

func (x *FixtureReply) ProtoReflect() protoreflect.Message {
	mi := &file_fixtures_test_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FixtureReply.ProtoReflect.Descriptor instead.
func (*FixtureReply) Descriptor() ([]byte, []int) {
	return file_fixtures_test_proto_rawDescGZIP(), []int{1}
}

func (x *FixtureReply) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_fixtures_test_proto protoreflect.FileDescriptor

var file_fixtures_test_proto_rawDesc = []byte{
	0x0a, 0x13, 0x66, 0x69, 0x78, 0x74, 0x75, 0x72, 0x65, 0x73, 0x5f, 0x74, 0x65, 0x73, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x67, 0x72, 0x70, 0x63, 0x22, 0x24, 0x0a, 0x0e, 0x46,
	0x69, 0x78, 0x74, 0x75, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x22, 0x28, 0x0a, 0x0c, 0x46, 0x69, 0x78, 0x74, 0x75, 0x72, 0x65, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x7b, 0x0a, 0x07, 0x46,
	0x69, 0x78, 0x74, 0x75, 0x72, 0x65, 0x12, 0x32, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x14,
	0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x46, 0x69, 0x78, 0x74, 0x75, 0x72, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x46, 0x69, 0x78, 0x74,
	0x75, 0x72, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x3c, 0x0a, 0x0a, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x14, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x46, 0x69, 0x78, 0x74, 0x75, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12,
	0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x46, 0x69, 0x78, 0x74, 0x75, 0x72, 0x65, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x6c, 0x0a, 0x19, 0x69, 0x6f, 0x2e, 0x67,
	0x72, 0x70, 0x63, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x74, 0x65, 0x73,
	0x74, 0x67, 0x72, 0x70, 0x63, 0x42, 0x0d, 0x54, 0x65, 0x73, 0x74, 0x47, 0x52, 0x50, 0x43, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3e, 0x67, 0x6f, 0x70, 0x6b, 0x67, 0x2e, 0x69, 0x6e,
	0x2f, 0x44, 0x61, 0x74, 0x61, 0x44, 0x6f, 0x67, 0x2f, 0x64, 0x64, 0x2d, 0x74, 0x72, 0x61, 0x63,
	0x65, 0x2d, 0x67, 0x6f, 0x2e, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x69, 0x62, 0x2f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2e, 0x6f, 0x72,
	0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_fixtures_test_proto_rawDescOnce sync.Once
	file_fixtures_test_proto_rawDescData = file_fixtures_test_proto_rawDesc
)

func file_fixtures_test_proto_rawDescGZIP() []byte {
	file_fixtures_test_proto_rawDescOnce.Do(func() {
		file_fixtures_test_proto_rawDescData = protoimpl.X.CompressGZIP(file_fixtures_test_proto_rawDescData)
	})
	return file_fixtures_test_proto_rawDescData
}

var file_fixtures_test_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_fixtures_test_proto_goTypes = []interface{}{
	(*FixtureRequest)(nil), // 0: grpc.FixtureRequest
	(*FixtureReply)(nil),   // 1: grpc.FixtureReply
}
var file_fixtures_test_proto_depIdxs = []int32{
	0, // 0: grpc.Fixture.Ping:input_type -> grpc.FixtureRequest
	0, // 1: grpc.Fixture.StreamPing:input_type -> grpc.FixtureRequest
	1, // 2: grpc.Fixture.Ping:output_type -> grpc.FixtureReply
	1, // 3: grpc.Fixture.StreamPing:output_type -> grpc.FixtureReply
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_fixtures_test_proto_init() }
func file_fixtures_test_proto_init() {
	if File_fixtures_test_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_fixtures_test_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FixtureRequest); i {
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
		file_fixtures_test_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FixtureReply); i {
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
			RawDescriptor: file_fixtures_test_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_fixtures_test_proto_goTypes,
		DependencyIndexes: file_fixtures_test_proto_depIdxs,
		MessageInfos:      file_fixtures_test_proto_msgTypes,
	}.Build()
	File_fixtures_test_proto = out.File
	file_fixtures_test_proto_rawDesc = nil
	file_fixtures_test_proto_goTypes = nil
	file_fixtures_test_proto_depIdxs = nil
}