// Code generated by protoc-gen-go. DO NOT EDIT.
// plugins:
// 	protoc-gen-go
// 	protoc
// source: temporal/server/api/testservice/v1/service.proto

package testservice

import (
	reflect "reflect"
	unsafe "unsafe"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_temporal_server_api_testservice_v1_service_proto protoreflect.FileDescriptor

const file_temporal_server_api_testservice_v1_service_proto_rawDesc = "" +
	"\n" +
	"0temporal/server/api/testservice/v1/service.proto\x12\"temporal.server.api.testservice.v1\x1a9temporal/server/api/testservice/v1/request_response.proto2\x89\x01\n" +
	"\vTestService\x12z\n" +
	"\tSendHello\x124.temporal.server.api.testservice.v1.SendHelloRequest\x1a5.temporal.server.api.testservice.v1.SendHelloResponse\"\x00B6Z4go.temporal.io/server/api/testservice/v1;testserviceb\x06proto3"

var file_temporal_server_api_testservice_v1_service_proto_goTypes = []any{
	(*SendHelloRequest)(nil),  // 0: temporal.server.api.testservice.v1.SendHelloRequest
	(*SendHelloResponse)(nil), // 1: temporal.server.api.testservice.v1.SendHelloResponse
}
var file_temporal_server_api_testservice_v1_service_proto_depIdxs = []int32{
	0, // 0: temporal.server.api.testservice.v1.TestService.SendHello:input_type -> temporal.server.api.testservice.v1.SendHelloRequest
	1, // 1: temporal.server.api.testservice.v1.TestService.SendHello:output_type -> temporal.server.api.testservice.v1.SendHelloResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_temporal_server_api_testservice_v1_service_proto_init() }
func file_temporal_server_api_testservice_v1_service_proto_init() {
	if File_temporal_server_api_testservice_v1_service_proto != nil {
		return
	}
	file_temporal_server_api_testservice_v1_request_response_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_temporal_server_api_testservice_v1_service_proto_rawDesc), len(file_temporal_server_api_testservice_v1_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_temporal_server_api_testservice_v1_service_proto_goTypes,
		DependencyIndexes: file_temporal_server_api_testservice_v1_service_proto_depIdxs,
	}.Build()
	File_temporal_server_api_testservice_v1_service_proto = out.File
	file_temporal_server_api_testservice_v1_service_proto_goTypes = nil
	file_temporal_server_api_testservice_v1_service_proto_depIdxs = nil
}
