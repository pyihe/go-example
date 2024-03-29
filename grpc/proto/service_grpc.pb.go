// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.17.2
// source: service.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UnaryClient is the client API for Unary service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UnaryClient interface {
	Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoResponse, error)
}

type unaryClient struct {
	cc grpc.ClientConnInterface
}

func NewUnaryClient(cc grpc.ClientConnInterface) UnaryClient {
	return &unaryClient{cc}
}

func (c *unaryClient) Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoResponse, error) {
	out := new(EchoResponse)
	err := c.cc.Invoke(ctx, "/proto.Unary/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UnaryServer is the server API for Unary service.
// All implementations must embed UnimplementedUnaryServer
// for forward compatibility
type UnaryServer interface {
	Echo(context.Context, *EchoRequest) (*EchoResponse, error)
	mustEmbedUnimplementedUnaryServer()
}

// UnimplementedUnaryServer must be embedded to have forward compatible implementations.
type UnimplementedUnaryServer struct {
}

func (UnimplementedUnaryServer) Echo(context.Context, *EchoRequest) (*EchoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Echo not implemented")
}
func (UnimplementedUnaryServer) mustEmbedUnimplementedUnaryServer() {}

// UnsafeUnaryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UnaryServer will
// result in compilation errors.
type UnsafeUnaryServer interface {
	mustEmbedUnimplementedUnaryServer()
}

func RegisterUnaryServer(s grpc.ServiceRegistrar, srv UnaryServer) {
	s.RegisterService(&Unary_ServiceDesc, srv)
}

func _Unary_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UnaryServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Unary/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UnaryServer).Echo(ctx, req.(*EchoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Unary_ServiceDesc is the grpc.ServiceDesc for Unary service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Unary_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Unary",
	HandlerType: (*UnaryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _Unary_Echo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}

// ServerStreamClient is the client API for ServerStream service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServerStreamClient interface {
	Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (ServerStream_EchoClient, error)
}

type serverStreamClient struct {
	cc grpc.ClientConnInterface
}

func NewServerStreamClient(cc grpc.ClientConnInterface) ServerStreamClient {
	return &serverStreamClient{cc}
}

func (c *serverStreamClient) Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (ServerStream_EchoClient, error) {
	stream, err := c.cc.NewStream(ctx, &ServerStream_ServiceDesc.Streams[0], "/proto.ServerStream/Echo", opts...)
	if err != nil {
		return nil, err
	}
	x := &serverStreamEchoClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ServerStream_EchoClient interface {
	Recv() (*EchoResponse, error)
	grpc.ClientStream
}

type serverStreamEchoClient struct {
	grpc.ClientStream
}

func (x *serverStreamEchoClient) Recv() (*EchoResponse, error) {
	m := new(EchoResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ServerStreamServer is the server API for ServerStream service.
// All implementations must embed UnimplementedServerStreamServer
// for forward compatibility
type ServerStreamServer interface {
	Echo(*EchoRequest, ServerStream_EchoServer) error
	mustEmbedUnimplementedServerStreamServer()
}

// UnimplementedServerStreamServer must be embedded to have forward compatible implementations.
type UnimplementedServerStreamServer struct {
}

func (UnimplementedServerStreamServer) Echo(*EchoRequest, ServerStream_EchoServer) error {
	return status.Errorf(codes.Unimplemented, "method Echo not implemented")
}
func (UnimplementedServerStreamServer) mustEmbedUnimplementedServerStreamServer() {}

// UnsafeServerStreamServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServerStreamServer will
// result in compilation errors.
type UnsafeServerStreamServer interface {
	mustEmbedUnimplementedServerStreamServer()
}

func RegisterServerStreamServer(s grpc.ServiceRegistrar, srv ServerStreamServer) {
	s.RegisterService(&ServerStream_ServiceDesc, srv)
}

func _ServerStream_Echo_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EchoRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ServerStreamServer).Echo(m, &serverStreamEchoServer{stream})
}

type ServerStream_EchoServer interface {
	Send(*EchoResponse) error
	grpc.ServerStream
}

type serverStreamEchoServer struct {
	grpc.ServerStream
}

func (x *serverStreamEchoServer) Send(m *EchoResponse) error {
	return x.ServerStream.SendMsg(m)
}

// ServerStream_ServiceDesc is the grpc.ServiceDesc for ServerStream service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ServerStream_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ServerStream",
	HandlerType: (*ServerStreamServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Echo",
			Handler:       _ServerStream_Echo_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "service.proto",
}

// ClientStreamClient is the client API for ClientStream service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClientStreamClient interface {
	Echo(ctx context.Context, opts ...grpc.CallOption) (ClientStream_EchoClient, error)
}

type clientStreamClient struct {
	cc grpc.ClientConnInterface
}

func NewClientStreamClient(cc grpc.ClientConnInterface) ClientStreamClient {
	return &clientStreamClient{cc}
}

func (c *clientStreamClient) Echo(ctx context.Context, opts ...grpc.CallOption) (ClientStream_EchoClient, error) {
	stream, err := c.cc.NewStream(ctx, &ClientStream_ServiceDesc.Streams[0], "/proto.ClientStream/Echo", opts...)
	if err != nil {
		return nil, err
	}
	x := &clientStreamEchoClient{stream}
	return x, nil
}

type ClientStream_EchoClient interface {
	Send(*EchoRequest) error
	CloseAndRecv() (*EchoResponse, error)
	grpc.ClientStream
}

type clientStreamEchoClient struct {
	grpc.ClientStream
}

func (x *clientStreamEchoClient) Send(m *EchoRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *clientStreamEchoClient) CloseAndRecv() (*EchoResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(EchoResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ClientStreamServer is the server API for ClientStream service.
// All implementations must embed UnimplementedClientStreamServer
// for forward compatibility
type ClientStreamServer interface {
	Echo(ClientStream_EchoServer) error
	mustEmbedUnimplementedClientStreamServer()
}

// UnimplementedClientStreamServer must be embedded to have forward compatible implementations.
type UnimplementedClientStreamServer struct {
}

func (UnimplementedClientStreamServer) Echo(ClientStream_EchoServer) error {
	return status.Errorf(codes.Unimplemented, "method Echo not implemented")
}
func (UnimplementedClientStreamServer) mustEmbedUnimplementedClientStreamServer() {}

// UnsafeClientStreamServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClientStreamServer will
// result in compilation errors.
type UnsafeClientStreamServer interface {
	mustEmbedUnimplementedClientStreamServer()
}

func RegisterClientStreamServer(s grpc.ServiceRegistrar, srv ClientStreamServer) {
	s.RegisterService(&ClientStream_ServiceDesc, srv)
}

func _ClientStream_Echo_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ClientStreamServer).Echo(&clientStreamEchoServer{stream})
}

type ClientStream_EchoServer interface {
	SendAndClose(*EchoResponse) error
	Recv() (*EchoRequest, error)
	grpc.ServerStream
}

type clientStreamEchoServer struct {
	grpc.ServerStream
}

func (x *clientStreamEchoServer) SendAndClose(m *EchoResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *clientStreamEchoServer) Recv() (*EchoRequest, error) {
	m := new(EchoRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ClientStream_ServiceDesc is the grpc.ServiceDesc for ClientStream service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClientStream_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ClientStream",
	HandlerType: (*ClientStreamServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Echo",
			Handler:       _ClientStream_Echo_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "service.proto",
}

// BidStreamClient is the client API for BidStream service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BidStreamClient interface {
	Echo(ctx context.Context, opts ...grpc.CallOption) (BidStream_EchoClient, error)
}

type bidStreamClient struct {
	cc grpc.ClientConnInterface
}

func NewBidStreamClient(cc grpc.ClientConnInterface) BidStreamClient {
	return &bidStreamClient{cc}
}

func (c *bidStreamClient) Echo(ctx context.Context, opts ...grpc.CallOption) (BidStream_EchoClient, error) {
	stream, err := c.cc.NewStream(ctx, &BidStream_ServiceDesc.Streams[0], "/proto.BidStream/Echo", opts...)
	if err != nil {
		return nil, err
	}
	x := &bidStreamEchoClient{stream}
	return x, nil
}

type BidStream_EchoClient interface {
	Send(*EchoRequest) error
	Recv() (*EchoResponse, error)
	grpc.ClientStream
}

type bidStreamEchoClient struct {
	grpc.ClientStream
}

func (x *bidStreamEchoClient) Send(m *EchoRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *bidStreamEchoClient) Recv() (*EchoResponse, error) {
	m := new(EchoResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BidStreamServer is the server API for BidStream service.
// All implementations must embed UnimplementedBidStreamServer
// for forward compatibility
type BidStreamServer interface {
	Echo(BidStream_EchoServer) error
	mustEmbedUnimplementedBidStreamServer()
}

// UnimplementedBidStreamServer must be embedded to have forward compatible implementations.
type UnimplementedBidStreamServer struct {
}

func (UnimplementedBidStreamServer) Echo(BidStream_EchoServer) error {
	return status.Errorf(codes.Unimplemented, "method Echo not implemented")
}
func (UnimplementedBidStreamServer) mustEmbedUnimplementedBidStreamServer() {}

// UnsafeBidStreamServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BidStreamServer will
// result in compilation errors.
type UnsafeBidStreamServer interface {
	mustEmbedUnimplementedBidStreamServer()
}

func RegisterBidStreamServer(s grpc.ServiceRegistrar, srv BidStreamServer) {
	s.RegisterService(&BidStream_ServiceDesc, srv)
}

func _BidStream_Echo_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BidStreamServer).Echo(&bidStreamEchoServer{stream})
}

type BidStream_EchoServer interface {
	Send(*EchoResponse) error
	Recv() (*EchoRequest, error)
	grpc.ServerStream
}

type bidStreamEchoServer struct {
	grpc.ServerStream
}

func (x *bidStreamEchoServer) Send(m *EchoResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *bidStreamEchoServer) Recv() (*EchoRequest, error) {
	m := new(EchoRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BidStream_ServiceDesc is the grpc.ServiceDesc for BidStream service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BidStream_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.BidStream",
	HandlerType: (*BidStreamServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Echo",
			Handler:       _BidStream_Echo_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "service.proto",
}
