package repositories

import (
	"context"
	"errors"
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

func (r *TasksRepositoryMongo) Update(ctx context.Context, taskId string, name string, description string) error {
	id, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return errors.New("wrong id format")
	}
	one, err := r.db.Collection(TasksCollection).UpdateOne(ctx, bson.M{"_id": id}, bson.D{{"$set", bson.D{{"name", name}, {"description", description}}}})
	if err != nil {
		return err
	}
	if one.ModifiedCount == 0 {
		return errors.New("no records updated")
	}
	return nil
}

func (r *TasksRepositoryMongo) FindById(ctx context.Context, taskId string) (*entities.Task, error) {
	id, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return nil, errors.New("wrong id format")
	}
	var task entities.Task
	err = r.db.Collection(TasksCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err == mongo.ErrNoDocuments {
		return nil, repositories.NoTaskWithIdError
	}
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TasksRepositoryMongo) AddFile(ctx context.Context, userId, taskId, fileId, fileName string) error {
	task, err := r.FindById(ctx, taskId)
	if err != nil {
		return err
	}
	if task.UserId != userId {
		return repositories.NoTaskWithIdError
	}

	fileIdAsObjectId, err := primitive.ObjectIDFromHex(fileId)
	task.Files = append(task.Files, entities.File{
		ID:   fileIdAsObjectId,
		Name: fileName,
	})
	return r.Save(ctx, *task)
}

func (r *TasksRepositoryMongo) Save(ctx context.Context, task entities.Task) error {
	one, err := r.db.Collection(TasksCollection).UpdateOne(ctx, bson.M{"_id": task.ID}, bson.D{{"$set", task}})
	if err != nil {
		return err
	}
	if one.ModifiedCount == 0 {
		return errors.New("no records updated")
	}
	return nil
}
