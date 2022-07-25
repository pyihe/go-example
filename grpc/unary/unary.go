package unary

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/pyihe/go-example/grpc/proto"
	"github.com/pyihe/go-pkg/syncs"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedUnaryServer
	wg         syncs.WgWrapper
	grpcServer *grpc.Server
}

func Run(addr string, opts ...grpc.ServerOption) (io.Closer, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	srv := &server{
		UnimplementedUnaryServer: proto.UnimplementedUnaryServer{},
		grpcServer:               grpc.NewServer(opts...),
	}
	srv.wg.Wrap(func() {
		proto.RegisterUnaryServer(srv.grpcServer, srv)
		srv.grpcServer.Serve(ln)
	})
	return srv, nil
}

func (s *server) Close() error {
	s.grpcServer.GracefulStop()
	s.wg.Wait()
	return nil
}

func (s *server) Echo(ctx context.Context, in *proto.EchoRequest) (resp *proto.EchoResponse, err error) {
	fmt.Printf("unary服务器收到消息: %v\n", in.Message)
	resp = &proto.EchoResponse{Message: in.Message}
	return
}
