package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"headless-todo-tasks-service/internal/adapters/endpoints"
	"log"
	"net/http"
)

func main() {
	initConfig()

	client, closeConnection := ConnectMongo()
	nc, sc, closeNats := ConnectNats()
	defer func() {
		closeConnection()
		closeNats()
	}()

	c := Init(client, nc, sc)

	endpoints.CreateAddFileHandler(c)
	http.Handle("/create-task", endpoints.CreateTaskHandler(c))
	http.Handle("/get-tasks", endpoints.GetTasksHandler(c))
	http.Handle("/update", endpoints.UpdateTaskHandler(c))
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
