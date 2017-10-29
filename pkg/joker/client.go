package joker

import (
	"context"
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
func NewClient(cfg *ClientConfig) *Client {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(cfg.GothamHost, opts...)
	if err != nil {
		log.Fatalf("could not dial server: %v", err)
	}

	// TODO - Investigate deferring the closing of the connection
	//defer conn.Close()

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

// connect connects the DraftClient to the DraftServer.
// func connect(server *Client) (conn *grpc.ClientConn, err error) {
// 	if conn, err = grpc.Dial(server.Config.GothamHost, grpc.WithInsecure()); err != nil {
// 		return nil, fmt.Errorf("failed to dial %q: %v", server.Config.GothamHost, err)
// 	}
// 	return conn, nil
// }
