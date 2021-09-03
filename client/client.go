package main

import (
	"bufio"
	"flag"
	"fmt"
	"go-cache/internal"
	"go-cache/internal/codec"
	"go-cache/internal/protocol"
	"go-cache/log"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
)

const InputFlag = ">"

type client struct {
	conn       net.Conn
	closedFlag int32
	sig        chan os.Signal
}

func newClient(conn net.Conn) *client {
	c := &client{
		conn:       conn,
		closedFlag: 0,
		sig:        make(chan os.Signal, 1),
	}
	signal.Notify(c.sig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	return c
}

func (c *client) Close() {
	close(c.sig)
	atomic.StoreInt32(&(c.closedFlag), 1)
	c.conn.Close()
}

func (c *client) isClosed() bool {
	return atomic.LoadInt32(&(c.closedFlag)) == int32(1)
}

func main() {

	port := flag.String("--port", "4399", "your listening port")

	conn, err := net.Dial("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		panic(err)
	}
	Print("connected successfully....")

	c := newClient(conn)

	go func() {
		defer func() {
			c.Close()

		}()
		select {
		case <-c.sig:
			return
		}
	}()

	// 回显
	go func() {
		buffer := bufio.NewReader(conn)
		for !c.isClosed() {
			msg, err := protocol.Unpack(protocol.PlainType, buffer)
			if err == io.EOF {
				c.Close()
				break
			}
			if err != nil {
				Print(err)
				continue
			}
			if msg.StatusCode == protocol.StatusErr {
				Print(msg.Error)
				continue
			}
			data := &internal.Payload{}
			err = codec.GetCodec(msg.CodeType).Decode(msg.Data, data)
			if err != nil {
				Print(err)
				continue
			}

			Print(data.Value)

		}

	}()

	// 读
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() && !c.isClosed() {

		text := scanner.Text()
		text = removeSpace(text)
		args := strings.SplitN(text, " ", 3)

		if len(args) < 2 {
			Print("invalid input")
			continue
		}
		args = append(args, "")
		payload := internal.Payload{
			Command: args[0],
			Key:     args[1],
			Value:   args[2],
		}

		data, err := codec.GetCodec(codec.Msg).Encode(payload)
		if err != nil {
			log.Logger.Error(err)
			continue
		}

		msg := protocol.NewMessage()
		msg.Data = data
		msg.CodeType = codec.Msg
		msg.MsgType = protocol.MessageRequest

		bytes, err := protocol.Pack(protocol.PlainType, msg)
		if err != nil {
			log.Logger.Error(err)
			continue
		}
		conn.Write(bytes)

	}

}

func Print(args interface{}) {
	fmt.Print(args, "\n", InputFlag)
}

func removeSpace(s string) string {
	s = strings.TrimSpace(s)
	bytes := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' && i >= 1 && s[i-1] == ' ' {
			continue
		}
		bytes = append(bytes, s[i])
	}
	return string(bytes)

}
