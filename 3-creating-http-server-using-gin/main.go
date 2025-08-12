package main

import (
	"errors"
	"net/http"
	"strings"

	"slices"

	"github.com/gin-gonic/gin"
)

// define the Todo struct
type Todo struct {
	ID        string `json:"id"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

// define a few todos
var TODOS = []Todo{
	{ID: "1", Item: "Buy groceries", Completed: false},
	{ID: "2", Item: "Read a book", Completed: false},
	{ID: "3", Item: "Write code", Completed: false},
	{ID: "4", Item: "Go for a walk", Completed: false},
}

// containsIgnoreCase is a helper function that checks if a string contains another string
func containsIgnoreCase(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}

// getTodoById is a helper function that retrieves a todo by its ID
// it takes an ID as a string and returns the corresponding Todo or an error if not found
func getTodoById(id string) (*Todo, error) {
	for i := range TODOS {
		if TODOS[i].ID == id {
			return &TODOS[i], nil
		}
	}
	return nil, errors.New("Todo not found")
}

// getTodos handles the GET request to retrieve all todos
// it filters todos based on query parameters: ?query= and ?status=
// it returns the filtered todos in the response body with a 200 OK status
func getTodos(context *gin.Context) {
	query := context.Query("query")
	status := context.Query("status")

	var filtered []Todo

	for _, todo := range TODOS {
		// Match query if provided
		if query != "" && !containsIgnoreCase(todo.Item, query) {
			continue
		}

		// Match status if provided
		if status == "true" && !todo.Completed {
			continue
		}
		if status == "false" && todo.Completed {
			continue
		}

		filtered = append(filtered, todo)
	}

	context.IndentedJSON(http.StatusOK, filtered)
}

// addTodo handles the POST request to add a new todo
// it takes context as an argument, context holds the info about the HTTP request
// it binds the JSON from the request body to a new Todo struct and appends it to the TODOS slice
// it returns a 201 Created status with the new todo in the response body
// if there is an error binding the JSON, it simply returns without doing anything
func addTodo(context *gin.Context) {
	var newTodo Todo

	if err := context.BindJSON(&newTodo); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	TODOS = append(TODOS, newTodo)

	context.IndentedJSON(http.StatusCreated, newTodo)
}

// getTodo handles the GET request to retrieve a specific todo by its ID
// it retrieves the ID from the URL parameters, calls getTodoById to find the todo,
// and returns it in the response body with a 200 OK status
// if the todo is not found, it returns a 404 Not Found status with an error
func getTodo(context *gin.Context) {
	id := context.Param("id")
	todo, err := getTodoById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
		return
	}

	context.IndentedJSON(http.StatusOK, todo)
}

// updateTodoStatus handles the PATCH request to toggle the completed status of a todo
// it retrieves the ID from the URL parameters, calls getTodoById to find the todo,
// toggles its Completed status, and returns the updated todo in the response body with a
// 200 OK status
// if the todo is not found, it returns a 404 Not Found status with an error message
func updateTodoStatus(context *gin.Context) {
	id := context.Param("id")
	todo, err := getTodoById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
		return
	}

	todo.Completed = !todo.Completed // toggle the completed status

	context.IndentedJSON(http.StatusOK, todo)
}

// deleteTodo handles the DELETE request to remove a todo by its ID
// it retrieves the ID from the URL parameters, searches for the todo in the TODOS slice,
// and removes it if found using the slices.Delete function from the slices package
func deleteTodo(context *gin.Context) {
	id := context.Param("id")
	for i, todo := range TODOS {
		if todo.ID == id {
			TODOS = slices.Delete(TODOS, i, i+1)
			context.IndentedJSON(http.StatusOK, gin.H{"message": "Todo deleted"})
			return
		}
	}
	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
}

// main function sets up the Gin router and defines the routes for the todo API
func main() {
	router := gin.Default()
	router.GET("/todos", getTodos)
	router.POST("/todos", addTodo)
	router.GET("/todos/:id", getTodo)
	router.PATCH("/todos/:id", updateTodoStatus)
	router.DELETE("/todos/:id", deleteTodo)
	router.Run("localhost:8080")
}
