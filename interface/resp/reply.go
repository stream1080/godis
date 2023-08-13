package resp

type Reply interface {
	ToBytes() []byte // 转换为 byte 数组
}
