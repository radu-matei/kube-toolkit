package main

import (
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/joker/pkg/joker"
	"github.com/spf13/cobra"
)

var (
	gcloudUsage   = "interacts with the Azure CLI pod"
	gcloudCommand string
)

type gcloudCmd struct {
	out    io.Writer
	client *joker.Client
}

func newGCloudCmd(out io.Writer) *cobra.Command {
	gcloud := &gcloudCmd{
		out: out,
	}

	cmd := &cobra.Command{
		Use:   "gcloud",
		Short: azUsage,
		Long:  azUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			setupConnection()

			return gcloud.run()
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&gcloudCommand, "execute", "", "command to execute agains the az CLI")

	return cmd
}

func (gcloud *gcloudCmd) run() error {
	log.Debugf("making request to: %s using gcloud command: %s", gothamHost, gcloudCommand)

	return nil
}
