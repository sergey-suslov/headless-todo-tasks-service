package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func ConnectMongo() (*mongo.Client, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://user:password@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	return client, func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}
}
