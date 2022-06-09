// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: pkg/workout/pb/workout/workout.proto

package workout

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

// WorkoutServiceClient is the client API for WorkoutService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WorkoutServiceClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error)
	One(ctx context.Context, in *OneRequest, opts ...grpc.CallOption) (*OneResponse, error)
	All(ctx context.Context, in *AllRequest, opts ...grpc.CallOption) (*AllResponse, error)
}

type workoutServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWorkoutServiceClient(cc grpc.ClientConnInterface) WorkoutServiceClient {
	return &workoutServiceClient{cc}
}

func (c *workoutServiceClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, "/WorkoutService.WorkoutService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workoutServiceClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, "/WorkoutService.WorkoutService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workoutServiceClient) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error) {
	out := new(UpdateResponse)
	err := c.cc.Invoke(ctx, "/WorkoutService.WorkoutService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workoutServiceClient) One(ctx context.Context, in *OneRequest, opts ...grpc.CallOption) (*OneResponse, error) {
	out := new(OneResponse)
	err := c.cc.Invoke(ctx, "/WorkoutService.WorkoutService/One", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workoutServiceClient) All(ctx context.Context, in *AllRequest, opts ...grpc.CallOption) (*AllResponse, error) {
	out := new(AllResponse)
	err := c.cc.Invoke(ctx, "/WorkoutService.WorkoutService/All", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WorkoutServiceServer is the server API for WorkoutService service.
// All implementations must embed UnimplementedWorkoutServiceServer
// for forward compatibility
type WorkoutServiceServer interface {
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	Delete(context.Context, *DeleteRequest) (*DeleteResponse, error)
	Update(context.Context, *UpdateRequest) (*UpdateResponse, error)
	One(context.Context, *OneRequest) (*OneResponse, error)
	All(context.Context, *AllRequest) (*AllResponse, error)
	mustEmbedUnimplementedWorkoutServiceServer()
}

// UnimplementedWorkoutServiceServer must be embedded to have forward compatible implementations.
type UnimplementedWorkoutServiceServer struct {
}

func (UnimplementedWorkoutServiceServer) Create(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedWorkoutServiceServer) Delete(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedWorkoutServiceServer) Update(context.Context, *UpdateRequest) (*UpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedWorkoutServiceServer) One(context.Context, *OneRequest) (*OneResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method One not implemented")
}
func (UnimplementedWorkoutServiceServer) All(context.Context, *AllRequest) (*AllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method All not implemented")
}
func (UnimplementedWorkoutServiceServer) mustEmbedUnimplementedWorkoutServiceServer() {}

// UnsafeWorkoutServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WorkoutServiceServer will
// result in compilation errors.
type UnsafeWorkoutServiceServer interface {
	mustEmbedUnimplementedWorkoutServiceServer()
}

func RegisterWorkoutServiceServer(s grpc.ServiceRegistrar, srv WorkoutServiceServer) {
	s.RegisterService(&WorkoutService_ServiceDesc, srv)
}

func _WorkoutService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkoutServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/WorkoutService.WorkoutService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkoutServiceServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkoutService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkoutServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/WorkoutService.WorkoutService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkoutServiceServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkoutService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkoutServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/WorkoutService.WorkoutService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkoutServiceServer).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkoutService_One_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OneRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkoutServiceServer).One(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/WorkoutService.WorkoutService/One",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkoutServiceServer).One(ctx, req.(*OneRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkoutService_All_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkoutServiceServer).All(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/WorkoutService.WorkoutService/All",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkoutServiceServer).All(ctx, req.(*AllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// WorkoutService_ServiceDesc is the grpc.ServiceDesc for WorkoutService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WorkoutService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "WorkoutService.WorkoutService",
	HandlerType: (*WorkoutServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _WorkoutService_Create_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _WorkoutService_Delete_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _WorkoutService_Update_Handler,
		},
		{
			MethodName: "One",
			Handler:    _WorkoutService_One_Handler,
		},
		{
			MethodName: "All",
			Handler:    _WorkoutService_All_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/workout/pb/workout/workout.proto",
}
