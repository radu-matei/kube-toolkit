package main

import (
	"fmt"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kubernetes/helm/pkg/kube"
	"github.com/radu-matei/kube-toolkit/pkg/k8s"
	"github.com/radu-matei/kube-toolkit/pkg/ktk"
	"github.com/radu-matei/kube-toolkit/pkg/portforwarder"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	globalUsage = "ktk - the client-side component of your awesome Kubernetes tool"

	flagDebug  bool
	kubeconfig string
	ktkdHost   string

	ktkdTunnel *kube.Tunnel
)

func newRootCmd(out io.Writer, in io.Reader) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "ktk",
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
	flags.StringVar(&ktkdHost, "host", "", "address of ktkd server")

	cmd.AddCommand(
		newInitCmd(out),
		newProxyCmd(out),
		newResetCmd(out),
		newVersionCmd(out),
		newServerStreamCmd(out),
	)

	return cmd
}

func setupConnection(port int) error {
	if ktkdHost == "" {
		clientset, config, err := k8s.GetKubeClient(kubeconfig)
		if err != nil {
			return err
		}

		ktkdTunnel, err = portforwarder.New(clientset, config, "default", port)
		if err != nil {
			return err
		}

		ktkdHost = fmt.Sprintf("localhost:%d", ktkdTunnel.Local)
		log.Debugf("Created tunnel using local port: '%d'", ktkdTunnel.Local)
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
	// if ktkdTunnel != nil {
	// 	log.Debugf("Tearing down tunnel connection to ktkd...")
	// 	ktkdTunnel. .Close()
	// }
}

func ensureKTKClient(client *ktk.Client, conn *grpc.ClientConn) *ktk.Client {
	cfg := &ktk.ClientConfig{
		KTKDHost: ktkdHost,
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
	}

	return ktk.NewClient(cfg, conn)
}

func getEnvVarOrExit(varName string) string {
	value := os.Getenv(varName)
	if value == "" {
		log.Fatalf("missing environment variable %s\n", varName)
	}

	return value
}
