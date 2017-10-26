package main

import (
	"context"
	"fmt"
	"io"

	"github.com/radu-matei/joker/pkg/joker"
	"github.com/spf13/cobra"
)

var (
	versionUsage = "prints the Joker and Gotham version information"
)

type versionCmd struct {
	out    io.Writer
	client *joker.Client
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
			versionCmd.client = ensureJokerClient(versionCmd.client)
			return versionCmd.run()
		},
	}

	return cmd
}

func (cmd *versionCmd) run() error {
	gothamVersion, err := cmd.client.GetVersion(context.Background())
	if err != nil {
		return fmt.Errorf("cannot get Gotham version: %v", err)
	}

	fmt.Printf("\nGotham version:  %s Git Commit: %s", gothamVersion.GetSemVer(), gothamVersion.GetGitCommit())
	return nil
}
