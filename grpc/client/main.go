package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/pyihe/go-example/grpc/bidstream"
	"github.com/pyihe/go-example/grpc/clientstream"
	"github.com/pyihe/go-example/grpc/proto"
	"github.com/pyihe/go-example/grpc/serverstream"
	"github.com/pyihe/go-example/grpc/unary"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 开启各个类型的服务
	// unary
	unaryCloser, err := unary.Run(":8801")
	if err != nil {
		fmt.Printf("unary run err: %v\n", err)
		return
	}
	defer unaryCloser.Close()

	// server stream
	ssCloser, err := serverstream.Run(":8802")
	if err != nil {
		fmt.Printf("server stream run err: %v\n", err)
		return
	}
	defer ssCloser.Close()

	// client stream
	ccCloser, err := clientstream.Run(":8803")
	if err != nil {
		fmt.Printf("client stream run err: %v\n", err)
		return
	}
	defer ccCloser.Close()

	// bid stream
	bidCloser, err := bidstream.Run(":8804")
	if err != nil {
		fmt.Printf("bid stream err: %v\n", err)
		return
	}
	defer bidCloser.Close()

	//unaryClient()
	//serverStream()
	//clientStream()
	bidStream()
}

func unaryClient() {
	conn, err := grpc.Dial(":8801", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("unary dial err: %v\n", err)
		return
	}
	client := proto.NewUnaryClient(conn)
	resp, err := client.Echo(context.Background(), &proto.EchoRequest{Message: "Hello World!"})
	if err != nil {
		fmt.Printf("unary echo err: %v\n", err)
		return
	}
	fmt.Printf("unary客户端收到回复: %v\n", resp.Message)
	time.Sleep(1 * time.Second)
	conn.Close()
}

func serverStream() {
	conn, err := grpc.Dial(":8802", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("server stream dial err: %v\n", err)
		return
	}
	defer conn.Close()

	// 客户端调用API后会返回一个Stream, 然后从Stream中读取服务器返回的响应
	client := proto.NewServerStreamClient(conn)
	stream, err := client.Echo(context.Background(), &proto.EchoRequest{Message: "Hello World!"})
	if err != nil {
		fmt.Printf("ServerStream echo err: %v\n", err)
		return
	}
	// 需要处理stream返回的数据
	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					fmt.Printf("server stream err: %v\n", err)
				}
				return
			}
			fmt.Printf("ServerStream客户端收到回复: %v\n", resp.Message)
		}

	}()
	time.Sleep(1 * time.Second)
	stream.CloseSend()
}

func clientStream() {
	conn, err := grpc.Dial(":8803", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("client stream dial err: %v\n", err)
		return
	}
	defer conn.Close()
	client := proto.NewClientStreamClient(conn)
	stream, err := client.Echo(context.Background())
	if err != nil {
		fmt.Printf("client stream err: %v\n", err)
		return
	}
	stream.Send(&proto.EchoRequest{Message: "Hello World!"})

	// 处理服务器返回的响应
	resp, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("CloseAndRecv() err: %v\n", err)
		return
	}
	fmt.Printf("ClientStream客户端收到回复: %v\n", resp.Message)

	time.Sleep(1 * time.Second)
	stream.CloseSend()
}

func bidStream() {
	conn, err := grpc.Dial(":8804", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("bid stream err: %v\n", err)
		return
	}
	defer conn.Close()

	client := proto.NewBidStreamClient(conn)
	stream, err := client.Echo(context.Background())
	if err != nil {
		fmt.Printf("new bid stream err: %v\n", err)
		return
	}

	go func() {
		for {
			resp := proto.EchoResponse{}
			err = stream.RecvMsg(&resp)
			if err != nil {
				if err == io.EOF {
					fmt.Printf("EOF")
					return
				}
				fmt.Printf("bid stream recv err: %v\n", err)
				return
			}
			fmt.Printf("bid stream收到回复: %v\n", resp.Message)
		}
	}()

	stream.SendMsg(&proto.EchoRequest{Message: "Hello World!"})
	time.Sleep(1 * time.Second)
	//stream.CloseSend()
}
