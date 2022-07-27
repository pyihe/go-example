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
	"github.com/pyihe/go-pkg/rands"
	"github.com/pyihe/go-pkg/syncs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	testUnary()
	testServerStream()
	testClientStream()
	testBidStream()
}

func testUnary() {
	// 运行服务端
	closer, err := unary.Run(":8801")
	if err != nil {
		fmt.Printf("unary: start server fail: %v\n", err)
		return
	}
	defer closer.Close()

	// 开启客户端
	// 先创建grpc连接
	grpcConn, err := grpc.Dial(":8801", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("unary: grpc dial fail(%v)\n", err)
		return
	}
	defer grpcConn.Close()

	// 通过连接创建Unary服务的客户端
	unaryClient := proto.NewUnaryClient(grpcConn)

	// 如果需要携带metadata
	md := metadata.Pairs("token", rands.String(10))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	request := &proto.EchoRequest{Message: "I'm Unary!"}
	fmt.Printf("unary: sending request, md(%v), request(%v)\n", md, request.String())

	// header, trailer用于接收服务器返回的metadata
	var header, trailer metadata.MD
	resp, err := unaryClient.Echo(ctx, request, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		fmt.Printf("unary: echo fail(%v)\n", err)
		return
	}
	fmt.Printf("unary: receive response, header(%v), trailer(%v), response(%v)\n\n", header, trailer, resp.String())
}

func testServerStream() {
	// 先开启服务端
	closer, err := serverstream.Run(":8802")
	if err != nil {
		fmt.Printf("serverstream: start server fail(%v)\n", err)
		return
	}
	defer closer.Close()

	// 创建grpc连接
	grpcConn, err := grpc.Dial(":8802", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("serverstream: grpc dial fail(%v)\n", err)
		return
	}
	defer grpcConn.Close()

	// 创建grpc client
	client := proto.NewServerStreamClient(grpcConn)

	md := metadata.Pairs("token", rands.String(10))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	request := &proto.EchoRequest{Message: "I'm serverstream!"}
	fmt.Printf("serverstream: sending request, md(%v), request(%v)\n", md, request.String())

	// 发送请求同时获取服务器返回的stream
	stream, err := client.Echo(ctx, request)
	if err != nil {
		fmt.Printf("serverstream: echo fail(%v)\n", err)
		return
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				fmt.Printf("serverstream: recv fail(%v)\n", err)
			}
			break
		}
		fmt.Printf("serverstream: receive response, response(%v)\n", resp.String())
	}
	// 读取回复，因为执行Echo后，服务器已经处理完并返回，所以读取header的顺序没有影响（可以在读取回复前或者后）
	header, err := stream.Header()
	if err != nil {
		fmt.Printf("serverstream: stream header fail(%v)\n", err)
		return
	}
	fmt.Printf("serverstream: header(%v)\n", header)

	//但是trailer必须得等到所有回复完毕后才能读取
	trailer := stream.Trailer()
	fmt.Printf("serverstream: trailer(%v)\n\n", trailer)
}

func testClientStream() {
	// 先运行服务端
	closer, err := clientstream.Run(":8803")
	if err != nil {
		fmt.Printf("clientstream: start server fail(%v)\n", err)
		return
	}
	defer closer.Close()

	// 创建grpc连接
	grpcConn, err := grpc.Dial(":8803", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("clientstream: grpc dial fail(%v)\n", err)
		return
	}
	defer grpcConn.Close()

	// 创建grpc客户端
	client := proto.NewClientStreamClient(grpcConn)

	// 初始化请求和metadata
	request := &proto.EchoRequest{Message: "I'm clientstream!"}
	// 创建客户端的流
	md := metadata.Pairs("token", rands.String(16))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// 创建用于发送请求的stream
	stream, err := client.Echo(ctx)
	if err != nil {
		fmt.Printf("clientstream: create stream fail(%v)\n", err)
		return
	}
	fmt.Printf("clientstream: send, md(%v)\n", md)

	// 如果要在读取回复之前读取header，那么服务端需要在回复之前就调用stream.SendHeader
	//header, err := stream.Header()
	//if err != nil {
	//	fmt.Printf("clientstream: stream header fail(%v)\n", err)
	//	return
	//}
	//fmt.Printf("clientstream: header(%v)\n", header)

	// 发送请求
	for i := 0; i < 10; i++ {
		err = stream.Send(request)
		if err != nil {
			fmt.Printf("clientstream: send fail(%v)\n", err)
			return
		}
	}

	// 读取回复
	resp, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("clientstream: recv fail(%v)\n", err)
		return
	}
	// 如果是在读取回复之后再读取header，那么服务端只需要在发送回复之前写入header即可
	header, err := stream.Header()
	if err != nil {
		fmt.Printf("clientstream: stream header fail(%v)\n", err)
		return
	}
	fmt.Printf("clientstream: header(%v)\n", header)
	trailer := stream.Trailer()
	fmt.Printf("clientstream: trailer(%v), response(%v)\n\n", trailer, resp.String())
}

func testBidStream() {
	// 开启服务
	closer, err := bidstream.Run(":8804")
	if err != nil {
		fmt.Printf("bidstream: start server fail(%v)\n", err)
		return
	}
	defer closer.Close()

	// 创建grpc连接
	grpcConn, err := grpc.Dial(":8804", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("bidstream: grpc dial fail(%v)\n", err)
		return
	}
	defer grpcConn.Close()

	counter := new(syncs.AtomicInt32)
	wg := syncs.WgWrapper{}
	// 创建grpc客户端
	client := proto.NewBidStreamClient(grpcConn)

	md := metadata.Pairs("token", rands.String(16))
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := client.Echo(ctx)
	if err != nil {
		fmt.Printf("bidstream: create stream fail(%v)\n", err)
		return
	}

	// 这里用两个协程分别读/写grpc stream
	// 读
	wg.Wrap(func() {
		var header metadata.MD
		var resp *proto.EchoResponse

		header, err = stream.Header()
		if err != nil {
			fmt.Printf("bidstream client: header fail(%v)\n", err)
			return
		}
		fmt.Printf("bidstream client: header(%v)\n", header)
	loop:
		for {
			resp, err = stream.Recv()
			switch err {
			case io.EOF:
				break loop
			case nil:
				fmt.Printf("bidstream client: receive response(%v)\n", resp.String())
				resp.Reset()
				counter.Inc(1)
			default:
				fmt.Printf("bidstream client: receive fail(%v)\n", err)
				break loop
			}
		}
		trailer := stream.Trailer()
		fmt.Printf("bidstream client: trailer(%v)\n", trailer)
	})

	// 写，写入10此数据，然后发送close
	wg.Wrap(func() {
		for i := 0; i < 10; i++ {
			stream.Send(&proto.EchoRequest{Message: fmt.Sprintf("I'm bidstream!——%d", i)})
		}
		var v int32
		for v != 10 {
			time.Sleep(50 * time.Millisecond)
			v = counter.Value()
		}
		stream.CloseSend()
	})
	wg.Wait()
}
