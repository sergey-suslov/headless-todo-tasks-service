package repositories

import (
	"context"
	"errors"
	"headless-todo-tasks-service/internal/entities"
)

var NoTaskWithIdError = errors.New("no task with this id")

type TasksRepository interface {
	Create(context.Context, string, string, string) (*entities.Task, error)
	FindById(context.Context, string) (*entities.Task, error)
	FindByUserId(context.Context, string, int64, int64) ([]entities.Task, error)
	Update(context.Context, string, string, string) error
	AddFile(ctx context.Context, userId, taskId, fileId, fileName string) error
	Save(ctx context.Context, task entities.Task) error
}
