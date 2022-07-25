package clientstream

import (
	"fmt"
	"io"
	"net"

	"github.com/pyihe/go-example/grpc/proto"
	"github.com/pyihe/go-pkg/syncs"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedClientStreamServer
	wg         syncs.WgWrapper
	grpcServer *grpc.Server
}

func Run(addr string, opts ...grpc.ServerOption) (io.Closer, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	srv := &server{
		UnimplementedClientStreamServer: proto.UnimplementedClientStreamServer{},
		grpcServer:                      grpc.NewServer(opts...),
	}
	srv.wg.Wrap(func() {
		proto.RegisterClientStreamServer(srv.grpcServer, srv)
		srv.grpcServer.Serve(ln)
	})
	return srv, nil
}

func (s *server) Close() error {
	s.grpcServer.GracefulStop()
	s.wg.Wait()
	return nil
}

func (s *server) Echo(stream proto.ClientStream_EchoServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return err
		}
		fmt.Printf("ClientStream服务器收到消息: %v\n", in.Message)
		if err = stream.SendMsg(&proto.EchoResponse{
			Message: in.Message,
		}); err != nil {
			fmt.Printf("ClientStream SendMsg err: %v\n", err)
		}
	}
}
