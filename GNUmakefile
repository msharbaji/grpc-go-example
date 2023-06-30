default: genProto

.PHONY: genProto
genProto:
	@protoc \
     -I=./api/proto \
     --go_out=. \
     --go_opt=module=github.com/msharbaji/grpc-go-example \
     --go-grpc_out=. \
     --go-grpc_opt=module=github.com/msharbaji/grpc-go-example \
     ./api/proto/v1/*.proto

.PHONY: run-server
run-server:
	@go run ./cmd/server/main.go
