package conn

import (
	"net"
	"sync"
	"time"

	"github.com/stream1080/godis/lib/sync/wait"
)

type Connection struct {
	conn         net.Conn
	waitingReply wait.Wait
	mu           sync.Mutex
	selecteDB    int
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connection) Write(bytes []byte) error {
	return nil
}

func (c *Connection) GetDBIndex() int {
	return c.selecteDB
}

func (c *Connection) SelectDB(dbNum int) {
	c.selecteDB = dbNum
}

func (c *Connection) Close() error {
	c.waitingReply.WaitWithTimeout(10 * time.Second)
	return c.conn.Close()
}
