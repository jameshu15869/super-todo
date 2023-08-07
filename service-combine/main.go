package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	pb "supertodo/combine/pb"
	"supertodo/combine/pkg/combineservice"
	"supertodo/combine/pkg/constants"
	"supertodo/combine/pkg/dbclient"

	"google.golang.org/grpc"
)

var (
	port = 4003
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
	pb.RegisterCombineServiceServer(grpcServer, combineservice.NewServer())
	grpcServer.Serve(lis)
}
