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
	cfg *ClientConfig
	rpc rpc.JokerClient
}

func newClient(cfg *ClientConfig) *Client {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(cfg.GothamHost, opts...)
	if err != nil {
		log.Fatalf("could not dial server: %v", err)
	}
	defer conn.Close()

	return &Client{
		cfg: cfg,
		rpc: rpc.NewJokerClient(conn),
	}
}

// GetVersion returns the Gotham version
func (client *Client) GetVersion(ctx context.Context) (*rpc.Version, error) {

	// TODO - remove this once google.protobuf.empty is used
	empty := &rpc.Empty{}

	return client.rpc.GetVersion(ctx, empty)
}
