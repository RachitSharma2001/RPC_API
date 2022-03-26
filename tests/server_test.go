package main

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	pb "fake.com/GoRPCApi/protobuf"
	. "fake.com/GoRPCApi/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tests Suite")
}

var _ = Describe("Given the client", func() {
	server := Server{}

	Describe("Given that we want to delete a user", func() {
		Context("if the user exists in the db", func() {
			userToDelete := &pb.User{Id: 99645, Password: "randompass", Email: "ronald99645@gmail.com"}
			response, errFromDelete := server.DeleteUser(context.Background(), &pb.DeleteUserRequest{Email: userToDelete.Email})
			It("no error should be returned", func() {
				Expect(errFromDelete).NotTo(HaveOccurred())
			})
			It("Response should contain the correct user", func() {
				Expect(response.GetUser()).To(Equal(userToDelete))
			})
		})

		Context("if the user does not exist in the db", func() {
			_, errFromDelete := server.DeleteUser(context.Background(), &pb.DeleteUserRequest{Email: "nobodywhoishere@nobody"})
			It("an error should be returned", func() {
				Expect(errFromDelete).To(HaveOccurred())
			})
		})
	})

	Describe("Given that we want to fetch a user", func() {
		Context("if the user exists in the db", func() {
			userEmail := "james@gmail.com"
			expectedUserId := int32(5)
			response, errFromFetch := server.FetchUser(context.Background(), &pb.FetchUserRequest{Email: userEmail})
			It("no error should be returned", func() {
				Expect(errFromFetch).NotTo(HaveOccurred())
			})
			It("the response should contain the correct user id", func() {
				observedUserId := response.GetUser().Id
				Expect(observedUserId).To(Equal(expectedUserId))
			})
		})

		Context("if user does not exist in the db", func() {
			userEmail := "nobodywhoexists@nothing.com"
			_, errFromFetch := server.FetchUser(context.Background(), &pb.FetchUserRequest{Email: userEmail})
			It("an error should have occurred", func() {
				Expect(errFromFetch).To(HaveOccurred())
			})
		})
	})

	Describe("Given that we want to update a user", func() {
		Context("if the user exists in the db", func() {
			randomNum := createRandomNum()
			newPassword := fmt.Sprintf("somethingnew%d", randomNum)
			user := &pb.User{Id: 5, Password: newPassword, Email: "james@gmail.com"}
			response, errFromUpdate := server.UpdateUser(context.Background(), &pb.UpdateUserRequest{User: user})
			It("no error should be returned", func() {
				Expect(errFromUpdate).NotTo(HaveOccurred())
			})
			It("response should contain correctly updated user", func() {
				Expect(response.GetUser()).To(Equal(user))
			})
		})

		Context("if the user does not exit in the db", func() {
			nonExistentUserId := -10
			user := &pb.User{Id: int32(nonExistentUserId), Password: "random", Email: "adfsf@gmail.com"}
			_, errFromUpdate := server.UpdateUser(context.Background(), &pb.UpdateUserRequest{User: user})
			It("an error should be returned", func() {
				Expect(errFromUpdate).To(HaveOccurred())
			})
		})
	})

	Describe("Given that we want to add a user", func() {

		Context("if the user is not already in the db", func() {
			user := &pb.User{Id: 99645, Password: "randompass", Email: "ronald99645@gmail.com"}
			response, errFromAdd := server.CreateUser(context.Background(), &pb.CreateUserRequest{User: user})
			It("no error should be returned", func() {
				Expect(errFromAdd).NotTo(HaveOccurred())
			})
			It("the response should contain the user", func() {
				Expect(response.GetUser()).To(Equal(user))
			})
		})

		Context("if the user is already in the db", func() {
			user := &pb.User{Id: 5, Password: "somethingnew", Email: "james@gmail.com"}
			_, errFromAdd := server.CreateUser(context.Background(), &pb.CreateUserRequest{User: user})
			It("an error should be returned", func() {
				Expect(errFromAdd).To(HaveOccurred())
			})
		})
	})
})

func createRandomNum() int32 {
	source := rand.NewSource(time.Now().UnixNano())
	return int32(rand.New(source).Intn(100000))
}
