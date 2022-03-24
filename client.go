package main

import (
	"context"
	"log"

	pb "fake.com/GoRPCApi/protobuf"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":5000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to port 5000: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	log.Printf("Creating a new user")
	user := &pb.User{Id: 1231, Email: "somedude@gmail.com", Password: "sfsadf"}
	response, err := client.CreateUser(context.Background(), &pb.CreateUserRequest{User: user})
	if err != nil {
		log.Fatalf("Error while creating user: %v", err)
	}
	log.Printf("Response: %v", response)
}
