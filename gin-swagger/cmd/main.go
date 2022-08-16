package main

import (
	"time"

	"github.com/pyihe/go-example/gin-swagger/cmd/docs"
	"github.com/pyihe/go-example/gin-swagger/route"
	"github.com/pyihe/go-example/gin-swagger/service"
	"github.com/pyihe/go-pkg/https/http_api"
	"github.com/pyihe/go-pkg/tools"
	"github.com/pyihe/plogs"
)

func main() {
	docs.SwaggerInfo.BasePath = ""

	opts := []plogs.Option{
		plogs.WithName("swagger"),
		plogs.WithFileOption(plogs.WriteByLevelMerged),
		plogs.WithLogPath("./logs"),
		plogs.WithStdout(true),
		plogs.WithLogLevel(plogs.LevelInfo | plogs.LevelDebug | plogs.LevelWarn | plogs.LevelError | plogs.LevelFatal | plogs.LevelPanic),
		plogs.WithMaxAge(24 * time.Hour),
		plogs.WithMaxSize(60 * 1024 * 1024),
	}

	logger := plogs.NewLogger(opts...)
	defer logger.Close()

	httpConfig := http_api.Config{
		Swagger:     true,
		Name:        "gin-swagger",
		Addr:        ":8080",
		RoutePrefix: "",
		CertFile:    "",
		KeyFile:     "",
	}
	httpServer := http_api.NewHTTPServer(httpConfig)
	defer httpServer.Stop()

	// 初始化service
	userService := service.NewUserService()

	// 初始化route
	userRoute := route.NewUserRoute(userService)

	httpServer.AddHandler(userRoute)

	httpServer.Run()

	tools.Wait()
}
