package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pyihe/go-example/fsnotify/service"
	"github.com/pyihe/go-pkg/https/http_api"
)

type FileRouter struct {
	service *service.FileService
}

func NewFileRouter(s *service.FileService) *FileRouter {
	return &FileRouter{service: s}
}

func (f *FileRouter) Handle(r http_api.IRouter) {
	// 静态文件服务器
	r.StaticFS("/files", http.Dir(f.service.GetConfigFilePath()))

	//
	r.POST("/list", http_api.WrapHandler(f.list))
}

func (f *FileRouter) list(c *gin.Context) (result interface{}, err error) {
	return f.service.List()
}
