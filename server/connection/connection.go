package connection

type Connection interface {
	Close()
	Write([]byte) error
}
