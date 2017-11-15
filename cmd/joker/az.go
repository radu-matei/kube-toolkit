package main

import (
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/joker/pkg/joker"
	"github.com/spf13/cobra"
)

var (
	azUsage   = "interacts with the Azure CLI pod"
	azCommand string
)

type azCmd struct {
	out    io.Writer
	client *joker.Client
}

func newAzCmd(out io.Writer) *cobra.Command {
	az := &azCmd{
		out: out,
	}

	cmd := &cobra.Command{
		Use:   "az",
		Short: azUsage,
		Long:  azUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			setupConnection()

			return az.run()
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&azCommand, "execute", "", "command to execute agains the az CLI")

	return cmd
}

func (az *azCmd) run() error {
	log.Debugf("making request to: %s using az command: %s", gothamHost, azCommand)

	return nil
}
