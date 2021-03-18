package server

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryMetric() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		st := time.Now()
		defer func() {
			duration := time.Since(st)
			fmt.Printf("rpc duration: %v\n", duration)
		}()

		return handler(ctx, req)
	})
}

func UnaryAuthorization(authorization Authorization) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = authorization.Do(ctx)
		if err != nil {
			return
		}

		return handler(ctx, req)
	})
}

func UnaryCrashHandler() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if p := recover(); p != nil {
				err = status.Errorf(codes.Internal, "%+v", p)
			}
		}()

		return handler(ctx, req)
	})
}

func StreamMetric() grpc.ServerOption {
	return grpc.ChainStreamInterceptor(func(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		st := time.Now()
		defer func() {
			duration := time.Since(st)
			fmt.Printf("rpc duration: %v\n", duration)
		}()

		return handler(srv, ss)
	})
}

func StreamAuthorization(authorization Authorization) grpc.ServerOption {
	return grpc.ChainStreamInterceptor(func(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := authorization.Do(ss.Context())
		if err != nil {
			return err
		}

		return handler(srv, ss)
	})
}

func StreamCrashHandler() grpc.ServerOption {
	return grpc.ChainStreamInterceptor(func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = status.Errorf(codes.Internal, "%+v", p)
			}
		}()

		return handler(srv, ss)
	})
}
