package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pyihe/go-example/grpc/bidstream"
	"github.com/pyihe/go-example/grpc/proto"
	"github.com/pyihe/go-example/grpc/unary"
	"github.com/pyihe/go-pkg/rands"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func testWithAuthAndTLS() {
	cert, err := tls.LoadX509KeyPair("server_cert.pem", "server_key.pem")
	if err != nil {
		fmt.Printf("load certifacate fail: %v\n", err)
		return
	}
	opts := []grpc.ServerOption{
		// TLS加密
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),

		// 简单RPC的拦截器，用于验证token，验证不通过返回错误，否则执行handler
		grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			//
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				err = status.Errorf(codes.InvalidArgument, "want metadata but not exist")
				return
			}
			if !validToken(md.Get("authorization")) {
				err = status.Errorf(codes.Unauthenticated, "invalid token")
				return
			}
			return handler(ctx, req)
		}),
		grpc.StreamInterceptor(func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
			md, ok := metadata.FromIncomingContext(ss.Context())
			if !ok {
				err = status.Errorf(codes.InvalidArgument, "want metadata but not exist")
				return
			}
			if !validToken(md.Get("authorization")) {
				err = status.Errorf(codes.Unauthenticated, "invalid token")
				return
			}
			return handler(srv, ss)
		}),
	}

	uCloser, err := unary.Run(":8801", opts...)
	if err != nil {
		fmt.Printf("unary run fail: %v\n", err)
		return
	}
	defer uCloser.Close()

	sCloser, err := bidstream.Run(":8802", opts...)
	if err != nil {
		fmt.Printf("stream run fail: %v\n", err)
		return
	}
	defer sCloser.Close()

	sendRequest()
}

func sendRequest() {
	// 发送请求
	perRPC := oauth.NewOauthAccess(&oauth2.Token{AccessToken: "your-authorization-token"})
	creds, err := credentials.NewClientTLSFromFile("ca_cert.pem", "x.test.example.com")
	if err != nil {
		fmt.Printf("load client TLS fail: %v\n", err)
		return
	}

	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(perRPC),
		grpc.WithTransportCredentials(creds),
	}

	// unary
	conn, err := grpc.Dial(":8801", opts...)
	if err != nil {
		fmt.Printf("dial unary fail: %v\n", err)
		return
	}
	sendUnaryRequest(conn)

	// stream
	conn, err = grpc.Dial(":8802", opts...)
	if err != nil {
		fmt.Printf("dial stream fail: %v\n", err)
		return
	}
	sendStreamRequest(conn)
}

func sendUnaryRequest(conn *grpc.ClientConn) {
	defer conn.Close()

	md := metadata.Pairs("token", rands.String(16))
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	client := proto.NewUnaryClient(conn)
	resp, err := client.Echo(ctx, &proto.EchoRequest{Message: "I'm unary with TLS&Auth!"})
	if err != nil {
		// handle error
		//status.FromError(err)
		fmt.Printf("unary echo fail: %v\n", err)
		return
	}

	fmt.Printf("unary resp: %v\n\n", resp.String())
}

func sendStreamRequest(conn *grpc.ClientConn) {
	defer conn.Close()

	waiter := sync.WaitGroup{}
	md := metadata.Pairs("token", rands.String(16))
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	client := proto.NewBidStreamClient(conn)
	stream, err := client.Echo(ctx)
	if err != nil {
		fmt.Printf("stream echo fail: %v\n", err)
		return
	}
	waiter.Add(2)
	go func(s proto.BidStream_EchoClient, wg *sync.WaitGroup) {
		defer wg.Done()

		//header, err := s.Header()
		//if err != nil {
		//	fmt.Printf("stream client: header fail(%v)\n", err)
		//	return
		//}
		//fmt.Printf("stream client: header(%v)\n", header)

		for {
			m, err := s.Recv()
			if err != nil {
				break
			}
			fmt.Printf("stream client recv: %v\n", m.String())
		}
	}(stream, &waiter)

	go func(s proto.BidStream_EchoClient, wg *sync.WaitGroup) {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			s.Send(&proto.EchoRequest{Message: "I'm stream with TLS&Auth!"})
		}

		time.Sleep(1 * time.Second)
		s.CloseSend()
	}(stream, &waiter)
	waiter.Wait()
}

func validToken(infos []string) bool {
	if len(infos) == 0 {
		return false
	}

	return strings.TrimPrefix(infos[0], "Bearer ") == "your-authorization-token"
}
