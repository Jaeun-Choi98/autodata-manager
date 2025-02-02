// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.0
// source: bcnet.proto

package bcnet

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	BlockChainNetworkService_SendMessage_FullMethodName = "/bcnet.BlockChainNetworkService/SendMessage"
)

// BlockChainNetworkServiceClient is the client API for BlockChainNetworkService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BlockChainNetworkServiceClient interface {
	SendMessage(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResponse, error)
}

type blockChainNetworkServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBlockChainNetworkServiceClient(cc grpc.ClientConnInterface) BlockChainNetworkServiceClient {
	return &blockChainNetworkServiceClient{cc}
}

func (c *blockChainNetworkServiceClient) SendMessage(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MessageResponse)
	err := c.cc.Invoke(ctx, BlockChainNetworkService_SendMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BlockChainNetworkServiceServer is the server API for BlockChainNetworkService service.
// All implementations must embed UnimplementedBlockChainNetworkServiceServer
// for forward compatibility.
type BlockChainNetworkServiceServer interface {
	SendMessage(context.Context, *MessageRequest) (*MessageResponse, error)
	mustEmbedUnimplementedBlockChainNetworkServiceServer()
}

// UnimplementedBlockChainNetworkServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBlockChainNetworkServiceServer struct{}

func (UnimplementedBlockChainNetworkServiceServer) SendMessage(context.Context, *MessageRequest) (*MessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedBlockChainNetworkServiceServer) mustEmbedUnimplementedBlockChainNetworkServiceServer() {
}
func (UnimplementedBlockChainNetworkServiceServer) testEmbeddedByValue() {}

// UnsafeBlockChainNetworkServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BlockChainNetworkServiceServer will
// result in compilation errors.
type UnsafeBlockChainNetworkServiceServer interface {
	mustEmbedUnimplementedBlockChainNetworkServiceServer()
}

func RegisterBlockChainNetworkServiceServer(s grpc.ServiceRegistrar, srv BlockChainNetworkServiceServer) {
	// If the following call pancis, it indicates UnimplementedBlockChainNetworkServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&BlockChainNetworkService_ServiceDesc, srv)
}

func _BlockChainNetworkService_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockChainNetworkServiceServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BlockChainNetworkService_SendMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockChainNetworkServiceServer).SendMessage(ctx, req.(*MessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BlockChainNetworkService_ServiceDesc is the grpc.ServiceDesc for BlockChainNetworkService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BlockChainNetworkService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bcnet.BlockChainNetworkService",
	HandlerType: (*BlockChainNetworkServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMessage",
			Handler:    _BlockChainNetworkService_SendMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "bcnet.proto",
}
