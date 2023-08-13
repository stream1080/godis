package parser

import (
	"bufio"
	"errors"
	"io"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/stream1080/godis/interface/resp"
	"github.com/stream1080/godis/lib/logger"
	"github.com/stream1080/godis/resp/reply"
)

// 客户端载荷
type Payload struct {
	Data resp.Reply
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

func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse0(reader, ch)
	return ch
}

func parse0(reader io.Reader, ch chan<- *Payload) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	bufReader := bufio.NewReader(reader)
	var ioErr bool
	var state readState
	var err error
	var msg []byte

	for {
		msg, ioErr, err = readLine(bufReader, &state)
		if err != nil {
			if ioErr {
				ch <- &Payload{Err: err}
				close(ch)
				return
			}
			ch <- &Payload{Err: err}
			state = readState{}
			continue
		}

		// 多行解析
		if !state.readingMultiLine {
			if msg[0] == '*' {
				if err := parseMultiBulkHeader(msg, &state); err != nil {
					ch <- &Payload{Err: errors.New("reply error: " + string(msg))}
					state = readState{}
					continue
				}
				if state.expectedArgsCount == 0 {
					ch <- &Payload{Data: &reply.EmptyMultiBulkReply{}}
					state = readState{}
					continue
				}
				// $3\r\n
			} else if msg[0] == '$' {
				err = parseMultiBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{Err: errors.New("reply error: " + string(msg))}
					state = readState{}
					continue
				}
				if state.expectedArgsCount == -1 {
					ch <- &Payload{Data: &reply.NullBulkReply{}}
					state = readState{}
					continue
				}
			} else {
				result, err := parseSingleLineReply(msg)
				ch <- &Payload{Data: result, Err: err}
				state = readState{}
				continue
			}

		} else {
			if err := readBody(msg, &state); err != nil {
				ch <- &Payload{Err: errors.New("reply error: " + string(msg))}
				state = readState{}
				continue
			}
			if state.finished() {
				var result resp.Reply
				if state.msgType == '*' {
					result = reply.MakeMultiBulkReply(state.args)
				} else if state.msgType == '$' {
					result = reply.MakeBulkReply(state.args[0])
				}
				ch <- &Payload{Data: result, Err: err}
				state = readState{}
			}
		}
	}
}

// readLine 读取一行
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
			return nil, false, errors.New("reply error: " + string(msg))
		}
	} else {
		// 读到 $ 数字
		msg = make([]byte, state.bulkLen+2)
		_, err = io.ReadFull(bufReader, msg)
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' || msg[len(msg)-1] != '\n' {
			return nil, false, errors.New("reply error: " + string(msg))
		}
		state.bulkLen = 0
	}
	return msg, false, nil
}

func parseMultiBulkHeader(msg []byte, state *readState) error {

	expectedLine, err := strconv.ParseUint(string(msg[1:len(msg)-2]), 10, 32)
	if err != nil {
		return errors.New("reply error: " + string(msg))
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
		return errors.New("reply error: " + string(msg))
	}
}

// +OK\r\n -err\r\n :5\r\n
func parseSingleLineReply(msg []byte) (resp.Reply, error) {
	str := strings.TrimSuffix(string(msg), "\r\n")
	var result resp.Reply
	switch msg[0] {
	case '+':
		result = reply.MakeStatusReply(str[1:])
	case '-':
		result = reply.MakeErrReply(str[1:])
	case ':':
		val, err := strconv.ParseInt(str[1:], 10, 64)
		if err != nil {
			return nil, errors.New("reply error: " + string(msg))
		}
		result = reply.MakeIntReply(val)
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
			return errors.New("reply error: " + string(msg))
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
