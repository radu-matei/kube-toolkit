package gotham

import (
	"fmt"
	"net"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/radu-matei/joker/pkg/rpc"
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
	server := new(Server)
	server.Config = cfg
	server.RPC = grpc.NewServer()

	return server

	// return &Server{
	// 	Config: cfg,
	// 	RPC:    grpc.NewServer(),
	// }
}

// Serve starts the server and listens on ListenAddress
func (server *Server) Serve() error {
	fmt.Println("serve function")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 10000))
	if err != nil {
		return fmt.Errorf("failed to start listening: %v", err)
	}

	rpc.RegisterJokerServer(server.RPC, server)

	err = server.RPC.Serve(lis)
	fmt.Printf("%v", err)

	return err

	//_, cancel := context.WithCancel(ctx)
	//var wg sync.WaitGroup
	//errc := make(chan error, 1)
	/*
		//wg.Add(1)
		go func() {
			errc <- grpcServer.Serve(lis)
			fmt.Println("actual serve")
			close(errc)
			//wg.Done()
		}()

		defer func() {
			grpcServer.Stop()
			cancel()
			//wg.Wait()
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errc:
			return err
		}
	*/
}

// GetVersion returns the current version of the server.
func (server *Server) GetVersion(ctx context.Context, _ *rpc.Empty) (*rpc.Version, error) {
	fmt.Println("executing gotham version")
	return &rpc.Version{
		SemVer:    "v0.1-experimental",
		GitCommit: "git commit"}, nil
}
