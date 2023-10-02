package cluster

import (
	"github.com/stream1080/godis/interface/resp"
	"github.com/stream1080/godis/resp/reply"
)

// rename k1 k2
func Rename(cluster *ClusterDatabases, c resp.Connection, args [][]byte) resp.Reply {
	if len(args) != 3 {
		return reply.MakeErrReply("ERR wrong number args")
	}

	src := string(args[1])
	dst := string(args[2])

	srcPeer := cluster.peerPick.PickNode(src)
	dstPeer := cluster.peerPick.PickNode(dst)

	if srcPeer != dstPeer {
		return reply.MakeErrReply("ERR rename must within peer")
	}

	return cluster.relay(srcPeer, c, args)
}
