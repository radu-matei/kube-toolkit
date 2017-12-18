plugins = grpc
target = go
protoc_location = pkg/rpc


GIT_COMMIT = $(shell git rev-parse HEAD)
SEMVER = "v0.1"

KTK_CMD_PATH = cmd/ktk
KTKD_CMD_PATH = cmd/ktkd

KTK_BINARY_NAME = ktk
KTKD_BINARY_NAME = ktkd

KTK_LINUX_BINARY = ktk-linux
KTKD_LINUX_BINARY = ktkd-linux

OUTPUT_DIR = bin
VERSION_PACKAGE = github.com/radu-matei/kube-toolkit/pkg/version
LDFLAGS += -X $(VERSION_PACKAGE).GitCommit=${GIT_COMMIT}
LDFLAGS += -X $(VERSION_PACKAGE).SemVer=${SEMVER}

.PHONY: rpc
rpc:
	protoc -I $(protoc_location) --$(target)_out=plugins=$(plugins):$(protoc_location) $(protoc_location)/*.proto

.PHONY: bin
bin:
	$(MAKE) ktk && \
	$(MAKE) ktkd

.PHONY: ktk
ktk:
	cd $(KTK_CMD_PATH) && \
	go build -ldflags '$(LDFLAGS)' -o ../../$(OUTPUT_DIR)/$(KTK_BINARY_NAME)

.PHONY: ktkd
ktkd:
	cd $(KTKD_CMD_PATH) && \
	go build -ldflags '$(LDFLAGS)' -o ../../$(OUTPUT_DIR)/$(KTKD_BINARY_NAME)

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: ktk-linux
ktk-linux:
	cd $(KTK_CMD_PATH) && \
	GOOS=linux go build -ldflags '$(LDFLAGS)' -o ../../$(OUTPUT_DIR)/$(KTK_LINUX_BINARY)

.PHONY: ktkd-linux
ktkd-linux:
	cd $(KTKD_CMD_PATH) && \
	GOOS=linux go build -ldflags '$(LDFLAGS)' -o ../../$(OUTPUT_DIR)/$(KTKD_LINUX_BINARY)
	