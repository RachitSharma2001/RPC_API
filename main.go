package main

import (
	"net"

	"fake.com/GoRPCApi/errhelp"
	pb "fake.com/GoRPCApi/protobuf"
	. "fake.com/GoRPCApi/server"
	"google.golang.org/grpc"
)

func main() {
	listener := listenAtPort(":5000")
	server := registerServer()
	connectServerToListener(listener, server)
}

func listenAtPort(port string) net.Listener {
	lis, err := net.Listen("tcp", port)
	if errhelp.ErrorExists(err) {
		errhelp.ThrowPortListenErr(err)
	}
	return lis
}

func registerServer() *grpc.Server {
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &Server{})
	return grpcServer
}

func connectServerToListener(listener net.Listener, server *grpc.Server) {
	err := server.Serve(listener)
	if errhelp.ErrorExists(err) {
		errhelp.ThrowServeErr(err)
	}
}
