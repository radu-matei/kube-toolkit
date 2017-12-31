# uses unzip at the moment
apt-get update
apt-get install -y unzip

# download, extract and add to path the protobuf compiler
cd /home
wget https://github.com/google/protobuf/releases/download/v3.5.1/protoc-3.5.1-linux-x86_64.zip
unzip protoc-3.5.1-linux-x86_64.zip
export PATH=$PATH:/home/bin

# go get the protobuf dependencies needed
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -u github.com/golang/protobuf/protoc-gen-go