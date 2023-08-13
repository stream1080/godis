package database

import (
	"github.com/stream1080/godis/datastruct/dict"
	"github.com/stream1080/godis/interface/resp"
	"github.com/stream1080/godis/resp/reply"
	"strings"
)

type DB struct {
	index int
	data  dict.Dict
}

type (
	ExecFunc func(db *DB, args [][]byte) resp.Reply
	CmdLine  = [][]byte
)

func MakeDB() *DB {
	return &DB{
		data: dict.MakeSyncDict(),
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
