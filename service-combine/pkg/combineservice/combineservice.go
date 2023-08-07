package combineservice

import (
	"context"
	"log"
	"supertodo/combine/pb"
	"supertodo/combine/pkg/dbclient"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
)

type CombineServer struct {
	pb.UnimplementedCombineServiceServer
}

func NewServer() *CombineServer {
	s := &CombineServer{}
	return s
}

func (s *CombineServer) PingHealth(ctx context.Context, empty *emptypb.Empty) (*pb.Health, error) {
	return &pb.Health{Status: "healthy"}, nil
}

func (s *CombineServer) GetAllCombined(ctx context.Context, empty *emptypb.Empty) (*pb.MultipleCombined, error) {
	start := time.Now()
	allCombined, err := dbclient.QueryAllCombined()
	if err != nil {
		log.Printf("Error occurred in GetAllCombined: %v\n", err)
		return nil, err
	}
	log.Println("GetAllCombined() - ", time.Since(start))
	return &pb.MultipleCombined{Combines: allCombined}, nil
}

func (s *CombineServer) GetTodosFromUserId(ctx context.Context, in *pb.SingleUserIDRequest) (*pb.TodosFromUserIdResponse, error) {
	start := time.Now()
	todo_ids, err := dbclient.QueryTodosFromUserId(in.UserId)
	if err != nil {
		log.Printf("Error occurred in GetTodosFromUserId: %v\n", err)
		return nil, err
	}
	log.Println("GetTodosFromUserId() - ", time.Since(start))
	return &pb.TodosFromUserIdResponse{TodoIds: todo_ids}, nil
}

func (s *CombineServer) GetUsersFromTodoId(ctx context.Context, in *pb.SingleTodoIDRequest) (*pb.UsersFromTodoIdResponse, error) {
	start := time.Now()
	user_ids, err := dbclient.QueryUsersFromTodoId(in.TodoId)
	if err != nil {
		log.Printf("Error occurred in GetUsersFromTodoId: %v\n", err)
		return nil, err
	}
	log.Println("GetUsersFromTodoId() - ", time.Since(start))
	return &pb.UsersFromTodoIdResponse{UserIds: user_ids}, nil
}

func (s *CombineServer) AddCombined(ctx context.Context, in *pb.CombinedRequest) (*pb.Combined, error) {
	start := time.Now()
	insertedCombined, err := dbclient.AddCombined(in.UserId, in.TodoId)
	if err != nil {
		log.Printf("Error occurred in AddCombined: %v\n", err)
		return nil, err
	}
	log.Println("AddCombined() - ", time.Since(start))
	return insertedCombined, nil
}

func (s *CombineServer) AddCombinedArray(ctx context.Context, in *pb.CombinedArrayRequest) (*pb.MultipleCombined, error) {
	start := time.Now()
	insertedCombines, err := dbclient.AddCombinedArray(in.TodoId, in.UserIds)
	if err != nil {
		log.Printf("Error occurred in AddCombinedArray: %v\n", err)
		return nil, err
	}
	log.Println("AddCombinedArray() - ", time.Since(start))
	return &pb.MultipleCombined{Combines: insertedCombines}, nil
}

func (s *CombineServer) DeleteCombined(ctx context.Context, in *pb.CombinedRequest) (*pb.Combined, error) {
	start := time.Now()
	removedCombined, err := dbclient.DeleteCombined(in.UserId, in.TodoId)
	if err != nil {
		log.Printf("Error occurred in DeleteCombined: %v\n", err)
		return nil, err
	}
	log.Println("DeleteCombined() - ", time.Since(start))
	return removedCombined, nil
}

func (s *CombineServer) DeleteTodo(ctx context.Context, in *pb.SingleTodoIDRequest) (*pb.MultipleCombined, error) {
	start := time.Now()
	removedCombined, err := dbclient.DeleteTodo(in.TodoId)
	if err != nil {
		log.Printf("Error occurred in DeleteTodo: %v\n", err)
		return nil, err
	}
	log.Println("DeleteTodo() - ", time.Since(start))
	return &pb.MultipleCombined{Combines: removedCombined}, nil
}

func (s *CombineServer) DeleteUser(ctx context.Context, in *pb.SingleUserIDRequest) (*pb.MultipleCombined, error) {
	start := time.Now()
	removedCombined, err := dbclient.DeleteUser(in.UserId)
	if err != nil {
		log.Printf("Error occurred in DeleteUser: %v\n", err)
		return nil, err
	}
	log.Println("DeleteUser() - ", time.Since(start))
	return &pb.MultipleCombined{Combines: removedCombined}, nil
}

func (s *CombineServer) UpdateTodo(ctx context.Context, in *pb.CombinedArrayRequest) (*pb.MultipleCombined, error) {
	start := time.Now()
	_, err := dbclient.DeleteTodo(in.TodoId)
	if err != nil {
		log.Printf("Error occurred in UpdateTodo: %v\n", err)
		return nil, err
	}

	insertedCombined, err := dbclient.AddCombinedArray(in.TodoId, in.UserIds)
	if err != nil {
		log.Printf("Error occurred in UpdateTodo: %v\n", err)
		return nil, err
	}
	log.Println("UpdateUser() - ", time.Since(start))
	return &pb.MultipleCombined{Combines: insertedCombined}, nil
}
