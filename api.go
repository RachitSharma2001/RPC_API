package main

import (
	"context"
	"errors"
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
	userInDb, err := findUserInDbByEmail(userEmail)
	protobufUser := convertToProtobufUser(userInDb)
	return &pb.FetchUserResponse{User: protobufUser}, err
}

func (s *Server) DeleteUser(ctx context.Context, request *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	emailOfUserToDelete := request.GetEmail()
	deletedUser, err := deleteFromDb(emailOfUserToDelete)
	protobufUser := convertToProtobufUser(deletedUser)
	return &pb.DeleteUserResponse{User: protobufUser}, err
}

func deleteFromDb(emailOfUserToDelete string) (User, error) {
	userToDelete, userFindErr := findUserInDbByEmail(emailOfUserToDelete)
	if errorExists(userFindErr) {
		return User{}, errors.New("User not found")
	} else {
		db.Table("enduser").Where("email=?", emailOfUserToDelete).Delete(&User{})
		return userToDelete, nil
	}
}

func (s *Server) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user := convertToDbUser(*request.GetUser())
	err := updateInDb(user)
	return &pb.UpdateUserResponse{User: request.GetUser()}, err
}

func updateInDb(user User) error {
	_, userFindErr := findUserInDbById(user.Id)
	if errorExists(userFindErr) {
		return errors.New("User not found")
	} else {
		db.Table("enduser").Where("id=?", user.Id).Updates(user)
		return nil
	}
}

func findUserInDbByEmail(userEmail string) (User, error) {
	foundUser := User{}
	resultOfFind := db.Table("enduser").Where("email=?", userEmail).Take(&foundUser)
	return foundUser, resultOfFind.Error
}

func findUserInDbById(userId int32) (User, error) {
	foundUser := User{}
	resultOfFind := db.Table("enduser").Where("id=?", userId).Take(&foundUser)
	return foundUser, resultOfFind.Error
}

func convertToDbUser(protobufUser pb.User) User {
	return User{Id: protobufUser.GetId(), Email: protobufUser.GetEmail(), Password: protobufUser.GetPassword()}
}

func convertToProtobufUser(dbUser User) *pb.User {
	return &pb.User{Id: dbUser.Id, Email: dbUser.Email, Password: dbUser.Password}
}
