package serverstream

import (
	"fmt"
	"io"
	"net"

	"github.com/pyihe/go-example/grpc/proto"
	"github.com/pyihe/go-pkg/syncs"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedServerStreamServer
	wg         syncs.WgWrapper
	grpcServer *grpc.Server
}

func Run(addr string, opts ...grpc.ServerOption) (io.Closer, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	srv := &server{
		UnimplementedServerStreamServer: proto.UnimplementedServerStreamServer{},
		grpcServer:                      grpc.NewServer(opts...),
	}

	srv.wg.Wrap(func() {
		proto.RegisterServerStreamServer(srv.grpcServer, srv)
		srv.grpcServer.Serve(ln)
	})
	return srv, nil
}

func (s *server) Close() error {
	s.grpcServer.GracefulStop()
	s.wg.Wait()
	return nil
}

func (s *server) Echo(in *proto.EchoRequest, stream proto.ServerStream_EchoServer) error {
	fmt.Printf("ServerStream服务器收到消息: %v\n", in.Message)
	return stream.Send(&proto.EchoResponse{Message: in.Message})
}
