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

type User struct {
	Id       int32
	Email    string
	Password string
}

func convertToDbUser(protobufUser pb.User) User {
	return User{Id: protobufUser.GetId(), Email: protobufUser.GetEmail(), Password: protobufUser.GetPassword()}
}

func convertToProtobufUser(dbUser User) *pb.User {
	return &pb.User{Id: dbUser.Id, Email: dbUser.Email, Password: dbUser.Password}
}

type Server struct {
	pb.UnimplementedUserServiceServer
}

func (s *Server) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	protobufUser := request.GetUser()
	userToAdd := convertToDbUser(*protobufUser)
	err := addUserToDb(userToAdd)
	return &pb.CreateUserResponse{User: protobufUser}, err
}

func addUserToDb(userToAdd User) error {
	resultFromAdding := db.Table("enduser").Create(&userToAdd)
	return resultFromAdding.Error
}

func (s *Server) FetchUser(ctx context.Context, request *pb.FetchUserRequest) (*pb.FetchUserResponse, error) {
	userEmail := request.GetEmail()
	userInDb, err := findUserInDb(userEmail)
	protobufUser := convertToProtobufUser(userInDb)
	return &pb.FetchUserResponse{User: protobufUser}, err
}

func findUserInDb(userEmail string) (User, error) {
	foundUser := User{}
	resultOfFind := db.Table("enduser").Where("email=?", userEmail).Take(&foundUser)
	return foundUser, resultOfFind.Error
}
