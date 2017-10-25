package main

import (
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
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

func main() {
	cmd := newRootCmd(os.Stdout, os.Stdin)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func teardown() {

}

func defaultGothamHost() string {
	return os.Getenv(hostEnvVar)
}
