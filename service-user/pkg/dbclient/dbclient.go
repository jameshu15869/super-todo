package dbclient

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	pb "supertodo/user/pb"
	"supertodo/user/pkg/constants"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		constants.EnvValues.POSTGRES_USER,
		constants.EnvValues.POSTGRES_PASSWORD,
		constants.EnvValues.POSTGRES_HOST,
		constants.EnvValues.POSTGRES_PORT,
		constants.EnvValues.POSTGRES_DBNAME,
		constants.EnvValues.POSTGRES_SSLMODE)
	openDB, err := sql.Open("postgres", connStr)
	for err != nil {
		log.Printf("Failed to open DB connection: %v\n", err)
		time.Sleep(1 * time.Second)
		openDB, err = sql.Open("postgres", connStr)
	}

	db = openDB
	err = db.Ping()
	if err != nil {
		log.Printf("Failed to ping DB: %v\n", err)
		time.Sleep(1 * time.Second)
		err = db.Ping()
	}

	initRedis()

	log.Println("Successfully connected to DB")
}

func userFromRow(rows *sql.Rows) (*pb.User, error) {
	var id int32
	var username string
	if err := rows.Scan(&id, &username); err != nil {
		log.Printf("Error occurred in userFromRow: Scanning Rows: %v\n", err)
		return nil, err
	}
	return &pb.User{Id: id, Username: username}, nil
}

// func queryRowWithSingleReturn(queryString string, functionName string) (*pb.User, error) {
	// var insertedId int
	// var insertedUsername string
	// if err := db.QueryRow(queryString).Scan(&insertedId, &insertedUsername); err != nil {
	// 	log.Printf("Error occurred in %s: %v\n", functionName, err)
	// 	return nil, err
	// }

	// insertedUser := &pb.User{Id: int32(insertedId), Username: insertedUsername}
	// return insertedUser, nil
// }

func queryWithDollars(queryString string, functionName string, args ...interface{}) (*pb.User, error) {
	var insertedId int
	var insertedUsername string
	if err := db.QueryRow(queryString, args...).Scan(&insertedId, &insertedUsername); err != nil {
		log.Printf("Error occurred in %s: %v\n", functionName, err)
		return nil, err
	}
	insertedUser := &pb.User{Id: int32(insertedId), Username: insertedUsername}
	return insertedUser, nil
}

func generateMultipleIdSetString(ids []int32) string {
	result := "("
	for i, id := range ids {
		if i == len(ids)-1 {
			result += strconv.Itoa(int(id))
		} else {
			result += (strconv.Itoa(int(id)) + ", ")
		}
	}
	result += ")"
	return result
}

func AddUser(user *pb.User) (*pb.User, error) {
	insertedUser, err := queryWithDollars(constants.ADD_USER, "AddUser", user.Username)
	if err != nil {
		log.Printf("Error occurred in AddUser: %v\n", err)
		return nil, err
	}
	return insertedUser, nil
}

func DeleteUser(userId int32) (*pb.User, error) {
	insertedUser, err := queryWithDollars(constants.DELETE_USER, "DeleteUser", userId)
	if err != nil {
		log.Printf("Error occurred in DeleteUser: %v\n", err)
		return nil, err
	}
	return insertedUser, nil
}

func UpdateUser(user *pb.User) (*pb.User, error) {
	insertedUser, err := queryWithDollars(constants.UPDATE_USER, "UpdateUser", user.Username, user.Id)
	if err != nil {
		log.Printf("Error occurred in UpdateUser: %v\n", err)
		return nil, err
	}
	return insertedUser, nil
}

func QueryAllUsers() ([]*pb.User, error) {
	queryString := constants.QUERY_ALL_USERS

	cacheVal, err := CacheQuery[[]*pb.User](queryString)

	if err != nil {
		// Value not found in cache
		log.Println("QueryAllUsers: Value not found in cache, querying DB")
		rows, err := db.Query(queryString)
		if err != nil {
			log.Printf("Error occurred in QueryAllUsers: %v\n", err)
			return nil, err
		}
		defer rows.Close()

		var users []*pb.User
		for rows.Next() {
			currentUser, err := userFromRow(rows)
			if err != nil {
				log.Printf("Error occurred in QueryAllUsers: Scanning Rows: %v\n", err)
				continue
			}
			users = append(users, currentUser)
		}
		CacheSet(queryString, users)
		return users, nil
	}

	log.Println("QueryAllUsers: Cache: HIT")
	return cacheVal, nil
}

func QueryUserById(id int32) (*pb.User, error) {
	queryString := fmt.Sprintf(constants.QUERY_USER_BY_ID, id)

	cacheVal, err := CacheQuery[*pb.User](queryString)

	if err != nil {
		// Value not found in cache
		log.Println("QueryUserById: Value not found in cache, querying DB")
		rows, err := db.Query(queryString)
		if err != nil {
			log.Printf("Error occurred in QueryUserById: %v\n", err)
			return nil, err
		}
		defer rows.Close()

		var user *pb.User
		for rows.Next() {
			currentUser, err := userFromRow(rows)
			if err != nil {
				log.Printf("Error occurred in QueryUserById: Scanning Rows: %v\n", err)
				continue
			}
			user = currentUser
			break
		}
		CacheSet(queryString, user)
		return user, nil
	}

	log.Println("QueryUserById: Cache: HIT")
	return cacheVal, nil
}

func QueryUsersByMultipleIds(userIds []int32) ([]*pb.User, error) {
	if len(userIds) == 0 {
		return []*pb.User{}, nil
	}
	setString := generateMultipleIdSetString(userIds)
	queryString := fmt.Sprintf(constants.QUERY_USERS_BY_MULTIPLE_IDS, setString)
	log.Println(queryString)

	cacheVal, err := CacheQuery[[]*pb.User](queryString)

	if err != nil {
		// Value not found in cache
		log.Println("QueryUsersByMultipleIds: Value not found in cache, querying DB")
		rows, err := db.Query(queryString)
		if err != nil {
			log.Printf("Error occurred in QueryUsersByMultipleIds: %v\n", err)
			return nil, err
		}
		defer rows.Close()

		var users []*pb.User
		for rows.Next() {
			currentUser, err := userFromRow(rows)
			if err != nil {
				log.Printf("Error occurred in QueryUsersByMultipleIds: Scanning Rows: %v\n", err)
				continue
			}
			users = append(users, currentUser)
		}
		CacheSet(queryString, users)
		return users, nil
	}

	log.Println("QueryUsersByMultipleIds: Cache: HIT")
	return cacheVal, nil
}

func QueryUserByUsername(username string) (*pb.User, error) {
	queryString := fmt.Sprintf(constants.QUERY_USER_BY_USERNAME, username)

	cacheVal, err := CacheQuery[*pb.User](queryString)

	if err != nil {
		// Value not found in cache
		log.Println("QueryUserByUsername: Value not found in cache, querying DB")
		rows, err := db.Query(queryString)
		if err != nil {
			log.Printf("Error occurred in QueryUserByUsername: %v\n", err)
			return nil, err
		}
		defer rows.Close()

		var user *pb.User
		for rows.Next() {
			currentUser, err := userFromRow(rows)
			if err != nil {
				log.Printf("Error occurred in QueryUserByUsername: Scanning Rows: %v\n", err)
				continue
			}
			user = currentUser
			break
		}
		CacheSet(queryString, user)
		return user, nil
	}

	log.Println("QueryUserByUsername: Cache: HIT")
	return cacheVal, nil
}
