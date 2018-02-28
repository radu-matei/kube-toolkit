package main

import (
	"fmt"
	"io"

	"github.com/radu-matei/kube-toolkit/example/scaffold/pkg/client"
	"github.com/radu-matei/kube-toolkit/pkg/k8s"
	"github.com/spf13/cobra"
)

var (
	resetUsage = "deletes the server deployment from the Kubernetes cluster"
)

type resetCmd struct {
	out    io.Writer
	client *client.Client
}

func newResetCmd(out io.Writer) *cobra.Command {
	initCmd := &resetCmd{
		out: out,
	}

	cmd := &cobra.Command{
		Use:   "reset",
		Short: resetUsage,
		Long:  resetUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			return initCmd.run()
		},
	}

	flags := cmd.PersistentFlags()
	flags.StringVar(&deploymentName, "name", "kube-toolkit", "docker image to use for the web gatewayw")

	return cmd
}

func (cmd *resetCmd) run() error {

	err := k8s.DeleteDeployment(kubeconfig, deploymentName)
	if err != nil {
		return fmt.Errorf("cannot delete deployment: %v", err)
	}
	fmt.Println("Deleted server deployment")

	return nil
}
