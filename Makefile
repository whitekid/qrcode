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

	@$(MAKE) tidy

	@go mod edit -replace github.com/whitekid/goxp@v0.0.10=github.com/whitekid/goxp@master
	@$(MAKE) tidy

tidy:
	@go mod tidy -v

cadl:
	@cd cadl && cadl compile .

