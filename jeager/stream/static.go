package main

import (
	"net"

	"github.com/pyihe/go-example/jeager/proto"
	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()

	proto.RegisterStaticServiceServer(grpcServer, &service{})

	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		return
	}
	grpcServer.Serve(l)
}

type service struct {
	proto.UnimplementedStaticServiceServer
}

func (s *service) Stream(server proto.StaticService_StreamServer) error {
	for {
		request, err := server.Recv()
		if err != nil {
			return err
		}
		err = server.Send(&proto.StaticResponse{Body: request.Body})
		if err != nil {
			return err
		}
	}
	return nil
}
