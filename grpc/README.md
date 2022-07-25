### GRPC

[grpc-go]()实现服务之间的调用，分别实现GRPC四种模式: 

1. [Unary RPC](https://grpc.io/docs/what-is-grpc/core-concepts/): 简单RPC调用, 客户端向服务器发送单个请求并返回单个响应，阻塞式的调用，客户端会一直等待直到服务器返回响应。
2. [Server streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/): 服务器流式调用, 客户端向服务器发送请求并获取流以读回一系列消息。客户端从返回的流中读取，直到没有更多消息为止。gRPC保证单个RPC调用中的消息顺序
3. [Client streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/): 客户端流式调用, 客户端写入一系列消息并将它们发送到服务器，再次使用提供的流。一旦客户端完成了消息的写入，它就会等待服务器读取它们并返回它的响应。gRPC再次保证单个RPC调用中的消息排序
4. [Bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/): 双向流式调用, 客户端和服务器均使用stream发送消息, 类似于TCP连接，全双工连接。

GRPC可以携带metadata(元数据), 通过[metadata](https://github.com/grpc/grpc-go/tree/master/metadata)包中的[NewIncomingContext](https://github.com/grpc/grpc-go/blob/master/metadata/metadata.go#L151)与[NewOutgoingContext](https://github.com/grpc/grpc-go/blob/master/metadata/metadata.go#L158)方法从客户端/服务器中写入或者读取元数据
```go
md := metadata.Pairs("key", value)
ctx := metadata.NewOutgoingContext(context.Background(), md)
_ = ctx
```

```go
md := metadata.MD{}
metadata.NewOutgoingContext(ctx, md)
```