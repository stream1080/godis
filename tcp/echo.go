package tcp

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/stream1080/godis/lib/logger"
	"github.com/stream1080/godis/lib/sync/atomic"
	"github.com/stream1080/godis/lib/sync/wait"
)

// 客户端
type EchoClient struct {
	Conn    net.Conn  // 连接
	Waiting wait.Wait // 正在处理业务
}

// 关闭客户端
func (c *EchoClient) Close() error {
	c.Waiting.WaitWithTimeout(10 * time.Second)
	c.Conn.Close()
	return nil
}

// 服务器
type EchoHandler struct {
	activeConn sync.Map       // 活跃连接
	closing    atomic.Boolean // 是否关闭
}

func MakeEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

// 处理业务
func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		_ = conn.Close()
		return
	}

	client := &EchoClient{
		Conn: conn,
	}
	h.activeConn.Store(client, struct{}{})

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("conn closed")
				h.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		client.Waiting.Add(1)

		// 转换数据，回写
		b := []byte(msg)
		_, _ = conn.Write(b)
		client.Waiting.Done()
	}
}

// 关闭 stop 的 handler
func (h *EchoHandler) Close() error {
	logger.Info("handler shutting down")
	// 标记关闭
	h.closing.Set(true)
	h.activeConn.Range(func(key, value any) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		return true
	})
	return nil
}
