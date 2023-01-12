TARGET=bin/qrcodeapi
SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "*_test.go")
GOPATH=$(shell go env GOPATH)
BUILD_FLAGS?=-v

.PHONY: clean test get tidy cadl

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

cadl:
	@cd cadl && cadl compile .

