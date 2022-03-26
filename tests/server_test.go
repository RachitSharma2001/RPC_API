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
	Describe("Given that we want to add a user", func() {

		Context("if the user is not already in the db", func() {
			id := createRandomId()
			email := fmt.Sprintf("Some%d@gmail.com", id)
			user := &pb.User{Id: id, Password: "pass", Email: email}
			_, errFromAdd := server.CreateUser(context.Background(), &pb.CreateUserRequest{User: user})
			It("the user should be successfully added", func() {
				Expect(errFromAdd).NotTo(HaveOccurred())
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
})

func createRandomId() int32 {
	source := rand.NewSource(time.Now().UnixNano())
	return int32(rand.New(source).Intn(100000))
}
