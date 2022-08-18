package tcp

import (
	"net"
	"sync/atomic"
	"time"

	"github.com/pyihe/go-pkg/buffers"
	"github.com/pyihe/go-pkg/bytes"
)

const (
	open   = 0
	closed = 1
)

// Conn 底层连接
type Conn interface {
	// GetID 获取连接唯一标识
	GetID() int64

	//RemoteAddr 获取客户端地址
	RemoteAddr() net.Addr

	// SendMsg 发送消息, 可变参数表示一次性可以发送多条消息(消息顺序与传递一致), 但每条消息仍然会以单独的数据包在TCP连接中被发送
	SendMsg(...[]byte) error

	// SendMsgWithTimeout 带超时机制发送消息
	SendMsgWithTimeout(...[]byte) error

	// Close 服务器主动关闭连接
	Close() error
}

type tcpConn struct {
	closeTag int32      // 是否关闭: 避免重复关闭
	id       int64      // 唯一ID
	conn     net.Conn   // 底层连接
	server   *tcpServer // 属于哪个服务器
}

func newTCPConn(conn net.Conn, s *tcpServer) *tcpConn {
	return &tcpConn{
		id:     s.handler.GenerateID(),
		conn:   conn,
		server: s,
	}
}

func (c *tcpConn) GetID() int64 {
	return c.id
}

func (c *tcpConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *tcpConn) SendMsg(message ...[]byte) (err error) {
	buf := buffers.Get()
	for _, m := range message {
		data, err := c.server.pkt.Packet(m)
		if err != nil {
			break
		}
		buf.Write(data)
	}
	m := bytes.Copy(buf.Bytes())
	buffers.Put(buf)
	if err == nil {
		_, err = c.conn.Write(m)
	}
	return
}

func (c *tcpConn) SendMsgWithTimeout(message ...[]byte) (err error) {
	buf := buffers.Get()
	for _, m := range message {
		data, err := c.server.pkt.Packet(m)
		if err != nil {
			break
		}
		buf.Write(data)
	}
	m := bytes.Copy(buf.Bytes())
	buffers.Put(buf)
	if err == nil {
		if err = c.setWriteDeadline(c.server.config.WriteTimeout); err != nil {
			return
		}
		_, err = c.conn.Write(m)
	}
	return
}

func (c *tcpConn) Close() error {
	if atomic.CompareAndSwapInt32(&c.closeTag, open, closed) == false {
		return nil
	}
	// 关闭连接
	return c.conn.Close()
}

func (c *tcpConn) setWriteDeadline(timeout time.Duration) error {
	if timeout > 0 {
		return c.conn.SetWriteDeadline(time.Now().Add(timeout))
	} else {
		return c.conn.SetWriteDeadline(time.Time{})
	}
}
