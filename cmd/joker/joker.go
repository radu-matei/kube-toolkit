package main

import (
	"fmt"
	"io"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	log "github.com/Sirupsen/logrus"
	"github.com/kubernetes/helm/pkg/kube"
	"github.com/radu-matei/joker/pkg/joker"
	"github.com/radu-matei/joker/pkg/portforwarder"
	"github.com/spf13/cobra"
)

const (
	hostEnvVar = "GOTHAM_HOST"
)

var (
	globalUsage = "Joker - the cloud service deployment for Kubernetes"

	flagDebug   bool
	kubeContext string
	gothamHost  string

	ghostTunnel *kube.Tunnel
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
	flags.StringVar(&gothamHost, "host", defaultGothamHost(), "address of Gotham server")

	cmd.AddCommand(newVersionCmd(out))

	return cmd
}

func setupConnection(c *cobra.Command, args []string) error {
	if gothamHost == "" {
		clientset, config, err := getKubeClient(kubeContext)
		if err != nil {
			return err
		}
		clientConfig, err := config.ClientConfig()
		if err != nil {
			return err
		}

		tunnel, err := portforwarder.New(clientset, clientConfig, "default")
		if err != nil {
			return err
		}

		gothamHost = fmt.Sprintf("localhost:%d", tunnel.Local)
		log.Debugf("Created tunnel using local port: '%d'", tunnel.Local)
	}

	log.Debugf("SERVER: %q", gothamHost)
	return nil
}

// getKubeClient is a convenience method for creating kubernetes config and client
// for a given kubeconfig context
func getKubeClient(context string) (*kubernetes.Clientset, clientcmd.ClientConfig, error) {
	config := kube.GetConfig(context)
	clientConfig, err := config.ClientConfig()
	if err != nil {
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

}

func defaultGothamHost() string {
	return "127.0.0.1:10000"
}

func ensureJokerClient(client *joker.Client) *joker.Client {

	cfg := joker.ClientConfig{
		GothamHost: gothamHost,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
	}
	return joker.NewClient(&cfg)

}
