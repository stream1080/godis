package cluster

import (
	"github.com/stream1080/godis/interface/resp"
	"github.com/stream1080/godis/resp/reply"
)

func Del(cluster *ClusterDatabases, c resp.Connection, args [][]byte) resp.Reply {
	replies := cluster.broadcast(c, args)
	var errReply reply.ErrorReply
	var deleted int64 = 0
	for _, r := range replies {
		if reply.IsErrorReply(r) {
			errReply = r.(reply.ErrorReply)
			return reply.MakeErrReply("error: " + errReply.Error())
		}

		intReply, ok := r.(*reply.IntReply)
		if !ok {
			return reply.MakeErrReply("error: " + errReply.Error())
		}

		deleted += intReply.Code
	}

	return reply.MakeIntReply(deleted)
}
