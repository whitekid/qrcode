// Original file: ../proto/v1alpha1.proto

import type * as grpc from '@grpc/grpc-js'
import type { MethodDefinition } from '@grpc/proto-loader'
import type { Empty as _google_protobuf_Empty, Empty__Output as _google_protobuf_Empty__Output } from '../../google/protobuf/Empty';
import type { Request as _api_v1alpha1_Request, Request__Output as _api_v1alpha1_Request__Output } from '../../api/v1alpha1/Request';
import type { Response as _api_v1alpha1_Response, Response__Output as _api_v1alpha1_Response__Output } from '../../api/v1alpha1/Response';
import type { StringValue as _google_protobuf_StringValue, StringValue__Output as _google_protobuf_StringValue__Output } from '../../google/protobuf/StringValue';

export interface QRCodeClient extends grpc.Client {
  generate(argument: _api_v1alpha1_Request, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_api_v1alpha1_Response__Output>): grpc.ClientUnaryCall;
  generate(argument: _api_v1alpha1_Request, metadata: grpc.Metadata, callback: grpc.requestCallback<_api_v1alpha1_Response__Output>): grpc.ClientUnaryCall;
  generate(argument: _api_v1alpha1_Request, options: grpc.CallOptions, callback: grpc.requestCallback<_api_v1alpha1_Response__Output>): grpc.ClientUnaryCall;
  generate(argument: _api_v1alpha1_Request, callback: grpc.requestCallback<_api_v1alpha1_Response__Output>): grpc.ClientUnaryCall;
  generate(argument: _api_v1alpha1_Request, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_api_v1alpha1_Response__Output>): grpc.ClientUnaryCall;
  generate(argument: _api_v1alpha1_Request, metadata: grpc.Metadata, callback: grpc.requestCallback<_api_v1alpha1_Response__Output>): grpc.ClientUnaryCall;
  generate(argument: _api_v1alpha1_Request, options: grpc.CallOptions, callback: grpc.requestCallback<_api_v1alpha1_Response__Output>): grpc.ClientUnaryCall;
  generate(argument: _api_v1alpha1_Request, callback: grpc.requestCallback<_api_v1alpha1_Response__Output>): grpc.ClientUnaryCall;
  
  version(argument: _google_protobuf_Empty, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_google_protobuf_StringValue__Output>): grpc.ClientUnaryCall;
  version(argument: _google_protobuf_Empty, metadata: grpc.Metadata, callback: grpc.requestCallback<_google_protobuf_StringValue__Output>): grpc.ClientUnaryCall;
  version(argument: _google_protobuf_Empty, options: grpc.CallOptions, callback: grpc.requestCallback<_google_protobuf_StringValue__Output>): grpc.ClientUnaryCall;
  version(argument: _google_protobuf_Empty, callback: grpc.requestCallback<_google_protobuf_StringValue__Output>): grpc.ClientUnaryCall;
  version(argument: _google_protobuf_Empty, metadata: grpc.Metadata, options: grpc.CallOptions, callback: grpc.requestCallback<_google_protobuf_StringValue__Output>): grpc.ClientUnaryCall;
  version(argument: _google_protobuf_Empty, metadata: grpc.Metadata, callback: grpc.requestCallback<_google_protobuf_StringValue__Output>): grpc.ClientUnaryCall;
  version(argument: _google_protobuf_Empty, options: grpc.CallOptions, callback: grpc.requestCallback<_google_protobuf_StringValue__Output>): grpc.ClientUnaryCall;
  version(argument: _google_protobuf_Empty, callback: grpc.requestCallback<_google_protobuf_StringValue__Output>): grpc.ClientUnaryCall;
  
}

export interface QRCodeHandlers extends grpc.UntypedServiceImplementation {
  generate: grpc.handleUnaryCall<_api_v1alpha1_Request__Output, _api_v1alpha1_Response>;
  
  version: grpc.handleUnaryCall<_google_protobuf_Empty__Output, _google_protobuf_StringValue>;
  
}

export interface QRCodeDefinition extends grpc.ServiceDefinition {
  generate: MethodDefinition<_api_v1alpha1_Request, _api_v1alpha1_Response, _api_v1alpha1_Request__Output, _api_v1alpha1_Response__Output>
  version: MethodDefinition<_google_protobuf_Empty, _google_protobuf_StringValue, _google_protobuf_Empty__Output, _google_protobuf_StringValue__Output>
}
