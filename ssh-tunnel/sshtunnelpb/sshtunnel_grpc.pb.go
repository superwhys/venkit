// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: sshtunnelpb/sshtunnel.proto

package sshtunnelpb

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
	SshTunnel_Forward_FullMethodName     = "/SshTunnel/Forward"
	SshTunnel_Reverse_FullMethodName     = "/SshTunnel/Reverse"
	SshTunnel_ListConnect_FullMethodName = "/SshTunnel/ListConnect"
	SshTunnel_Disconnect_FullMethodName  = "/SshTunnel/Disconnect"
)

// SshTunnelClient is the client API for SshTunnel service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SshTunnelClient interface {
	Forward(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (*ForwardReply, error)
	Reverse(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (*ReverseReply, error)
	ListConnect(ctx context.Context, in *ListConnectRequest, opts ...grpc.CallOption) (*ListConnectReply, error)
	Disconnect(ctx context.Context, in *DisconnectRequest, opts ...grpc.CallOption) (*DisconnectReply, error)
}

type sshTunnelClient struct {
	cc grpc.ClientConnInterface
}

func NewSshTunnelClient(cc grpc.ClientConnInterface) SshTunnelClient {
	return &sshTunnelClient{cc}
}

func (c *sshTunnelClient) Forward(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (*ForwardReply, error) {
	out := new(ForwardReply)
	err := c.cc.Invoke(ctx, SshTunnel_Forward_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sshTunnelClient) Reverse(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (*ReverseReply, error) {
	out := new(ReverseReply)
	err := c.cc.Invoke(ctx, SshTunnel_Reverse_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sshTunnelClient) ListConnect(ctx context.Context, in *ListConnectRequest, opts ...grpc.CallOption) (*ListConnectReply, error) {
	out := new(ListConnectReply)
	err := c.cc.Invoke(ctx, SshTunnel_ListConnect_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sshTunnelClient) Disconnect(ctx context.Context, in *DisconnectRequest, opts ...grpc.CallOption) (*DisconnectReply, error) {
	out := new(DisconnectReply)
	err := c.cc.Invoke(ctx, SshTunnel_Disconnect_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SshTunnelServer is the server API for SshTunnel service.
// All implementations must embed UnimplementedSshTunnelServer
// for forward compatibility
type SshTunnelServer interface {
	Forward(context.Context, *ConnectRequest) (*ForwardReply, error)
	Reverse(context.Context, *ConnectRequest) (*ReverseReply, error)
	ListConnect(context.Context, *ListConnectRequest) (*ListConnectReply, error)
	Disconnect(context.Context, *DisconnectRequest) (*DisconnectReply, error)
	mustEmbedUnimplementedSshTunnelServer()
}

// UnimplementedSshTunnelServer must be embedded to have forward compatible implementations.
type UnimplementedSshTunnelServer struct {
}

func (UnimplementedSshTunnelServer) Forward(context.Context, *ConnectRequest) (*ForwardReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Forward not implemented")
}
func (UnimplementedSshTunnelServer) Reverse(context.Context, *ConnectRequest) (*ReverseReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reverse not implemented")
}
func (UnimplementedSshTunnelServer) ListConnect(context.Context, *ListConnectRequest) (*ListConnectReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListConnect not implemented")
}
func (UnimplementedSshTunnelServer) Disconnect(context.Context, *DisconnectRequest) (*DisconnectReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Disconnect not implemented")
}
func (UnimplementedSshTunnelServer) mustEmbedUnimplementedSshTunnelServer() {}

// UnsafeSshTunnelServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SshTunnelServer will
// result in compilation errors.
type UnsafeSshTunnelServer interface {
	mustEmbedUnimplementedSshTunnelServer()
}

func RegisterSshTunnelServer(s grpc.ServiceRegistrar, srv SshTunnelServer) {
	s.RegisterService(&SshTunnel_ServiceDesc, srv)
}

func _SshTunnel_Forward_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SshTunnelServer).Forward(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SshTunnel_Forward_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SshTunnelServer).Forward(ctx, req.(*ConnectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SshTunnel_Reverse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SshTunnelServer).Reverse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SshTunnel_Reverse_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SshTunnelServer).Reverse(ctx, req.(*ConnectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SshTunnel_ListConnect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListConnectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SshTunnelServer).ListConnect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SshTunnel_ListConnect_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SshTunnelServer).ListConnect(ctx, req.(*ListConnectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SshTunnel_Disconnect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DisconnectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SshTunnelServer).Disconnect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SshTunnel_Disconnect_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SshTunnelServer).Disconnect(ctx, req.(*DisconnectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SshTunnel_ServiceDesc is the grpc.ServiceDesc for SshTunnel service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SshTunnel_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "SshTunnel",
	HandlerType: (*SshTunnelServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Forward",
			Handler:    _SshTunnel_Forward_Handler,
		},
		{
			MethodName: "Reverse",
			Handler:    _SshTunnel_Reverse_Handler,
		},
		{
			MethodName: "ListConnect",
			Handler:    _SshTunnel_ListConnect_Handler,
		},
		{
			MethodName: "Disconnect",
			Handler:    _SshTunnel_Disconnect_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sshtunnelpb/sshtunnel.proto",
}