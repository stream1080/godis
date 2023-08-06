package redis

type Connection interface {
	Write([]byte) (int, error) // 写入数据
	GetDBIndex() int
	SelectDB(int)
}
