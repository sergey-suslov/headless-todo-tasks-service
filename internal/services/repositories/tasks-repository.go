package repositories

import (
	"context"
	"headless-todo-tasks-service/internal/entities"
)

type TasksRepository interface {
	Create(context.Context, entities.Task) (*entities.Task, error)
}
