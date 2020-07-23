package services

import (
	"context"
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
	return service.tasksRepository.Create(ctx, name, description, userId)
}
