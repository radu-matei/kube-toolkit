package joker

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"

	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/joker/pkg/rpc"
)

// ClientConfig contains all configuration for the Joker client
type ClientConfig struct {
	GothamHost string
	Stdout     io.Writer
	Stderr     io.Writer
}

// Client contains all necessary information to
// connect to the Gotham server
type Client struct {
	Config *ClientConfig
	RPC    rpc.JokerClient
}

// NewClient returns a new instance of the Joker client
func NewClient(cfg *ClientConfig, conn *grpc.ClientConn) *Client {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	return &Client{
		Config: cfg,
		RPC:    rpc.NewJokerClient(conn),
	}
}

// GetVersion returns the Gotham version
func (client *Client) GetVersion(ctx context.Context) (*rpc.Version, error) {

	// TODO - remove this once google.protobuf.empty is used
	empty := &rpc.Empty{}

	return client.RPC.GetVersion(ctx, empty)
}

// InitializeCloud initializes a cloud
func (client *Client) InitializeCloud(ctx context.Context, cfg *rpc.CloudConfig, opts ...grpc.CallOption) error {
	log.Debugf("called InitializeCloud client method...")

	stream, err := client.RPC.InitializeCloud(ctx, cfg)
	if err != nil {
		log.Fatalf("cannot initialize cloud: %v", err)
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

//GetGRPCConnection returns a new grpc connection
func GetGRPCConnection(gothamHost string) (conn *grpc.ClientConn, err error) {
	if conn, err = grpc.Dial(gothamHost, grpc.WithInsecure()); err != nil {
		return nil, fmt.Errorf("failed to dial %q: %v", gothamHost, err)
	}
	return conn, nil
}
