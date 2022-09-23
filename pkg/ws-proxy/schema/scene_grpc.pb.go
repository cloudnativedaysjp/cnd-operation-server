// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.6
// source: pkg/ws-proxy/schema/scene.proto

package schema

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SceneServiceClient is the client API for SceneService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SceneServiceClient interface {
	ListScene(ctx context.Context, in *ListSceneRequest, opts ...grpc.CallOption) (*ListSceneResponse, error)
	MoveSceneToNext(ctx context.Context, in *MoveSceneToNextRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type sceneServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSceneServiceClient(cc grpc.ClientConnInterface) SceneServiceClient {
	return &sceneServiceClient{cc}
}

func (c *sceneServiceClient) ListScene(ctx context.Context, in *ListSceneRequest, opts ...grpc.CallOption) (*ListSceneResponse, error) {
	out := new(ListSceneResponse)
	err := c.cc.Invoke(ctx, "/schema.SceneService/ListScene", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sceneServiceClient) MoveSceneToNext(ctx context.Context, in *MoveSceneToNextRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/schema.SceneService/MoveSceneToNext", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SceneServiceServer is the server API for SceneService service.
// All implementations must embed UnimplementedSceneServiceServer
// for forward compatibility
type SceneServiceServer interface {
	ListScene(context.Context, *ListSceneRequest) (*ListSceneResponse, error)
	MoveSceneToNext(context.Context, *MoveSceneToNextRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedSceneServiceServer()
}

// UnimplementedSceneServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSceneServiceServer struct {
}

func (UnimplementedSceneServiceServer) ListScene(context.Context, *ListSceneRequest) (*ListSceneResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListScene not implemented")
}
func (UnimplementedSceneServiceServer) MoveSceneToNext(context.Context, *MoveSceneToNextRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MoveSceneToNext not implemented")
}
func (UnimplementedSceneServiceServer) mustEmbedUnimplementedSceneServiceServer() {}

// UnsafeSceneServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SceneServiceServer will
// result in compilation errors.
type UnsafeSceneServiceServer interface {
	mustEmbedUnimplementedSceneServiceServer()
}

func RegisterSceneServiceServer(s grpc.ServiceRegistrar, srv SceneServiceServer) {
	s.RegisterService(&SceneService_ServiceDesc, srv)
}

func _SceneService_ListScene_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSceneRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SceneServiceServer).ListScene(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/schema.SceneService/ListScene",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SceneServiceServer).ListScene(ctx, req.(*ListSceneRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SceneService_MoveSceneToNext_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MoveSceneToNextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SceneServiceServer).MoveSceneToNext(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/schema.SceneService/MoveSceneToNext",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SceneServiceServer).MoveSceneToNext(ctx, req.(*MoveSceneToNextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SceneService_ServiceDesc is the grpc.ServiceDesc for SceneService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SceneService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "schema.SceneService",
	HandlerType: (*SceneServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListScene",
			Handler:    _SceneService_ListScene_Handler,
		},
		{
			MethodName: "MoveSceneToNext",
			Handler:    _SceneService_MoveSceneToNext_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/ws-proxy/schema/scene.proto",
}