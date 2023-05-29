# grpc-go-example
This a grpc show the example for grpc-go, include server and client.

## Environment variables
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
```shell
chmod +x ./scripts/gen-proto.sh
./scripts/gen-proto.sh
```

## Run server
```shell
go run ./cmd/server/main.go
```