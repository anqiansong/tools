package rpcx

import (
	"net"
	"time"

	"github.com/anqiansong/tools/rpcx/server"

	"google.golang.org/grpc"
)

const netWorkTcp = "tcp"

type ServerConfig struct {
	ListenOn string
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
	options := []grpc.ServerOption{
		server.MetricServerOption,
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

// WithTimeOutServerOption returns a ServerOption type to control timeout
func WithTimeOutServerOption(timeout time.Duration) ServerOption {
	return func(s Server) {
		s.AddUnaryServerInterceptor(server.TimeOut(timeout))
	}
}

// WithDeadlineServerOption returns a ServerOption type to control deadline
func WithDeadlineServerOption(at time.Time) ServerOption {
	return func(s Server) {
		s.AddUnaryServerInterceptor(server.Deadline(at))
	}
}
