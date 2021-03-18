package rpcx

import (
	"log"
	"time"

	"github.com/anqiansong/tools/rpcx/client"
	"google.golang.org/grpc"
)

type ClientConfig struct {
	EndPoint string
	Token    string
	TimeOut  int
}

type Client interface {
	grpc.ClientConnInterface
}

type ClientOption func(client Client)

type defaultClient struct {
	conf    *ClientConfig
	options []grpc.DialOption
	*grpc.ClientConn
}

func NewClient(conf *ClientConfig, options ...ClientOption) Client {
	dlt := &defaultClient{
		conf: conf,
	}

	for _, opt := range options {
		opt(dlt)
	}

	dlt.buildConn()

	return dlt
}

func (c *defaultClient) buildConn() {
	auth := client.NewAuthorization(c.conf.Token)
	options := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),

		client.UnaryCrashHandler(),
		client.UnaryAuthorization(auth),
		client.UnaryMetric(),

		client.StreamCrashHandler(),
		client.StreamAuthorization(auth),
		client.StreamMetric(),
	}

	if c.conf.TimeOut > 0 {
		options = append(options, client.StreamTimeout(time.Duration(c.conf.TimeOut)*time.Millisecond), client.UnaryTimeout(time.Duration(c.conf.TimeOut)*time.Millisecond))
	}

	options = append(options, c.options...)
	conn, err := grpc.Dial(c.conf.EndPoint, options...)
	if err != nil {
		log.Fatal(err)
	}

	c.ClientConn = conn
}
