package handler

import (
	"bufio"
	"context"
	"go-cache/internal"
	"go-cache/internal/codec"
	"go-cache/internal/protocol"
	"go-cache/log"
	"go-cache/server/connection"
	"go-cache/server/db"
	"io"
	"net"
	"sync"
	"sync/atomic"
)

type Deal struct {
	conns  sync.Map
	status uint32 // 0 开启 1 关闭
	db     *db.DB
}

func (d *Deal) Close() {
	atomic.StoreUint32(&(d.status), 1)

	// 将所有的连接一一关闭
	d.conns.Range(func(key, value interface{}) bool {
		c := key.(connection.Connection)
		c.Close()
		return true
	})

}
func (d *Deal) IsClosed() bool {
	return atomic.LoadUint32(&(d.status)) == uint32(1)
}

func (d *Deal) Handle(ctx context.Context, conn net.Conn) {

	if d.IsClosed() {
		conn.Close()
		return
	}

	client := connection.NewPlainConnection(conn)
	d.conns.Store(client, struct{}{})

	// 解析client
	buffer := bufio.NewReader(conn)
	for {
		msg, err := protocol.Unpack(protocol.PlainType, buffer)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Logger.Error(err.Error())
			h := protocol.Header{
				MsgType:    protocol.MessageResponse,
				CodeType:   codec.Msg,
				StatusCode: protocol.StatusErr,
				Error:      err.Error(),
			}
			data := &protocol.Message{
				Header: h,
			}
			bytes, _ := protocol.Pack(protocol.PlainType, data)
			conn.Write(bytes)
			continue
		}
		value := &internal.Payload{}
		codec.GetCodec(msg.CodeType).Decode(msg.Data, value)
		log.Logger.Info(value)

		res := d.exec(value.Command, value.Key, value.Value)
		// 返回
		bytes, _ := protocol.Pack(protocol.PlainType, res)
		conn.Write(bytes)
	}
	d.conns.Delete(client)
	log.Logger.Info(conn.RemoteAddr().String(), " logged out")

}

func NewDeal() *Deal {
	return &Deal{
		conns:  sync.Map{},
		status: 0,
		db:     db.NewDB(),
	}
}

func (d *Deal) exec(command string, key string, value interface{}) *protocol.Message {
	p := db.ExecCmd(d.db, command, key, value)

	bytes, _ := codec.GetCodec(codec.Msg).Encode(*p)

	h := protocol.Header{
		MsgType:  protocol.MessageResponse,
		CodeType: codec.Msg,
	}
	m := &protocol.Message{
		Header: h,
		Data:   bytes,
	}
	return m
}
