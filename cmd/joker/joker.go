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

	gothamTunnel *kube.Tunnel
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
	flags.BoolVar(&flagDebug, "debug", true, "enable verbose output")
	flags.StringVar(&kubeContext, "kube-context", "", "kubeconfig context to use")
	flags.StringVar(&gothamHost, "host", "", "address of Gotham server")

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

		gothamTunnel, err := portforwarder.New(clientset, clientConfig, "default")
		if err != nil {
			return err
		}

		gothamHost = fmt.Sprintf("localhost:%d", gothamTunnel.Local)
		log.Debugf("Created tunnel using local port: '%d'", gothamTunnel.Local)
	}

	// clientset, _, err := getKubeClient(kubeContext)
	// if err != nil {
	// 	return err
	// }
	// pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	// if err != nil {
	// 	panic(err.Error())
	// }
	// fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

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
		log.Debug("cannot get clientConfig: %v", err)
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
	log.Debug("Tearing down tunnel connection to Gotham...")
	if gothamTunnel != nil {
		gothamTunnel.Close()
	}
}

func ensureJokerClient(client *joker.Client) *joker.Client {
	log.Debugf("passing gothamHost: %s", gothamHost)
	cfg := joker.ClientConfig{
		GothamHost: gothamHost,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
	}
	return joker.NewClient(&cfg)

}
