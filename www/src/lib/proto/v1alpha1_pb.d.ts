import * as jspb from 'google-protobuf'

import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';
import * as google_protobuf_wrappers_pb from 'google-protobuf/google/protobuf/wrappers_pb';


export class Request extends jspb.Message {
  getContent(): string;
  setContent(value: string): Request;

  getUrl(): string;
  setUrl(value: string): Request;

  getWidth(): number;
  setWidth(value: number): Request;

  getHeight(): number;
  setHeight(value: number): Request;

  getAccept(): string;
  setAccept(value: string): Request;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Request.AsObject;
  static toObject(includeInstance: boolean, msg: Request): Request.AsObject;
  static serializeBinaryToWriter(message: Request, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Request;
  static deserializeBinaryFromReader(message: Request, reader: jspb.BinaryReader): Request;
}

export namespace Request {
  export type AsObject = {
    content: string,
    url: string,
    width: number,
    height: number,
    accept: string,
  }
}

export class Response extends jspb.Message {
  getContentType(): string;
  setContentType(value: string): Response;

  getWidth(): number;
  setWidth(value: number): Response;

  getHeight(): number;
  setHeight(value: number): Response;

  getImage(): Uint8Array | string;
  getImage_asU8(): Uint8Array;
  getImage_asB64(): string;
  setImage(value: Uint8Array | string): Response;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Response.AsObject;
  static toObject(includeInstance: boolean, msg: Response): Response.AsObject;
  static serializeBinaryToWriter(message: Response, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Response;
  static deserializeBinaryFromReader(message: Response, reader: jspb.BinaryReader): Response;
}

export namespace Response {
  export type AsObject = {
    contentType: string,
    width: number,
    height: number,
    image: Uint8Array | string,
  }
}

