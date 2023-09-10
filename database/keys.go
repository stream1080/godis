package database

import (
	"github.com/stream1080/godis/interface/resp"
	"github.com/stream1080/godis/lib/utils"
	"github.com/stream1080/godis/lib/wildcard"
	"github.com/stream1080/godis/resp/reply"
)

func init() {
	RegisterCommand("DEL", ExecDel, -2)          // DEL k1,k2,k3
	RegisterCommand("EXISTS", ExecExists, -2)    // EXISTS k1,k2,k3
	RegisterCommand("flushdb", ExecFlushDB, -1)  // FLUSHDB a,b,a
	RegisterCommand("type", ExecType, 2)         // TYPE k1
	RegisterCommand("RENAME", ExecRename, 3)     // RENAME k1,k2
	RegisterCommand("RENAMENX", ExecRenameNx, 3) // RENAMENX k1,k2
	RegisterCommand("KEYS", ExecKeys, 2)         // KEYS *
}

func ExecDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}
	deleted := db.Removes(keys...)
	if deleted > 0 {
		db.addAof(utils.ToCmdLine2("del", args...))
	}

	return reply.MakeIntReply(int64(deleted))
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
	db.addAof(utils.ToCmdLine2("flushdb", args...))
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
	db.addAof(utils.ToCmdLine2("rename", args...))

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
	db.addAof(utils.ToCmdLine2("renamenx", args...))

	return reply.MakeIntReply(1)
}

func ExecKeys(db *DB, args [][]byte) resp.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0)
	db.data.ForEach(func(key string, val interface{}) bool {
		if pattern.IsMatch(key) {
			result = append(result, []byte(key))
		}
		return true
	})

	return reply.MakeMultiBulkReply(result)
}
