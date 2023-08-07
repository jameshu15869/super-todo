package dbclient

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	pb "supertodo/todo/pb"
	"supertodo/todo/pkg/constants"
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

func todoFromRow(rows *sql.Rows) (*pb.Todo, error) {
	var id int32
	var title string
	var todo_date string
	var body string
	if err := rows.Scan(&id, &title, &todo_date, &body); err != nil {
		log.Printf("Error occurred in Scanning Rows: %v\n", err)
		return nil, err
	}
	return &pb.Todo{Id: id, Title: title, TodoDate: todo_date, Body: body}, nil
}

func queryRowWithSingleReturn(queryString string, functionName string, args ...interface{}) (*pb.Todo, error) {
	var insertedId int
	var insertedTitle string
	var insertedTodoDate string
	var insertedBody string
	if err := db.QueryRow(queryString, args...).Scan(&insertedId, &insertedTitle, &insertedTodoDate, &insertedBody); err != nil {
		log.Printf("Error occurred in %s: %v\n", functionName, err)
		return nil, err
	}
	return &pb.Todo{Id: int32(insertedId), Title: insertedTitle, TodoDate: insertedTodoDate, Body: insertedBody}, nil
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

func QueryAllTodos() ([]*pb.Todo, error) {
	queryString := constants.QUERY_ALL_TODOS

	cacheVal, err := CacheQuery[[]*pb.Todo](queryString)
	if err != nil {
		// Value not found in cache
		log.Println("QueryAllTodos: Value not found in cache, querying DB")
		rows, err := db.Query(queryString)
		if err != nil {
			log.Printf("Error occurred in QueryAllTodos: %v\n", err)
			return nil, err
		}
		defer rows.Close()

		var todos []*pb.Todo
		for rows.Next() {
			currentTodo, err := todoFromRow(rows)
			if err != nil {
				log.Printf("Error occurred in QueryAllTodos: Scanning Rows: %v\n", err)
				continue
			}
			todos = append(todos, currentTodo)
		}
		CacheSet(queryString, todos)
		return todos, nil
	}

	log.Println("QueryAllTodos: Cache: HIT")
	return cacheVal, nil
}

func QueryTodoById(id int32) (*pb.Todo, error) {
	// queryString := `SELECT * FROM users WHERE id=$1 LIMIT 1;`
	cacheKey := fmt.Sprintf(constants.QUERY_TODO_BY_ID_CACHE_KEY, id)

	cacheVal, err := CacheQuery[*pb.Todo](cacheKey)

	if err != nil {
		// Value not found in cache
		log.Println("QueryTodoById: Value not found in cache, querying DB")
		rows, err := db.Query(constants.QUERY_TODO_BY_ID, id)
		if err != nil {
			log.Printf("Error occurred in QueryTodoById: %v\n", err)
			return nil, err
		}
		defer rows.Close()

		var todo *pb.Todo
		for rows.Next() {
			currentTodo, err := todoFromRow(rows)
			if err != nil {
				log.Printf("Error occurred in QueryTodoById: Scanning Rows: %v\n", err)
				continue
			}
			todo = currentTodo
			break
		}
		CacheSet(cacheKey, todo)
		return todo, nil
	}

	log.Println("QueryTodoById: Cache: HIT")
	return cacheVal, nil
}

func QueryTodosByMultipleIds(todoIds []int32) ([]*pb.Todo, error) {
	setString := generateMultipleIdSetString(todoIds)
	var queryString string
	if len(todoIds) > 0 {
		queryString = fmt.Sprintf(constants.QUERY_TODOS_BY_MULTIPLE_IDS, setString)
	} else {
		return []*pb.Todo{}, nil
	}
	log.Println("QueryMultiple String: ", queryString)

	cacheVal, err := CacheQuery[[]*pb.Todo](queryString)
	if err != nil {
		// Value not found in cache
		log.Println("QueryTodosByMultipleIds: Value not found in cache, querying DB")
		rows, err := db.Query(queryString)
		if err != nil {
			log.Printf("Error occurred in QueryTodosByMultipleIds: %v\n", err)
			return nil, err
		}
		defer rows.Close()

		var todos []*pb.Todo
		for rows.Next() {
			currentTodo, err := todoFromRow(rows)
			if err != nil {
				log.Printf("Error occurred in QueryTodosByMultipleIds: Scanning Rows: %v\n", err)
				continue
			}
			todos = append(todos, currentTodo)
		}
		CacheSet(queryString, todos)
		return todos, nil
	}

	log.Println("QueryTodosByMultipleIds: Cache: HIT")
	return cacheVal, nil
}

func AddTodo(todo *pb.Todo) (*pb.Todo, error) {
	updatedTodo, err := queryRowWithSingleReturn(constants.ADD_TODO, "AddTodo", todo.Title, todo.TodoDate, todo.Body)
	if err != nil {
		log.Printf("Error occurred in AddTodo: %v\n", err)
		return nil, err
	}
	return updatedTodo, nil
}

func UpdateTodoById(todo *pb.Todo) (*pb.Todo, error) {
	updatedTodo, err := queryRowWithSingleReturn(constants.UPDATE_TODO_BY_ID, "UpdateTodoById", todo.Title, todo.TodoDate, todo.Body, todo.Id)
	if err != nil {
		log.Printf("Error occurred in UpdateTodoById: %v\n", err)
		return nil, err
	}
	return updatedTodo, nil
}

func DeleteTodo(todoId int32) (*pb.Todo, error) {
	updatedTodo, err := queryRowWithSingleReturn(constants.DELETE_TODO, "DeleteTodo", todoId)
	if err != nil {
		log.Printf("Error occurred in DeleteTodo: %v\n", err)
		return nil, err
	}
	return updatedTodo, nil
}
