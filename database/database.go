package database

import (
	"strconv"
	"strings"

	"github.com/stream1080/godis/aof"
	"github.com/stream1080/godis/config"
	"github.com/stream1080/godis/interface/resp"
	"github.com/stream1080/godis/lib/logger"
	"github.com/stream1080/godis/resp/reply"
)

type Database struct {
	dbSet      []*DB
	aofHandler *aof.AofHandler
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

	if config.Properties.AppendOnly {
		aofHandler, err := aof.NewAofHandler(database)
		if err != nil {
			panic(err)
		}
		database.aofHandler = aofHandler
		for _, db := range database.dbSet {
			sdb := db
			sdb.addAof = func(line CmdLine) {
				database.aofHandler.AddAof(sdb.index, line)
			}
		}
	}

	return database
}

func (database *Database) Exec(client resp.Connection, args [][]byte) resp.Reply {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	cmd := strings.ToLower(string(args[0]))
	if cmd == "select" {
		if len(args) != 2 {
			return reply.MakeArgNumErrReply("select")
		}
		return execSelect(client, database, args[1:])
	}

	db := database.dbSet[client.GetDBIndex()]
	return db.Exec(client, args)
}

func (database *Database) Close() {

}

func (db *Database) AfterClientClose(c resp.Connection) {

}

func execSelect(c resp.Connection, database *Database, args [][]byte) resp.Reply {

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
