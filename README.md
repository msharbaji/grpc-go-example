# grpc-go-example
This a grpc show the example for grpc-go

## How to run
```shell
# run server
go run server/main.go
```

## How to generate code
```shell
# generate code
protoc --go_out=plugins=grpc:. ./proto/helloworld.proto
```
