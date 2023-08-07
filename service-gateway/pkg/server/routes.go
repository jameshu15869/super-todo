package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	pb "supertodo/gateway/pb"
	"supertodo/gateway/pkg/grpcclient"
	"supertodo/gateway/pkg/redisclient"
	"supertodo/gateway/pkg/sse"
	"supertodo/gateway/pkg/utils"
	"sync"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg" : "pong"})
}

func GetBackendHealth(c *gin.Context) {
	var healthStatusMap map[string]string = make(map[string]string)

	var wg sync.WaitGroup
	var mapLock sync.RWMutex = sync.RWMutex{}

	for _, ms := range Microservices {
		wg.Add(1)
		go func(currentMicroservice *Microservice) {
			defer wg.Done()
			response, err := currentMicroservice.GetHealth()
			var serviceHealth string
			if err != nil {
				log.Printf("An error occurred in GetHealth(): %v", err)
				serviceHealth = "unhealthy"
			} else {
				serviceHealth = response.Status
			}
			mapLock.Lock()
			healthStatusMap[currentMicroservice.Name] = serviceHealth
			mapLock.Unlock()
		}(ms)
	}

	wg.Wait()

	c.JSON(http.StatusOK, gin.H{"status": healthStatusMap})
}

func GetUsers(c *gin.Context) {
	users, err := grpcclient.GetUsers()
	if err != nil {
		log.Printf("Error occured in grpcclient.GetUsers(): %v\n", err)
	}

	if users == nil {
		users = []*pb.User{}
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func GetUserById(c *gin.Context) {
	id := c.Params.ByName("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("Error occured in GetUserById: ParseInt: %v\n", err)
		c.AbortWithStatusJSON(http.StatusNotFound, utils.FormatJsonError("GetUserById", err))
		return
	}
	user, err := grpcclient.GetUserById(int32(intId))
	if err != nil {
		log.Printf("Error occurred in GetUserById: grpcclient.GetUserById: %v\n", err)
		c.AbortWithStatusJSON(http.StatusNotFound, utils.FormatJsonError("GetUserById", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func GetUserWithTodos(c *gin.Context) {
	id := c.Params.ByName("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("Error occured in GetUserWithTodos: ParseInt: %v\n", err)
		c.AbortWithStatusJSON(http.StatusNotFound, utils.FormatJsonError("GetUserWithTodos", err))
		return
	}

	user, err := grpcclient.GetUserWithTodos(int32(intId))
	if err != nil {
		log.Printf("Error occured in GetUserWithTodos: grpcclient.GetUserWithTodos: %v\n", err)
		c.AbortWithStatusJSON(http.StatusNotFound, utils.FormatJsonError("GetUserWithTodos", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"userWithTodos": user})
}

func GetAllUsersWithTodos(c *gin.Context) {
	users, err := grpcclient.GetAllUsersWithTodos()
	if err != nil {
		log.Printf("Error occured in GetAllUsersWithTodos: grpcclient.GetUsersWithTodos: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func GetAllTodosWithUsers(c *gin.Context) {
	todos, err := grpcclient.GetAllTodosWithUsers()
	if err != nil {
		log.Printf("Error occurred in GetAllTodosWithUsers: grpcclient.GetTodosWithUsers: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{"todos": todos})
}

func PostAddUser(c *gin.Context) {
	var newUser pb.PutUser
	if err := c.BindJSON(&newUser); err != nil {
		log.Printf("Error occurred in PostAddUser: %v\n", err)
	}
	log.Printf("String: %s\n", newUser.Username)
	if len(newUser.Username) == 0 {
		log.Printf("Username cannot be empty")
		c.AbortWithError(http.StatusNotAcceptable, errors.New("username cannot be empty"))
		return
	}
	addedUser, err := grpcclient.AddUser(&pb.User{Id: -1, Username: newUser.Username})
	if err != nil {
		log.Printf("Error occurred in PostAddUser: grpcclient.AddUser: %v\n", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.FormatJsonError("GetUserWithTodos", err))
		return
	}
	jsonString, err := json.Marshal(gin.H{"user": addedUser})
	if err != nil {
		log.Printf("Error occurred in PostAddUser: json.Marshal: %v\n", err)
		return
	}
	redisclient.PublishSseMessage(c, &sse.Message{MessageType: "add-user", Content: string(jsonString)})
	c.JSON(http.StatusOK, gin.H{"user": addedUser})
}

func PostDeleteUser(c *gin.Context) {
	strId := c.Params.ByName("id")
	parsedId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		log.Printf("Error occurred in PostDeleteUser: ParseInt: %v\n", err)
	}

	deletedUser, err := grpcclient.DeleteUser(int32(parsedId))
	if err != nil {
		log.Printf("Error occurred in PostDeleteUser: grpcclient.DeleteUser: %v\n", err)
	}
	c.JSON(http.StatusOK, gin.H{"user": deletedUser})
}

func PostUpdateUser(c *gin.Context) {
	strId := c.Params.ByName("id")
	parsedId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		log.Printf("Error occurred in PostDeleteUser: ParseInt: %v\n", err)
	}

	var putUser pb.PutUser
	if err := c.BindJSON(&putUser); err != nil {
		log.Printf("Error occurred in PostUpdateUser: BindJson: %v\n", err)
	}
	updatedUser, err := grpcclient.UpdateUser(&pb.User{Id: int32(parsedId), Username: putUser.Username})
	if err != nil {
		log.Printf("Error occurred in PostUpdateUser: grpcclient.UpdateUser: %v\n", err)
	}
	c.JSON(http.StatusOK, gin.H{"user": updatedUser})
}

func GetTodos(c *gin.Context) {
	todos, err := grpcclient.GetTodos()
	if err != nil {
		log.Printf("Error occurred in grpcclient.GetTodos(): %v\n", err)
	}
	if todos == nil {
		todos = []*pb.Todo{}
	}
	c.JSON(http.StatusOK, gin.H{"todos": todos})
}

func GetTodoById(c *gin.Context) {
	id := c.Params.ByName("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("Error occured in GetTodoById: ParseInt: %v\n", err)
	}
	user, err := grpcclient.GetTodoById(int32(intId))
	if err != nil {
		log.Printf("Error occurred in GetTodoById: grpcclient.GetTodoById: %v\n", err)
	}
	c.JSON(http.StatusOK, gin.H{"todo": user})
}

func GetTodoWithUsers(c *gin.Context) {
	id := c.Params.ByName("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("Error occured in GetTodoWithUsers: ParseInt: %v\n", err)
	}

	todoChan := make(chan *pb.Todo)
	usersChan := make(chan []*pb.User)

	go func() {
		todo, err := grpcclient.GetTodoById(int32(intId))
		if err != nil {
			log.Printf("Error occurred in GetTodoWithUsers: grpcclient.GetTodoById: %v\n", err)
		}
		todoChan <- todo
	}()

	go func() {
		users, err := grpcclient.GetUsersFromTodoId(int32(intId))
		if err != nil {
			log.Printf("Error occurred in GetTodoWithUsers: grpcclient.GetUsersFromTodoId: %v\n", err)
		}
		usersChan <- users
	}()

	todo := <-todoChan
	users := <-usersChan

	if todo == nil {
		c.AbortWithError(http.StatusNotFound, errors.New("todo was null"))
		return
	}

	if users == nil {
		c.AbortWithError(http.StatusNotFound, errors.New("users was null"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"todo": todo, "users": users})
}

func PostAddTodo(c *gin.Context) {
	var newTodo pb.PutTodo

	if err := c.BindJSON(&newTodo); err != nil {
		log.Printf("Error occurred in PostAddTodo: %v\n", err)
		return
	}
	addedTodo, err := grpcclient.AddTodo(&pb.Todo{Id: -1, Title: newTodo.Title, TodoDate: newTodo.TodoDate, Body: newTodo.Body})
	if err != nil {
		log.Printf("Error occurred in PostAddTodo: grpcclient.AddTodo: %v\n", err)
		return
	}

	addedCombines, err := grpcclient.AddCombinedArray(&pb.CombinedArrayRequest{TodoId: addedTodo.Id, UserIds: newTodo.UserIds})
	if err != nil {
		log.Printf("Error occurred in PostAddTodo: grpcclient.AddCombinedArray: %v\n", err)
		return
	}

	jsonString, err := json.Marshal(gin.H{"todo": addedTodo, "combined": addedCombines})
	if err != nil {
		log.Printf("Error occurred in PostAddTodo: json.Marshal: %v\n", err)
		return
	}
	// sseServer.PushMessage(&sse.Message{MessageType: "add-todo", Content: string(jsonString)})
	redisclient.PublishSseMessage(c, &sse.Message{MessageType: "add-todo", Content: string(jsonString)})
	c.JSON(http.StatusOK, gin.H{"todo": addedTodo, "combined": addedCombines})
}

func PostDeleteTodo(c *gin.Context) {
	strId := c.Params.ByName("id")
	parsedId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		log.Printf("Error occurred in PostDeleteTodo: ParseInt: %v\n", err)
	}

	deletedTodo, err := grpcclient.DeleteTodo(int32(parsedId))
	if err != nil {
		log.Printf("Error occurred in PostDeleteTodo: grpcclient.DeleteTodo: %v\n", err)
	}

	jsonString, err := json.Marshal(deletedTodo)
	if err != nil {
		log.Printf("Error occurred in PostDeleteTodo: json.Marshal: %v\n", err)
		return
	}
	// sseServer.PushMessage(&sse.Message{MessageType: "delete-todo", Content: string(jsonString)})
	redisclient.PublishSseMessage(c, &sse.Message{MessageType: "delete-todo", Content: string(jsonString)})

	c.JSON(http.StatusOK, gin.H{"todo": deletedTodo})
}

func PostUpdateTodo(c *gin.Context) {
	strId := c.Params.ByName("id")
	parsedId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		log.Printf("Error occurred in PostUpdateTodo: ParseInt: %v\n", err)
	}

	var newTodo pb.PutTodo
	if err := c.BindJSON(&newTodo); err != nil {
		log.Printf("Error occurred in PostAddTodo: %v\n", err)
	}

	todoChan := make(chan *pb.Todo)
	combinedChan := make(chan []*pb.Combined)

	go func() {
		updatedTodo, err := grpcclient.UpdateTodo(&pb.Todo{Id: int32(parsedId), Title: newTodo.Title, TodoDate: newTodo.TodoDate, Body: newTodo.Body})
		if err != nil {
			log.Printf("Error occurred in PostUpdateTodo: grpcclient.UpdateTodo: %v\n", err)
			todoChan <- nil
			return
		}
		todoChan <- updatedTodo
	}()

	go func() {
		updatedCombined, err := grpcclient.UpdateTodoCombined(int32(parsedId), newTodo.UserIds)
		if err != nil {
			log.Printf("Error occurred in PostUpdateTodo: grpcclient.UpdateTodoCombined: %v\n", err)
			combinedChan <- nil
			return
		}
		combinedChan <- updatedCombined
	}()

	todo := <-todoChan
	combined := <-combinedChan
	if todo == nil {
		c.AbortWithError(http.StatusNotFound, errors.New("todo was null"))
		return
	}
	if combined == nil {
		c.AbortWithError(http.StatusNotFound, errors.New("combined was null"))
		return
	}

	jsonString, err := json.Marshal(gin.H{"todo": todo, "combined": combined})
	if err != nil {
		log.Printf("Error occurred in PostUpdateTodo: json.Marshal: %v\n", err)
		return
	}
	// sseServer.PushMessage(&sse.Message{MessageType: "update-todo", Content: string(jsonString)})
	redisclient.PublishSseMessage(c, &sse.Message{MessageType: "update-todo", Content: string(jsonString)})

	c.JSON(http.StatusOK, gin.H{"todo": todo})
}

func GetCombined(c *gin.Context) {
	combines, err := grpcclient.GetAllCombined()
	if err != nil {
		log.Printf("Error occurred in grpcclient.GetAllCombined(): %v\n", err)
	}

	if combines == nil {
		combines = []*pb.Combined{}
	}
	
	c.JSON(http.StatusOK, gin.H{"combines": combines})
}

func PostAddCombined(c *gin.Context) {
	var combinedRequest pb.CombinedRequest
	if err := c.BindJSON(&combinedRequest); err != nil {
		log.Printf("Error occurred in PostAddCombined: %v\n", err)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, utils.FormatJsonError("PostAddCombined", err))
		return
	}

	addedCombined, err := grpcclient.AddCombined(&combinedRequest)
	if err != nil {
		log.Printf("Error occurred in PostAddCombined: %v\n", err)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, utils.FormatJsonError("PostAddCombined", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"combined": addedCombined})
}

func PostAddCombinedArray(c *gin.Context) {
	var combinedArrayRequest pb.CombinedArrayRequest
	if err := c.BindJSON(&combinedArrayRequest); err != nil {
		log.Printf("Error occurred in PostAddCombined: %v\n", err)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, utils.FormatJsonError("PostAddCombinedArray", err))
		return
	}

	addedCombines, err := grpcclient.AddCombinedArray(&combinedArrayRequest)
	if err != nil {
		log.Printf("Error occurred in PostAddCombinedArray: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.FormatJsonError("PostAddCombined", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"combined": addedCombines})
}

func PostDeleteCombined(c *gin.Context) {
	var combinedRequest pb.CombinedRequest
	if err := c.BindJSON(&combinedRequest); err != nil {
		log.Printf("Error occurred in PostDeleteCombined: %v\n", err)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, utils.FormatJsonError("PostDeleteCombined", err))
		return
	}

	deletedCombined, err := grpcclient.DeleteCombined(&combinedRequest)
	if err != nil {
		log.Printf("Error occurred in PostDeleteCombined: %v\n", err)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, utils.FormatJsonError("PostDeleteCombined", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"combined": deletedCombined})
}
