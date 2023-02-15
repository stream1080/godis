package tcp

import (
	"context"
	"net"
)

// 处理 tcp 服务器应用
type Handler interface {
	Handle(ctx context.Context, conn net.Conn) // 处理逻辑
	Close() error                              // 关闭
}
