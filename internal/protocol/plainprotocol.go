package protocol

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
)

type Plain struct {
}

func (p *Plain) Pack(message *Message) ([]byte, error) {

	topBytes := []byte{magic, version}

	headerBytes, err := json.Marshal(message.Header)
	if err != nil {
		return nil, err
	}

	totalLen := 4 + len(headerBytes) + len(message.Data)
	totalLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(totalLenBytes, uint32(totalLen))
	headerLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(headerLenBytes, uint32(len(headerBytes)))

	// 6: magic + version + totallen
	data := make([]byte, 6+totalLen)

	copy(data[:2], topBytes)
	copy(data[2:6], totalLenBytes)
	copy(data[6:10], headerLenBytes)
	copy(data[10:10+len(headerBytes)], headerBytes)
	copy(data[10+len(headerBytes):], message.Data)
	return data, nil
}

func (p *Plain) Unpack(r io.Reader) (*Message, error) {
	top2Bytes := make([]byte, 2)

	_, err := io.ReadFull(r, top2Bytes)
	if err != nil {
		return nil, err
	}
	if !check(top2Bytes) {
		return nil, errors.New("magic and version is not right")
	}

	totalLenBytes := make([]byte, 4)
	_, err = io.ReadFull(r, totalLenBytes)
	if err != nil {
		return nil, err
	}

	totalLen := int(binary.BigEndian.Uint32(totalLenBytes))
	if totalLen < 4 {
		return nil, errors.New("invalid length")
	}

	data := make([]byte, totalLen)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return nil, err
	}

	headerLen := int(binary.BigEndian.Uint32(data[:4]))
	headerBytes := data[4 : 4+headerLen]
	h := &Header{}
	err = json.Unmarshal(headerBytes, h)
	if err != nil {
		return nil, err
	}
	msg := new(Message)
	msg.Header = *h
	msg.Data = data[4+headerLen:]
	return msg, nil

}

func check(v []byte) bool {

	return v[0] == magic && v[1] == version
}
