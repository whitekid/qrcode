import { QRCodeClient as Client } from '@/lib/proto/V1alpha1ServiceClientPb';
import { Request } from '@/lib/proto/v1alpha1_pb';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

export class QRCodeClient {
  s: Client;

  constructor() {
    this.s = new Client(process.env.NEXT_PUBLIC_API_ENDPOINT!);
  }

  async version() {
    const r = await this.s.version(new Empty(), null);
    return r.toObject().value;
  }

  async text(value: string) {
    const req = new Request();
    req.setWidth(200);
    req.setHeight(200);
    req.setContent(value);
    const r = await this.s.generate(req, null);
    return r.toObject();
  }
}
