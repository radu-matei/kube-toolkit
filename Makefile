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
VERSION_PACKAGE = github.com/radu-matei/joker/pkg/version
LDFLAGS += -X $(VERSION_PACKAGE).GitCommit=${GIT_COMMIT}
LDFLAGS += -X $(VERSION_PACKAGE).SemVer=${SEMVER}

.PHONY: rpc
rpc:
	protoc -I $(protoc_location) --$(target)_out=plugins=$(plugins):$(protoc_location) $(protoc_location)/*.proto

.PHONY: bin
bin:
	$(MAKE) joker && \
	$(MAKE) gotham

.PHONY: joker
joker:
	cd $(JOKER_CMD_PATH) && \
	go build -ldflags '$(LDFLAGS)' -o ../../$(OUTPUT_DIR)/$(JOKER_BINARY_NAME)

.PHONY: gotham
gotham:
	cd $(GOTHAM_CMD_PATH) && \
	go build -ldflags '$(LDFLAGS)' -o ../../$(OUTPUT_DIR)/$(GOTHAM_BINARY_NAME)

.PHONY: clean
clean:
	rm -rf bin/