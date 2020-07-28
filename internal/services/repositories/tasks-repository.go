package repositories

import (
	"context"
	"headless-todo-tasks-service/internal/entities"
)

type TasksRepository interface {
	Create(context.Context, string, string, string) (*entities.Task, error)
	FindById(context.Context, string) (*entities.Task, error)
	FindByUserId(context.Context, string, int64, int64) ([]entities.Task, error)
	Update(context.Context, string, string, string) error
}
