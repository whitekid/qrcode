'use client';

import { QRCodeClient } from '@/lib/grpcwebclient';
import Image from 'next/image';
import { ChangeEvent, useState } from 'react';

const client = new QRCodeClient();

export default function Panel() {
  const [value, setValue] = useState('');
  const [image, setImage] = useState('');
  const [width, setWidth] = useState(0);
  const [height, setHeight] = useState(0);

  return (
    <>
      Text <input
        value={value}
        onChange={(v: ChangeEvent<HTMLInputElement>) => {
          setValue(v.currentTarget.value);
          client.text(v.currentTarget.value).then((r) => {
            setWidth(r.width);
            setHeight(r.height);
            setImage(`data:${r.contentType};base64, ${r.image}`);
          });
        }}
      />
      <Image alt="preview" width={width} height={height} src={image} />
    </>
  );
}
