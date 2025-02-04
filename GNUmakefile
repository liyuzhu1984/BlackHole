TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
GFLAGS=-ldflags "-w -s -buildid=${BUILD_ID}"

UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

# 设置默认值
ifeq ($(UNAME_S),Linux)
    GOOS := linux
else ifeq ($(UNAME_S),Darwin)
    GOOS := darwin
else ifeq ($(UNAME_S),Windows_NT)
    GOOS := windows
else
    GOOS := unknown
endif

ifeq ($(UNAME_M),x86_64)
    GOARCH := amd64
else ifeq ($(UNAME_M),arm64)
    GOARCH := arm64
else ifeq ($(UNAME_M),aarch64)
    GOARCH := arm64
else
    GOARCH := unknown
endif

default: build

build: all

fmt:
	@if command -v gofmt >/dev/null 2>&1; then \
		echo "Running gofmt..."; \
		gofmt -w $(GOFMT_FILES); \
	else \
		echo "gofmt not found, skipping..."; \
	fi
	@if command -v goimports >/dev/null 2>&1; then \
		echo "Running goimports..."; \
		goimports -w $(GOFMT_FILES); \
	else \
		echo "goimports not found, skipping..."; \
	fi

init:
	go mod init BlackHole

env:
	@go env -w GO111MODULE=on
	@go env -w GOPROXY=https://goproxy.cn,direct

vendor:
	@go mod tidy
	@go mod vendor

.PHONY: build env fmt vendor

all: fmt vendor stash voidengine

dev: clean fmt mac copy

copy:
	tar -xvf bin/BlackHole_darwin-amd64.tgz

clean:
	rm -rf bin/*
	go clean -modcache

tools:
	GOOS=${GOOS} GOARCH=${GOARCH} go build ${GFLAGS} -o bin/tools cmd/tools/main.go

swagger-generator: fmt vendor
	GOOS=${GOOS} GOARCH=${GOARCH} go build ${GFLAGS} -o bin/swagger-generator cmd/swagger-generator/main.go
	./bin/swagger-generator -source=./internal/voidengine/voidengine.go -output=./docs/api/voidengine

stash: fmt vendor
	GOOS=${GOOS} GOARCH=${GOARCH} go build ${GFLAGS} -o bin/stash cmd/stash/main.go

voidengine: fmt vendor swagger-generator
	GOOS=${GOOS} GOARCH=${GOARCH} go build ${GFLAGS} -o bin/voidengine cmd/voidengine/main.go
