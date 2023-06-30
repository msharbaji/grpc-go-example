# grpc-go-example
This is an example repository showcasing the usage of gRPC with Go, including server and client implementations.

## Environment variables
The following environment variables can be used to configure the application:

| Name          | Description          | Default         | Required |
|---------------|----------------------|-----------------|----------|
| GRPC_ENDPOINT | grpc server endpoint | localhost:50051 | false    |
| KEY_ID        | hmac key id          | 1               | false    |
| SECRET_KEY    | hmac secret key      | 123456          | false    |


## Set environment variables
```shell
export GRPC_ENDPOINT=localhost:50051
export KEY_ID=my-key-id
export SECRET_KEY=my-secret-key
```

## Generate proto code
You can generate the protobuf code using one of the following methods:

### Using go generate
```shell
go generate ./...
```
### Using protoc
```shell
protoc \
     -I=./api/proto \
     --go_out=. \
     --go_opt=module=github.com/msharbaji/grpc-go-example \
     --go-grpc_out=. \
     --go-grpc_opt=module=github.com/msharbaji/grpc-go-example \
     ./api/proto/v1/*.proto
```

### Using `buf.build`
```shell
buf generate
```


## Run server
To run the server, execute the following command:

```shell
make run-server
```