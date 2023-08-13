package database

import (
	"github.com/stream1080/godis/interface/resp"
	"github.com/stream1080/godis/resp/reply"
)

func Ping(db *DB, args [][]byte) resp.Reply {
	return &reply.PongReply{}
}

func init() {
	RegisterCommand("ping", Ping, 1)
}
