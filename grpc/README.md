### GRPC

[grpc-go](https://github.com/grpc/grpc-go)实现服务之间的调用，分别实现GRPC四种模式: 

1. [Unary RPC](https://grpc.io/docs/what-is-grpc/core-concepts/): 简单RPC调用, request-response一对一模式, 一个请求对应一个响应。客户端向服务器发送单个请求并返回单个响应，阻塞式的调用，客户端会一直等待直到服务器返回响应。
2. [Server streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/): 服务器流式调用, request-response一对多模式, 一个请求对应多个响应。客户端向服务器发送请求并获取流以读回一系列消息。客户端从返回的流中读取，直到没有更多消息为止。gRPC保证单个RPC调用中的消息顺序
3. [Client streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/): 客户端流式调用, request-response多对一模式, 客户端写入一系列消息并将它们发送到服务器，同时使用返回的流读取服务器回复。一旦客户端完成了消息的写入，它就会等待服务器读取它们并返回它的响应。gRPC保证单个RPC调用中的消息排序。
4. [Bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/): 双向流式调用, request-response多对多模式, 客户端和服务器均使用stream发送消息, 类似于TCP连接，全双工连接。

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

### Protocol Buffers 3
1. 编号1-15对应出现频率高的字段，所以此区间的编号最好预留给需要频繁出现的字段。
2. 编号19000-19999为Protocol Buffers预留字段，不可用。
3. 为了proto文件在使用期间的前后兼容，中途需要删除或者废弃的字段，建议使用`reserved`关键字进行保留，`reserved`关键字既可以作用于字段编号也可以作用于字段名；消息中不能再使用被`reserved`标记的编号或者字段名。 
4. 字段默认值
   * `string`(字符串类型): 空字符串
   * `bytes`(字节切片类型): `nil`
   * `bool`(布尔类型): `false`
   * 数值类型: 0
   * 枚举类型: 默认值为定义的第一个枚举值(必须为0)
   * repeated类型: nil
   * 其他message类型: nil
5. 枚举类型的第一个常量元素对应的值必须为0，如果想让同一个枚举类型中的两个常量元素具有相同的值，可以通过设置`allow_alias`为`true`来实现：
   ```go
   enum Sports {
    option allow_alias = true;
    Unknown = 0; // 第一个常量元素值必须为0
    Basketball = 1; // 常量值为1
    Soccer = 1; // 常量值为1
   }
   ```
6. 数据类型

| proto类型  | Go类型    | 备注                                            |
|:---------|:--------|:----------------------------------------------|
| bool     | bool    ||
| double   | float64 ||
| float    | float32 ||
| int32    | int32   | 使用可变长度编码，如果字段值包含负数，建议使用sint32，因为int32编码负数的效率低 |
| int64    | int64   | 使用可变长度编码，如果字段值包含负数，建议使用sint64，因为int64编码负数的效率低 |
| uint32   | uint32  | 使用可变长度进行编码                                    |
| uint64   | uint64  | 使用可变长度进行编码                                    |
| sint32   | int32   | 使用可变长度进行编码，有符号整型，编码负数的效率比int32类型高             |
| sint64   | int64   | 使用可变长度进行编码，有符号整型，编码负数的效率比int64类型高             |
| fixed32  | uint32  | 编码长度总是4字节，如果字段值经常比2^28大，fixed32编码效率将比uint32高  |
| fixed64  | uint64  | 编码长度总是8字节，如果字段值经常比2^56大，fixed64编码效率比uint64高   |
| sfixed32 | int32   | 编码长度总是4字节                                     |
| sfixed64 | int64   | 编码长度总是8字节                                     |
| string   | string  | 值必须包含UTF-8编码或者7位ASCII文本，值长度不能超过2^32           |
| bytes    | []byte  | 长度不能超过2^32                                    |