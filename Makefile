.DEFAULT_GOAL:=help

BIN_DIR       ?= bin
TOOLS_DIR     := tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin
GOLANGCI_LINT := $(TOOLS_BIN_DIR)/golangci-lint

GOFLAGS = -gcflags "all=-trimpath=$(PWD)" -asmflags "all=-trimpath=$(PWD)"

GO_BUILD_ENV_VARS := GO111MODULE=on CGO_ENABLED=0

##@ Building

.PHONY: all
all: kv-api kv-grpc-server ## Build all Go binaries

.PHONY: kv-api
kv-api: ## Build the kv-api Go binary
	$(GO_BUILD_ENV_VARS) go build -o bin/kv-api $(GOFLAGS) -ldflags '$(LDFLAGS)' ./cmd/kv-api

.PHONY: kv-grpc-server
kv-grpc-server: ## Build the kv-grpc-server Go binary
	$(GO_BUILD_ENV_VARS) go build -o bin/kv-grpc-server $(GOFLAGS) -ldflags '$(LDFLAGS)' ./cmd/kv-grpc-server


##@ Testing

.PHONY: lint
lint: $(GOLANGCI_LINT) ## Lint codebase
	$(GOLANGCI_LINT) run -v

.PHONY: test
test: ## Run go unit tests
	@go test -v ./internal/...


##@ Generate protobuf

.PHONY: generate
generate: ## Generate Go protobuf files
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		--go-grpc_opt=require_unimplemented_servers=false \
		--plugin protoc-gen-go="$(shell go env GOPATH)/bin/protoc-gen-go" \
		internal/kvpb/kv.proto

.PHONY: generate-install
generate-install: ## Install Go protobuf plugins
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1


##@ Helpers

.PHONY: help clean

$(GOLANGCI_LINT): $(TOOLS_DIR)/go.mod # Build golangci-lint from tools folder.
	cd $(TOOLS_DIR); go build -o $(BIN_DIR)/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

clean: ## Cleanup the project folders
	@rm -f $(BIN_DIR)/*
	@rm -f $(TOOLS_BIN_DIR)/*

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
