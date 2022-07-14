package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pyihe/go-example/tcp"
	"github.com/pyihe/go-pkg/snowflakes"
	"github.com/pyihe/go-pkg/syncs"
)

func main() {
	config := tcp.Config{
		Ticker:     true,
		HeaderSize: 4,
		MaxMsgSize: 4 * 1024,
		//ReadBuffer:  1024,
		Port:        9999,
		IP:          "localhost",
		ReadTimeout: 5 * time.Second,
		//TLSConfig: &tcp.TLSConfig{
		//	ServerCert: "../certs/server.pem",
		//	ServerKey:  "../certs/server.key",
		//	RootCa: "../certs/client.pem", // 这里使用客户端的证书作为根证书, 实际使用中应该是自己计算机内置的根证书
		//},
	}

	closer, err := tcp.Run(&config, &EchoServer{})
	if err != nil {
		fmt.Printf("run err: %v\n", err)
		return
	}
	defer closer.Close()

	//
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT)
	<-interrupt
}

// EchoServer echo服务器
type EchoServer struct {
	count syncs.AtomicInt64 // 记录当前的连接数
}

// OnMessage 收到消息时执行
func (s *EchoServer) OnMessage(conn tcp.Conn, data []byte) {
	conn.SendMsg(data)
}

// OnTick 定时任务
func (s *EchoServer) OnTick() (time.Duration, bool) {
	now := time.Now().Format("15:04:05")
	println(now, s.count.Value())
	return 1 * time.Second, false
}

// OnConnect 新连接建立时执行
func (s *EchoServer) OnConnect(conn tcp.Conn) {
	s.count.Inc(1)
}

// OnClose 连接关闭时执行
func (s *EchoServer) OnClose(conn tcp.Conn) {
	s.count.Inc(-1)
}

// NewUniqueID 获取全局唯一ID
func (s *EchoServer) NewUniqueID() int64 {
	return snowflakes.NewWorker(1).GetInt64()
}
