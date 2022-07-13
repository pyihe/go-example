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

	"github.com/pyihe/go-pkg/errors"
	"github.com/pyihe/go-pkg/packets"
	"github.com/pyihe/go-pkg/syncs"
)

// Handler 服务Handler, 由上层实现
type Handler interface {
	OnMessage(Conn, []byte)        // 收到消息时调用
	OnConnect(Conn)                // 建立连接时调用
	OnClose(Conn)                  // 连接断开时调用
	OnTick() (time.Duration, bool) // 定时任务
	NewUniqueID() int64            // 获取全局唯一ID
}

type tcpServer struct {
	ctx      context.Context    // 上下文
	cancel   context.CancelFunc // 取消函数
	wg       syncs.WgWrapper    // waitgroup
	conns    sync.Map           // 保存所有的连接
	listener net.Listener       // net listener
	pkt      packets.IPacket    // 封包、拆包
	handler  Handler            // 服务器的handler
	config   *Config            // 服务器配置
}

func Run(config *Config, handler Handler) (io.Closer, error) {
	var err error
	var address = fmt.Sprintf("%s:%d", config.IP, config.Port)
	var s = &tcpServer{
		wg:      syncs.WgWrapper{},
		conns:   sync.Map{},
		pkt:     packets.NewPacket(config.MsgHeader, config.MsgSize),
		handler: handler,
		config:  config,
	}

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
	return s, nil
}

func (s *tcpServer) Close() error {
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

func (s *tcpServer) start() error {
	if s.config.Ticker {
		go s.tick()
	}

	var wg syncs.WgWrapper
	for {
		clientConn, err := s.listener.Accept()
		if err != nil {
			if !isServerClose(err) {
				return fmt.Errorf("TCP Accept err(%v)", err)
			}
			break
		}
		wg.Wrap(func() {
			s.handleConn(clientConn)
		})
	}
	wg.Wait()
	return nil
}

func (s *tcpServer) handleConn(conn net.Conn) {
	//s.logger.Infof("Accept New Connection: %s\n", conn.RemoteAddr())

	client := newTCPConn(conn, s)
	s.handler.OnConnect(client)
	s.conns.Store(conn.RemoteAddr(), client)

	s.ioLoop(client)

	s.handler.OnClose(client)
	s.conns.Delete(conn.RemoteAddr())
	conn.Close()
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
			}
			break
		}
		s.handler.OnMessage(tc, data)
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
