# Makefile for spacedock build and test targets

.PHONY: all build build-windows-amd64 test test-windows-amd64 clean

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X github.com/spacedock-dev/spacedock/internal/cli.Version=$(VERSION)

all: build

build:
	CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o dist/spacedock ./cmd/spacedock

build-windows-amd64:
	@mkdir -p dist
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o dist/spacedock_windows_amd64.exe ./cmd/spacedock

test:
	go test ./...

test-windows-amd64:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build ./...

clean:
	rm -rf dist/
