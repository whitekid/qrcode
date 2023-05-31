// server side gprc
import { credentials, loadPackageDefinition } from '@grpc/grpc-js';
import { loadSync } from '@grpc/proto-loader';
import { promisify } from 'util';
import { QRCodeClient } from './proto/api/v1alpha1/QRCode';
import { ProtoGrpcType } from './proto/v1alpha1';

const packageDefinition = loadSync('../proto/v1alpha1.proto');
const grpcObject = loadPackageDefinition(packageDefinition);
const QRCodeService = (grpcObject as unknown as ProtoGrpcType).api.v1alpha1.QRCode;

export namespace api {
  export namespace v1alpha1 {
    export class QRCode {
      s: QRCodeClient;

      constructor() {
        this.s = new QRCodeService('localhost:8080', credentials.createInsecure());
      }

      public async version() {
        return promisify(this.s.version).bind(this.s)({});
      }
    }
  }
}
