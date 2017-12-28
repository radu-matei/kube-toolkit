plugins = grpc
target = go
protoc_location = pkg/rpc


GIT_COMMIT = $(shell git rev-parse HEAD)
SEMVER = "v0.3.1"

KTK_CMD_PATH = cmd/ktk
KTKD_CMD_PATH = cmd/ktkd
GATEWAY_CMD_PATH = cmd/gateway

KTK_BINARY_NAME = ktk
KTKD_BINARY_NAME = ktkd
GATEWAY_BINARY_NAME = gateway

KTK_LINUX_BINARY = ktk-linux
KTKD_LINUX_BINARY = ktkd-linux
GATEWAY_LINUX_BINARY = gateway-linux


OUTPUT_DIR = bin
VERSION_PACKAGE = github.com/radu-matei/kube-toolkit/pkg/version
LDFLAGS += -X $(VERSION_PACKAGE).GitCommit=${GIT_COMMIT}
LDFLAGS += -X $(VERSION_PACKAGE).SemVer=${SEMVER}

PROTOBUF_INCLUDE_DIR = /usr/local/include
GRPC_GATEWAY_PROTO = vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis

.PHONY: rpc
rpc:
	protoc -I $(protoc_location) -I $(PROTOBUF_INCLUDE_DIR) -I $(GRPC_GATEWAY_PROTO)  --$(target)_out=plugins=$(plugins):$(protoc_location) $(protoc_location)/*.proto --grpc-gateway_out=logtostderr=true:$(protoc_location)


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

.PHONY: gateway
gateway:
	cd $(GATEWAY_CMD_PATH) && \
	go build -o ../../$(OUTPUT_DIR)/$(GATEWAY_BINARY_NAME)

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
	
