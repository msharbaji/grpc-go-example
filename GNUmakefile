default: genProto

.PHONY: genProto
genProto:
	@protoc \
     -I=./api/proto \
     --go_out=./pkg \
     --go_opt module=${module} \
     --go-grpc_out=./pkg \
     --go-grpc_opt module=${module} \
     ./api/proto/*.proto

.PHONY: run-server
run-server:
	@go run ./cmd/server/main.go
