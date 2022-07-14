package websocket

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pyihe/go-pkg/bytes"
	"github.com/pyihe/go-pkg/errors"
	"github.com/pyihe/go-pkg/maps"
	"github.com/pyihe/go-pkg/syncs"
)

type Handler interface {
	// OnMessage 消息处理
	OnMessage(conn Conn, message []byte)

	// OnConnect 连接建立时调用
	OnConnect(conn Conn)

	// OnClose 连接关闭时调用
	OnClose(conn Conn)

	// OnTick 定时任务, bool值返回true的话定时器将被终止, 并不可恢复
	OnTick() (time.Duration, bool)

	// OnPing 收到心跳时调用
	OnPing(conn Conn)

	// GenerateID 生成全局唯一ID
	GenerateID() int64
}

type message struct {
	messageType int
	conn        *wsConn
	data        *bytes.ByteBuffer
}

type wsServer struct {
	closeTag  int32
	ctx       context.Context
	cancel    context.CancelFunc
	wg        syncs.WgWrapper
	msgPool   sync.Pool
	msgBuffer chan *message      // 消息缓冲区
	ln        net.Listener       // 连接监听器
	handler   Handler            // 服务器handler
	upgrader  websocket.Upgrader // 将请求升级为websocket
	conns     *maps.Map          // 所有连接的维护信息
	config    *Config            // 配置
}

func Run(config *Config, handler Handler) (io.Closer, error) {
	if config == nil {
		return nil, errors.New("config cannot be empty")
	}
	if handler == nil {
		return nil, errors.New("handler cannot be nil")
	}
	// 初始化消息缓冲区大小，默认64
	bufferSize := 64
	if config.ReadBufferSize > 0 {
		bufferSize = config.ReadBufferSize
	}
	s := &wsServer{
		closeTag:  open,
		wg:        syncs.WgWrapper{},
		msgPool:   sync.Pool{},
		msgBuffer: make(chan *message, bufferSize),
		handler:   handler,
		upgrader: websocket.Upgrader{
			HandshakeTimeout: config.ReadTimeout,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		conns:  maps.NewMap(),
		config: config,
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.upgrader = websocket.Upgrader{}

	// 初始化Listener
	if config.CertFile != "" && config.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
		if err != nil {
			return nil, err
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			NextProtos:   []string{"http/1.1"},
		}
		s.ln, err = tls.Listen("tcp", config.Addr, tlsConfig)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		s.ln, err = net.Listen("tcp", config.Addr)
		if err != nil {
			return nil, err
		}
	}
	httpServer := &http.Server{
		Addr:           config.Addr,
		Handler:        s,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: 1024,
	}

	s.wg.Wrap(func() {
		httpServer.Serve(s.ln)
	})
	s.wg.Wrap(func() {
		s.processMessage()
	})
	s.wg.Wrap(func() {
		s.tick()
	})
	return s, nil
}

// ServeHTTP 将收到的GET请求升级为websocket请求
func (s *wsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 升级协议必须是GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// 升级请求为websocket, 并返回一个websocket连接
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	// 设置允许读取的最大消息长度
	conn.SetReadLimit(int64(s.config.MaxMsgSize))

	// 判断是否超过最大连接数限制
	if maxConns := s.config.MaxConns; maxConns > 0 && s.countConn() >= maxConns {
		conn.Close()
		return
	}

	// 新建连接
	newConn := newWsConn(conn, s)
	s.handler.OnConnect(newConn)
	// 添加连接
	s.addConn(newConn)

	// 这里为每个新建立的连接开启一个写消息的协程
	s.wg.Wrap(func() {
		newConn.writeLoop()
	})

	s.wg.Add(1)
	// 监听连接上的数据(阻塞)
	s.ioLoop(newConn)

	// 连接断开后需要关闭并删除连接信息
	fmt.Printf("[%v]断开连接...\n", newConn.RemoteAddr())
	s.handler.OnClose(newConn)
	newConn.Close()
	s.removeConn(newConn)
	s.wg.Done()
}

func (s *wsServer) Close() error {
	// 如果已经关闭
	if atomic.CompareAndSwapInt32(&s.closeTag, open, closed) == false {
		return nil
	}

	s.cancel()
	close(s.msgBuffer)
	s.conns.LockRange(func(k interface{}, v interface{}) bool {
		c, ok := v.(*wsConn)
		if ok {
			c.Close()
		}
		return false
	})
	s.wg.Wait()
	return nil
}

func (s *wsServer) getMessage() *message {
	m, ok := s.msgPool.Get().(*message)
	if !ok {
		return &message{
			data: bytes.Get(),
		}
	}
	return m
}

func (s *wsServer) putMessage(m *message) {
	if m != nil {
		m.messageType = 0
		m.conn = nil
		m.data.Reset()
		bytes.Put(m.data)
		s.msgPool.Put(m)
	}
}

func (s *wsServer) addConn(wsConn *wsConn) {
	s.conns.Set(wsConn.conn.RemoteAddr(), wsConn)
}

func (s *wsServer) removeConn(wsConn *wsConn) {
	s.conns.Del(wsConn.conn.RemoteAddr())
}

func (s *wsServer) countConn() int {
	return s.conns.Len()
}

func (s *wsServer) processMessage() {
	for {
		select {
		case m, ok := <-s.msgBuffer:
			if !ok {
				break
			}
			switch m.messageType {
			case websocket.PongMessage:
			case websocket.PingMessage:
				s.handler.OnPing(m.conn)
			case websocket.TextMessage, websocket.BinaryMessage:
				s.handler.OnMessage(m.conn, m.data.Bytes())
			case websocket.CloseMessage:
				s.handler.OnClose(m.conn)
			}
			s.putMessage(m)
		}
	}
}

func (s *wsServer) ioLoop(wConn *wsConn) {
	for {
		messageType, b, err := wConn.conn.ReadMessage()
		if err != nil {
			if err != io.EOF {
				fmt.Printf("read websocket fail: %v\n", err)
			}
			break
		}
		if !s.isClosed() {
			m := s.getMessage()
			m.messageType = messageType
			m.conn = wConn
			m.data.Write(b)
			s.msgBuffer <- m
		}
	}
}

func (s *wsServer) tick() {
	if !s.config.Tick {
		return
	}
	var (
		delay time.Duration
		timer *time.Timer
		stop  bool
	)
	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()
	for {
		delay, stop = s.handler.OnTick()
		if stop {
			break
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

func (s *wsServer) isClosed() bool {
	return atomic.LoadInt32(&s.closeTag) == closed
}
