package main

import (
	kitlog "github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"headless-todo-tasks-service/internal/adapters/repositories"
	"headless-todo-tasks-service/internal/services"
	"log"
	"os"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Init(client *mongo.Client) *dig.Container {
	c := dig.New()

	err := c.Provide(func() *mongo.Client {
		return client
	})
	handleError(err)

	err = c.Provide(func() *mongo.Database {
		return client.Database("tasks")
	})
	handleError(err)

	err = c.Provide(repositories.NewTasksRepositoryMongo)
	handleError(err)

	err = c.Provide(services.NewTasksService)
	handleError(err)

	err = c.Provide(func() kitlog.Logger {
		logger := kitlog.NewLogfmtLogger(os.Stderr)
		return logger
	})
	handleError(err)

	return c
}
