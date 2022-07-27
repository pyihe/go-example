package clientstream

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/pyihe/go-example/grpc/proto"
	"github.com/pyihe/go-pkg/errors"
	"github.com/pyihe/go-pkg/syncs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

func (s *server) Echo(stream proto.ClientStream_EchoServer) (err error) {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		err = errors.New("unary: want metadata but not exist")
		return
	}
	fmt.Printf("clientstream: md(%v)\n", md)

	//header := metadata.Pairs("token", md.Get("token")[0])
	//stream.SendHeader(header)

	var builder strings.Builder
loop:
	for {
		in, err := stream.Recv()
		switch err {
		case io.EOF:
			break loop
		case nil:
			builder.WriteString(in.Message)
			fmt.Printf("clientstream: receive request(%v)\n", in.String())
		default:
			fmt.Printf("clientstream: receive fail(%v)\n", err)
			return err
		}
	}
	header := metadata.Pairs("token", md.Get("token")[0])
	stream.SetHeader(header)
	trailer := metadata.Pairs("token", md.Get("token")[0])
	stream.SetTrailer(trailer)
	return stream.SendMsg(&proto.EchoResponse{Message: builder.String()})
}
