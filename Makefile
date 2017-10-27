plugins = grpc
target = go
protoc_location = pkg/rpc


GIT_COMMIT = $(shell git rev-parse HEAD)
SEMVER = "v0.1"

JOKER_CMD_PATH = cmd/joker
GOTHAM_CMD_PATH = cmd/gotham

JOKER_BINARY_NAME = joker
GOTHAM_BINARY_NAME = gotham

OUTPUT_DIR = bin

LDFLAGS += -X github.com/radu-matei/joker/pkg/version.GitCommit=${GIT_COMMIT}
LDFLAGS += -X github.com/radu-matei/joker/pkg/version.SemVer=${SEMVER}

.PHONY: rpc
rpc:
	protoc -I $(protoc_location) --$(target)_out=plugins=$(plugins):$(protoc_location) $(protoc_location)/*.proto

.PHONY: joker
joker:
	cd $(JOKER_CMD_PATH) && \
	go build -ldflags '$(LDFLAGS)' -o ../../$(OUTPUT_DIR)/$(JOKER_BINARY_NAME)


.PHONY: gotham
gotham:
	cd $(GOTHAM_CMD_PATH) && \
	go build -ldflags '$(LDFLAGS)' -o ../../$(OUTPUT_DIR)/$(GOTHAM_BINARY_NAME)