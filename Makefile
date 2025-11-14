
BIN?=diag
PKG=github.com/chigaji/diag-cli
# LD_FLAGS=-x $(PKG)/internal/version.Version=$(VERSION) -x $(PKG)/internal/version.Commit=$(COMMIT) -x $(PKG)/internal/version.BuildDate=$(BUILD_DATE)
LD_FLAGS=-X $(PKG)/internal/version.Version=$(VERSION) -X $(PKG)/internal/version.Commit=$(COMMIT) -X $(PKG)/internal/version.BuildDate=$(DATE)
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE?=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
# BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

all: build

build:
# 	@mkdir -p bin
	go build -ldflags "$(LD_FLAGS)" -o bin/$(BIN) ./cmd/diag

fmt: 
	go fmt ./...

lint:
	golangci-lint run || true 

test:
	go test ./...

run: 
	go run -ldflags "$(LD_FLAGS)" ./cmd/diag --help

release: 
	goreleaser build --snapshot --clean