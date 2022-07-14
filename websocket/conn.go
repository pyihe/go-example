package websocket

import "github.com/gorilla/websocket"

type Conn interface {
	// GetID 获取唯一ID
	GetID() int64

	// Close 关闭连接
	Close() error

	// SendMsg 发送消息
	SendMsg(...[]byte) error
}

type wsConn struct {
	conn *websocket.Conn
}

func newWsConn(c *websocket.Conn) *wsConn {
	return &wsConn{conn: c}
}

func (c *wsConn) GetID() int64 {
	// TODO
	return 0
}

func (c *wsConn) Close() error {
	// TODO
	return nil
}

func (c *wsConn) SendMsg(message ...[]byte) error {
	// TODO
	return nil
}
