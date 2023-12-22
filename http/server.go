package httpnet

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pyihe/go-pkg/syncs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config *config         // 服务配置
	engine *gin.Engine     // 路由
	server *http.Server    // HTTP服务
	wg     syncs.WgWrapper // waiter
}

func NewServer(addr string, opts ...Option) *Server {
	c := &config{
		addr: addr,
	}
	for _, op := range opts {
		op(c)
	}

	engine := gin.Default()
	s := &Server{}
	s.config = c
	s.engine = engine
	s.server = &http.Server{
		Addr:    c.addr,
		Handler: engine,
	}

	if c.swaggerURL != "" {
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL(c.swaggerURL)))
	}

	return s
}

func (s *Server) Use(middleware ...gin.HandlerFunc) {
	s.engine.Use(middleware...)
}

func (s *Server) AddHandler(handler APIHandler) {
	engin := s.engine
	prefix := s.config.routePrefix
	handler.Handle(engin.Group(prefix))
}

func (s *Server) Run() {
	s.wg.Wrap(func() {
		if s.config.certFile != "" && s.config.keyFile != "" {
			if err := s.server.ListenAndServeTLS(s.config.certFile, s.config.keyFile); err != nil {
				log.Printf("listen fail: %v\n", err)
			}
		} else {
			if err := s.server.ListenAndServe(); err != nil {
				log.Printf("listen fail: %v\n", err)
			}
		}
	})
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	s.wg.Wait()
	return nil
}
