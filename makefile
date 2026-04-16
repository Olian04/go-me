# All targets are phony (no file named "build", "test", etc. should shadow these).
.PHONY: build build-release format help lint test test-coverage

REV := $(shell git rev-parse HEAD 2>/dev/null || echo unknown)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

help:
	@echo "Usage: make <target>"
	@echo "Targets:"
	@printf '  %-20s %s\n' 'build' 'Build me (./cmd/me)'
	@printf '  %-20s %s\n' 'build-release' 'GoReleaser snapshot to dist/ (no publish)'
	@printf '  %-20s %s\n' 'format' 'go fmt + gofmt -w'
	@printf '  %-20s %s\n' 'lint' 'go vet, go mod verify, govulncheck, gosec, golangci-lint'
	@printf '  %-20s %s\n' 'test' 'Run all tests'
	@printf '  %-20s %s\n' 'test-coverage' 'Generate coverage.out + HTML report'

lint:
	go vet ./...
	go mod verify
	go tool govulncheck ./...
	go tool gosec -fmt text -stdout -quiet ./...
	golangci-lint run ./...

format:
	go fmt ./...
	gofmt -w .

build:
	go build -trimpath -ldflags "-s -w -X github.com/Olian04/go-me/cmd/me/version.Version=dev -X github.com/Olian04/go-me/cmd/me/version.Revision=$(REV) -X github.com/Olian04/go-me/cmd/me/version.BuildTime=$(BUILD_TIME)" -o me ./cmd/me

build-release:
	@command -v syft >/dev/null 2>&1 || { echo >&2 "syft not on PATH (install: https://github.com/anchore/syft#installation)"; exit 1; }
	@command -v goreleaser >/dev/null 2>&1 || { echo >&2 "goreleaser not on PATH (install: https://goreleaser.com/install/)"; exit 1; }
	ME_BUILD_TIME=$(BUILD_TIME) RELEASE_NAME=local-snapshot RELEASE_BODY='Local snapshot (not a production release).' goreleaser release --snapshot --clean --skip=publish,validate --config goreleaser/.goreleaser.yaml

test:
	go test -shuffle=on ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
