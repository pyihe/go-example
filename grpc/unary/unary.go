package unary

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/pyihe/go-example/grpc/proto"
	"github.com/pyihe/go-pkg/errors"
	"github.com/pyihe/go-pkg/syncs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// 简单模式，即常规的rpc request-response请求响应模式，每个客户端的请求对应服务器的一个响应
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
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = errors.New("unary: want metadata but not exist")
		return
	}

	fmt.Printf("unary: receive request, md(%v), request(%v)\n", md, in.String())

	header := metadata.Pairs("token", md.Get("token")[0])
	trailer := metadata.Pairs("token", md.Get("token")[0])

	// 响应请求携带header, trailer
	grpc.SetHeader(ctx, header)
	grpc.SetTrailer(ctx, trailer)
	resp = &proto.EchoResponse{Message: in.Message}
	return
}
