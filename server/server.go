package server

import (
	"context"
	"fmt"
	"go-cache/config"
	"go-cache/internal/syncx"
	"go-cache/server/handler"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type server struct {
	miss   int
	thresh int
	wg     syncx.WaitGroup
}

func (s *server) getMiss() int {
	return s.miss
}

func (s *server) addMiss(x int) {
	s.miss += x
}

func (s *server) resetMiss() {
	s.miss = 0
}

func Start(handler handler.Handler) {

	sig := make(chan os.Signal)

	signal.Notify(sig, syscall.SIGHUP, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGTERM)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.GlobalConfig.Address, config.GlobalConfig.Port))

	if err != nil {
		panic(err)
	}
	defer func() {
		listener.Close()
		handler.Close()
	}()

	go func() {
		select {
		case <-sig:
			listener.Close()
			handler.Close()
		}
	}()

	server := &server{
		thresh: 5,
	}

	for {
		conn, err := listener.Accept()
		// 过程中突然出现了一次连接失败,先不着急退出
		if err != nil {
			if server.getMiss() > server.thresh {
				break
			}
			server.addMiss(1)
			continue
		}
		// 没问题后
		server.resetMiss()

		// 开始处理
		server.wg.Add(1)

		go func() {
			defer func() {
				server.wg.Done()
			}()

			handler.Handle(context.TODO(), conn)
		}()
	}

	// 做这一步的意义在于说，虽然服务器要关了，但是要等正在连接中的连接处理完,限制5秒内
	server.wg.WaitTime(time.Second * 5)

}
