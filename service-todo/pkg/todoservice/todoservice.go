package todoservice

import (
	"context"
	"log"
	"supertodo/todo/pb"
	"supertodo/todo/pkg/dbclient"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
)

type TodoServer struct {
	pb.UnimplementedTodoServiceServer
}

func NewServer() *TodoServer {
	s := &TodoServer{}
	return s
}

func (s *TodoServer) PingHealth(ctx context.Context, empty *emptypb.Empty) (*pb.Health, error) {
	return &pb.Health{Status: "healthy"}, nil
}

func (s *TodoServer) GetAllTodos(ctx context.Context, empty *emptypb.Empty) (*pb.Todos, error) {
	start := time.Now()
	todos, err := dbclient.QueryAllTodos()
	if err != nil {
		log.Printf("Error in GetAllTodos: %v\n", err)
		return nil, err
	}
	log.Println("GetAllTodos() - ", time.Since(start))
	return &pb.Todos{Todos: todos}, nil
}

func (s *TodoServer) GetTodoById(ctx context.Context, in *pb.SingleTodoIDRequest) (*pb.Todo, error) {
	start := time.Now()
	todo, err := dbclient.QueryTodoById(in.TodoId)
	if err != nil {
		log.Printf("Error in GetTodoById: %v\n", err)
		return nil, err
	}
	log.Println("GetTodoById() - ", time.Since(start))
	return todo, nil
}

func (s *TodoServer) GetTodosByMultipleIds(ctx context.Context, in *pb.MultipleTodoIDRequest) (*pb.Todos, error) {
	start := time.Now()
	todos, err := dbclient.QueryTodosByMultipleIds(in.TodoIds)
	if err != nil {
		log.Printf("Error in GetTodosByMultipleId: %v\n", err)
		return nil, err
	}
	log.Println("GetTodosByMultipleId() - ", time.Since(start))
	return &pb.Todos{Todos: todos}, nil
}

func (s *TodoServer) AddTodo(ctx context.Context, in *pb.Todo) (*pb.Todo, error) {
	start := time.Now()
	todo, err := dbclient.AddTodo(in)
	if err != nil {
		log.Printf("Error in AddTodo: %v\n", err)
		return nil, err
	}
	log.Println("AddTodo() - ", time.Since(start))
	return todo, nil
}

func (s *TodoServer) UpdateTodoById(ctx context.Context, in *pb.Todo) (*pb.Todo, error) {
	start := time.Now()
	todo, err := dbclient.UpdateTodoById(in)
	if err != nil {
		log.Printf("Error in UpdateTodoById: %v\n", err)
		return nil, err
	}
	log.Println("UpdateTodoById() - ", time.Since(start))
	return todo, nil
}

func (s *TodoServer) DeleteTodo(ctx context.Context, in *pb.SingleTodoIDRequest) (*pb.Todo, error) {
	start := time.Now()
	todo, err := dbclient.DeleteTodo(in.TodoId)
	if err != nil {
		log.Printf("Error in DeleteTodo: %v\n", err)
		return nil, err
	}
	log.Println("DeleteTodo() - ", time.Since(start))
	return todo, nil
}
