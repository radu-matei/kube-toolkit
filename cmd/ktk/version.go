package main

import (
	"context"
	"fmt"
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/kube-toolkit/pkg/ktk"
	"github.com/radu-matei/kube-toolkit/pkg/version"
	"github.com/spf13/cobra"
)

var (
	versionUsage = "prints the ktk and ktkd version information"
)

type versionCmd struct {
	out    io.Writer
	client *ktk.Client
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
			err := setupConnection(10000)
			if err != nil {
				log.Fatalf("cannot setup connection: %v", err)
			}

			return versionCmd.run()
		},
	}

	return cmd
}

func (cmd *versionCmd) run() error {

	conn, err := ktk.GetGRPCConnection(ktkdHost)
	if err != nil {
		log.Fatalf("cannot create grpc connection: %v", err)
	}
	defer conn.Close()

	cmd.client = ensureKTKClient(cmd.client, conn)

	ktkdVersion, err := cmd.client.GetVersion(context.Background())
	if err != nil {
		return fmt.Errorf("cannot get ktkd version: %v", err)
	}
	fmt.Printf("ktk Semver: %s GitCommit: %s\n", version.SemVer, version.GitCommit)
	fmt.Printf("ktkd SemVer:  %s Git Commit: %s\n", ktkdVersion.SemVer, ktkdVersion.GitCommit)
	return nil
}
