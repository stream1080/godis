package tcp

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/stream1080/go-redis/interface/tcp"
	"github.com/stream1080/go-redis/lib/logger"
)

type Config struct {
	Address string // 服务地址
}

// 通过配置启动一个服务
func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {

	closeChan := make(chan struct{})
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	// 接收系统中的关闭信号
	go func() {
		switch <-signalChan {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	// 监听连接
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}

	logger.Info("start listen, addr:", cfg.Address)

	ListenAndServe(listener, handler, closeChan)

	return nil
}

// 通过 listener 启动服务
func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {
	var (
		ctx      = context.Background()
		waitDone sync.WaitGroup
	)

	// 关闭资源
	go func() {
		<-closeChan
		logger.Info("shutting down...")
		_ = listener.Close()
		_ = handler.Close()
	}()

	// 关闭连接
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()

	for {
		// 建立连接
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("listener accept")
		waitDone.Add(1)
		go func() {
			defer waitDone.Done()
			handler.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}
