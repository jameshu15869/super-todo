package dbclient

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	pb "supertodo/combine/pb"
	"supertodo/combine/pkg/constants"
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

func combinedFromRow(rows *sql.Rows) (*pb.Combined, error) {
	var id int32
	var user_id int32
	var todo_id int32
	if err := rows.Scan(&id, &user_id, &todo_id); err != nil {
		log.Printf("Error occurred in combinedFromRow: Scanning Rows: %v\n", err)
		return nil, err
	}
	return &pb.Combined{Id: id, UserId: user_id, TodoId: todo_id}, nil
}

func queryRowWithSingleReturn(queryString string, functionName string) (*pb.Combined, error) {
	var insertedId int
	var insertedUserId int
	var insertedTodoId int
	if err := db.QueryRow(queryString).Scan(&insertedId, &insertedUserId, &insertedTodoId); err != nil {
		log.Printf("Error occurred in %s: %v\n", functionName, err)
		return nil, err
	}

	insertedCombined := &pb.Combined{Id: int32(insertedId), UserId: int32(insertedUserId), TodoId: int32(insertedTodoId)}
	return insertedCombined, nil
}

func queryRowWithMultipleReturn(queryString string, functionName string) ([]*pb.Combined, error) {
	rows, err := db.Query(queryString)
	if err != nil {
		log.Printf("Error occured in %s: %v\n", functionName, err)
		return nil, err
	}
	var removedCombined []*pb.Combined
	for rows.Next() {
		var returnedId int
		var returnedUserId int
		var returnedTodoId int
		if err := rows.Scan(&returnedId, &returnedUserId, &returnedTodoId); err != nil {
			log.Printf("Error occured in %s Scan: %v\n", functionName, err)
			continue
		}
		removedCombined = append(removedCombined, &pb.Combined{Id: int32(returnedId), UserId: int32(returnedUserId), TodoId: int32(returnedTodoId)})
	}
	return removedCombined, nil
}

func QueryAllCombined() ([]*pb.Combined, error) {
	queryString := constants.QUERY_ALL_COMBINED

	cacheVal, err := CacheQuery[[]*pb.Combined](queryString)

	if err != nil {
		// Value not found in cache
		log.Println("QueryAllCombined: Value not found in cache, querying DB")
		rows, err := db.Query(queryString)
		if err != nil {
			log.Printf("Error occurred in QueryAllCombined: %v\n", err)
			return nil, err
		}
		defer rows.Close()

		var allCombined []*pb.Combined
		for rows.Next() {
			currentCombined, err := combinedFromRow(rows)
			if err != nil {
				log.Printf("Error occurred in QueryAllCombined: Scanning Rows: %v\n", err)
				continue
			}
			allCombined = append(allCombined, currentCombined)
		}
		CacheSet(queryString, allCombined)
		return allCombined, nil
	}

	log.Println("QueryAllCombined: Cache: HIT")
	return cacheVal, nil
}

func QueryTodosFromUserId(user_id int32) ([]int32, error) {
	queryString := fmt.Sprintf(constants.QUERY_TODOS_FROM_USER_ID, user_id)

	cacheVal, err := CacheQuery[[]int32](queryString)

	if err != nil {
		// Value not found in cache
		log.Println("QueryTodosFromUserId: Value not found in cache, querying DB")
		rows, err := db.Query(queryString)
		if err != nil {
			log.Printf("Error occurred in QueryTodosFromUserId: %v\n", err)
			return nil, err
		}
		defer rows.Close()

		var todo_ids []int32
		for rows.Next() {
			currentCombined, err := combinedFromRow(rows)
			if err != nil {
				log.Printf("Error occurred in QueryAllCombined: Scanning Rows: %v\n", err)
				continue
			}
			todo_ids = append(todo_ids, currentCombined.TodoId)
		}
		CacheSet(queryString, todo_ids)
		return todo_ids, nil
	}

	log.Println("QueryAllCombined: Cache: HIT")
	return cacheVal, nil
}

func QueryUsersFromTodoId(todo_id int32) ([]int32, error) {
	queryString := fmt.Sprintf(constants.QUERY_USERS_FROM_TODO_ID, todo_id)

	cacheVal, err := CacheQuery[[]int32](queryString)

	if err != nil {
		// Value not found in cache
		log.Println("QueryUsersFromTodoId: Value not found in cache, querying DB")
		rows, err := db.Query(queryString)
		if err != nil {
			log.Printf("Error occurred in QueryUsersFromTodoId: %v\n", err)
			return nil, err
		}
		defer rows.Close()

		var user_ids []int32
		for rows.Next() {
			currentCombined, err := combinedFromRow(rows)
			if err != nil {
				log.Printf("Error occurred in QueryUsersFromTodoId: Scanning Rows: %v\n", err)
				continue
			}
			user_ids = append(user_ids, currentCombined.UserId)
		}
		CacheSet(queryString, user_ids)
		return user_ids, nil
	}

	log.Println("QueryUsersFromTodoId: Cache: HIT")
	return cacheVal, nil
}

func AddCombined(user_id int32, todo_id int32) (*pb.Combined, error) {
	dynamicString := fmt.Sprintf(constants.ADD_COMBINED, user_id, todo_id)

	insertedCombined, err := queryRowWithSingleReturn(dynamicString, "AddCombined")
	if err != nil {
		log.Printf("Error occurred in AddCombined: %v\n", err)
		return nil, err
	}
	return insertedCombined, nil
}

func AddCombinedArray(todo_id int32, user_ids []int32) ([]*pb.Combined, error) {
	queryString := strings.Clone(constants.ADD_COMBINED_ARRAY_PREFIX)
	for index, user_id := range user_ids {
		queryString += fmt.Sprintf(`(%d, %d)`, user_id, todo_id)
		if index < len(user_ids)-1 {
			queryString += ","
		}
	}
	queryString += constants.ADD_COMBINED_ARRAY_POSTFIX

	addedCombines, err := queryRowWithMultipleReturn(queryString, "AddCombinedArray")
	if err != nil {
		log.Printf("Error occurred in AddCombinedArray: %v\n", err)
		return nil, err
	}
	return addedCombines, nil
}

func DeleteCombined(user_id int32, todo_id int32) (*pb.Combined, error) {
	dynamicString := fmt.Sprintf(constants.DELETE_COMBINED, user_id, todo_id)

	removedCombined, err := queryRowWithSingleReturn(dynamicString, "DeleteCombined")
	if err != nil {
		log.Printf("Error occurred in DeleteCombined: %v\n", err)
		return nil, err
	}
	return removedCombined, nil
}

func DeleteTodo(todo_id int32) ([]*pb.Combined, error) {
	dynamicString := fmt.Sprintf(constants.DELETE_TODO, todo_id)

	removedCombined, err := queryRowWithMultipleReturn(dynamicString, "DeleteTodo")
	if err != nil {
		log.Printf("Error occurred in DeleteTodo: %v\n", err)
		return nil, err
	}
	return removedCombined, nil
}

func DeleteUser(user_id int32) ([]*pb.Combined, error) {
	dynamicString := fmt.Sprintf(constants.DELETE_USER, user_id)

	removedCombined, err := queryRowWithMultipleReturn(dynamicString, "DeleteUser")
	if err != nil {
		log.Printf("Error occurred in DeleteUser: %v\n", err)
		return nil, err
	}
	return removedCombined, nil
}
