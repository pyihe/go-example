package main

import (
	_ "github.com/pyihe/go-example/gin-swagger/docs"
	"github.com/pyihe/go-example/gin-swagger/route"
	"github.com/pyihe/go-example/gin-swagger/service"
	"github.com/pyihe/go-example/pkg"
	"github.com/pyihe/go-pkg/https/http_api"
	"github.com/pyihe/go-pkg/tools"
)

// @title          这里填写文档名称(必填项): Gin-Swagger在线文档测试项目
// @version        0.0.1(必填项)
// description 	   这里填写项目描述: 本项目用于演示swagger在gin框架中的应用, 项目完全虚构，用于演示
// @termsOfService 这里填写服务条款: http://swagger.io/terms/

// @host     localhost:8080(服务器运行地址)
// @BasePath /api(API基本路径)

// contact.name API联系人
// contact.url 联系人的URL信息
// contact.email 联系人Email

// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization
// @description                API安全验证
func main() {
	logger := pkg.InitLogger("swagger")
	defer logger.Close()

	httpConfig := http_api.Config{
		SwaggerURL:  "http://localhost:8080/swagger/doc.json",
		Name:        "gin-swagger",
		Addr:        ":8080",
		RoutePrefix: "/api",
	}
	httpServer := http_api.NewHTTPServer(httpConfig)
	defer httpServer.Stop()

	// 初始化service
	userService := service.NewUserService()
	departService := service.NewDepartmentService()

	// 初始化route
	userRouter := route.NewUserRoute(userService)
	departRouter := route.NewDepartmentRouter(departService)

	httpServer.AddHandler(userRouter)
	httpServer.AddHandler(departRouter)

	httpServer.Run()

	tools.Wait()
}
