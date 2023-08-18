package bookClient

import (
	"context"
	"time"

	proto_book "github.com/dwnGnL/testWork/proto/book"
	"github.com/pkg/errors"
	"go.elastic.co/apm/module/apmgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	bookClient proto_book.BookServiceClient
}

type Config struct {
	Host              string
	ConnectionTimeout time.Duration
}

func New(cfg *Config) (*Client, error) {
	grpcOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			apmgrpc.NewUnaryClientInterceptor(),
		),
		grpc.WithChainStreamInterceptor(
			apmgrpc.NewStreamClientInterceptor(),
		),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: cfg.ConnectionTimeout,
		}),
	}

	grpcClientConn, err := grpc.Dial(cfg.Host, grpcOpts...)
	if err != nil {
		return nil, errors.WithMessage(err, "grpc.Dial error")
	}
	return &Client{
		bookClient: proto_book.NewBookServiceClient(grpcClientConn),
	}, nil
}

func (c Client) GetBooks(ctx context.Context, name string) (*proto_book.BookList, error) {
	return c.bookClient.GetBooks(ctx, &proto_book.BookFilter{Name: name})
}
