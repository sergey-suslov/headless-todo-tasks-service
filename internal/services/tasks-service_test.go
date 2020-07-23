package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"headless-todo-tasks-service/internal/entities"
	"testing"
)

type MockedTasksRepository struct {
	mock.Mock
}

func (m *MockedTasksRepository) Create(ctx context.Context, name, description, userId string) (*entities.Task, error) {
	args := m.Called(ctx, name, description, userId)
	return args.Get(0).(*entities.Task), nil
}

func TestTasksService_Create_Success(t *testing.T) {
	mockedTaskRepository := new(MockedTasksRepository)
	testTask := &entities.Task{
		ID:          primitive.ObjectID{1},
		Name:        "task 1",
		Description: "description 1",
		UserId:      "2",
	}
	mockedTaskRepository.On("Create", context.Background(), testTask.Name, testTask.Description, testTask.UserId).Return(testTask)
	service := &TasksService{mockedTaskRepository}

	task, _ := service.Create(context.Background(), testTask.Name, testTask.Description, testTask.UserId)

	mockedTaskRepository.AssertCalled(t, "Create", context.Background(), testTask.Name, testTask.Description, testTask.UserId)
	assert.Equal(t, task, testTask)
}
