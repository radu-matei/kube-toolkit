package gotham

import (
	"fmt"
	"net"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/joker/pkg/rpc"
	"github.com/radu-matei/joker/pkg/version"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// ServerConfig contains all configuration for the Gotham server
type ServerConfig struct {
	ListenAddress string
}

// Server contains all methods and config for the Gotham server
type Server struct {
	Config *ServerConfig
	RPC    *grpc.Server
}

// NewServer returns a new instance of the Gotham server
func NewServer(cfg *ServerConfig) *Server {
	return &Server{
		Config: cfg,
		RPC:    grpc.NewServer(),
	}
}

// Serve starts the server and listens on ListenAddress
func (server *Server) Serve(ctx context.Context) error {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 10000))
	if err != nil {
		return fmt.Errorf("failed to start listening: %v", err)
	}

	rpc.RegisterJokerServer(server.RPC, server)

	_, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	errc := make(chan error, 1)

	wg.Add(1)
	go func() {
		errc <- server.RPC.Serve(lis)
		log.Debugf("starting to serve...")
		close(errc)
		wg.Done()
	}()

	defer func() {
		server.RPC.Stop()
		log.Debugf("stopping the server")
		cancel()
		wg.Wait()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errc:
		return err
	}

}

// GetVersion returns the current version of the server.
func (server *Server) GetVersion(ctx context.Context, _ *rpc.Empty) (*rpc.Version, error) {
	log.Debugf("executing gotham version")
	return &rpc.Version{
		SemVer:    version.SemVer,
		GitCommit: version.GitCommit}, nil
}

// InitializeCloud initializes a cloud
func (server *Server) InitializeCloud(cfg *rpc.CloudConfig, stream rpc.Joker_InitializeCloudServer) error {
	log.Debugf("received InitializeCloud server method with cfg: %s", cfg.CloudProvider.String())
	for i := 0; i < 5; i++ {
		err := stream.Send(&rpc.CloudInitStream{
			Message: fmt.Sprintf("Sending stream back to client, iteration: %d", i),
		})
		if err != nil {
			return err
		}

		time.Sleep(2 * time.Second)
	}

	return nil
}
