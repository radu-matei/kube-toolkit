package server

import (
	"fmt"
	"net"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/radu-matei/kube-toolkit/example/scaffold/pkg/rpc"
	"github.com/radu-matei/kube-toolkit/example/scaffold/pkg/version"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Config contains all configuration for the server
type Config struct {
	ListenAddress string
}

// Server contains all methods and config for the server
type Server struct {
	Config *Config
	RPC    *grpc.Server
}

// NewServer returns a new instance of the server
func NewServer(cfg *Config) *Server {
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

	rpc.RegisterGRPCServer(server.RPC, server)

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
func (server *Server) GetVersion(ctx context.Context, _ *google_protobuf.Empty) (*rpc.Version, error) {
	log.Debugf("executing GetVersion")
	return &rpc.Version{
		SemVer:    version.SemVer,
		GitCommit: version.GitCommit}, nil
}

// ServerStream starts a new stream from the server
func (server *Server) ServerStream(_ *google_protobuf.Empty, stream rpc.GRPC_ServerStreamServer) error {
	log.Debugf("received server stream command")
	for i := 0; i < 5; i++ {
		err := stream.Send(&rpc.Message{
			Message: fmt.Sprintf("Sending stream back to client, iteration: %d", i),
		})
		if err != nil {
			return err
		}

		time.Sleep(2 * time.Second)
	}

	return nil
}
