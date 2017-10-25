plugins = grpc
target = go
protoc_location = pkg/rpc

.PHONY: rpc
rpc:
	protoc -I $(protoc_location) --$(target)_out=plugins=$(plugins):$(protoc_location) $(protoc_location)/*.proto