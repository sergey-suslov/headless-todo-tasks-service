package main

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"headless-todo-tasks-service/internal/adapters/repositories"
	"headless-todo-tasks-service/internal/services"
	"log"
)

func Init(client *mongo.Client) *dig.Container {
	c := dig.New()

	err := c.Provide(func() *mongo.Client {
		return client
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.Provide(func() *mongo.Database {
		return client.Database("tasks")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.Provide(repositories.NewTasksRepositoryMongo)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Provide(services.NewTasksService)
	if err != nil {
		log.Fatal(err)
	}
	return c
}
