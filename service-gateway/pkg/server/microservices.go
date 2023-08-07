package server

import (
	"supertodo/gateway/pb"
	"supertodo/gateway/pkg/grpcclient"
)

type Microservice struct {
	Name      string
	GetHealth func() (*pb.Health, error)
}

var Microservices []*Microservice = []*Microservice{
	{
		Name:      "user",
		GetHealth: grpcclient.GetUserHealth,
	},
	{
		Name:      "todo",
		GetHealth: grpcclient.GetTodoHealth,
	},
	{
		Name:      "combine",
		GetHealth: grpcclient.GetCombineHealth,
	},
}
