package etcd

import (
	"fmt"

	"github.com/coreos/etcd/client"
)

// NewKeysAPI returns a new etcd keys api based on a config
func NewKeysAPI(cfg client.Config) (client.KeysAPI, error) {
	c, err := client.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot get etcd client: %v", err)
	}

	return client.NewKeysAPI(c), nil
}
