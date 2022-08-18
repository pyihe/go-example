package main

import (
	"time"

	"github.com/pyihe/go-example/fsnotify/route"
	"github.com/pyihe/go-example/fsnotify/service"
	"github.com/pyihe/go-pkg/https/http_api"
	"github.com/pyihe/go-pkg/tools"
	"github.com/pyihe/plogs"
)

func main() {
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

	config := http_api.Config{
		Name: "file_server",
		Addr: ":8080",
	}
	server := http_api.NewHTTPServer(config)
	defer server.Stop()

	// service
	fService := service.NewFileService("./files")
	defer fService.Close()

	// router
	fRouter := route.NewFileRouter(fService)

	// add handler
	server.AddHandler(fRouter)

	server.Run()

	tools.Wait()
}
