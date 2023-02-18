package protocol

type ErrorReply interface {
	Error() error
	ToBytes() []byte
}
