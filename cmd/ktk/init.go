package main

import (
	"fmt"
	"io"

	"github.com/radu-matei/kube-toolkit/pkg/k8s"
	"github.com/radu-matei/kube-toolkit/pkg/ktk"
	"github.com/spf13/cobra"
)

var (
	initUsage = "deploys the ktkd server to your cluster"

	dockerImage string
)

type initCmd struct {
	out    io.Writer
	client *ktk.Client
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
	flags.StringVar(&dockerImage, "docker-image", "", "docker image to use for ktkd deployment")

	return cmd
}

func (cmd *initCmd) run() error {

	err := k8s.CreateDeployment(kubeconfig, dockerImage, "ktkd")
	if err != nil {
		return fmt.Errorf("cannot create deployment: %v", err)
	}
	fmt.Printf("Created ktkd deployment")

	return nil
}
