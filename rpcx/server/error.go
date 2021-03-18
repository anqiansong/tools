package server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var grpcUnauthenticated = status.Error(codes.Unauthenticated, "not authenticated")
