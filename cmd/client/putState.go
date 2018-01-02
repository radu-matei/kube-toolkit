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
	putStateUsage = "creates a new state entry in etcd"
	value         string
)

type putStateCmd struct {
	out    io.Writer
	client *client.Client
}

func newPutStateCmd(out io.Writer) *cobra.Command {
	putState := &putStateCmd{
		out: out,
	}

	cmd := &cobra.Command{
		Use:   "put",
		Short: putStateUsage,
		Long:  putStateUsage,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := setupConnection(remoteServerPort, localRandomPort)
			if err != nil {
				log.Fatalf("cannot setup connection: %v", err)
			}

			return putState.run()
		},
	}

	flags := cmd.PersistentFlags()

	flags.StringVar(&key, "key", "", "key to put to etcd")
	flags.StringVar(&value, "value", "", "value to put to etcd")

	return cmd
}

func (cmd *putStateCmd) run() error {

	conn, err := client.GetGRPCConnection(serverHost)
	if err != nil {
		log.Fatalf("cannot create grpc connection: %v", err)
	}
	defer conn.Close()
	cmd.client = ensureGRPCClient(cmd.client, conn)

	s := &rpc.StateMessage{
		Key:   key,
		Value: value,
	}

	v, err := cmd.client.PutValue(context.Background(), s)
	if err != nil {
		return fmt.Errorf("cannot get value for key %v: %v", key, err)
	}
	fmt.Printf("%v: %v", v.Key, v.Value)

	return nil
}
