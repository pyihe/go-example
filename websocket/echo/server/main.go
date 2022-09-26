package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pyihe/go-example/websocket"
	"github.com/pyihe/go-pkg/snowflakes"
	"github.com/pyihe/go-pkg/syncs"
)

func main() {
	config := &websocket.Config{
		Tick:            true,
		Addr:            "127.0.0.1:8888",
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    5 * time.Second,
		MaxMsgSize:      4 * 1024,
		MinMsgSize:      0,
		MaxConns:        1024,
		ReadBufferSize:  128,
		WriteBufferSize: 128,
		//CertFile:        "../certs/server.pem",
		//KeyFile:         "../certs/server.key",
	}
	s, err := websocket.Run(config, &echoServer{})
	if err != nil {
		fmt.Printf("run server fail: %v\n", err)
		return
	}
	defer s.Close()

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-interrupt
}

type echoServer struct {
	count syncs.AtomicInt64
}

func (s *echoServer) OnMessage(conn websocket.Conn, message []byte) {
	fmt.Printf("[%v][%v]OnMessage: %v\n", conn.GetID(), conn.RemoteAddr(), string(message))
	conn.SendTextMsg(message)
}

func (s *echoServer) OnConnect(conn websocket.Conn) {
	fmt.Printf("[%v] OnCennect\n", conn.RemoteAddr())
	s.count.Inc(1)
}

func (s *echoServer) OnClose(conn websocket.Conn) {
	fmt.Printf("[%v] OnClose\n", conn.RemoteAddr())
	s.count.Inc(-1)
}

func (s *echoServer) OnTick() (time.Duration, bool) {
	fmt.Println(time.Now().Format("15:04:05"), s.count.Value())
	return 1 * time.Second, false
}

func (s *echoServer) OnPing(conn websocket.Conn, pingMessage string) {
	fmt.Printf("[%v] OnPing\n", conn.RemoteAddr())
}

func (s *echoServer) GenerateID() int64 {
	return snowflakes.NewWorker(1).GetInt64()
}
