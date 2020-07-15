package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"headless-todo-tasks-service/internal/entities"
	"headless-todo-tasks-service/internal/services/repositories"
)

const TasksCollection = "tasks"

type TasksRepositoryMongo struct {
	db *mongo.Database
}

func NewTasksRepositoryMongo(db *mongo.Database) repositories.TasksRepository {
	return &TasksRepositoryMongo{db}
}

func (r *TasksRepositoryMongo) Create(ctx context.Context, task entities.Task) (*entities.Task, error) {
	result, err := r.db.Collection(TasksCollection).InsertOne(ctx, task)
	if err != nil {
		return nil, err
	}
	task.ID = result.InsertedID.(primitive.ObjectID)
	return &task, nil
}
