package main

import (
	"context"

	gwr "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	pb "github.com/vitalik-mironov/go-labs/grpc/grpc-gateway/pkg/echo_api"
)

type echoService struct{}

var (
	grpcEndpoint = "localhost:8081"
)

func (s echoService) EchoV1(ctx context.Context, req *pb.EchoV1Request) (resp *pb.EchoV1Response, err error) {
	return &pb.EchoV1Response{}, nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	server := grpc.NewServer()
	pb.RegisterEchoAPIServer(server, *new(echoService))

	mux := gwr.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterYourServiceHandlerFromEndpoint(ctx, mux, server, opts)
	if err != nil {
		return err
	}
}
