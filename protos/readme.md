# protos for farseer 
This code is a modification of the original protos from `farcasterxyz/hub-monorepo`.

```
# How to compile the protos normally
protoc -I=protos --go_out=protos protos/*.proto
```

```
# How to compile the protos for GRPC
protoc -I=protos --go-grpc_opt=paths=source_relative --go_out=protos --go-grpc_out=protos protos/rpc.proto
```