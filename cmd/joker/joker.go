package main

import (
	"fmt"
	"io"
	"os"

	"google.golang.org/grpc"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	log "github.com/Sirupsen/logrus"
	"github.com/kubernetes/helm/pkg/kube"
	"github.com/radu-matei/joker/pkg/joker"
	"github.com/radu-matei/joker/pkg/portforwarder"
	"github.com/spf13/cobra"
)

var (
	globalUsage = "Joker - the cloud service deployment for Kubernetes"

	flagDebug   bool
	kubeContext string
	gothamHost  string

	//gothamTunnel *kube.Tunnel
)

func newRootCmd(out io.Writer, in io.Reader) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "joker",
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
	flags.StringVar(&kubeContext, "kube-context", "", "kubeconfig context to use")
	flags.StringVar(&gothamHost, "host", "", "address of Gotham server")

	cmd.AddCommand(
		newVersionCmd(out),
		newCloudInitCmd(out),
	)

	return cmd
}

func setupConnection() error {
	if gothamHost == "" {
		clientset, config, err := getKubeClient(kubeContext)
		if err != nil {
			return err
		}

		clientConfig, err := config.ClientConfig()
		if err != nil {
			return err
		}

		gothamTunnel, err := portforwarder.New(clientset, clientConfig, "default")
		if err != nil {
			return err
		}

		gothamHost = fmt.Sprintf("localhost:%d", gothamTunnel.Local)
		log.Debugf("Created tunnel using local port: '%d'", gothamTunnel.Local)
	}

	log.Debugf("SERVER: %q", gothamHost)
	return nil
}

// getKubeClient is a convenience method for creating kubernetes config and client
// for a given kubeconfig context
func getKubeClient(context string) (*kubernetes.Clientset, clientcmd.ClientConfig, error) {
	log.Debugf("getting kube client using Kubernetes context: %v", context)

	config := kube.GetConfig(context)
	clientConfig, err := config.ClientConfig()
	if err != nil {
		log.Debugf("cannot get clientConfig: %v", err)
		return nil, nil, fmt.Errorf("could not get kubernetes config for context '%s': %s", context, err)
	}

	client, err := kubernetes.NewForConfig(clientConfig)
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

func teardown() {
	// if gothamTunnel != nil {
	// 	log.Debugf("Tearing down tunnel connection to Gotham...")
	// 	gothamTunnel.Close()
	// }
}

func ensureJokerClient(client *joker.Client, conn *grpc.ClientConn) *joker.Client {
	log.Debugf("passing gothamHost: %s", gothamHost)
	cfg := &joker.ClientConfig{
		GothamHost: gothamHost,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
	}

	return joker.NewClient(cfg, conn)
}
