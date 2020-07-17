package main

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/golang/glog"
	gwr "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	pb "github.com/vitalik-mironov/go-labs/pkg/echo_api"
)

type echoService struct{}

var (
	endpoint = ":8081"
)

func (s echoService) EchoV1(ctx context.Context, req *pb.EchoV1Request) (resp *pb.EchoV1Response, err error) {
	return &pb.EchoV1Response{Text: req.Text}, nil
}

func main() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		io.Copy(w, strings.NewReader(pb.Swagger))
	})

	var grpcEchoService echoService

	gwmux := gwr.NewServeMux()
	err := pb.RegisterEchoAPIHandlerServer(ctx, gwmux, grpcEchoService)
	if err != nil {
		glog.Fatal(err)
	}

	mux.Handle("/", gwmux)

	httpServer := &http.Server{
		Addr:    endpoint,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		glog.Infof("Shutting down the http gateway server")
		if err := httpServer.Shutdown(context.Background()); err != nil {
			glog.Errorf("Failed to shutdown http gateway server: %v", err)
		}
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		glog.Errorf("Failed to listen and serve: %v", err)
	}

	// server := grpc.NewServer()
	// pb.RegisterEchoAPIServer(server, *new(echoService))

}
