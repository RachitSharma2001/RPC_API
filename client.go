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
	userEmail := fmt.Sprintf("somedude%d@gmail.com", randomId)
	user := &pb.User{Id: randomId, Email: userEmail, Password: "sfsadf"}
	createUserResp, err := client.CreateUser(context.Background(), &pb.CreateUserRequest{User: user})
	if err != nil {
		log.Fatalf("Error while creating user: %v", err)
	}
	log.Printf("Response: %v", createUserResp)
	log.Println("----------------------------------------------")

	log.Printf("Fetching an existent user")
	fetchUserResp, err := client.FetchUser(context.Background(), &pb.FetchUserRequest{Email: "ronald@gmail.com"})
	if err != nil {
		log.Fatalf("Error while fetching user: %v", err)
	}
	log.Printf("Response: %v", fetchUserResp)
	log.Println("----------------------------------------------")

	log.Printf("Fetching an non-existent user")
	fetchUserResp, err = client.FetchUser(context.Background(), &pb.FetchUserRequest{Email: "asfsdfs@gmail.com"})
	if err == nil {
		log.Fatalf("Recieved no err, instead got response: %v", fetchUserResp)
	}
	log.Printf("Correctly received error: %v", err)
	log.Println("----------------------------------------------")

	log.Printf("Updating an existent user")
	userToUpdate := pb.User{Id: 5, Email: "james@gmail.com", Password: "somethingnew"}
	updateUserResp, err := client.UpdateUser(context.Background(), &pb.UpdateUserRequest{User: &userToUpdate})
	if err != nil {
		log.Fatalf("Received error while updating: %v", err)
	}
	log.Printf("Response: %v", updateUserResp)
	log.Println("----------------------------------------------")

	log.Printf("Updating a non-existent user")
	nonExistentUser := pb.User{Id: 241241, Email: userEmail, Password: "somethingnew"}
	updateUserResp, err = client.UpdateUser(context.Background(), &pb.UpdateUserRequest{User: &nonExistentUser})
	if err == nil {
		log.Fatalf("Recieved no err, instead got response: %v", updateUserResp)
	}
	log.Printf("Correctly received error: %v", err)
	log.Println("----------------------------------------------")

	log.Printf("Deleting an existent user")
	deleteUserResp, err := client.DeleteUser(context.Background(), &pb.DeleteUserRequest{Email: userEmail})
	if err != nil {
		log.Fatalf("Received error while deleting: %v", err)
	}
	log.Printf("Response: %v", deleteUserResp)
	log.Println("----------------------------------------------")

	log.Printf("Deleting a non-existent user")
	deleteUserResp, err = client.DeleteUser(context.Background(), &pb.DeleteUserRequest{Email: "asfsdfs@gmail.com"})
	if err == nil {
		log.Fatalf("Recieved no err, instead got response: %v", fetchUserResp)
	}
	log.Printf("Correctly received error: %v", err)
	log.Println("----------------------------------------------")
}

func createRandomId() int32 {
	source := rand.NewSource(time.Now().UnixNano())
	return int32(rand.New(source).Intn(100000))
}
