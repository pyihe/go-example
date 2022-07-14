package websocket

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
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
	OnTick() (time.Duration, error)

	// OnPing 收到心跳时调用
	OnPing(conn Conn)

	// GenerateID 生成全局唯一ID
	GenerateID() int64
}

type wsServer struct {
	ctx      context.Context
	cancel   context.CancelFunc
	wg       syncs.WgWrapper
	ln       net.Listener       // 连接监听器
	handler  Handler            // 服务器handler
	upgrader websocket.Upgrader // 将请求升级为websocket
	conns    *maps.Map          // 所有连接的维护信息
	config   *Config            // 配置
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

	s.wg.Add(1)
	defer s.wg.Done()

	// 判断是否超过最大连接数限制
	if maxConns := s.config.MaxConns; maxConns > 0 && s.countConn() >= maxConns {
		conn.Close()
		return
	}

	// 新建连接
	wsConn := newWsConn(conn)
	s.handler.OnConnect(wsConn)
	// 添加连接
	s.addConn(wsConn)
	// 监听连接上的数据
	s.ioLoop(wsConn)
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

func (s *wsServer) ioLoop(conn *wsConn) {

}
