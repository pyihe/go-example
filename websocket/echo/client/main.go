package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pyihe/go-pkg/maps"
)

var conns = maps.NewMap()

func main() {
	var closeChan = make(chan *client, 10) // 传递被服务器关闭的连接
	var interrupt = make(chan os.Signal)   // 接收中断信号
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT)
	// 开启连接
	for i := 0; i < 100; i++ {
		c := dial("ws://127.0.0.1:8888")
		if c != nil {
			conns.Set(c, struct{}{})
			go func(cli *client) {
				cli.send()
			}(c)
			go func(cli *client) {
				cli.read()
				closeChan <- cli
			}(c)
		}
	}

	select {
	case <-interrupt: // 如果是客户端被终止，则断开所有连接并删除
		conns.LockRange(func(k interface{}, v interface{}) bool {
			c, ok := k.(*client)
			if ok {
				c.close()
				conns.UnsafeDel(c)
			}
			return false
		})

	case c := <-closeChan: // 如果是被服务器断开，删除连接即可
		if c != nil {
			conns.Del(c)
		}
		if conns.Len() == 0 {
			return
		}
	}
}

type client struct {
	ctx    context.Context
	cancel context.CancelFunc
	conn   *websocket.Conn
}

func dial(addr string) *client {
	c, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		fmt.Printf("dial err: %v\n", err)
		return nil
	}
	cc := &client{
		conn: c,
	}
	cc.ctx, cc.cancel = context.WithCancel(context.Background())
	return cc
}

func (c *client) send() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			//err := c.conn.WriteMessage(websocket.BinaryMessage, []byte("BinaryMessage"))
			//if err != nil {
			//	return
			//}
			err := c.conn.WriteMessage(websocket.TextMessage, []byte("TextMessage"))
			if err != nil {
				return
			}
			//err := c.conn.WriteMessage(websocket.PingMessage, []byte("PingMessage"))
			//if err != nil {
			//	return
			//}
			//err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			//if err != nil {
			//	return
			//}

			time.Sleep(1 * time.Second)
		}
	}
}

func (c *client) read() {
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Printf("read err: %v\n", err)
			break
		}
		fmt.Printf("[%v]: %v\n", c.conn.LocalAddr(), string(data))
	}
}

func (c *client) close() {
	c.cancel()
	c.conn.Close()
}
