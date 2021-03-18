package server

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

func TimeOut(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		timeCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return handler(timeCtx, req)
	}
}

func Deadline(at time.Time) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		timeCtx, cancel := context.WithDeadline(ctx, at)
		defer cancel()

		return handler(timeCtx, req)
	}
}

var MetricServerOption = grpc.ChainUnaryInterceptor(Metric())

func Metric() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		st := time.Now()
		defer func() {
			duration := time.Since(st)
			fmt.Printf("rpc duration: %v\n", duration)
		}()

		return handler(ctx, req)
	}
}
