package cluster

import "github.com/stream1080/godis/interface/resp"

func defaultFunc(cluster *ClusterDatabases, c resp.Connection, args [][]byte) resp.Reply {
	peer := cluster.peerPick.PickNode(string(args[1]))
	return cluster.relay(peer, c, args)
}
