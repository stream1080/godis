package parser

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/stream1080/godis/interface/redis"
	"github.com/stream1080/godis/redis/protocol"
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

func parseMultiBulkHeader(msg []byte, state *readState) error {

	expectedLine, err := strconv.ParseUint(string(msg[1:len(msg)-2]), 10, 32)
	if err != nil {
		return errors.New("protocol error: " + string(msg))
	}

	if expectedLine == 0 {
		state.expectedArgsCount = 0
		return nil
	} else if expectedLine > 0 {
		state.msgType = msg[0]
		state.readingMultiLine = true
		state.expectedArgsCount = int(expectedLine)
		state.args = make([][]byte, 0, expectedLine)
		return nil
	} else {
		return errors.New("protocol error: " + string(msg))
	}
}

// +OK\r\n -err\r\n :5\r\n
func parseSingleLineReply(msg []byte) (redis.Reply, error) {
	str := strings.TrimSuffix(string(msg), "\r\n")
	var result redis.Reply
	switch msg[0] {
	case '+':
		result = protocol.MakeStatusReply(str[1:])
	case '-':
		result = protocol.MakeErrReply(str[1:])
	case ':':
		val, err := strconv.ParseInt(str[1:], 10, 64)
		if err != nil {
			return nil, errors.New("protocol error: " + string(msg))
		}
		result = protocol.MakeIntReply(val)
	}

	return result, nil
}

// PING\r\n
func readBody(msg []byte, state *readState) error {
	line := msg[0 : len(msg)-2]
	var err error

	// $3
	if line[0] == '$' {
		state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return errors.New("protocol error: " + string(msg))
		}
		// $0\r\n
		if state.bulkLen <= 0 {
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	} else {
		state.args = append(state.args, line)
	}

	return nil
}
