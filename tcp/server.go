package tcp

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"sync"
	"time"

	"github.com/pyihe/go-pkg/bytes"
	"github.com/pyihe/go-pkg/errors"
	"github.com/pyihe/go-pkg/packets"
	"github.com/pyihe/go-pkg/syncs"
)

// Handler 服务Handler, 由上层实现
type Handler interface {
	// OnMessage 收到消息时调用, 需要注意的是如果需要改变[]byte, 请拷贝并另行处理, 切勿操作原始的[]byte
	OnMessage(Conn, []byte)

	// OnConnect 建立连接时调用
	OnConnect(Conn)

	// OnClose 连接断开时调用, 这里同时包含客户端主动关闭以及服务器关闭连接
	OnClose(Conn)

	// OnTick 定时任务
	OnTick() (time.Duration, bool)

	// NewUniqueID 获取全局唯一ID
	NewUniqueID() int64
}

type message struct {
	conn *tcpConn
	data *bytes.ByteBuffer
}

type tcpServer struct {
	closed     bool               // 服务器是否主动关闭
	ctx        context.Context    // 上下文
	cancel     context.CancelFunc // 取消函数
	wg         syncs.WgWrapper    // waitgroup
	conns      sync.Map           // 保存所有的连接
	listener   net.Listener       // net listener
	pkt        packets.IPacket    // 封包、拆包
	msgPool    sync.Pool          // message pool
	readBuffer chan *message      // 处理收到消息的缓冲区
	handler    Handler            // 服务器的handler
	config     *Config            // 服务器配置
}

func Run(config *Config, handler Handler) (io.Closer, error) {
	var bufferSize = 64
	var err error
	var address = fmt.Sprintf("%s:%d", config.IP, config.Port)
	var packetOpts = []packets.Option{
		packets.WithHeaderSize(config.HeaderSize),
		packets.WithMaxMsgSize(config.MaxMsgSize),
		packets.WithMinMsgSize(config.MinMsgSize),
	}
	var s = &tcpServer{
		wg:      syncs.WgWrapper{},
		conns:   sync.Map{},
		pkt:     packets.NewPacket(packetOpts...),
		msgPool: sync.Pool{},
		handler: handler,
		config:  config,
	}

	if config.ReadBuffer > 0 {
		bufferSize = config.ReadBuffer
	}

	s.readBuffer = make(chan *message, bufferSize)
	s.ctx, s.cancel = context.WithCancel(context.Background())

	if tlsConfig := config.TLSConfig; tlsConfig == nil {
		s.listener, err = net.Listen("tcp", address)
		if err != nil {
			return nil, err
		}
	} else {
		conf, err := s.loadTLSConfig()
		if err != nil {
			return nil, err
		}
		s.listener, err = tls.Listen("tcp", address, conf)
		if err != nil {
			return nil, err
		}
	}

	s.wg.Wrap(func() {
		s.start()
	})
	s.wg.Wrap(func() {
		s.process()
	})
	return s, nil
}

func (s *tcpServer) Close() error {
	if s.closed {
		return nil
	}

	s.closed = true

	// 关闭缓冲区
	close(s.readBuffer)

	// 关闭Listener
	s.listener.Close()

	// 执行取消函数
	s.cancel()

	// 关闭所有连接
	s.conns.Range(func(k, v any) bool {
		c, ok := v.(*tcpConn)
		if ok {
			c.conn.Close()
		}
		return true
	})

	s.wg.Wait()

	return nil
}

func (s *tcpServer) process() {
	for {
		select {
		case msg, ok := <-s.readBuffer:
			if !ok {
				return
			}
			if msg != nil {
				s.handler.OnMessage(msg.conn, msg.data.Bytes())
				s.putMessage(msg)
			}
		}
	}
}

func (s *tcpServer) loadTLSConfig() (*tls.Config, error) {
	tlsConfig := s.config.TLSConfig
	serverCert, err := tls.LoadX509KeyPair(tlsConfig.ServerCert, tlsConfig.ServerKey)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(tlsConfig.RootCa)
	if err != nil {
		return nil, err
	}
	clientCertPool := x509.NewCertPool()
	if ok := clientCertPool.AppendCertsFromPEM(b); !ok {
		return nil, errors.New("failed to parse root certificate")
	}
	return &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool, // 用于验证客户端证书是否合法
	}, nil
}

func (s *tcpServer) newUniqueID() int64 {
	return s.handler.NewUniqueID()
}

func (s *tcpServer) getMessage() *message {
	data := s.msgPool.Get()
	if data == nil {
		return &message{
			data: bytes.Get(),
		}
	}
	return data.(*message)
}

func (s *tcpServer) putMessage(m *message) {
	if m != nil {
		m.data.Reset()
		bytes.Put(m.data)
		m.data = nil
		m.conn = nil
	}
}

func (s *tcpServer) start() {
	if s.config.Ticker {
		go s.tick()
	}

	var wg syncs.WgWrapper
	for {
		clientConn, err := s.listener.Accept()
		if err != nil {
			if !isServerClose(err) {
				fmt.Printf("TCP Accept err(%v)\n", err)
			}
			break
		}
		wg.Wrap(func() {
			s.handleConn(clientConn)
		})
	}
	wg.Wait()
}

func (s *tcpServer) handleConn(conn net.Conn) {
	client := newTCPConn(conn, s)
	s.handler.OnConnect(client)
	s.conns.Store(conn.RemoteAddr(), client)

	s.ioLoop(client)

	s.handler.OnClose(client)
	client.Close()
	s.conns.Delete(conn.RemoteAddr())
}

func (s *tcpServer) ioLoop(tc *tcpConn) {
	var pkt = s.pkt
	var reader = bufio.NewReader(tc.conn)

	for {
		err := s.setReadTimeout(tc)
		if err != nil {
			break
		}
		data, err := pkt.UnPacket(reader)
		if err != nil {
			if !isClientClose(err) && !isServerClose(err) {
				fmt.Printf("read tcp fail: %v\n", err)
			}
			break
		}
		if !s.closed {
			msg := s.getMessage()
			msg.conn = tc
			msg.data.Write(data)
			s.readBuffer <- msg
		}
	}
}

func (s *tcpServer) setReadTimeout(c *tcpConn) error {
	timeout := s.config.ReadTimeout
	if timeout > 0 {
		return c.conn.SetReadDeadline(time.Now().Add(timeout))
	} else {
		return c.conn.SetReadDeadline(time.Time{})
	}
}

func (s *tcpServer) tick() {
	var (
		stop  bool
		delay time.Duration
		timer *time.Timer
	)
	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()

	for {
		delay, stop = s.handler.OnTick()
		if stop {
			return
		}
		if timer == nil {
			timer = time.NewTimer(delay)
		} else {
			timer.Reset(delay)
		}
		select {
		case <-s.ctx.Done():
			return
		case <-timer.C:
		}
	}
}
