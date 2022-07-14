package websocket

import "time"

// Config websocket服务配置
type Config struct {
	Tick            bool          // 是否开启定时器
	Addr            string        // ws服务地址
	ReadTimeout     time.Duration // 读超时
	WriteTimeout    time.Duration // 写超时
	MaxMsgSize      int           // 消息体最大长度
	MinMsgSize      int           // 消息体最小长度
	MaxConns        int           // 最大连接数
	ReadBufferSize  int           // 读缓冲区大小
	WriteBufferSize int           // 写缓冲区大小
	CertFile        string        // 证书路径
	KeyFile         string        // 证书密钥路径
}
