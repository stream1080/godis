package database

import (
	"github.com/stream1080/godis/interface/database"
	"github.com/stream1080/godis/interface/resp"
	"github.com/stream1080/godis/resp/reply"
)

func init() {
	RegisterCommand("Get", ExecGet, 2)       // get k1
	RegisterCommand("Set", ExecSet, 3)       // set k v
	RegisterCommand("SetNX", ExecSetNX, 3)   // setnx k v
	RegisterCommand("GetSet", ExecGetSet, 3) // getset k v
	RegisterCommand("StrLen", ExecStrLen, 2) // strlen k
}

func ExecGet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeNullBulkReply()
	}

	return reply.MakeBulkReply(entity.Data.([]byte))
}

func ExecSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity := &database.DataEntity{
		Data: args[1],
	}

	db.PutEntity(key, entity)

	return reply.MakeOkReply()
}

func ExecSetNX(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity := &database.DataEntity{
		Data: args[1],
	}

	return reply.MakeIntReply(int64(db.PutIfAbsen(key, entity)))
}

func ExecGetSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	db.PutEntity(key, &database.DataEntity{Data: args[1]})
	if !exists {
		return reply.MakeNullBulkReply()
	}

	return reply.MakeBulkReply(entity.Data.([]byte))
}

func ExecStrLen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeNullBulkReply()
	}

	return reply.MakeIntReply(int64(len(entity.Data.([]byte))))
}
