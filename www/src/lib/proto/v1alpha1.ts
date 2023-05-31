import type * as grpc from '@grpc/grpc-js';
import type { MessageTypeDefinition } from '@grpc/proto-loader';

import type { QRCodeClient as _api_v1alpha1_QRCodeClient, QRCodeDefinition as _api_v1alpha1_QRCodeDefinition } from './api/v1alpha1/QRCode';

type SubtypeConstructor<Constructor extends new (...args: any) => any, Subtype> = {
  new(...args: ConstructorParameters<Constructor>): Subtype;
};

export interface ProtoGrpcType {
  api: {
    v1alpha1: {
      QRCode: SubtypeConstructor<typeof grpc.Client, _api_v1alpha1_QRCodeClient> & { service: _api_v1alpha1_QRCodeDefinition }
      Request: MessageTypeDefinition
      Response: MessageTypeDefinition
    }
  }
  google: {
    protobuf: {
      BoolValue: MessageTypeDefinition
      BytesValue: MessageTypeDefinition
      DoubleValue: MessageTypeDefinition
      Empty: MessageTypeDefinition
      FloatValue: MessageTypeDefinition
      Int32Value: MessageTypeDefinition
      Int64Value: MessageTypeDefinition
      StringValue: MessageTypeDefinition
      UInt32Value: MessageTypeDefinition
      UInt64Value: MessageTypeDefinition
    }
  }
}

