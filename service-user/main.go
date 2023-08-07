package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"supertodo/user/pb"
	"supertodo/user/pkg/constants"
	"supertodo/user/pkg/dbclient"
	"supertodo/user/pkg/userservice"

	"google.golang.org/grpc"
)

var (
	port = 4001
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
	pb.RegisterUserServiceServer(grpcServer, userservice.NewServer())
	grpcServer.Serve(lis)
}
