package handler

import (
	"context"
	"io"
	"net"
	"strings"
	"sync"

	databaseface "github.com/stream1080/godis/interface/database"
	"github.com/stream1080/godis/lib/logger"
	"github.com/stream1080/godis/lib/sync/atomic"
	connection "github.com/stream1080/godis/redis/conn"
	"github.com/stream1080/godis/redis/parser"
	"github.com/stream1080/godis/redis/protocol"
)

type RespHandler struct {
	activeConn sync.Map              // 活跃连接
	db         databaseface.Database // database
	closing    atomic.Boolean        // 是否关闭
}

func (r *RespHandler) closeClient(client *connection.Connection) {
	_ = client.Close()
	// TODO: r.db.AfterClientClose(client)
	r.db.AfterClientClose(nil)
	r.activeConn.Delete(client)
}

func (r *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if r.closing.Get() {
		_ = conn.Close()
		return
	}
	client := connection.NewConn(conn)
	r.activeConn.Store(client, struct{}{})

	ch := parser.ParseStream(conn)
	for payload := range ch {
		// handler error
		if payload.Err != nil {
			if payload.Err == io.EOF ||
				payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				r.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			// protocol error
			errReply := protocol.MakeErrReply(payload.Err.Error())
			err := client.Write(errReply.ToBytes())
			if err != nil {
				r.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			continue
		}
		// exec

	}

}

func (r *RespHandler) Close() error {
	logger.Info("handler shutting down")
	r.closing.Set(true)
	r.activeConn.Range(
		func(key, value any) bool {
			client := key.(*connection.Connection)
			_ = client.Close()
			return true
		},
	)
	r.db.Close()
	return nil
}
