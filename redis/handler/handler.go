package handler

import (
	"sync"

	databaseface "github.com/stream1080/godis/interface/database"
	"github.com/stream1080/godis/lib/logger"
	"github.com/stream1080/godis/lib/sync/atomic"
	"github.com/stream1080/godis/redis/conn"
)

type RespHandler struct {
	activeConn sync.Map              // 活跃连接
	db         databaseface.Database // database
	closing    atomic.Boolean        // 是否关闭
}

func (r *RespHandler) Close() error {
	logger.Info("handler shutting down")
	r.closing.Set(true)
	r.activeConn.Range(
		func(key, value any) bool {
			client := key.(*conn.Connection)
			_ = client.Close()
			return true
		},
	)
	r.db.Close()
	return nil
}
