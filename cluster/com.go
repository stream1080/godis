package cluster

import (
	"context"
	"errors"
	"github.com/stream1080/godis/interface/resp"
	"github.com/stream1080/godis/lib/utils"
	"github.com/stream1080/godis/resp/client"
	"github.com/stream1080/godis/resp/reply"
	"strconv"
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

func (cluster *ClusterDatabases) returnPeerClient(peer string, client *client.Client) error {
	pool, ok := cluster.peerConn[peer]
	if !ok {
		return errors.New("conn not found")
	}

	return pool.ReturnObject(context.Background(), client)
}

func (cluster *ClusterDatabases) relay(peer string, c resp.Connection, args [][]byte) resp.Reply {
	if peer == cluster.self {
		return cluster.db.Exec(c, args)
	}

	peerClient, err := cluster.getPeerClient(peer)
	if err != nil {
		return reply.MakeErrReply(err.Error())
	}

	defer func() {
		_ = cluster.returnPeerClient(peer, peerClient)
	}()
	peerClient.Send(utils.ToCmdLine("SELECT", strconv.Itoa(c.GetDBIndex())))

	return peerClient.Send(args)
}

func (cluster *ClusterDatabases) broadcast(c resp.Connection, args [][]byte) map[string]resp.Reply {
	results := make(map[string]resp.Reply)
	for _, node := range cluster.nodes {
		results[node] = cluster.relay(node, c, args)
	}

	return results
}
