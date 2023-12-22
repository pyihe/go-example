package httpnet

import "github.com/gin-gonic/gin"

type IRouter = gin.IRouter

type APIHandler interface {
	Handle(IRouter)
}

type Option func(*config)

type config struct {
	name        string // 服务名称
	addr        string // 服务地址
	keyFile     string // 密钥文件
	certFile    string // 证书文件
	routePrefix string // 路由前缀
	swaggerURL  string // swagger文档地址，如果填写则开启swagger
}

func WithName(name string) Option {
	return func(c *config) {
		c.name = name
	}
}

func WithTLS(key, cert string) Option {
	return func(c *config) {
		c.keyFile = key
		c.certFile = cert
	}
}

func WithPrefix(prefix string) Option {
	return func(c *config) {
		c.routePrefix = prefix
	}
}

func WithSwagger(swaggerUrl string) Option {
	return func(c *config) {
		c.swaggerURL = swaggerUrl
	}
}
