package main

import (
	"headless-todo-tasks-service/internal/adapters/endpoints"
	"log"
	"net/http"
)

func main() {

	client, closeConnection := ConnectMongo()
	defer func() { closeConnection() }()

	c := Init(client)

	http.Handle("/create-task", endpoints.CreateTaskHandler(c))
	http.Handle("/get-tasks", endpoints.GetTasksHandler(c))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
