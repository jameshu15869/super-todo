package server

import (
	"supertodo/gateway/pkg/constants"
	"supertodo/gateway/pkg/grpcclient"
	"supertodo/gateway/pkg/redisclient"
	"supertodo/gateway/pkg/sse"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var sseServer *sse.SseServer

func Start() {
	constants.InitEnvConstants();

	grpcclient.InitGrpc()
	router := gin.Default()
	router.Use(cors.Default())
	sseServer = sse.NewSseServer()
	redisclient.InitRedis()

	go redisclient.Subscribe(sseServer)

	router.GET("/health", GetBackendHealth)
	router.GET("/ping", Ping)

	router.GET("/users", GetUsers)
	router.GET("/userswithtodos", GetAllUsersWithTodos)
	router.GET("/users/:id", GetUserWithTodos)
	router.POST("/users/add", PostAddUser)
	router.POST("/users/:id/delete", PostDeleteUser)
	router.POST("/users/:id/update", PostUpdateUser)

	router.GET("/todos", GetTodos)
	router.GET("/todoswithusers", GetAllTodosWithUsers)
	router.POST("/todos/add", PostAddTodo)
	router.GET("/todos/:id", GetTodoWithUsers)
	router.POST("/todos/:id/delete", PostDeleteTodo)
	router.POST("/todos/:id/update", PostUpdateTodo)

	router.GET("/combined", GetCombined)
	router.POST("/combined/add", PostAddCombined)
	router.POST("/combined/addarray", PostAddCombinedArray)
	router.POST("/combined/delete", PostDeleteCombined)

	router.GET("/sse", sse.SseMiddleware(), sseServer.InitializeConnection(), sse.Stream())

	router.Run("0.0.0.0:4000")
}
