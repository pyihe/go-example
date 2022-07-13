package tcp

import "time"

type Config struct {
	Ticker      bool          // 是否开启tick
	MsgHeader   int           // 消息头部长度
	MsgSize     int           // 消息体最大长度
	Port        int           // 服务器端口号
	IP          string        // 服务器IP
	ReadTimeout time.Duration // 读超时
	TLSConfig   *TLSConfig    // TLS配置
}

type TLSConfig struct {
	ServerCert string // 服务器证书
	ServerKey  string // 服务器密钥
	ClientCert string // 客户端证书
}
