SHELL:=/usr/bin/env bash

BIN_NAME:=wdproc
BIN_VERSION:=$(shell ./.version.sh)

default: help
.PHONY: help
help: ## Print help
	@grep -E '^[a-zA-Z_-\/]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: all
all: clean build-darwin-arm64 ## Build for macOS (arm64)

.PHONY: clean
clean: ## Remove build products (./out)
	rm -rf ./out

.PHONY: build
build: ## Build for the current platform & architecture to ./out
	mkdir -p out
	env CGO_ENABLED=0 go build -ldflags="-X main.version=${BIN_VERSION}" -o ./out/${BIN_NAME} .

.PHONY: build-darwin-arm64
build-darwin-arm64: ## Build for macOS/arm64 to ./out
	mkdir -p out
	env GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-X main.version=${BIN_VERSION}" -o ./out/${BIN_NAME}-${BIN_VERSION}-darwin-arm64 .

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run
