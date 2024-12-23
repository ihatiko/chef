LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

compile-proto:
	protoc -I . \
        --plugin=protoc-gen-go=./bin/protoc-gen-go --go_out . --go_opt paths=source_relative \
        --plugin=protoc-gen-go-grpc=./bin/protoc-gen-go-grpc --go-grpc_out . --go-grpc_opt paths=source_relative --go-grpc_opt=require_unimplemented_servers=false \
        ./pkg/protoc/*/*.proto

vendor-proto:
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi