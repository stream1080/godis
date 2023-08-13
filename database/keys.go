package database

import (
	"github.com/stream1080/godis/interface/resp"
	"github.com/stream1080/godis/resp/reply"
)

func init() {
	RegisterCommand("DEL", ExecDel, -2)
	RegisterCommand("EXISTS", ExecExists, -2)
	RegisterCommand("flushdb", ExecFlushDB, -1)
	RegisterCommand("type", ExecType, 2)
	RegisterCommand("RENAME", ExecRename, 3)
	RegisterCommand("RENAMENX", ExecRenameNx, 3)
}

func ExecDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}

	return reply.MakeIntReply(int64(db.Removes(keys...)))
}

func ExecExists(db *DB, args [][]byte) resp.Reply {
	result := int64(0)
	for _, arg := range args {
		key := string(arg)
		_, exists := db.GetEntity(key)
		if exists {
			result++
		}
	}

	return reply.MakeIntReply(result)
}

func ExecFlushDB(db *DB, args [][]byte) resp.Reply {
	db.Flush()
	return reply.MakeOkReply()
}

func ExecType(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeStatusReply("none")
	}

	switch entity.Data.(type) {
	case []byte:
		return reply.MakeStatusReply("string")
	}
	// TODO
	return &reply.UnknownErrReply{}
}

func ExecRename(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dest := string(args[1])

	entity, exists := db.GetEntity(src)
	if !exists {
		return reply.MakeErrReply("no such key")
	}

	db.PutEntity(dest, entity)
	db.Remove(src)

	return reply.MakeOkReply()
}

func ExecRenameNx(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dest := string(args[1])

	if _, exists := db.GetEntity(dest); exists {
		return reply.MakeIntReply(0)
	}

	entity, exists := db.GetEntity(src)
	if !exists {
		return reply.MakeErrReply("no such key")
	}

	db.PutEntity(dest, entity)
	db.Remove(src)

	return reply.MakeIntReply(1)
}
