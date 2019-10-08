package client

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"time"
)

type RpcConnection struct {
	GrpcCon *grpc.ClientConn
}

type Optional struct {
	GrpcServerAdd string
}

var GrpcConn RpcConnection

func GrpcClientConn(opts *Optional) {
	var err error

	grpcOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100 * time.Millisecond)),
		grpc_retry.WithCodes(codes.NotFound, codes.Aborted, codes.ResourceExhausted, codes.Unavailable),
	}

	if len(opts.GrpcServerAdd) > 0 {
		GrpcConn.GrpcCon, err = grpc.Dial(opts.GrpcServerAdd, grpc.WithInsecure(),
			grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
				grpc_prometheus.StreamClientInterceptor,
				grpc_retry.StreamClientInterceptor(grpcOpts...),
			)),
			grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
				grpc_prometheus.UnaryClientInterceptor,
				grpc_retry.UnaryClientInterceptor(grpcOpts...),
			)))
		if err != nil {
			fmt.Println("Unable to initialize connection to IMS GRPC")
		}
	}
}
