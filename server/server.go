package server

import (
	"context"

	database "fake.com/GoRPCApi/database"
	"fake.com/GoRPCApi/errhelp"
	pb "fake.com/GoRPCApi/protobuf"
	"gorm.io/gorm"
)

var db *gorm.DB

const (
	userTableInDb   = "enduser"
	whatIsUserId    = "id=?"
	whatIsUserEmail = "email=?"
)

func init() {
	db = database.InitDB()
}

type DbUser struct {
	Id       int32
	Email    string
	Password string
}

type Server struct {
	pb.UnimplementedUserServiceServer
}

func (s *Server) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	userInProtoForm := request.GetUser()
	userInDbForm := convertToUserInDbForm(*userInProtoForm)
	err := addUserToDb(userInDbForm)
	return &pb.CreateUserResponse{User: userInProtoForm}, err
}

func (s *Server) DeleteUser(ctx context.Context, request *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	emailOfUserToDelete := request.GetEmail()
	deletedUser, err := deleteUserFromDb(emailOfUserToDelete)
	userInProtoForm := convertToUserInProtoForm(deletedUser)
	return &pb.DeleteUserResponse{User: userInProtoForm}, err
}

func (s *Server) FetchUser(ctx context.Context, request *pb.FetchUserRequest) (*pb.FetchUserResponse, error) {
	emailOfUserToFetch := request.GetEmail()
	userInDbForm, err := findUserInDbByEmail(emailOfUserToFetch)
	userInProtoForm := convertToUserInProtoForm(userInDbForm)
	return &pb.FetchUserResponse{User: userInProtoForm}, err
}

func (s *Server) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	userInProtoForm := request.GetUser()
	userInDbForm := convertToUserInDbForm(*userInProtoForm)
	err := updateUserInDb(userInDbForm)
	return &pb.UpdateUserResponse{User: userInProtoForm}, err
}

func addUserToDb(userToAdd DbUser) error {
	resultFromAdding := db.Table(userTableInDb).Create(&userToAdd)
	return resultFromAdding.Error
}

func deleteUserFromDb(emailOfUserToDelete string) (DbUser, error) {
	userToDelete, userFindErr := findUserInDbByEmail(emailOfUserToDelete)
	if errhelp.ErrorExists(userFindErr) {
		return DbUser{}, errhelp.ErrUnableToDeleteUser
	} else {
		db.Table(userTableInDb).Where(whatIsUserEmail, emailOfUserToDelete).Delete(&DbUser{})
		return userToDelete, nil
	}
}

func updateUserInDb(user DbUser) error {
	_, userFindErr := findUserInDbById(user.Id)
	if errhelp.ErrorExists(userFindErr) {
		return errhelp.ErrUnableToUpdateUser
	} else {
		db.Table(userTableInDb).Where(whatIsUserId, user.Id).Updates(user)
		return nil
	}
}

func findUserInDbByEmail(userEmail string) (DbUser, error) {
	foundUser := DbUser{}
	resultOfFind := db.Table(userTableInDb).Where(whatIsUserEmail, userEmail).Take(&foundUser)
	return foundUser, resultOfFind.Error
}

func findUserInDbById(userId int32) (DbUser, error) {
	foundUser := DbUser{}
	resultOfFind := db.Table(userTableInDb).Where(whatIsUserId, userId).Take(&foundUser)
	return foundUser, resultOfFind.Error
}

func convertToUserInDbForm(userInProtoForm pb.User) DbUser {
	return DbUser{Id: userInProtoForm.GetId(), Email: userInProtoForm.GetEmail(), Password: userInProtoForm.GetPassword()}
}

func convertToUserInProtoForm(dbUser DbUser) *pb.User {
	return &pb.User{Id: dbUser.Id, Email: dbUser.Email, Password: dbUser.Password}
}
