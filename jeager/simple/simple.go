package main

import (
	"context"
	"net"

	"github.com/pyihe/go-example/jeager/proto"
	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()
	proto.RegisterSimpleServiceServer(grpcServer, &service{})

	l, err := net.Listen("tcp", ":8082")
	if err != nil {
		return
	}

	grpcServer.Serve(l)
}

type service struct {
	proto.UnimplementedSimpleServiceServer
}

func (s *service) Simple(ctx context.Context, in *proto.SimpleRequest) (rsp *proto.SimpleResponse, err error) {
	rsp.Body = in.Body
	return
}
