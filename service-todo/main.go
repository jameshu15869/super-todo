package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	pb "supertodo/todo/pb"
	"supertodo/todo/pkg/constants"
	"supertodo/todo/pkg/dbclient"
	"supertodo/todo/pkg/todoservice"

	"google.golang.org/grpc"
)

var (
	port = 4002
)

func main() {
	constants.InitEnvConstants()
	dbclient.InitDB()

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Listening on 0.0.0.0:%d\n", port)

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterTodoServiceServer(grpcServer, todoservice.NewServer())
	grpcServer.Serve(lis)
}
