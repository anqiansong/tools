package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

func UnaryAuthorization(authorization Authorization) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		newCtx := authorization.Do(ctx)
		return invoker(newCtx, method, req, reply, cc, opts...)
	})
}

func UnaryMetric() grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		st := time.Now()
		defer func() {
			fmt.Printf("duration: %v\n", time.Since(st))
		}()

		return invoker(ctx, method, req, reply, cc, opts...)
	})
}

func UnaryCrashHandler() grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = fmt.Errorf("%+v", p)
			}
		}()

		return invoker(ctx, method, req, reply, cc, opts...)
	})
}

func UnaryTimeout(duration time.Duration) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		if duration > 0 {
			newCtx, cancel := context.WithTimeout(ctx, duration)
			defer cancel()
			return invoker(newCtx, method, req, reply, cc, opts...)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	})
}

func StreamAuthorization(authorization Authorization) grpc.DialOption {
	return grpc.WithChainStreamInterceptor(func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		newCtx := authorization.Do(ctx)
		return streamer(newCtx, desc, cc, method, opts...)
	})
}

func StreamMetric() grpc.DialOption {
	return grpc.WithChainStreamInterceptor(func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		st := time.Now()
		defer func() {
			fmt.Printf("duration: %v\n", time.Since(st))
		}()

		return streamer(ctx, desc, cc, method, opts...)
	})
}

func StreamCrashHandler() grpc.DialOption {
	return grpc.WithChainStreamInterceptor(func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (s grpc.ClientStream, err error) {
		defer func() {
			if p := recover(); p != nil {
				err = fmt.Errorf("%+v", p)
			}
		}()

		return streamer(ctx, desc, cc, method, opts...)
	})
}

func StreamTimeout(duration time.Duration) grpc.DialOption {
	return grpc.WithChainStreamInterceptor(func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (s grpc.ClientStream, err error) {
		if duration > 0 {
			newCtx, cancel := context.WithTimeout(ctx, duration)
			defer cancel()
			return streamer(newCtx, desc, cc, method, opts...)
		}

		return streamer(ctx, desc, cc, method, opts...)
	})
}
