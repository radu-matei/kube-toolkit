package main

import (
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	gw "github.com/radu-matei/kube-toolkit/pkg/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	serverEndpoint = flag.String("backend", "localhost:10000", "endpoint of server")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterGRPCHandlerFromEndpoint(ctx, gwMux, *serverEndpoint, opts)
	if err != nil {
		return err
	}

	m := http.NewServeMux()

	m.HandleFunc("/", gwMux.ServeHTTP)

	serveStatic(m)

	return http.ListenAndServe(":8080", m)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}

func serveStatic(mux *http.ServeMux) {
	fileServer := http.FileServer(http.Dir("web"))
	prefix := "/ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}
