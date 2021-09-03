package protocol

import (
	"go-cache/internal/codec"
	"io"
	"sync"
)

// 发送数据格式定义

//-------------------------------------------------------------------------------------------------
//|1byte|1byte  |4byte       |4byte        | header length |(total length - header length - 4byte)|
//-------------------------------------------------------------------------------------------------
//|magic|version|total length|header length|     header    |                    body              |
//-------------------------------------------------------------------------------------------------

type Protocol interface {
	Pack(*Message) ([]byte, error)
	Unpack(r io.Reader) (*Message, error)
}

var protocols = map[ProtocolType]Protocol{PlainType: &Plain{}}
var lock = sync.RWMutex{}

func register(key ProtocolType, value Protocol) {
	protocols[key] = value
}

type Header struct {
	Seq        uint64            // 序列号，标识唯一请求
	MsgType    MessageType       // 消息的类型
	CodeType   codec.CodeType    // 数据的编码类型
	StatusCode StatusCode        // 请求状态码
	Error      string            //
	MetaData   map[string]string //其他属性
}

type Message struct {
	Header
	Data []byte
}

func NewMessage() *Message {
	return &Message{}
}

func Pack(t ProtocolType, m *Message) ([]byte, error) {
	lock.RLock()
	p := protocols[t]
	lock.RUnlock()
	return p.Pack(m)
}

func Unpack(t ProtocolType, r io.Reader) (*Message, error) {
	lock.RLock()
	p := protocols[t]
	lock.RUnlock()
	return p.Unpack(r)
}
