package main

import (
	"context"
	"fmt"
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/kube-toolkit/pkg/client"
	"github.com/radu-matei/kube-toolkit/pkg/rpc"
	"github.com/spf13/cobra"
)

var (
	getStateUsage = "prints value of an existing key in etcd"
	key           string
)

type getStateCmd struct {
	out    io.Writer
	client *client.Client
}

func newGetStateCmd(out io.Writer) *cobra.Command {
	getState := &getStateCmd{
		out: out,
	}

	cmd := &cobra.Command{
		Use:   "get",
		Short: getStateUsage,
		Long:  getStateUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := setupConnection(remoteServerPort, localRandomPort)
			if err != nil {
				log.Fatalf("cannot setup connection: %v", err)
			}

			return getState.run()
		},
	}

	flags := cmd.PersistentFlags()
	flags.StringVar(&key, "key", "", "key to get from etcd")

	return cmd
}

func (cmd *getStateCmd) run() error {

	conn, err := client.GetGRPCConnection(serverHost)
	if err != nil {
		log.Fatalf("cannot create grpc connection: %v", err)
	}
	defer conn.Close()
	cmd.client = ensureGRPCClient(cmd.client, conn)

	s := &rpc.StateMessage{
		Key: key,
	}

	v, err := cmd.client.GetValue(context.Background(), s)
	if err != nil {
		return fmt.Errorf("cannot get value for key %v: %v", key, err)
	}
	fmt.Printf("%v: %v", v.Key, v.Value)

	return nil
}
