package main

import (
	"fmt"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/kube-toolkit/pkg/client"
	"github.com/radu-matei/kube-toolkit/pkg/k8s"
	"github.com/radu-matei/kube-toolkit/pkg/portforwarder"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	remoteServerPort  = 10000
	remoteGatewayPort = 8080
	localRandomPort   = 0
)

var (
	globalUsage = "the client-side component of your awesome Kubernetes tool"

	flagDebug bool

	kubeconfig string
	serverHost string

	deploymentName string

	kubeTunnel *k8s.Tunnel
)

func newRootCmd(out io.Writer, in io.Reader) *cobra.Command {
	cmd := &cobra.Command{
		Short:        globalUsage,
		Long:         globalUsage,
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if flagDebug {
				log.SetLevel(log.DebugLevel)
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			teardown()
		},
	}

	kubeconfig = getEnvVarOrExit("KUBECONFIG")

	flags := cmd.PersistentFlags()
	flags.BoolVar(&flagDebug, "debug", false, "enable verbose output")
	flags.StringVar(&serverHost, "host", "", "address of the server")
	flags.StringVar(&deploymentName, "name", "kube-toolkit", "kubernetes deployment name")

	cmd.AddCommand(
		newInitCmd(out),
		newProxyCmd(out),
		newResetCmd(out),
		newVersionCmd(out),
		newServerStreamCmd(out),
	)

	return cmd
}

func setupConnection(remotePort, localPort int) error {
	if serverHost == "" {
		clientset, config, err := k8s.GetKubeClient(kubeconfig)
		if err != nil {
			return err
		}

		kubeTunnel, err = portforwarder.New(clientset, config, "default", deploymentName, remotePort, localPort)
		if err != nil {
			return err
		}

		serverHost = fmt.Sprintf("localhost:%d", kubeTunnel.Local)
		log.Debugf("Created tunnel using local port: '%d'", kubeTunnel.Local)
	}

	return nil
}

func main() {
	cmd := newRootCmd(os.Stdout, os.Stdin)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// the tunnel appears to be already closed here
func teardown() {
	// if kubeTunnel != nil {
	// 	log.Debugf("Tearing down tunnel connection to server...")
	// 	kubeTunnel.Close()
	// }
}

func ensureGRPCClient(c *client.Client, conn *grpc.ClientConn) *client.Client {
	cfg := &client.Config{
		ServerHost: serverHost,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
	}

	return client.NewClient(cfg, conn)
}

func getEnvVarOrExit(varName string) string {
	value := os.Getenv(varName)
	if value == "" {
		log.Fatalf("missing environment variable %s\n", varName)
	}

	return value
}
