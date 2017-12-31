plugins = grpc
target = go
protoc_location = pkg/rpc


GIT_COMMIT = $(shell git rev-parse HEAD)
SEMVER = "v0.3.1"

CLIENT_CMD_PATH = cmd/client
SERVER_CMD_PATH = cmd/server
GATEWAY_CMD_PATH = gateway

CLIENT_BINARY_NAME = ktk
SERVER_BINARY_NAME = ktkd
GATEWAY_BINARY_NAME = gateway

CLIENT_LINUX_BINARY = $(CLIENT_BINARY_NAME)-linux
SERVER_LINUX_BINARY = $(SERVER_BINARY_NAME)-linux
GATEWAY_LINUX_BINARY = $(GATEWAY_BINARY_NAME)-linux

OUTPUT_DIR = bin

VERSION_PACKAGE = github.com/radu-matei/kube-toolkit/pkg/version
LDFLAGS += -X $(VERSION_PACKAGE).GitCommit=${GIT_COMMIT}
LDFLAGS += -X $(VERSION_PACKAGE).SemVer=${SEMVER}

PROTOBUF_INCLUDE_DIR = vendor/protobuf-include
GRPC_GATEWAY_PROTO = vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis
SWAGGER_ANNOTATIONS = vendor/github.com/grpc-ecosystem/grpc-gateway/


.PHONY: rpc
rpc:
	protoc -I $(protoc_location) -I $(PROTOBUF_INCLUDE_DIR) -I $(GRPC_GATEWAY_PROTO) -I $(SWAGGER_ANNOTATIONS) --$(target)_out=plugins=$(plugins):$(protoc_location) $(protoc_location)/*.proto --grpc-gateway_out=logtostderr=true:$(protoc_location) --swagger_out=logtostderr=true:gateway/web && \
	cd gateway/web && mkdir src/generated-client && docker run --rm -v ${PWD}/gateway/web:/local swaggerapi/swagger-codegen-cli generate     -i /local/rpc.swagger.json     -l typescript-angular      -c /local/swagger.config.json     -o /local/src/generated-client


.PHONY: bin
bin:
	$(MAKE) client && \
	$(MAKE) server

.PHONY: client
client:
	cd $(CLIENT_CMD_PATH) && \
	go build -ldflags '$(LDFLAGS)' -o ../../$(OUTPUT_DIR)/$(CLIENT_BINARY_NAME)

.PHONY: server
server:
	cd $(SERVER_CMD_PATH) && \
	go build -ldflags '$(LDFLAGS)' -o ../../$(OUTPUT_DIR)/$(SERVER_BINARY_NAME)

.PHONY: gateway
gateway:
	cd $(GATEWAY_CMD_PATH) && \
	go build -o ../$(OUTPUT_DIR)/$(GATEWAY_BINARY_NAME)

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: client-linux
client-linux:
	cd $(CLIENT_CMD_PATH) && \
	GOOS=linux go build -ldflags '$(LDFLAGS)' -o ../../$(OUTPUT_DIR)/$(CLIENT_LINUX_BINARY)

.PHONY: server-linux
server-linux:
	cd $(SERVER_CMD_PATH) && \
	GOOS=linux go build -ldflags '$(LDFLAGS)' -o ../../$(OUTPUT_DIR)/$(SERVER_LINUX_BINARY)
	
.PHONY: gateway-linux
gateway-linux:
	cd $(GATEWAY_CMD_PATH) && \
	GOOS=linux go build -o ../$(OUTPUT_DIR)/$(GATEWAY_LINUX_BINARY)
