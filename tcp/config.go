package tcp

import "time"

// Config TCP服务器参数配置
type Config struct {
	Ticker       bool          // 是否开启tick
	HeaderSize   int           // 消息头部长度
	MaxMsgSize   int           // 消息体最大长度
	MinMsgSize   int           // 消息体最小长度
	Port         int           // 服务器端口号
	ReadBuffer   int           // 读取消息的通道缓冲区大小, 默认64
	IP           string        // 服务器IP
	ReadTimeout  time.Duration // 读超时
	WriteTimeout time.Duration // 写超时
	TLSConfig    *TLSConfig    // TLS配置
}

// TLSConfig TLS证书配置
type TLSConfig struct {
	ServerCert string // 服务器证书
	ServerKey  string // 服务器密钥
	RootCa     string // 根证书
}
