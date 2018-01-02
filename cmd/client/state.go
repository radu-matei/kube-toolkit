package main

import (
	"io"

	"github.com/spf13/cobra"
)

func newStateCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use: "state",
	}

	return cmd
}
