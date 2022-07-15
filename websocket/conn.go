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

type Conn interface {
	// GetID 获取唯一ID
	GetID() int64

	// Close 关闭连接
	Close() error

	// SendTextMsg 发送文本消息
	SendTextMsg(...[]byte) error

	// SendBinaryMsg 发送二进制消息
	SendBinaryMsg(...[]byte) error

	// RemoteAddr 获取客户端地址
	RemoteAddr() net.Addr
}

type wsConn struct {
	minMsgSize  int
	maxMsgSize  int
	closedTag   int32
	id          int64
	conn        *websocket.Conn
	writeBuffer chan *message
}

func newWsConn(id int64, c *websocket.Conn, config *Config) *wsConn {
	writeBufferSize := 64
	if config.WriteBufferSize > 0 {
		writeBufferSize = config.WriteBufferSize
	}
	return &wsConn{
		minMsgSize:  config.MinMsgSize,
		maxMsgSize:  config.MaxMsgSize,
		closedTag:   open,
		id:          id,
		conn:        c,
		writeBuffer: make(chan *message, writeBufferSize),
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

func (c *wsConn) SendTextMsg(message ...[]byte) error {
	if c.isClosed() {
		return ErrSendOnClosedConn
	}
	for _, m := range message {
		size := len(m)
		if c.maxMsgSize > 0 && size > c.maxMsgSize {
			return ErrMessageTooLong
		}
		if c.minMsgSize > 0 && size < c.minMsgSize {
			return ErrMessageTooShort
		}
		msg := getMessage()
		msg.mType = websocket.TextMessage
		msg.conn = c
		msg.write(m)
		c.writeBuffer <- msg
	}
	return nil
}

func (c *wsConn) SendBinaryMsg(message ...[]byte) error {
	if c.isClosed() {
		return ErrSendOnClosedConn
	}
	for _, m := range message {
		size := len(m)
		if c.maxMsgSize > 0 && size > c.maxMsgSize {
			return ErrMessageTooLong
		}
		if c.minMsgSize > 0 && size < c.minMsgSize {
			return ErrMessageTooShort
		}
		msg := getMessage()
		msg.mType = websocket.BinaryMessage
		msg.write(m)
		c.writeBuffer <- msg
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
		err := c.conn.WriteMessage(m.mType, m.data)
		putMessage(m)
		if err != nil {
			break
		}
	}
}

func (c *wsConn) isClosed() (b bool) {
	return atomic.LoadInt32(&c.closedTag) == closed
}
