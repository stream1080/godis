package cluster

import "github.com/stream1080/godis/interface/resp"

func Ping(cluster *ClusterDatabases, c resp.Connection, args [][]byte) resp.Reply {
	return cluster.db.Exec(c, args)
}
