package main

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var (
	versionUsage = "prints the Joker and Gotham version information"
)

type versionCmd struct {
	out io.Writer
}

func newVersionCmd(out io.Writer) *cobra.Command {
	version := &versionCmd{
		out: out,
	}

	cmd := &cobra.Command{
		Use:   "version",
		Short: versionUsage,
		Long:  versionUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			return version.run()
		},
	}

	return cmd
}

func (version *versionCmd) run() error {
	fmt.Println("Version of the Joker client")
	return nil
}
