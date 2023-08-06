package database

import (
	"github.com/stream1080/godis/interface/redis"
)

type CmdLine = [][]byte

type Database interface {
	Exec(client redis.Connection, args []byte) redis.Reply
	Close()
	AfterClientClose(c redis.Connection)
}

type DataEntity struct {
	Data interface{}
}
