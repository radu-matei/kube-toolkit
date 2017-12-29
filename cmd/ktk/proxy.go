package main

import (
	"io"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var (
	proxyUsage = "starts a proxy to your web client gateway"

	port int
)

type proxyCmd struct {
	out io.Writer
}

func newProxyCmd(out io.Writer) *cobra.Command {
	proxyCommand := &proxyCmd{
		out: out,
	}

	cmd := &cobra.Command{
		Use:   "proxy",
		Short: initUsage,
		Long:  initUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := setupConnection(8080)
			if err != nil {
				log.Fatalf("cannot setup connection: %v", err)
			}
			return proxyCommand.run()
		},
	}

	flags := cmd.PersistentFlags()
	flags.IntVar(&port, "port", 8081, "local port to start the proxy")

	return cmd
}

func (cmd *proxyCmd) run() error {

	time.Sleep(10 * time.Second)
	return nil
}
