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
	// POSTGRES_DBNAME   = "supertodo_users"
	// POSTGRES_SSLMODE  = "disable"

	// REDIS_HOST     = "localhost"
	// REDIS_PORT     = 6379
	// REDIS_PASSWORD = ""
	// REDIS_DB       = 0

	// REDIS_TIMEOUT  = 5 * time.Second
	// CACHE_DURATION = 10 * time.Second

	QUERY_ALL_USERS             = `SELECT * FROM "users";`
	QUERY_USER_BY_ID            = `SELECT * FROM users WHERE id=%d LIMIT 1;`
	QUERY_USER_BY_USERNAME      = `SELECT * FROM users WHERE username='%s' LIMIT 1;`
	QUERY_USERS_BY_MULTIPLE_IDS = `SELECT * FROM users WHERE id IN %s;`
	ADD_USER                    = `INSERT INTO "users"(username) VALUES($1) RETURNING *;`
	DELETE_USER                 = `DELETE FROM users WHERE id=$1 RETURNING *;`
	UPDATE_USER                 = `UPDATE users SET username=$1 WHERE id=$2 RETURNING *;`
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
	EnvValues.POSTGRES_DBNAME = getGenericEnv("POSTGRES_DBNAME", "supertodo_users")
	EnvValues.POSTGRES_SSLMODE = getGenericEnv("POSTGRES_SSLMODE", "disable")

	EnvValues.REDIS_HOST = getGenericEnv("REDIS_HOST", "localhost")
	EnvValues.REDIS_PORT = getGenericEnv("REDIS_PORT", 6379)
	EnvValues.REDIS_PASSWORD = getGenericEnv("REDIS_PASSWORD", "")
	EnvValues.REDIS_DB = getGenericEnv("REDIS_DB", 0)

	EnvValues.REDIS_TIMEOUT = time.Duration(getGenericEnv("REDIS_TIMEOUT", 5)) * time.Second
	EnvValues.CACHE_DURATION = time.Duration(getGenericEnv("CACHE_DURATION", 10)) * time.Second
}
