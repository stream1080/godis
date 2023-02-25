package parser

import (
	"bufio"
	"errors"
	"io"

	"github.com/stream1080/godis/interface/redis"
)

// 客户端载荷
type Payload struct {
	Data redis.Reply
	Err  error
}

type readState struct {
	readingMultiLine  bool     // 是否为多行
	expectedArgsCount int      // 参数个数
	msgType           byte     // 消息类型
	args              [][]byte // 参数
	bulkLen           int64    // 数组长度
}

// 是否解析完成
func (r *readState) finished() bool {
	return r.expectedArgsCount > 0 && r.expectedArgsCount == len(r.args)
}

func parseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse0(reader, ch)
	return ch
}

func parse0(reader io.Reader, ch chan<- *Payload) {

}

// 读取一行
func readLine(bufReader *bufio.Reader, state *readState) ([]byte, bool, error) {
	var msg []byte
	var err error

	// 以 \r\n 切分
	if state.bulkLen == 0 {
		msg, err = bufReader.ReadBytes('\n')
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\n' {
			return nil, false, errors.New("protocol error: " + string(msg))
		}
	} else {
		// 读到 $ 数字
		msg = make([]byte, state.bulkLen+2)
		_, err = io.ReadFull(bufReader, msg)
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' || msg[len(msg)-1] != '\n' {
			return nil, false, errors.New("protocol error: " + string(msg))
		}
		state.bulkLen = 0
	}
	return msg, false, nil
}
