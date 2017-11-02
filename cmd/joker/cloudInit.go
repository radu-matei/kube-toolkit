package main

import (
	"context"
	"io"
	"log"

	"github.com/radu-matei/joker/pkg/joker"
	"github.com/radu-matei/joker/pkg/rpc"
	"github.com/spf13/cobra"
)

var (
	cloudInitUsage = "prints the Joker and Gotham version information"
)

type cloudInitCmd struct {
	out    io.Writer
	client *joker.Client
}

func newCloudInitCmd(out io.Writer) *cobra.Command {
	cloudInitCmd := &cloudInitCmd{
		out: out,
	}

	cmd := &cobra.Command{
		Use:   "init",
		Short: cloudInitUsage,
		Long:  cloudInitUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			setupConnection(cmd, args)

			cloudInitCmd.client = ensureJokerClient(cloudInitCmd.client)
			return cloudInitCmd.run()
		},
	}

	return cmd
}

func (cmd *cloudInitCmd) run() error {

	cfg := &rpc.CloudConfig{
		CloudName:      "awesome-cloud",
		ContainerImage: "docker-hub-image",
	}
	err := cmd.client.InitializeCloud(context.Background(), cfg)

	if err != nil {
		log.Fatalf("cannot initialize cloud: %v", err)
		return err
	}

	return nil
}
