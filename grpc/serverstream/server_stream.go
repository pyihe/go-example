package serverstream

import (
	"fmt"
	"io"
	"net"

	"github.com/pyihe/go-example/grpc/proto"
	"github.com/pyihe/go-pkg/errors"
	"github.com/pyihe/go-pkg/syncs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//server stream: 客户端采用简单的rpc模式，服务端采用流式，客户端没发送一个请求，服务器可以发送一系列回复
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

func (s *server) Echo(in *proto.EchoRequest, stream proto.ServerStream_EchoServer) (err error) {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		err = errors.New("unary: want metadata but not exist")
		return
	}
	fmt.Printf("serverstream: receive request, md(%v), request(%v)\n", md, in.Message)

	header := metadata.Pairs("token", md.Get("token")[0])
	trailer := metadata.Pairs("token", md.Get("token")[0])
	stream.SetHeader(header)
	stream.SetTrailer(trailer)

	for i := 0; i < 10; i++ {
		err = stream.Send(&proto.EchoResponse{Message: fmt.Sprintf("%s——%d", in.Message, i+1)})
		if err != nil {
			break
		}
	}
	return
}
