FROM golang:1.9.2 as builder

WORKDIR /go/src/github.com/radu-matei/kube-toolkit

COPY . .

RUN ["chmod", "+x", "prerequisites.sh"]

RUN ./prerequisites.sh

RUN make server-linux

# starting from ubuntu right now, there's an issue starting from alpine/scratch
FROM ubuntu

COPY --from=builder /go/src/github.com/radu-matei/kube-toolkit/bin /app

WORKDIR /app
EXPOSE 10000

ENTRYPOINT ["./ktkd-linux"]
CMD ["start", "--debug"]




####
# original version of Dockerfile, before addinng prerequisites.sh
# here for future reference
####
# FROM golang:1.9.2 as builder

# WORKDIR /go/src/github.com/radu-matei/kube-toolkit

# COPY . .

# RUN apt-get update
# RUN apt-get install -y unzip


# RUN cd /home && \
#     wget https://github.com/google/protobuf/releases/download/v3.5.1/protoc-3.5.1-linux-x86_64.zip && \
#     unzip protoc-3.5.1-linux-x86_64.zip && \
#     export PATH=$PATH:/home/bin

# RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway && \
#     go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger && \
#     go get -u github.com/golang/protobuf/protoc-gen-go

# RUN make server-linux

# # starting from ubuntu right now, there's an issue starting from alpine/scratch
# FROM ubuntu

# COPY --from=builder /go/src/github.com/radu-matei/kube-toolkit/bin /app

# WORKDIR /app
# EXPOSE 10000

# ENTRYPOINT ["./ktkd-linux"]
# CMD ["start", "--debug"]
