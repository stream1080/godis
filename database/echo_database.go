package database

import (
	"github.com/stream1080/godis/interface/resp"
	"github.com/stream1080/godis/resp/reply"
)

type EchoDatabase struct {
}

func NewEchoDatabase() *EchoDatabase {
	return &EchoDatabase{}
}

func (e *EchoDatabase) Exec(client redis.Connection, args [][]byte) redis.Reply {
	return protocol.MakeMultiBulkReply(args)
}

func (e *EchoDatabase) Close() {

}

func (e *EchoDatabase) AfterClientClose(c redis.Connection) {

}
