package rpcx

import (
	"net"

	"github.com/anqiansong/tools/rpcx/server"

	"google.golang.org/grpc"
)

const netWorkTcp = "tcp"

type ServerConfig struct {
	ListenOn string
	Auth     bool
}

type Register func(server *grpc.Server)

type Server interface {
	AddServerOption(options ...grpc.ServerOption)
	AddUnaryServerInterceptor(interceptors ...grpc.UnaryServerInterceptor)
	AddStreamServerInterceptor(interceptors ...grpc.StreamServerInterceptor)
	Serve(register Register) error
}

type ServerOption func(s Server)

type defaultServer struct {
	conf    *ServerConfig
	options []grpc.ServerOption
}

func NewServer(conf *ServerConfig, options ...ServerOption) Server {
	dft := &defaultServer{
		conf: conf,
	}
	for _, opt := range options {
		opt(dft)
	}

	return dft
}

func (s *defaultServer) AddServerOption(options ...grpc.ServerOption) {
	s.options = append(s.options, options...)
}

func (s *defaultServer) AddUnaryServerInterceptor(interceptors ...grpc.UnaryServerInterceptor) {
	s.options = append(s.options, grpc.ChainUnaryInterceptor(interceptors...))
}

func (s *defaultServer) AddStreamServerInterceptor(interceptors ...grpc.StreamServerInterceptor) {
	s.options = append(s.options, grpc.ChainStreamInterceptor(interceptors...))
}

func (s *defaultServer) Serve(register Register) error {
	authorization := server.NewAuthorization(s.conf.Auth)
	options := []grpc.ServerOption{
		server.UnaryCrashHandler(),
		server.UnaryAuthorization(authorization),
		server.UnaryMetric(),

		server.StreamCrashHandler(),
		server.StreamAuthorization(authorization),
		server.StreamMetric(),
	}

	options = append(options, s.options...)
	serv := grpc.NewServer(options...)
	register(serv)
	lis, err := net.Listen(netWorkTcp, s.conf.ListenOn)
	if err != nil {
		return err
	}

	return serv.Serve(lis)
}
