# nextjs app

    npx create-next-app@latest --typescript --eslint --tailwind --app --src-dir --use-yarn www

## GRPC supports

- `grpcweb`은 [browser의 XHR에 기반하고 있음](https://github.com/grpc/grpc-web/discussions/1161). 서버사이드에서 사용하기 위해서는 `grpc-js`를 사용함
- `grpcweb`을 `use client`를 이용해 사용하려고 해도 오류가 발생함

    yarn add @grpc/grpc-js

- `grpc-js`는 `makeUnaryRequest()`를 통해서 요청을 전달하는데, 이러면 typescript와 grpc의 strong type의 잇접을 가지지 못한다.

generating typescript types

    yarn add @grpc/proto-loader
    yarn proto-loader-gen-types --longs=String --enums=String --defaults --oneofs --grpcLib=@grpc/grpc-js --outDir=src/lib/proto/ ../proto/*.proto

<https://github.com/grpc/grpc-node/tree/master/packages/grpc-js>
