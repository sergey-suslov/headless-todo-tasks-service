package middlewares

import (
	"context"
	"github.com/go-kit/kit/log"
	"headless-todo-tasks-service/internal/entities"
	"headless-todo-tasks-service/internal/services"
	"time"
)

type LoggerMiddleware struct {
	Logger log.Logger
	Next   services.TasksService
}

func (l *LoggerMiddleware) Update(ctx context.Context, userId, taskId, name, description string) (err error) {
	defer func(begin time.Time) {
		_ = l.Logger.Log(
			"method", "Update",
			"name", name,
			"description", description,
			"userId", taskId,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return l.Next.Update(ctx, userId, taskId, name, description)
}

func (l *LoggerMiddleware) Create(ctx context.Context, name, description, userId string) (output *entities.Task, err error) {
	defer func(begin time.Time) {
		_ = l.Logger.Log(
			"method", "Create",
			"name", name,
			"description", description,
			"userId", userId,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return l.Next.Create(ctx, name, description, userId)
}

func (l *LoggerMiddleware) GetByUserId(ctx context.Context, userId string, limit, offset int64) (output []entities.Task, err error) {
	defer func(begin time.Time) {
		_ = l.Logger.Log(
			"method", "GetByUserId",
			"userId", userId,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return l.Next.GetByUserId(ctx, userId, limit, offset)
}
