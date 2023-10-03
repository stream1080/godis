package cluster

import (
	"context"
	"strings"

	"github.com/stream1080/godis/config"
	database2 "github.com/stream1080/godis/database"
	"github.com/stream1080/godis/interface/database"
	"github.com/stream1080/godis/interface/resp"
	"github.com/stream1080/godis/lib/consistenthash"
	"github.com/stream1080/godis/lib/logger"
	"github.com/stream1080/godis/resp/reply"

	pool "github.com/jolestar/go-commons-pool/v2"
)

type ClusterDatabases struct {
	self     string
	nodes    []string
	peerPick *consistenthash.NodeMap
	peerConn map[string]*pool.ObjectPool
	db       database.Database
}

func MakeClusterDatabases() *ClusterDatabases {
	cluster := &ClusterDatabases{
		self:     config.Properties.Self,
		db:       database2.NewStandaloneDatabase(),
		peerPick: consistenthash.NewNodeMap(nil),
		peerConn: make(map[string]*pool.ObjectPool),
	}
	nodes := make([]string, 0, len(config.Properties.Peers)+1)
	for _, peer := range config.Properties.Peers {
		nodes = append(nodes, peer)
	}
	nodes = append(nodes, config.Properties.Self)
	cluster.peerPick.AddNode(nodes...)
	ctx := context.Background()
	for _, peer := range config.Properties.Peers {
		cluster.peerConn[peer] = pool.NewObjectPoolWithDefaultConfig(ctx, &ConnFactory{
			Peer: peer,
		})
	}
	cluster.nodes = nodes
	return cluster
}

type CmdFunc func(cluster *ClusterDatabases, c resp.Connection, args [][]byte) resp.Reply

var router = MakeRouter()

func (cluster *ClusterDatabases) Exec(client resp.Connection, args [][]byte) (result resp.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
			result = &reply.UnknownErrReply{}
		}
	}()

	cmdName := strings.ToLower(string(args[0]))
	cmdFunc, ok := router[cmdName]
	if !ok {
		return reply.MakeErrReply("not supported cmd")
	}

	return cmdFunc(cluster, client, args)
}

func (cluster *ClusterDatabases) Close() {
	cluster.db.Close()
}

func (cluster *ClusterDatabases) AfterClientClose(conn resp.Connection) {
	cluster.db.AfterClientClose(conn)
}
