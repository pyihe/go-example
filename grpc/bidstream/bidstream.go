package bidstream

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

type session struct {
	request *proto.EchoRequest
	stream  proto.BidStream_EchoServer
}

type server struct {
	proto.UnimplementedBidStreamServer
	wg          syncs.WgWrapper
	ln          net.Listener
	requestChan chan *session
	errChan     chan error
	s           *grpc.Server
}

func Run(addr string, opts ...grpc.ServerOption) (io.Closer, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	s := &server{
		UnimplementedBidStreamServer: proto.UnimplementedBidStreamServer{},
		ln:                           ln,
		requestChan:                  make(chan *session, 64),
		errChan:                      make(chan error, 1),
		s:                            grpc.NewServer(opts...),
	}

	s.wg.Wrap(func() {
		s.readLoop()
	})
	s.wg.Wrap(func() {
		s.exitLoop()
	})

	s.wg.Wrap(func() {
		proto.RegisterBidStreamServer(s.s, s)
		s.s.Serve(ln)
	})
	return s, nil
}

func (s *server) Close() error {
	for len(s.requestChan) > 0 {
	}
	s.s.GracefulStop()
	close(s.requestChan)
	close(s.errChan)
	s.wg.Wait()
	return nil
}

func (s *server) exitLoop() {
	for {
		select {
		case err, ok := <-s.errChan:
			if ok {
				fmt.Printf("exit with err: %v\n", err)
			}
			return
		}
	}
}

func (s *server) readLoop() {
	for {
		select {
		case sess, ok := <-s.requestChan:
			if !ok {
				return
			}
			if err := sess.stream.SendMsg(&proto.EchoResponse{Message: sess.request.Message}); err != nil {
				s.errChan <- err
			}
		}
	}
}

func (s *server) Echo(stream proto.BidStream_EchoServer) (err error) {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		err = errors.New("unary: want metadata but not exist")
		return
	}
	fmt.Printf("bidstream server: recv md(%v)\n", md)
	header := metadata.Pairs("token", md.Get("token")[0])
	stream.SendHeader(header)

	defer func() {
		trailer := metadata.Pairs("token", md.Get("token")[0])
		stream.SetTrailer(trailer)
	}()

	for {
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return err
		}
		fmt.Printf("bidstream server: recv(%v)\n", in.String())
		sess := &session{
			request: in,
			stream:  stream,
		}
		s.requestChan <- sess
	}
}
