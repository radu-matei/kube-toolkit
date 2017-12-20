package main

import (
	"fmt"
	"io"

	"github.com/radu-matei/kube-toolkit/pkg/k8s"
	"github.com/radu-matei/kube-toolkit/pkg/ktk"
	"github.com/spf13/cobra"
)

var (
	reset = "deletes the ktkd deployment from the Kubernetes cluster"
)

type resetCmd struct {
	out    io.Writer
	client *ktk.Client
}

func newResetCmd(out io.Writer) *cobra.Command {
	initCmd := &resetCmd{
		out: out,
	}

	cmd := &cobra.Command{
		Use:   "reset",
		Short: initUsage,
		Long:  initUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			return initCmd.run()
		},
	}

	return cmd
}

func (cmd *resetCmd) run() error {

	err := k8s.DeleteDeployment(kubeconfig, "ktkd")
	if err != nil {
		return fmt.Errorf("cannot delete deployment: %v", err)
	}

	return nil
}
