package main

import (
	"github.com/pyihe/go-example/fsnotify/route"
	"github.com/pyihe/go-example/fsnotify/service"
	"github.com/pyihe/go-example/pkg"
	"github.com/pyihe/go-pkg/https/http_api"
	"github.com/pyihe/go-pkg/tools"
)

func main() {
	logger := pkg.InitLogger("file")
	defer logger.Close()

	config := http_api.Config{
		Name: "file_server",
		Addr: ":8080",
	}
	server := http_api.NewHTTPServer(config)
	defer server.Close()

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
