package websocket

import "time"

// Config websocket服务配置
type Config struct {
	Tick         bool          // 是否开启定时器
	Addr         string        // ws服务地址
	ReadTimeout  time.Duration // 读超时
	WriteTimeout time.Duration // 写超时
	MaxMsgSize   int           // 消息体最大长度
	MaxConns     int           // 最大连接数
}
