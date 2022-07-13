package tcp

import (
	"bufio"
	"net"
	"time"
)

// Conn 底层连接
type Conn interface {
	GetID() int64                                   // 获取连接唯一标识
	RemoteAddr() net.Addr                           // 获取客户端地址
	SendMsg([]byte) error                           // 发送消息
	SendMsgWithTimeout([]byte, time.Duration) error // 带超时机制发送消息
}

type tcpConn struct {
	id     int64         // 唯一ID
	conn   net.Conn      // 底层连接
	w      *bufio.Writer // writer
	server *tcpServer    // 属于哪个服务器
}

func newTCPConn(conn net.Conn, s *tcpServer) *tcpConn {
	return &tcpConn{
		id:     s.handler.NewUniqueID(),
		conn:   conn,
		w:      bufio.NewWriter(conn),
		server: s,
	}
}

func (c *tcpConn) GetID() int64 {
	return c.id
}

func (c *tcpConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *tcpConn) SendMsg(message []byte) (err error) {
	data, err := c.server.pkt.Packet(message)
	if err != nil {
		return err
	}
	if _, err = c.w.Write(data); err != nil {
		return
	}
	return c.w.Flush()
}

func (c *tcpConn) SendMsgWithTimeout(message []byte, timeout time.Duration) (err error) {
	if err = c.setWriteDeadline(timeout); err != nil {
		return
	}
	data, err := c.server.pkt.Packet(message)
	if err != nil {
		return err
	}
	if _, err = c.w.Write(data); err != nil {
		return
	}
	return c.w.Flush()
}

func (c *tcpConn) setWriteDeadline(timeout time.Duration) error {
	if timeout > 0 {
		return c.conn.SetWriteDeadline(time.Now().Add(timeout))
	} else {
		return c.conn.SetWriteDeadline(time.Time{})
	}
}
