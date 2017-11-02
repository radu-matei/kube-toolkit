package main

import (
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	globalUsage = "Gotham - the cloud service deployment server for Kubernetes"

	flagDebug bool
)

func newRootCmd(out io.Writer, in io.Reader) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "gotham",
		Short:        globalUsage,
		Long:         globalUsage,
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if flagDebug {
				log.SetLevel(log.DebugLevel)
			}
		},
	}

	flags := cmd.PersistentFlags()
	flags.BoolVar(&flagDebug, "debug", false, "enable verbose output")

	cmd.AddCommand(newStartCmd(out))
	return cmd
}

func main() {
	cmd := newRootCmd(os.Stdout, os.Stdin)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
