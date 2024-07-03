// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.27.1
// source: rpc.proto

package __

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

// HubServiceClient is the client API for HubService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HubServiceClient interface {
	SubmitMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error)
}

type hubServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewHubServiceClient(cc grpc.ClientConnInterface) HubServiceClient {
	return &hubServiceClient{cc}
}

func (c *hubServiceClient) SubmitMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/HubService/SubmitMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HubServiceServer is the server API for HubService service.
// All implementations must embed UnimplementedHubServiceServer
// for forward compatibility
type HubServiceServer interface {
	SubmitMessage(context.Context, *Message) (*Message, error)
	mustEmbedUnimplementedHubServiceServer()
}

// UnimplementedHubServiceServer must be embedded to have forward compatible implementations.
type UnimplementedHubServiceServer struct {
}

func (UnimplementedHubServiceServer) SubmitMessage(context.Context, *Message) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitMessage not implemented")
}
func (UnimplementedHubServiceServer) mustEmbedUnimplementedHubServiceServer() {}

// UnsafeHubServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HubServiceServer will
// result in compilation errors.
type UnsafeHubServiceServer interface {
	mustEmbedUnimplementedHubServiceServer()
}

func RegisterHubServiceServer(s grpc.ServiceRegistrar, srv HubServiceServer) {
	s.RegisterService(&HubService_ServiceDesc, srv)
}

func _HubService_SubmitMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HubServiceServer).SubmitMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/HubService/SubmitMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HubServiceServer).SubmitMessage(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

// HubService_ServiceDesc is the grpc.ServiceDesc for HubService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var HubService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "HubService",
	HandlerType: (*HubServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubmitMessage",
			Handler:    _HubService_SubmitMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc.proto",
}
