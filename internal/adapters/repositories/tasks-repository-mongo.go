package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *TasksRepositoryMongo) Create(ctx context.Context, name, description, userId string) (*entities.Task, error) {
	task := entities.NewTask(name, description, userId)
	result, err := r.db.Collection(TasksCollection).InsertOne(ctx, bson.M{"name": task.Name, "userId": task.UserId, "description": task.Description, "created": task.Created})
	if err != nil {
		return nil, err
	}
	task.ID = result.InsertedID.(primitive.ObjectID)
	return &task, nil
}

func (r *TasksRepositoryMongo) FindByUserId(ctx context.Context, userId string, limit int64, offset int64) ([]entities.Task, error) {
	var tasks []entities.Task
	cursor, err := r.db.Collection(TasksCollection).Find(ctx, bson.M{"userId": userId}, options.Find().SetSort(bson.D{{"created", -1}}).SetLimit(limit).SetSkip(offset))
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var task entities.Task
		err := cursor.Decode(&task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
