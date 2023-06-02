TARGET=bin/qrcodeapi
SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "*_test.go")
GOPATH=$(shell go env GOPATH)
BUILD_FLAGS?=-v -ldflags="-s -w"

.PHONY: clean test get tidy

all: build
build: $(TARGET)

$(TARGET): $(SRC)
	@mkdir -p bin
	go build -o bin/ ${BUILD_FLAGS} ./cmd/...

clean:
	rm -f ${TARGET}

test:
	go test -v ./...

# update modules & tidy
dep:
	@rm -f go.mod go.sum
	@go mod init qrcodeapi

	@go mod edit -replace github.com/chai2010/webp=github.com/chai2010/webp@v1.2-alpha1

	@$(MAKE) tidy

tidy:
	@go mod tidy -v

spec/tsp-output/@typespec/openapi3/openapi.yaml: spec/main.tsp spec/qrcodeapi_v1.tsp
	@cd spec && tsp compile .

proto/v1alpha1.pb.go: proto/v1alpha1.proto
	protoc -I=./proto \
      --go_out=./proto \
      --go_opt=paths=source_relative \
      --go-grpc_out=./proto \
      --go-grpc_opt=paths=source_relative \
      --js_out=import_style=commonjs:./www/src/lib/proto \
      --grpc-web_out=import_style=typescript,mode=grpcweb:./www/src/lib/proto \
	  proto/v1alpha1.proto

	  cd ./proto && mockery --name QRCodeClient
