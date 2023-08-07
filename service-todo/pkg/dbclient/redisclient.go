package dbclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"supertodo/todo/pb"
	"supertodo/todo/pkg/constants"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func initRedis() {
	connectedRedisDB := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", constants.EnvValues.REDIS_HOST, constants.EnvValues.REDIS_PORT),
		Password: constants.EnvValues.REDIS_PASSWORD,
		DB:       constants.EnvValues.REDIS_DB,
	})

	rdb = connectedRedisDB
}

func CacheQuery[K []*pb.Todo | *pb.Todo](queryString string) (K, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.REDIS_TIMEOUT)
	defer cancel()
	val, err := rdb.Get(ctx, queryString).Result()
	switch {
	case err == redis.Nil:
		log.Printf("Key does not exist: %v\n", err)
		return nil, err
	case err != nil:
		log.Printf("Redis CacheQuery Failed: %v\n", err)
		return nil, err
	}

	var response K
	if err := json.Unmarshal([]byte(val), &response); err != nil {
		log.Printf("Error occurred while unmarshaling JSON: %v\n", err)
		return nil, err
	}

	return response, nil
}

func CacheSet(key string, val interface{}) error {
	encoded, err := json.Marshal(val)
	if err != nil {
		log.Printf("Error occurred in CacheSet: Marshal: %v\n", err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), constants.EnvValues.REDIS_TIMEOUT)
	defer cancel()
	if err := rdb.Set(ctx, key, encoded, constants.EnvValues.CACHE_DURATION).Err(); err != nil {
		log.Printf("Error occurred in CacheSet: %v\n", err)
		return errors.New("error occurred while setting Redis cache")
	}
	return nil
}
