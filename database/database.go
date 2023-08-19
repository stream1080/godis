package database

import (
	"strconv"

	"github.com/stream1080/godis/config"
	"github.com/stream1080/godis/interface/resp"
	"github.com/stream1080/godis/resp/reply"
)

type Database struct {
	dbSet []*DB
}

func NewDatabase() *Database {
	database := &Database{}
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}

	database.dbSet = make([]*DB, config.Properties.Databases)
	for i := range database.dbSet {
		db := MakeDB()
		db.index = i
		database.dbSet[i] = db
	}

	return database
}

func (db *Database) Exec(client resp.Connection, args [][]byte) resp.Reply {
	return reply.MakeMultiBulkReply(args)
}

func (db *Database) Close() {

}

func (db *Database) AfterClientClose(c resp.Connection) {

}

func execSelect(c resp.Connection, database Database, args [][]byte) resp.Reply {

	index, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.MakeErrReply("ERR invalid DB index")
	}

	if index > len(database.dbSet) {
		return reply.MakeErrReply("ERR DB index is out of range")
	}

	c.SelectDB(index)

	return reply.MakeOkReply()
}
