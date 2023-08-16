package database

import (
	"github.com/stream1080/godis/interface/resp"
	"github.com/stream1080/godis/resp/reply"
)

func ExecGet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeNullBulkReply()
	}

	return reply.MakeBulkReply(entity.Data.([]byte))
}
