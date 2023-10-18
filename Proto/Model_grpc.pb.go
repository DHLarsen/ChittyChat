// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: Proto/Model.proto

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

const (
	Model_SendMessage_FullMethodName = "/proto.Model/SendMessage"
)

// ModelClient is the client API for Model service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ModelClient interface {
	SendMessage(ctx context.Context, opts ...grpc.CallOption) (Model_SendMessageClient, error)
}

type modelClient struct {
	cc grpc.ClientConnInterface
}

func NewModelClient(cc grpc.ClientConnInterface) ModelClient {
	return &modelClient{cc}
}

func (c *modelClient) SendMessage(ctx context.Context, opts ...grpc.CallOption) (Model_SendMessageClient, error) {
	stream, err := c.cc.NewStream(ctx, &Model_ServiceDesc.Streams[0], Model_SendMessage_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &modelSendMessageClient{stream}
	return x, nil
}

type Model_SendMessageClient interface {
	Send(*Message) error
	CloseAndRecv() (*Ack, error)
	grpc.ClientStream
}

type modelSendMessageClient struct {
	grpc.ClientStream
}

func (x *modelSendMessageClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *modelSendMessageClient) CloseAndRecv() (*Ack, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Ack)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ModelServer is the server API for Model service.
// All implementations must embed UnimplementedModelServer
// for forward compatibility
type ModelServer interface {
	SendMessage(Model_SendMessageServer) error
	mustEmbedUnimplementedModelServer()
}

// UnimplementedModelServer must be embedded to have forward compatible implementations.
type UnimplementedModelServer struct {
}

func (UnimplementedModelServer) SendMessage(Model_SendMessageServer) error {
	return status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedModelServer) mustEmbedUnimplementedModelServer() {}

// UnsafeModelServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ModelServer will
// result in compilation errors.
type UnsafeModelServer interface {
	mustEmbedUnimplementedModelServer()
}

func RegisterModelServer(s grpc.ServiceRegistrar, srv ModelServer) {
	s.RegisterService(&Model_ServiceDesc, srv)
}

func _Model_SendMessage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ModelServer).SendMessage(&modelSendMessageServer{stream})
}

type Model_SendMessageServer interface {
	SendAndClose(*Ack) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type modelSendMessageServer struct {
	grpc.ServerStream
}

func (x *modelSendMessageServer) SendAndClose(m *Ack) error {
	return x.ServerStream.SendMsg(m)
}

func (x *modelSendMessageServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Model_ServiceDesc is the grpc.ServiceDesc for Model service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Model_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Model",
	HandlerType: (*ModelServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SendMessage",
			Handler:       _Model_SendMessage_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "Proto/Model.proto",
}