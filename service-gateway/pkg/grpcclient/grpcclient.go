package grpcclient

import (
	"context"
	"log"
	pb "supertodo/gateway/pb"
	"supertodo/gateway/pkg/constants"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var userServiceClient pb.UserServiceClient
var todoServiceClient pb.TodoServiceClient
var combineServiceClient pb.CombineServiceClient

func InitGrpc() {
	var opts []grpc.DialOption = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	initUserServiceClient(opts)
	initTodoServiceClient(opts)
	initCombineServiceClient(opts)
}

func initUserServiceClient(opts []grpc.DialOption) {
	conn, err := grpc.Dial(constants.EnvValues.USER_ENDPOINT, opts...)
	if err != nil {
		log.Printf("Failed to dial user service: %v\n", err)
	}

	userServiceClient = pb.NewUserServiceClient(conn)
}

func initTodoServiceClient(opts []grpc.DialOption) {
	conn, err := grpc.Dial(constants.EnvValues.TODO_ENDPOINT, opts...)
	if err != nil {
		log.Printf("Failed to dial todo service: %v\n", err)
	}

	todoServiceClient = pb.NewTodoServiceClient(conn)
}

func initCombineServiceClient(opts []grpc.DialOption) {
	conn, err := grpc.Dial(constants.EnvValues.COMBINE_ENDPOINT, opts...)
	if err != nil {
		log.Printf("Failed to dial combine service: %v\n", err)
	}

	combineServiceClient = pb.NewCombineServiceClient(conn)
}

func GetUserHealth() (*pb.Health, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	healthStatus, err := userServiceClient.PingHealth(ctx, &emptypb.Empty{})

	if err != nil {
		log.Printf("userServiceClient.PingHealth Failed : %v\n", err)
		return nil, err
	}
	return healthStatus, nil
}

func GetTodoHealth() (*pb.Health, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	healthStatus, err := todoServiceClient.PingHealth(ctx, &emptypb.Empty{})

	if err != nil {
		log.Printf("todoServiceClient.PingHealth Failed : %v\n", err)
		return nil, err
	}
	return healthStatus, nil
}

func GetAllCombined() ([]*pb.Combined, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	allCombined, err := combineServiceClient.GetAllCombined(ctx, &emptypb.Empty{})

	if err != nil {
		log.Printf("combineServiceClient.GetAllCombined Failed : %v\n", err)
		return nil, err
	}
	return allCombined.Combines, nil
}

func GetCombineHealth() (*pb.Health, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	healthStatus, err := combineServiceClient.PingHealth(ctx, &emptypb.Empty{})

	if err != nil {
		log.Printf("combineServiceClient.PingHealth Failed : %v\n", err)
		return nil, err
	}
	return healthStatus, nil
}

func GetUsers() ([]*pb.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT*time.Second)
	defer cancel()
	users, err := userServiceClient.GetAllUsers(ctx, &emptypb.Empty{})

	if err != nil {
		log.Printf("userServiceClient.GetAllUsers( Failed: %v\n", err)
		return nil, err
	}
	return users.Users, nil
}

func goGetAllCombined(ctx context.Context, allCombinedChan chan []*pb.Combined) {
	allCombined, err := combineServiceClient.GetAllCombined(ctx, &emptypb.Empty{})

	if err != nil {
		log.Printf("combinedServiceClient.GetAllCombined( Failed: %v\n", err)
		allCombinedChan <- []*pb.Combined{}
		return
	}
	allCombinedChan <- allCombined.Combines
}

func goGetAllUsers(ctx context.Context, usersChan chan []*pb.User) {
	users, err := userServiceClient.GetAllUsers(ctx, &emptypb.Empty{})

	if err != nil {
		log.Printf("userServiceClient.GetAllUsers( Failed: %v\n", err)
		usersChan <- []*pb.User{}
		return
	}
	usersChan <- users.Users
}

func goGetAllTodos(ctx context.Context, todosChan chan []*pb.Todo) {
	todos, err := todoServiceClient.GetAllTodos(ctx, &emptypb.Empty{})

	if err != nil {
		log.Printf("todoServiceClient.GetAllTodos( Failed: %v\n", err)
		todosChan <- []*pb.Todo{} // Prevent hanging other services.
		return
	}

	todosChan <- todos.Todos
}

// mapAllUsersToTodos takes in all the users from service-user and all combined entries
// from service-combine and will find the Todos associated with each service. This allows
// us to only make one call to service-combine's DB.
func mapAllUsersToTodos(ctx context.Context, users []*pb.User, combined []*pb.Combined) []*pb.UserWithTodos {
	usersToTodos := make(map[int32][]int32)

	var wg sync.WaitGroup
	var mapLock sync.RWMutex

	for _, user := range users {
		usersToTodos[user.Id] = make([]int32, 0)
	}

	for _, currentCombined := range combined {
		wg.Add(1)
		go func(currCombined *pb.Combined) {
			defer wg.Done()
			mapLock.Lock()

			userTodos := usersToTodos[currCombined.UserId]
			userTodos = append(userTodos, currCombined.TodoId)
			usersToTodos[currCombined.UserId] = userTodos

			mapLock.Unlock()
		}(currentCombined)
	}
	wg.Wait()

	// TODO: Watch this section in case of goroutine channel blocking - find better solution?
	// 10 is arbitrary to prevent buffer blocking (?)
	usersWithTodosChan := make(chan *pb.UserWithTodos, len(users)+constants.EnvValues.CHAN_BUFFER)

	for _, user := range users {
		wg.Add(1)
		go func(currentUser *pb.User) {
			defer wg.Done()
			userTodos, err := todoServiceClient.GetTodosByMultipleIds(ctx,
				&pb.MultipleTodoIDRequest{TodoIds: usersToTodos[currentUser.Id]})
			if err != nil {
				log.Printf("Error occurred in GetUsersWithTodos: go GetTodosFromUserId: %v\n", err)
				usersWithTodosChan <- &pb.UserWithTodos{User: currentUser, Todos: &pb.Todos{Todos: make([]*pb.Todo, 1)}}
				return
			}
			usersWithTodosChan <- &pb.UserWithTodos{User: currentUser, Todos: userTodos}
		}(user)
	}

	wg.Wait()
	close(usersWithTodosChan)

	var usersWithTodos []*pb.UserWithTodos
	for user := range usersWithTodosChan {
		usersWithTodos = append(usersWithTodos, user)
	}

	return usersWithTodos
}

func mapAllTodosToUsers(ctx context.Context, todos []*pb.Todo, combined []*pb.Combined) []*pb.TodoWithUsers {
	todosToUsers := make(map[int32][]int32)

	var wg sync.WaitGroup
	var mapLock sync.RWMutex

	for _, todo := range todos {
		todosToUsers[todo.Id] = make([]int32, 0)
	}

	for _, currentCombined := range combined {
		wg.Add(1)
		go func(currCombined *pb.Combined) {
			defer wg.Done()

			mapLock.Lock()
			todoUsers := todosToUsers[currCombined.TodoId]
			mapLock.Unlock()

			todoUsers = append(todoUsers, currCombined.UserId)

			mapLock.Lock()
			todosToUsers[currCombined.TodoId] = todoUsers
			mapLock.Unlock()
		}(currentCombined)
	}
	wg.Wait()

	log.Printf("Combined: %v\n", combined)

	// TODO: Watch this section in case of goroutine channel blocking - find better solution?
	// 10 is arbitrary to prevent buffer blocking (?)
	todosWithUsersChan := make(chan *pb.TodoWithUsers, len(todos)+constants.EnvValues.CHAN_BUFFER)

	for _, todo := range todos {
		wg.Add(1)
		go func(currentTodo *pb.Todo) {
			defer wg.Done()

			todoUsers, err := userServiceClient.GetUsersByMultipleIds(ctx,
				&pb.MultipleUserIDRequest{UserIds: todosToUsers[currentTodo.Id]})
			if err != nil {
				log.Printf("Error occurred in GetUsersWithTodos: go GetTodosFromUserId: %v\n", err)
				todosWithUsersChan <- &pb.TodoWithUsers{Todo: currentTodo, Users: &pb.Users{Users: make([]*pb.User, 1)}}
				return
			}
			todosWithUsersChan <- &pb.TodoWithUsers{Todo: currentTodo, Users: todoUsers}
		}(todo)
	}

	wg.Wait()
	close(todosWithUsersChan)

	var todosWithUsers []*pb.TodoWithUsers
	for todo := range todosWithUsersChan {
		todosWithUsers = append(todosWithUsers, todo)
	}

	return todosWithUsers
}

func GetAllUsersWithTodos() ([]*pb.UserWithTodos, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT*time.Second)
	defer cancel()

	usersChan := make(chan []*pb.User)
	allCombinedChan := make(chan []*pb.Combined)
	go goGetAllUsers(ctx, usersChan)
	go goGetAllCombined(ctx, allCombinedChan)

	users := <-usersChan
	combined := <-allCombinedChan

	usersWithTodos := mapAllUsersToTodos(ctx, users, combined)
	return usersWithTodos, nil
}

func GetAllTodosWithUsers() ([]*pb.TodoWithUsers, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT*time.Second)
	defer cancel()

	todosChan := make(chan []*pb.Todo)
	allCombinedChan := make(chan []*pb.Combined)

	go goGetAllTodos(ctx, todosChan)
	go goGetAllCombined(ctx, allCombinedChan)

	todos := <-todosChan
	combined := <-allCombinedChan

	todosWithUsers := mapAllTodosToUsers(ctx, todos, combined)
	return todosWithUsers, nil
}

func GetUserById(id int32) (*pb.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	user, err := userServiceClient.GetUserById(ctx, &pb.SingleUserIDRequest{UserId: id})

	if err != nil {
		log.Printf("userServiceClient.GetUserById( Failed: %v\n", err)
		return nil, err
	}
	return user, nil
}

func GetUserWithTodos(id int32) (*pb.UserWithTodos, error) {
	userChan := make(chan *pb.User)
	todosChan := make(chan []*pb.Todo)

	go func() {
		user, err := GetUserById(int32(id))
		if err != nil {
			log.Printf("Error occurred in GetUserById: grpcclient.GetUserById: %v\n", err)
			userChan <- &pb.User{} // to prevent hanging up the other services
			return
		}
		userChan <- user
	}()

	go func() {
		todos, err := GetTodosFromUserId(int32(id))
		if err != nil {
			log.Printf("Error occurred in GetUserWithTodos: grpcclient.GetTodosFromUserId: %v\n", err)
			todosChan <- []*pb.Todo{}
			return
		}
		todosChan <- todos
	}()

	user := <-userChan
	todos := <-todosChan

	return &pb.UserWithTodos{User: user, Todos: &pb.Todos{Todos: todos}}, nil
}

func AddUser(user *pb.User) (*pb.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	user, err := userServiceClient.AddUser(ctx, user)

	if err != nil {
		log.Printf("userServiceClient.AddUser( Failed: %v\n", err)
		return nil, err
	}
	return user, nil
}

func UpdateUser(user *pb.User) (*pb.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	user, err := userServiceClient.UpdateUserById(ctx, user)

	if err != nil {
		log.Printf("userServiceClient.UpdateUser( Failed: %v\n", err)
		return nil, err
	}
	return user, nil
}

func DeleteUser(userId int32) (*pb.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	user, err := userServiceClient.DeleteUser(ctx, &pb.SingleUserIDRequest{UserId: userId})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		_, err := combineServiceClient.DeleteUser(ctx, &pb.SingleUserIDRequest{UserId: userId})
		if err != nil {
			log.Printf("Error in DeleteUser go: error: %v\n", err)
			return
		}
	}()

	wg.Wait()

	if err != nil {
		log.Printf("userServiceClient.DeleteUser( Failed: %v\n", err)
		return nil, err
	}
	return user, nil
}

func GetTodos() ([]*pb.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	todos, err := todoServiceClient.GetAllTodos(ctx, &emptypb.Empty{})

	if err != nil {
		log.Printf("todoServiceClient.GetAllTodos( Failed: %v\n", err)
		return nil, err
	}
	return todos.Todos, nil
}

func GetTodoById(id int32) (*pb.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	todo, err := todoServiceClient.GetTodoById(ctx, &pb.SingleTodoIDRequest{TodoId: id})

	if err != nil {
		log.Printf("todoServiceClient.GetTodoById( Failed: %v\n", err)
		return nil, err
	}
	return todo, nil
}

func AddTodo(todo *pb.Todo) (*pb.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	addedTodo, err := todoServiceClient.AddTodo(ctx, todo)

	if err != nil {
		log.Printf("todoServiceClient.AddTodo( Failed: %v\n", err)
		return nil, err
	}
	return addedTodo, nil
}

func UpdateTodo(todo *pb.Todo) (*pb.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	updatedTodo, err := todoServiceClient.UpdateTodoById(ctx, todo)

	if err != nil {
		log.Printf("todoServiceClient.UpdateTodo( Failed: %v\n", err)
		return nil, err
	}
	return updatedTodo, nil
}

func DeleteTodo(todoId int32) (*pb.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	deletedTodo, err := todoServiceClient.DeleteTodo(ctx, &pb.SingleTodoIDRequest{TodoId: todoId})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := combineServiceClient.DeleteTodo(ctx, &pb.SingleTodoIDRequest{TodoId: todoId})
		if err != nil {
			log.Printf("combineServiceClient.DeleteTodo: Failed: %v\n", err)
			return
		}
	}()
	wg.Wait()

	if err != nil {
		log.Printf("todoServiceClient.DeleteTodo( Failed: %v\n", err)
		return nil, err
	}
	return deletedTodo, nil
}

func GetTodosFromUserId(userId int32) ([]*pb.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	todo_ids, err := combineServiceClient.GetTodosFromUserId(ctx, &pb.SingleUserIDRequest{UserId: userId})

	if err != nil {
		log.Printf("combineServiceClient.GetTodosFromUserId( Failed: %v\n", err)
		return nil, err
	}

	// Important: Handles case when a user has no todos to prevent bad SQL query in service-todo
	if len(todo_ids.TodoIds) == 0 {
		return []*pb.Todo{}, nil
	}

	log.Printf("User's todos: %v\n", todo_ids)

	todos, err := todoServiceClient.GetTodosByMultipleIds(ctx, &pb.MultipleTodoIDRequest{TodoIds: todo_ids.TodoIds})
	if err != nil {
		log.Printf("todoServiceClient.GetTodosByMultipleIds( Failed: %v\n", err)
		return nil, err
	}
	return todos.Todos, nil
}

func GetUsersFromTodoId(todoId int32) ([]*pb.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	user_ids, err := combineServiceClient.GetUsersFromTodoId(ctx, &pb.SingleTodoIDRequest{TodoId: todoId})

	if err != nil {
		log.Printf("combineServiceClient.GetUsersFromTodoId( Failed: %v\n", err)
		return nil, err
	}
	users, err := userServiceClient.GetUsersByMultipleIds(ctx, &pb.MultipleUserIDRequest{UserIds: user_ids.UserIds})
	if err != nil {
		log.Printf("userServiceClient.GetUsersByMultipleIds( Failed: %v\n", err)
		return nil, err
	}
	return users.Users, nil
}

func AddCombined(combinedRequest *pb.CombinedRequest) (*pb.Combined, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	insertedCombined, err := combineServiceClient.AddCombined(ctx, combinedRequest)

	if err != nil {
		log.Printf("combineServiceClient.AddCombined( Failed: %v\n", err)
		return nil, err
	}
	return insertedCombined, nil
}

func AddCombinedArray(combinedArrayRequest *pb.CombinedArrayRequest) ([]*pb.Combined, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	insertedCombines, err := combineServiceClient.AddCombinedArray(ctx, combinedArrayRequest)

	if err != nil {
		log.Printf("combineServiceClient.AddCombined( Failed: %v\n", err)
		return nil, err
	}
	return insertedCombines.Combines, nil
}

func DeleteCombined(combinedRequest *pb.CombinedRequest) (*pb.Combined, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	deletedCombined, err := combineServiceClient.DeleteCombined(ctx, combinedRequest)

	if err != nil {
		log.Printf("combineServiceClient.DeleteCombined( Failed: %v\n", err)
		return nil, err
	}
	return deletedCombined, nil
}

func DeleteTodoFromCombined(todoId int32) ([]*pb.Combined, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	deletedCombined, err := combineServiceClient.DeleteTodo(ctx, &pb.SingleTodoIDRequest{TodoId: todoId})

	if err != nil {
		log.Printf("combineServiceClient.DeleteTodoFromCombined( Failed: %v\n", err)
		return nil, err
	}
	return deletedCombined.Combines, nil
}

func DeleteUserFromCombined(userId int32) ([]*pb.Combined, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	deletedCombined, err := combineServiceClient.DeleteUser(ctx, &pb.SingleUserIDRequest{UserId: userId})

	if err != nil {
		log.Printf("combineServiceClient.DeleteUserFromCombined( Failed: %v\n", err)
		return nil, err
	}
	return deletedCombined.Combines, nil
}

func UpdateTodoCombined(todoId int32, userIds []int32) ([]*pb.Combined, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.GRPC_TIMEOUT)
	defer cancel()
	updatedCombined, err := combineServiceClient.UpdateTodo(ctx, &pb.CombinedArrayRequest{TodoId: todoId, UserIds: userIds})

	if err != nil {
		log.Printf("combineServiceClient.UpdateTodoCombined() Failed: %v\n", err)
		return nil, err
	}
	return updatedCombined.Combines, nil
}
