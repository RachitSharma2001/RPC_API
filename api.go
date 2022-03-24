package main

import (
	"context"
	"log"
	"net"

	pb "fake.com/GoRPCApi/protobuf"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	db = InitDB()
}

func main() {
	lis, err := net.Listen("tcp", ":5000")
	if errorExists(err) {
		log.Fatalf("Unable to listen at port 5000: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &Server{})
	err = grpcServer.Serve(lis)
	if errorExists(err) {
		log.Fatalf("Unable to listen serve: %v", err)
	}
}

type Server struct {
	pb.UnimplementedUserServiceServer
}

func (s *Server) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	givenUser := request.GetUser()
	return &pb.CreateUserResponse{User: givenUser}, nil
}
