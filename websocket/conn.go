package websocket

import (
	"net"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/pyihe/go-pkg/errors"
)

var (
	ErrSendOnClosedConn = errors.New("send message on closed connection")
	ErrMessageTooLong   = errors.New("message too long")
	ErrMessageTooShort  = errors.New("message too short")
)

const (
	open   = 0
	closed = 1
)

type Message struct {
	MessageType int
	Message     []byte
}

type Conn interface {
	// GetID 获取唯一ID
	GetID() int64

	// Close 关闭连接
	Close() error

	// SendMsg 发送消息
	SendMsg(...*Message) error

	// RemoteAddr 获取客户端地址
	RemoteAddr() net.Addr
}

type wsConn struct {
	closedTag   int32
	id          int64
	conn        *websocket.Conn
	server      *wsServer
	writeBuffer chan *Message
}

func newWsConn(c *websocket.Conn, s *wsServer) *wsConn {
	writeBufferSize := 64
	if s.config.WriteBufferSize > 0 {
		writeBufferSize = s.config.WriteBufferSize
	}
	return &wsConn{
		closedTag:   open,
		id:          s.handler.GenerateID(),
		conn:        c,
		server:      s,
		writeBuffer: make(chan *Message, writeBufferSize),
	}
}

func (c *wsConn) GetID() int64 {
	return c.id
}

func (c *wsConn) Close() error {
	// 是否已经关闭
	if atomic.CompareAndSwapInt32(&c.closedTag, open, closed) == false {
		return nil
	}
	close(c.writeBuffer)
	return c.conn.Close()
}

func (c *wsConn) SendMsg(message ...*Message) error {
	if c.isClosed() {
		return ErrSendOnClosedConn
	}
	config := c.server.config
	for _, m := range message {
		m := m
		size := len(m.Message)
		if config.MaxMsgSize > 0 && size > config.MaxMsgSize {
			return ErrMessageTooLong
		}
		if config.MinMsgSize > 0 && size < config.MinMsgSize {
			return ErrMessageTooShort
		}
		c.writeBuffer <- m
	}
	return nil
}

func (c *wsConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *wsConn) writeLoop() {
	for m := range c.writeBuffer {
		if m == nil {
			break
		}
		err := c.conn.WriteMessage(m.MessageType, m.Message)
		if err != nil {
			break
		}
	}
}

func (c *wsConn) isClosed() (b bool) {
	return atomic.LoadInt32(&c.closedTag) == closed
}
