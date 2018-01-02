package client

import (
	"context"
	"fmt"
	"io"

	log "github.com/Sirupsen/logrus"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/radu-matei/kube-toolkit/pkg/rpc"
	"google.golang.org/grpc"
)

// Config contains all configuration for the client
type Config struct {
	ServerHost string
	Stdout     io.Writer
	Stderr     io.Writer
}

// Client contains all necessary information to
// connect to the server
type Client struct {
	Config *Config
	RPC    rpc.GRPCClient
}

// NewClient returns a new instance of the client
func NewClient(cfg *Config, conn *grpc.ClientConn) *Client {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	return &Client{
		Config: cfg,
		RPC:    rpc.NewGRPCClient(conn),
	}
}

//GetGRPCConnection returns a new grpc connection
func GetGRPCConnection(serverHost string) (conn *grpc.ClientConn, err error) {
	if conn, err = grpc.Dial(serverHost, grpc.WithInsecure()); err != nil {
		return nil, fmt.Errorf("failed to dial %q: %v", serverHost, err)
	}
	return conn, nil
}

// GetVersion returns the server version
func (client *Client) GetVersion(ctx context.Context) (*rpc.Version, error) {
	return client.RPC.GetVersion(ctx, &google_protobuf.Empty{})
}

// ServerStream starts a stream from the server
func (client *Client) ServerStream(ctx context.Context, opts ...grpc.CallOption) error {
	stream, err := client.RPC.ServerStream(ctx, &google_protobuf.Empty{})
	if err != nil {
		log.Fatalf("cannot start server stream: %v", err)
		return err
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error receiving stream: %v", err)
		}

		fmt.Println(message.Message)
	}

	return nil
}

// GetValue gets the value of an existing key in etcd
func (client *Client) GetValue(ctx context.Context, m *rpc.StateMessage) (*rpc.StateMessage, error) {
	return client.RPC.GetValue(ctx, m)
}

// PutValue creates a new entry in etcd with m.Key / m.Value
func (client *Client) PutValue(ctx context.Context, m *rpc.StateMessage, opts ...grpc.CallOption) (*rpc.StateMessage, error) {
	return client.RPC.PutValue(ctx, m)
}
