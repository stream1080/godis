package redis

type Connection interface {
	Write([]byte) (int, error) // 写入数据
	Close() error              //关闭
}
