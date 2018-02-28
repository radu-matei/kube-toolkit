package main

import (
	"fmt"
	"io"

	"github.com/radu-matei/kube-toolkit/pkg/k8s"
	"github.com/spf13/cobra"
)

var (
	initUsage = "deploys the server and dashboard to your cluster"

	serverImage    string
	dashboardImage string
)

type initCmd struct {
	out io.Writer
}

func newInitCmd(out io.Writer) *cobra.Command {
	initCmd := &initCmd{
		out: out,
	}

	cmd := &cobra.Command{
		Use:   "init",
		Short: initUsage,
		Long:  initUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			return initCmd.run()
		},
	}

	flags := cmd.PersistentFlags()
	flags.StringVar(&serverImage, "server-image", "", "docker image to use for the server deployment")
	flags.StringVar(&dashboardImage, "dashboard-image", "", "docker image to use for the dashboard")

	return cmd
}

func (cmd *initCmd) run() error {

	err := k8s.CreateDeployment(kubeconfig, serverImage, dashboardImage, deploymentName)
	if err != nil {
		return fmt.Errorf("cannot create deployment: %v", err)
	}
	fmt.Printf("Created the server deployment\n")

	return nil
}
