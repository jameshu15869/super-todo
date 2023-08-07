package constants

import (
	"log"
	"time"

	"github.com/joho/godotenv"
)

const (
	// POSTGRES_USER     = "postgres"
	// POSTGRES_PASSWORD = "password"
	// POSTGRES_HOST     = "localhost"
	// POSTGRES_PORT     = 5432
	// POSTGRES_DBNAME   = "supertodo_combined"
	// POSTGRES_SSLMODE  = "disable"

	// REDIS_HOST     = "localhost"
	// REDIS_PORT     = 6379
	// REDIS_PASSWORD = ""
	// REDIS_DB       = 0

	// REDIS_TIMEOUT  = 5 * time.Second
	// CACHE_DURATION = 10 * time.Second

	QUERY_ALL_COMBINED         = `SELECT * FROM combined;`
	QUERY_TODOS_FROM_USER_ID   = `SELECT * FROM combined WHERE user_id=%d;`
	QUERY_USERS_FROM_TODO_ID   = `SELECT * FROM combined WHERE todo_id=%d;`
	ADD_COMBINED               = `INSERT INTO combined(user_id, todo_id) VALUES (%d, %d) RETURNING *;`
	ADD_COMBINED_ARRAY_PREFIX  = `INSERT INTO combined(user_id, todo_id) VALUES`
	ADD_COMBINED_ARRAY_POSTFIX = ` RETURNING *;`
	DELETE_COMBINED            = `DELETE FROM combined WHERE user_id=%d AND todo_id=%d RETURNING *;`
	DELETE_TODO                = `DELETE FROM combined WHERE todo_id=%d RETURNING *;`
	DELETE_USER                = `DELETE FROM combined WHERE user_id=%d RETURNING *;`
)

type EnvConstants struct {
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_HOST     string
	POSTGRES_PORT	int
	POSTGRES_DBNAME   string
	POSTGRES_SSLMODE  string

	REDIS_HOST     string
	REDIS_PORT     int
	REDIS_PASSWORD string
	REDIS_DB       int

	REDIS_TIMEOUT  time.Duration
	CACHE_DURATION time.Duration
}

var EnvValues EnvConstants

func InitEnvConstants() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Env file not found: %v\n", err)
	}

	EnvValues.POSTGRES_USER = getGenericEnv("POSTGRES_USER", "postgres")
	EnvValues.POSTGRES_PASSWORD = getGenericEnv("POSTGRES_PASSWORD", "password")
	EnvValues.POSTGRES_HOST = getGenericEnv("POSTGRES_HOST", "localhost")
	EnvValues.POSTGRES_PORT = getGenericEnv("POSTGRES_PORT", 5432)
	EnvValues.POSTGRES_DBNAME = getGenericEnv("POSTGRES_DBNAME", "supertodo_combined")
	EnvValues.POSTGRES_SSLMODE = getGenericEnv("POSTGRES_SSLMODE", "disable")

	EnvValues.REDIS_HOST = getGenericEnv("REDIS_HOST", "localhost")
	EnvValues.REDIS_PORT = getGenericEnv("REDIS_PORT", 6379)
	EnvValues.REDIS_PASSWORD = getGenericEnv("REDIS_PASSWORD", "")
	EnvValues.REDIS_DB = getGenericEnv("REDIS_DB", 0)

	EnvValues.REDIS_TIMEOUT = time.Duration(getGenericEnv("REDIS_TIMEOUT", 5)) * time.Second
	EnvValues.CACHE_DURATION = time.Duration(getGenericEnv("CACHE_DURATION", 10)) * time.Second
}
