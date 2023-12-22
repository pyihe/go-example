package main

import (
	"context"
	"fmt"

	httpnet2 "github.com/pyihe/go-example/http"

	"github.com/gin-gonic/gin"
	"github.com/pyihe/go-example/http/middleware/jeager"
	"github.com/pyihe/go-example/jeager/proto"
	"github.com/pyihe/go-pkg/tools"
	"google.golang.org/grpc"
)

func main() {
	httpServer := httpnet2.NewServer(":8080")
	defer httpServer.Close()

	staticConn, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer staticConn.Close()

	simpleConn, err := grpc.Dial(":8082", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer simpleConn.Close()

	s := newService(proto.NewStaticServiceClient(staticConn), proto.NewSimpleServiceClient(simpleConn))
	httpServer.AddHandler(s)

	httpServer.Run()

	tools.Wait()
}

type service struct {
	staticClient proto.StaticService_StreamClient
	sampleClient proto.SimpleServiceClient
}

func newService(staticClient proto.StaticServiceClient, simpleClient proto.SimpleServiceClient) *service {
	static, err := staticClient.Stream(context.Background())
	if err != nil {
		panic(err)
	}
	s := &service{
		staticClient: static,
		sampleClient: simpleClient,
	}

	go s.streamLoop()

	return s
}

func (s *service) Handle(router httpnet2.IRouter) {
	router.GET("/gate", httpnet2.WrapHandler(s.gate), jeager.WithTracing())
}

func (s *service) gate(c *gin.Context) (result interface{}, err error) {
	// static
	s.staticClient.Send(&proto.StaticRequest{Body: "xxxxxxx"})

	// simple
	s.sampleClient.Simple(context.Background(), &proto.SimpleRequest{Body: "zzzzzzz"})

	return "这里是Gate", nil
}

func (s *service) streamLoop() {
	for {
		rsp, err := s.staticClient.Recv()
		if err != nil {
			return
		}
		fmt.Println(rsp)
	}
}
