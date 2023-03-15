```shell
export PATH="$PATH:$(go env GOPATH)/bin"
protoc --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. auth.proto
```

```shell
swag init -g **/**/*.go
```