package services

import (
	"context"
	"errors"
	"headless-todo-tasks-service/internal/entities"
	"headless-todo-tasks-service/internal/services/repositories"
)

type TasksService struct {
	tasksRepository repositories.TasksRepository
}

func NewTasksService(tasksRepository repositories.TasksRepository) *TasksService {
	return &TasksService{tasksRepository}
}

func (service *TasksService) Create(ctx context.Context, name, description, userId string) (*entities.Task, error) {
	if name == "" {
		return nil, errors.New("name must be present")
	}
	if description == "" {
		return nil, errors.New("description must be present")
	}
	if userId == "" {
		return nil, errors.New("userId must be present")
	}
	return service.tasksRepository.Create(ctx, name, description, userId)
}

func (service *TasksService) GetByUserId(ctx context.Context, userId string, limit, offset int64) ([]entities.Task, error) {
	if userId == "" {
		return nil, errors.New("userId must be present")
	}
	return service.tasksRepository.FindByUserId(ctx, userId, limit, offset)
}
