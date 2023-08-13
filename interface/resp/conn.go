package resp

type Connection interface {
	Write([]byte) error // 写入数据
	GetDBIndex() int
	SelectDB(int)
}
