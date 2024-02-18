// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: location/location.proto

package Proto

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

// LocationClient is the client API for Location service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LocationClient interface {
	// A simple RPC.
	//
	// Obtains the feature at a given position.
	//
	// A feature with an empty name is returned if there's no feature at the given
	// position.
	GetFeature(ctx context.Context, in *Point, opts ...grpc.CallOption) (*Feature, error)
}

type locationClient struct {
	cc grpc.ClientConnInterface
}

func NewLocationClient(cc grpc.ClientConnInterface) LocationClient {
	return &locationClient{cc}
}

func (c *locationClient) GetFeature(ctx context.Context, in *Point, opts ...grpc.CallOption) (*Feature, error) {
	out := new(Feature)
	err := c.cc.Invoke(ctx, "/location.location/GetFeature", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LocationServer is the server API for Location service.
// All implementations must embed UnimplementedLocationServer
// for forward compatibility
type LocationServer interface {
	// A simple RPC.
	//
	// Obtains the feature at a given position.
	//
	// A feature with an empty name is returned if there's no feature at the given
	// position.
	GetFeature(context.Context, *Point) (*Feature, error)
	mustEmbedUnimplementedLocationServer()
}

// UnimplementedLocationServer must be embedded to have forward compatible implementations.
type UnimplementedLocationServer struct {
}

func (UnimplementedLocationServer) GetFeature(context.Context, *Point) (*Feature, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFeature not implemented")
}
func (UnimplementedLocationServer) mustEmbedUnimplementedLocationServer() {}

// UnsafeLocationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LocationServer will
// result in compilation errors.
type UnsafeLocationServer interface {
	mustEmbedUnimplementedLocationServer()
}

func RegisterLocationServer(s grpc.ServiceRegistrar, srv LocationServer) {
	s.RegisterService(&Location_ServiceDesc, srv)
}

func _Location_GetFeature_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Point)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocationServer).GetFeature(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/location.location/GetFeature",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocationServer).GetFeature(ctx, req.(*Point))
	}
	return interceptor(ctx, in, info, handler)
}

// Location_ServiceDesc is the grpc.ServiceDesc for Location service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Location_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "location.location",
	HandlerType: (*LocationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFeature",
			Handler:    _Location_GetFeature_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "location/location.proto",
}
