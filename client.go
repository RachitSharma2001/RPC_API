package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

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
	randomId := createRandomId()
	user := &pb.User{Id: randomId, Email: fmt.Sprintf("somedude%d@gmail.com", randomId), Password: "sfsadf"}
	createUserResp, err := client.CreateUser(context.Background(), &pb.CreateUserRequest{User: user})
	if err != nil {
		log.Fatalf("Error while creating user: %v", err)
	}
	log.Printf("Response: %v", createUserResp)
	log.Println("----------------------------------------------")

	log.Printf("Fetching an existent user")
	fetchUserResp, err := client.FetchUser(context.Background(), &pb.FetchUserRequest{Email: "ronald@gmail.com"})
	if err != nil {
		log.Fatalf("Error while creating user: %v", err)
	}
	log.Printf("Response: %v", fetchUserResp)
	log.Println("----------------------------------------------")

	log.Printf("Fetching an non-existent user")
	fetchUserResp, err = client.FetchUser(context.Background(), &pb.FetchUserRequest{Email: "asfsdfs@gmail.com"})
	if err != nil {
		log.Fatalf("Error while creating user: %v", err)
	}
	log.Printf("Response: %v", fetchUserResp)
	log.Println("----------------------------------------------")
}

func createRandomId() int32 {
	source := rand.NewSource(time.Now().UnixNano())
	return int32(rand.New(source).Intn(100000))
}
