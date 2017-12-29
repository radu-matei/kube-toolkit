package main

import (
	"fmt"
	"io"

	"github.com/radu-matei/kube-toolkit/pkg/k8s"
	"github.com/spf13/cobra"
)

var (
	initUsage = "deploys the ktkd server to your cluster"

	serverImage  string
	gatewayImage string
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
	flags.StringVar(&serverImage, "docker-image", "", "docker image to use for ktkd deployment")
	flags.StringVar(&gatewayImage, "gateway-image", "", "docker image to use for the gatewayw")

	return cmd
}

func (cmd *initCmd) run() error {

	err := k8s.CreateDeployment(kubeconfig, serverImage, gatewayImage, "ktkd")
	if err != nil {
		return fmt.Errorf("cannot create deployment: %v", err)
	}
	fmt.Printf("Created ktkd deployment")

	return nil
}
