package client

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const authorizationKey = "token"

type Authorization interface {
	Do(ctx context.Context) context.Context
}

type defaultAuthorization struct {
	token string
}

func NewAuthorization(token string) Authorization {
	return &defaultAuthorization{
		token: token,
	}
}

func (a *defaultAuthorization) Do(ctx context.Context) context.Context {
	md := metadata.New(map[string]string{
		authorizationKey: a.token,
	})

	return metadata.NewOutgoingContext(ctx, md)
}
