package cluster

import "github.com/stream1080/godis/interface/resp"

func defaultFunc(cluster *ClusterDatabases, c resp.Connection, args [][]byte) resp.Reply {
	peer := cluster.peerPick.PickNode(string(args[1]))
	return cluster.relay(peer, c, args)
}

func MakeRouter() map[string]CmdFunc {
	routerMap := make(map[string]CmdFunc)

	routerMap["exists"] = defaultFunc
	routerMap["type"] = defaultFunc
	routerMap["set"] = defaultFunc
	routerMap["setnx"] = defaultFunc
	routerMap["get"] = defaultFunc
	routerMap["getset"] = defaultFunc

	routerMap["ping"] = Ping

	return routerMap
}
