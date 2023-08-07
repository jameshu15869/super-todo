package redisclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"supertodo/gateway/pkg/constants"
	"supertodo/gateway/pkg/sse"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func InitRedis() {
	connectedRedisDB := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", constants.EnvValues.REDIS_HOST, constants.EnvValues.REDIS_PORT),
		Password: constants.EnvValues.REDIS_PASSWORD,
		DB:       constants.EnvValues.REDIS_DB,
	})

	rdb = connectedRedisDB
}

func Publish(ctx context.Context, message string) {
	err := rdb.Publish(ctx, constants.EnvValues.REDIS_PUBSUB_CHANNEL, message).Err()
	if err != nil {
		log.Printf("Error occurred in Publishing: %v\n", err)
	}
}

func PublishSseMessage(ctx context.Context, message *sse.Message) {
	jsonString, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error occurred in PublishSseMessage: json.Marshal: %v\n", err)
		return
	}

	Publish(ctx, string(jsonString))
}

func Subscribe(s *sse.SseServer) {
	ctx := context.Background()
	pubsub := rdb.Subscribe(ctx, constants.EnvValues.REDIS_PUBSUB_CHANNEL)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		log.Println("Received pubsub message")
		log.Println(msg.Channel, msg.Payload)
		var sseMessage *sse.Message
		if err := json.Unmarshal([]byte(msg.Payload), &sseMessage); err != nil {
			log.Printf("Error occurred while unmarshaling JSON: %v\n", err)
			return
		}
		s.PushMessage(&sse.Message{MessageType: sseMessage.MessageType, Content: sseMessage.Content})
	}
}
