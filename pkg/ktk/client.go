package ktk

import (
	"context"
	"fmt"
	"io"

	log "github.com/Sirupsen/logrus"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/radu-matei/kube-toolkit/pkg/rpc"
	"google.golang.org/grpc"
)

// ClientConfig contains all configuration for the ktk client
type ClientConfig struct {
	KTKDHost string
	Stdout   io.Writer
	Stderr   io.Writer
}

// Client contains all necessary information to
// connect to the ktkd server
type Client struct {
	Config *ClientConfig
	RPC    rpc.KTKClient
}

// NewClient returns a new instance of the ktk client
func NewClient(cfg *ClientConfig, conn *grpc.ClientConn) *Client {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	return &Client{
		Config: cfg,
		RPC:    rpc.NewKTKClient(conn),
	}
}

// GetVersion returns the ktk version
func (client *Client) GetVersion(ctx context.Context) (*rpc.Version, error) {

	// TODO - remove this once google.protobuf.empty is used
	empty := &google_protobuf.Empty{}

	return client.RPC.GetVersion(ctx, empty)
}

// ServerStream starts a stream from the server
func (client *Client) ServerStream(ctx context.Context, opts ...grpc.CallOption) error {
	log.Debugf("called InitializeCloud client method...")
	empty := &google_protobuf.Empty{}

	stream, err := client.RPC.ServerStream(ctx, empty)
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

//GetGRPCConnection returns a new grpc connection
func GetGRPCConnection(ktkdHost string) (conn *grpc.ClientConn, err error) {
	if conn, err = grpc.Dial(ktkdHost, grpc.WithInsecure()); err != nil {
		return nil, fmt.Errorf("failed to dial %q: %v", ktkdHost, err)
	}
	return conn, nil
}
