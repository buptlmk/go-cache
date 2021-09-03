package connection

import "net"

type PlainConnection struct {
	conn net.Conn
}

func (p *PlainConnection) Close() {
	p.conn.Close()
}

func (p *PlainConnection) Write(bytes []byte) error {
	_, err := p.conn.Write(bytes)
	return err
}

func NewPlainConnection(conn net.Conn) *PlainConnection {
	return &PlainConnection{
		conn: conn,
	}
}
