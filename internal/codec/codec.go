package codec

import (
	"encoding/json"
	"github.com/vmihailenco/msgpack/v5"
	"sync"
)

type Codec interface {
	Encode(v interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
}

type CodeType byte

const (
	Json CodeType = iota
	Msg
)

var codecs = map[CodeType]Codec{
	Json: &JsonCodec{},
	Msg:  &MsgCodec{},
}
var lock = sync.RWMutex{}

type JsonCodec struct {
}

func (j *JsonCodec) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *JsonCodec) Decode(bytes []byte, value interface{}) error {
	return json.Unmarshal(bytes, value)
}

func GetCodec(key CodeType) Codec {
	lock.RLock()
	defer lock.RUnlock()
	return codecs[key]
}

type MsgCodec struct {
}

func (m *MsgCodec) Encode(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (m *MsgCodec) Decode(bytes []byte, v interface{}) error {
	return msgpack.Unmarshal(bytes, v)
}
