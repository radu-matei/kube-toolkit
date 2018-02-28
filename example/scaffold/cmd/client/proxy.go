package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
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
		Short: proxyUsage,
		Long:  proxyUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := setupConnection(dashboardPort, port)
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
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()

	for {
		fmt.Printf("serving proxy on localhost:%d...\n", port)
		time.Sleep(10 * time.Second)
	}
}
