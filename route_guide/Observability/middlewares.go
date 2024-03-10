package Observability

import (
	"context"
	"google.golang.org/grpc"
	"log"
)

// UnaryServerInterceptor implements a unary server interceptor middleware
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Inbound middleware called%s", info.FullMethod)
	// Perform pre-processing logic here
	resp, err := handler(ctx, req)
	// Perform post-processing logic here
	log.Printf("Inboud middleware Completed: %s", info.FullMethod)
	return resp, err
}

// UnaryClientInterceptor implements a unary client interceptor middleware
func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Printf("Outbound middleware Called: %s", method)
	// Perform pre-processing logic here
	err := invoker(ctx, method, req, reply, cc, opts...)
	// Perform post-processing logic here
	log.Printf("Outbound middleware Completed: %s", method)
	return err
}
