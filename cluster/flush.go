package cluster

import (
	"github.com/stream1080/godis/interface/resp"
	"github.com/stream1080/godis/resp/reply"
)

func FlushDB(cluster *ClusterDatabases, c resp.Connection, args [][]byte) resp.Reply {
	replies := cluster.broadcast(c, args)
	var errReply reply.ErrorReply
	for _, r := range replies {
		if reply.IsErrorReply(r) {
			errReply = r.(reply.ErrorReply)
			break
		}
	}

	if errReply != nil {
		return reply.MakeErrReply("error: " + errReply.Error())
	}

	return reply.MakeOkReply()
}
