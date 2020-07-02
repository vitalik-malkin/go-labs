

.PHONY: gen
gen:
	cd pkg/echo-api \
	&& protoc \
		-I../../api \
		--go_out=plugins=grpc:. \
		../../api/echo-api.proto \
	&& protoc \
		-I../../api \
		--grpc-gateway_out=logtostderr=true,paths=source_relative,grpc_api_configuration=../../api/echo-api.yaml:. \
		../../api/echo-api.proto \
	&& protoc \
		-I../../api \
		--swagger_out=logtostderr=true,grpc_api_configuration=../../api/echo-api.yaml:. \
		../../api/echo-api.proto