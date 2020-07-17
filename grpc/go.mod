module github.com/vitalik-mironov/go-labs

go 1.14

replace github.com/valyala/fasthttp => git.wildberries.ru/mironov.vitaliy3/fasthttp v1.14.1

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/grpc-gateway v1.14.6
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.0.0-beta.3
	github.com/valyala/fasthttp v1.14.0
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.25.0
)
