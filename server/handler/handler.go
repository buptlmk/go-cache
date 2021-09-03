package handler

import (
	"context"
	"net"
)

type Handler interface {
	Close()
	Handle(ctx context.Context, conn net.Conn)
}
