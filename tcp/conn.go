package tcp

import (
	"net"
	"time"

	"github.com/pyihe/go-pkg/bytes"
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
	id          int64             // 唯一ID
	conn        net.Conn          // 底层连接
	writeBuffer *bytes.ByteBuffer // 写缓冲区
	server      *tcpServer        // 属于哪个服务器
}

func newTCPConn(conn net.Conn, s *tcpServer) *tcpConn {
	return &tcpConn{
		id:     s.handler.NewUniqueID(),
		conn:   conn,
		server: s,
	}
}

func (c *tcpConn) getBuffer() *bytes.ByteBuffer {
	if c.writeBuffer == nil {
		c.writeBuffer = bytes.Get()
	}
	return c.writeBuffer
}

func (c *tcpConn) releaseBuffer() {
	if c.writeBuffer == nil {
		return
	}
	c.writeBuffer.Reset()
}

func (c *tcpConn) GetID() int64 {
	return c.id
}

func (c *tcpConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *tcpConn) SendMsg(message ...[]byte) (err error) {
	connBuffer := c.getBuffer()
	tempBuffer := bytes.Get()
	for _, m := range message {
		tempBuffer.B, err = c.server.pkt.Packet(m)
		if err != nil {
			break
		}
		if _, err = connBuffer.Write(tempBuffer.B); err != nil {
			break
		}
		tempBuffer.Reset()
	}
	if err == nil {
		_, err = c.conn.Write(connBuffer.B)
	}
	bytes.Put(tempBuffer)
	c.releaseBuffer()
	return
}

func (c *tcpConn) SendMsgWithTimeout(message ...[]byte) (err error) {
	connBuffer := c.getBuffer()
	tempBuffer := bytes.Get()
	for _, m := range message {
		tempBuffer.B, err = c.server.pkt.Packet(m)
		if err != nil {
			break
		}
		if _, err = connBuffer.Write(tempBuffer.B); err != nil {
			break
		}
		tempBuffer.Reset()
	}
	if err == nil {
		if err = c.setWriteDeadline(c.server.config.WriteTimeout); err != nil {
			goto end
		}
		_, err = c.conn.Write(connBuffer.B)
	}
end:
	bytes.Put(tempBuffer)
	c.releaseBuffer()
	return
}

func (c *tcpConn) Close() error {
	// 归还缓存
	bytes.Put(c.writeBuffer)
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
