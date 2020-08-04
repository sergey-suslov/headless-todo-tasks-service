package services

import (
	"context"
	"errors"
	"headless-todo-tasks-service/internal/entities"
	"headless-todo-tasks-service/internal/services/repositories"
)

type TasksService interface {
	Create(ctx context.Context, name, description, userId string) (*entities.Task, error)
	GetByUserId(ctx context.Context, userId string, limit, offset int64) ([]entities.Task, error)
	Update(ctx context.Context, userId, taskId, name, description string) error
	AddFile(ctx context.Context, taskId, fileId, fileName string) error
	GetById(ctx context.Context, taskId string) (*entities.Task, error)
}

type tasksService struct {
	tasksRepository repositories.TasksRepository
}

func NewTasksService(tasksRepository repositories.TasksRepository) TasksService {
	return &tasksService{tasksRepository}
}

func (service *tasksService) Create(ctx context.Context, name, description, userId string) (*entities.Task, error) {
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

func (service *tasksService) GetByUserId(ctx context.Context, userId string, limit, offset int64) ([]entities.Task, error) {
	if userId == "" {
		return nil, errors.New("userId must be present")
	}
	return service.tasksRepository.FindByUserId(ctx, userId, limit, offset)
}

func (service *tasksService) Update(ctx context.Context, userId, taskId, name, description string) error {
	task, err := service.tasksRepository.FindById(ctx, taskId)
	if err != nil {
		return err
	}
	if task == nil || task.UserId != userId {
		return errors.New("no task with the given id")
	}

	return service.tasksRepository.Update(ctx, taskId, name, description)
}

func (service *tasksService) AddFile(ctx context.Context, taskId, fileId, fileName string) error {
	if taskId == "" {
		return errors.New("taskId must be present")
	}
	if fileId == "" {
		return errors.New("fileId must be present")
	}
	if fileName == "" {
		return errors.New("fileName must be present")
	}

	return service.tasksRepository.AddFile(ctx, taskId, fileId, fileName)
}

func (service *tasksService) GetById(ctx context.Context, taskId string) (*entities.Task, error) {
	if taskId == "" {
		return nil, errors.New("taskId must be present")
	}
	return service.tasksRepository.FindById(ctx, taskId)
}
