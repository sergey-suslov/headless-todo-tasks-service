package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Task struct {
	ID          primitive.ObjectID
	Name        string
	Description string
	Created     primitive.Timestamp
}

func NewTask(name, description string) Task {
	return Task{Name: name, Description: description, Created: primitive.Timestamp{T: uint32(time.Now().Unix())}}
}
