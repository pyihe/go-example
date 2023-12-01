package http

import "github.com/gin-gonic/gin"

type Handler interface {
	Handle(ir gin.IRouter)
}

type Server struct {
	address string
	engin   *gin.Engine
}

func NewHTTPServer(addr string) *Server {
	s := &Server{}
	s.address = addr
	s.engin = gin.Default()
	s.engin.Use(gin.Recovery())

	return s
}

func (s *Server) Use(middleware ...gin.HandlerFunc) {
	s.engin.Use(middleware...)
}

func (s *Server) Run() {
	go s.engin.Run(s.address)
}

func (s *Server) AddHandler(h Handler) {
	if h != nil {
		h.Handle(s.engin.Group(""))
	}
}
