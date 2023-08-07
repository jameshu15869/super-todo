#!/bin/bash

echo 'Generating Protobuf...'

echo 'Generating Gateway: protoc -I protobuf --go_out=service-gateway/pb --go_opt=paths=source_relative --go-grpc_out=service-gateway/pb --go-grpc_opt=paths=source_relative super-todo.proto'
protoc -I protobuf --go_out=service-gateway/pb --go_opt=paths=source_relative --go-grpc_out=service-gateway/pb --go-grpc_opt=paths=source_relative super-todo.proto
echo 'Generating User: protoc -I protobuf --go_out=service-user/pb --go_opt=paths=source_relative --go-grpc_out=service-user/pb --go-grpc_opt=paths=source_relative super-todo.proto' 
protoc -I protobuf --go_out=service-user/pb --go_opt=paths=source_relative --go-grpc_out=service-user/pb --go-grpc_opt=paths=source_relative super-todo.proto
echo 'Generating Todo: protoc -I protobuf --go_out=service-todo/pb --go_opt=paths=source_relative --go-grpc_out=service-todo/pb --go-grpc_opt=paths=source_relative super-todo.proto'
protoc -I protobuf --go_out=service-todo/pb --go_opt=paths=source_relative --go-grpc_out=service-todo/pb --go-grpc_opt=paths=source_relative super-todo.proto
echo 'Generating Combine: protoc -I protobuf --go_out=service-combine/pb --go_opt=paths=source_relative --go-grpc_out=service-combine/pb --go-grpc_opt=paths=source_relative super-todo.proto'
protoc -I protobuf --go_out=service-combine/pb --go_opt=paths=source_relative --go-grpc_out=service-combine/pb --go-grpc_opt=paths=source_relative super-todo.proto
echo 'Generating Client: protoc -I=protobuf --ts_out=client/src/pb super-todo.proto'
protoc -I=protobuf --ts_out=client/pb super-todo.proto
