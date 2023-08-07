package userservice

import (
	"context"
	"log"
	pb "supertodo/user/pb"
	"supertodo/user/pkg/dbclient"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
}

func NewServer() *UserServer {
	s := &UserServer{}
	return s
}

func (s *UserServer) PingHealth(ctx context.Context, empty *emptypb.Empty) (*pb.Health, error) {
	return &pb.Health{Status: "healthy"}, nil
}

func (s *UserServer) GetAllUsers(ctx context.Context, empty *emptypb.Empty) (*pb.Users, error) {
	start := time.Now()
	users, err := dbclient.QueryAllUsers()
	if err != nil {
		log.Printf("Error occurred in GetAllUsers: %v\n", err)
		return nil, err
	}
	log.Println("GetAllUsers() - ", time.Since(start))
	return &pb.Users{Users: users}, nil
}

func (s *UserServer) GetUserById(ctx context.Context, in *pb.SingleUserIDRequest) (*pb.User, error) {
	start := time.Now()
	user, err := dbclient.QueryUserById(in.UserId)
	if err != nil {
		log.Printf("Error occurred in GetUserById: %v\n", err)
		return nil, err
	}
	log.Println("GetUserById() - ", time.Since(start))
	return user, nil
}

func (s *UserServer) GetUsersByMultipleIds(ctx context.Context, in *pb.MultipleUserIDRequest) (*pb.Users, error) {
	start := time.Now()
	users, err := dbclient.QueryUsersByMultipleIds(in.UserIds)
	if err != nil {
		log.Printf("Error occurred in GetUsersByMultipleIds: %v\n", err)
		return nil, err
	}
	log.Println("GetUsersByMultipleIds() - ", time.Since(start))
	return &pb.Users{Users: users}, nil
}

func (s *UserServer) AddUser(ctx context.Context, in *pb.User) (*pb.User, error) {
	start := time.Now()
	insertedUser, err := dbclient.AddUser(in)
	if err != nil {
		log.Printf("Error in AddUser: %v\n", err)
		return nil, err
	}
	log.Println("AddUser() - ", time.Since(start))
	return insertedUser, nil
}

func (s *UserServer) DeleteUser(ctx context.Context, in *pb.SingleUserIDRequest) (*pb.User, error) {
	start := time.Now()
	deletedUser, err := dbclient.DeleteUser(in.UserId)
	if err != nil {
		log.Printf("Error in DeleteUser: %v\n", err)
		return nil, err
	}
	log.Println("DeleteUser() - ", time.Since(start))
	return deletedUser, nil
}

func (s *UserServer) UpdateUserById(ctx context.Context, in *pb.User) (*pb.User, error) {
	start := time.Now()
	updatedUser, err := dbclient.UpdateUser(in)
	if err != nil {
		log.Printf("Error in UpdateUser: %v\n", err)
		return nil, err
	}
	log.Println("UpdateUser() - ", time.Since(start))
	return updatedUser, nil
}
