// Original file: ../proto/v1alpha1.proto


export interface Response {
  'contentType'?: (string);
  'width'?: (number);
  'height'?: (number);
  'image'?: (Buffer | Uint8Array | string);
}

export interface Response__Output {
  'contentType': (string);
  'width': (number);
  'height': (number);
  'image': (Buffer);
}
