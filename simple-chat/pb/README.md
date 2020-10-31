# Protocol Buffer for Simple Chat App
Protocol Buffer for Simple Chat app.  
gRPC Server and Client for Go language are pre-built.

# Build package
Run the following command to build gRPC Server and Client for Go language.

```shell
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative simple-chat.proto 
```