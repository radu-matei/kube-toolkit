package main

import (
	"io"

	"github.com/radu-matei/joker/pkg/gotham"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var (
	startUsage = "starts the Gotham server"

	listenAddress string
)

type startCmd struct {
	out           io.Writer
	listenAddress string
}

func newStartCmd(out io.Writer) *cobra.Command {
	startCmd := &startCmd{}

	cmd := &cobra.Command{
		Use:   "start",
		Short: startUsage,
		Long:  startUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			return startCmd.run()
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&listenAddress, "listen-addr", "0.0.0.0:10000", "the Gotham server listen address")

	return cmd
}

func (cmd *startCmd) run() error {

	cfg := &gotham.ServerConfig{
		ListenAddress: listenAddress,
	}

	return gotham.NewServer(cfg).Serve(context.Background())
}
