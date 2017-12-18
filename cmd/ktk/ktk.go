package main

import (
	"fmt"
	"io"
	"os"

	"google.golang.org/grpc"

	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	log "github.com/Sirupsen/logrus"
	"github.com/kubernetes/helm/pkg/kube"
	"github.com/radu-matei/kube-toolkit/pkg/ktk"
	"github.com/radu-matei/kube-toolkit/pkg/portforwarder"
	"github.com/spf13/cobra"
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

	flags := cmd.PersistentFlags()
	flags.BoolVar(&flagDebug, "debug", false, "enable verbose output")
	flags.StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig file to use")
	flags.StringVar(&ktkdHost, "host", "", "address of ktkd server")

	cmd.AddCommand(
		newVersionCmd(out),
		newServerStreamCmd(out),
	)

	return cmd
}

func setupConnection() error {
	if ktkdHost == "" {
		clientset, config, err := getKubeClient(kubeconfig)
		if err != nil {
			return err
		}

		ktkdTunnel, err = portforwarder.New(clientset, config, "default")
		if err != nil {
			return err
		}

		ktkdHost = fmt.Sprintf("localhost:%d", ktkdTunnel.Local)
		log.Debugf("Created tunnel using local port: '%d'", ktkdTunnel.Local)
	}

	log.Debugf("SERVER: %q", ktkdHost)
	return nil
}

// getKubeClient is a convenience method for creating kubernetes config and client
// for a given kubeconfig
func getKubeClient(kubeconfig string) (*kubernetes.Clientset, *restclient.Config, error) {
	log.Debugf("getting kube client using Kubernetes kubeconfig: %v", kubeconfig)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get kubernetes config from kubeconfig '%s': %v", kubeconfig, err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get kubernetes client: %s", err)
	}
	return client, config, nil
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
	log.Debugf("passing ktkdHost: %s", ktkdHost)
	cfg := &ktk.ClientConfig{
		KTKDHost: ktkdHost,
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
	}

	return ktk.NewClient(cfg, conn)
}
