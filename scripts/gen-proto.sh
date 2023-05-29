# protoc gen

module="github.com/msharbaji/grpc-go-example"

protoc \
-I=./api/proto \
--go_out=./api \
--go_opt module=${module} \
--go-grpc_out=./api \
--go-grpc_opt module=${module} \
./api/proto/*.proto
