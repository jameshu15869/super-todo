package constants

import (
	"log"
	"time"

	"github.com/joho/godotenv"
)

// const (
// 	USER_ENDPOINT    = "localhost:4001"
// 	TODO_ENDPOINT    = "localhost:4002"
// 	COMBINE_ENDPOINT = "localhost:4003"

// GRPC_TIMEOUT = 2 * time.Second

// 	CHAN_BUFFER = 10

// 	REDIS_HOST           = "localhost"
// 	REDIS_PORT           = 6379
// 	REDIS_PASSWORD       = ""
// 	REDIS_DB             = 0
// 	REDIS_PUBSUB_CHANNEL = "updateChannel"

// 	REDIS_TIMEOUT  = 5 * time.Second
// 	CACHE_DURATION = 10 * time.Second
// )

type EnvConstants struct {
	USER_ENDPOINT string
	TODO_ENDPOINT string
	COMBINE_ENDPOINT string

	GRPC_TIMEOUT time.Duration

	CHAN_BUFFER int

	REDIS_HOST string
	REDIS_PORT int
	REDIS_PASSWORD string
	REDIS_DB int
	REDIS_PUBSUB_CHANNEL string
}

var EnvValues EnvConstants;

func InitEnvConstants() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Env file not found: %v\n", err)
	}

	EnvValues.USER_ENDPOINT = getGenericEnv("USER_ENDPOINT", "localhost:4001")
	EnvValues.TODO_ENDPOINT = getGenericEnv("TODO_ENDPOINT", "localhost:4002")
	EnvValues.COMBINE_ENDPOINT = getGenericEnv("COMBINE_ENDPOINT", "localhost:4003")

	EnvValues.GRPC_TIMEOUT = time.Duration(getGenericEnv("GRPC_TIMEOUT", 2)) * time.Second

	EnvValues.CHAN_BUFFER = getGenericEnv("CHAN_BUFFER", 10)

	EnvValues.REDIS_HOST = getGenericEnv("REDIS_HOST", "localhost")
	EnvValues.REDIS_PORT = getGenericEnv("REDIS_PORT", 6379)
	EnvValues.REDIS_PASSWORD = getGenericEnv("REDIS_PASSWORD", "")
	EnvValues.REDIS_DB = getGenericEnv("REDIS_DB", 0)
	EnvValues.REDIS_PUBSUB_CHANNEL = getGenericEnv("REDIS_PUBSUB_CHANNEL", "updateChannel")
}
