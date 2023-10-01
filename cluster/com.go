package cluster

import (
	"context"
	"errors"
	"github.com/stream1080/godis/resp/client"
)

func (cluster *ClusterDatabases) getPeerClient(peer string) (*client.Client, error) {
	pool, ok := cluster.peerConn[peer]
	if !ok {
		return nil, errors.New("conn not found")
	}

	object, err := pool.BorrowObject(context.Background())
	if err != nil {
		return nil, err
	}

	c, ok := object.(*client.Client)
	if !ok {
		return nil, errors.New("wrong type")
	}

	return c, nil
}
