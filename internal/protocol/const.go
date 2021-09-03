package protocol

const magic byte = 0xab
const version byte = 0x01

type MessageType byte

const (
	MessageRequest MessageType = iota
	MessageResponse
)

type StatusCode byte

const (
	StatusOk StatusCode = iota
	StatusErr
)

type ProtocolType byte

const (
	PlainType = iota
)
