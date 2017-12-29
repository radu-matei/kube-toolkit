package main

import (
	"context"
	"fmt"
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/kube-toolkit/pkg/client"
	"github.com/radu-matei/kube-toolkit/pkg/version"
	"github.com/spf13/cobra"
)

var (
	versionUsage = "prints the client and server version information"
)

type versionCmd struct {
	out    io.Writer
	client *client.Client
}

func newVersionCmd(out io.Writer) *cobra.Command {
	versionCmd := &versionCmd{
		out: out,
	}

	cmd := &cobra.Command{
		Use:   "version",
		Short: versionUsage,
		Long:  versionUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := setupConnection(remoteServerPort, localRandomPort)
			if err != nil {
				log.Fatalf("cannot setup connection: %v", err)
			}

			return versionCmd.run()
		},
	}

	return cmd
}

func (cmd *versionCmd) run() error {

	conn, err := client.GetGRPCConnection(serverHost)
	if err != nil {
		log.Fatalf("cannot create grpc connection: %v", err)
	}
	defer conn.Close()

	cmd.client = ensureGRPCClient(cmd.client, conn)

	serverVersion, err := cmd.client.GetVersion(context.Background())
	if err != nil {
		return fmt.Errorf("cannot get server version: %v", err)
	}
	fmt.Printf("Client SemVer: %s Git Commit: %s\n", version.SemVer, version.GitCommit)
	fmt.Printf("Server SemVer:  %s Git Commit: %s\n", serverVersion.SemVer, serverVersion.GitCommit)
	return nil
}
