package database

import (
	"strings"

	"github.com/stream1080/godis/datastruct/dict"
	"github.com/stream1080/godis/interface/database"
	"github.com/stream1080/godis/interface/resp"
	"github.com/stream1080/godis/resp/reply"
)

type DB struct {
	index  int
	data   dict.Dict
	addAof func(CmdLine)
}

type (
	ExecFunc func(db *DB, args [][]byte) resp.Reply
	CmdLine  = [][]byte
)

func MakeDB() *DB {
	return &DB{
		data:   dict.MakeSyncDict(),
		addAof: func(line CmdLine) {},
	}
}
func (db *DB) Exec(c resp.Connection, cmdLine CmdLine) resp.Reply {
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return reply.MakeErrReply("ERR unknown command" + cmdName)
	}

	if !validateArity(cmd.arity, cmdLine) {
		return reply.MakeArgNumErrReply(cmdName)
	}

	return cmd.executor(db, cmdLine[1:])
}

func validateArity(arity int, cmdArgs [][]byte) bool {
	argNum := len(cmdArgs)
	if arity >= 0 {
		return argNum == arity
	}

	return argNum >= -arity
}

func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {
	raw, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}

	return raw.(*database.DataEntity), true
}

func (db *DB) PutEntity(key string, entity *database.DataEntity) int {
	return db.data.Put(key, entity)
}

func (db *DB) PutIfExists(key string, entity *database.DataEntity) int {
	return db.data.PutIfExists(key, entity)
}

func (db *DB) PutIfAbsen(key string, entity *database.DataEntity) int {
	return db.data.PutIfAbsent(key, entity)
}

func (db *DB) Remove(key string) {
	db.data.Remove(key)
}

func (db *DB) Removes(keys ...string) int {
	deleted := 0
	for _, key := range keys {
		_, exists := db.data.Get(key)
		if exists {
			db.Remove(key)
			deleted++
		}
	}

	return deleted
}

func (db *DB) Flush() {
	db.data.Clear()
}
