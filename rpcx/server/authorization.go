package server

import (
	"context"

	"google.golang.org/grpc/metadata"
)

type Authorization interface {
	Do(ctx context.Context) error
}

type defaultAuthorization struct {
	auth bool
}

func NewAuthorization(auth bool) Authorization {
	return &defaultAuthorization{
		auth: auth,
	}
}

func (a *defaultAuthorization) Do(ctx context.Context) error {
	if !a.auth {
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return grpcUnauthenticated
	}

	list, ok := md["token"]
	if !ok || len(list) == 0 {
		return grpcUnauthenticated
	}

	key := list[0]
	if key != "" {
		// TODO: auth logic
	}

	return nil
}
