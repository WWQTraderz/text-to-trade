package grpc

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func UnaryUserIdInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return invoker(context.WithValue(ctx, "user_id", "123"), method, req, reply, cc, opts...)
}

func UnaryLoggingInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	startTime := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("method=%s duration=%s error=%v", method, time.Since(startTime), err)
	return err
}
