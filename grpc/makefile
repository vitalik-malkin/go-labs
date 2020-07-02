PROJECT_STRUCTURE_CLIENTS=pkg
PROJECT_STRUCTURE_GO_GRPC_API_CLIENT=$(PROJECT_STRUCTURE_CLIENTS)/go-grpc-api
IMPLEMENTATION_TYPE_NAME="Implementation"

.PHONY: generate
generate:
	cd $(PROJECT_STRUCTURE_GO_GRPC_API_CLIENT) \
	 && protoc --plugin=protoc-gen-go \
		-I../../api/ \
		--go_out=plugins=grpc:. \
		../../api/go-grpc-api.proto
