

.PHONY: gen
gen:
	cd pkg/echo_api \
	&& protoc \
		-I../../api \
		--go_out=plugins=grpc,paths=source_relative:. \
		../../api/echo-api.proto \
	&& protoc \
		-I../../api \
		--grpc-gateway_out=logtostderr=true,paths=source_relative,grpc_api_configuration=../../api/echo-api.yaml:. \
		../../api/echo-api.proto \
	&& protoc \
		-I../../api \
		--swagger_out=logtostderr=true,grpc_api_configuration=../../api/echo-api.yaml:. \
		../../api/echo-api.proto