package Observability

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

// UnaryServerInterceptor implements a unary server interceptor middleware
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	strt := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	if md == nil {
		md = metadata.New(map[string]string{})
	}
	caller := "unknown"
	if len(md["caller"]) > 0 {
		caller = md["caller"][0]
	}
	log.Printf("Starting execution of request fom %s for method : %s ", caller, info.FullMethod[1:])
	md["fullmethod"] = []string{info.FullMethod[1:]}
	ctx = metadata.NewIncomingContext(ctx, md)
	// Perform pre-processing logic here
	resp, err := handler(ctx, req)
	// Perform post-processing logic here
	log.Printf("Execution completed for request from %s for method %s in %s", caller, info.FullMethod[1:], time.Now().UTC().Sub(strt).String())
	return resp, err
}

// UnaryClientInterceptor implements a unary client interceptor middleware
func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	strt := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	if md == nil {
		md = metadata.New(map[string]string{})
	}
	origCaller := "unknown"
	if len(md["caller"]) > 0 {
		origCaller = md["caller"][0]
	}
	newCaller := "gateway"
	if len(md["fullmethod"]) > 0 {
		newCaller = metadata.ValueFromIncomingContext(ctx, "fullmethod")[0]
	}
	md["caller"] = []string{newCaller}
	log.Printf("Call to %s invoked in call chain\n%s->%s->%s", method[1:], origCaller, newCaller, method[1:])
	// Perform pre-processing logic here
	ctx = metadata.NewOutgoingContext(ctx, md)
	err := invoker(ctx, method, req, reply, cc, opts...)
	// Perform post-processing logic here
	log.Printf("Outbound Call to %s in call chain\n%s->%s->%s completed in %s", method[1:], origCaller, newCaller, method[1:], time.Now().UTC().Sub(strt).String())
	return err
}
