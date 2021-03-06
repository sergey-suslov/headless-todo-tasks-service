package middlewares

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"headless-todo-tasks-service/internal/entities"
	"headless-todo-tasks-service/internal/services"
	"time"
)

type PrometheusMetrics struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	CountResult    metrics.Histogram
}

func NewPrometheusMetrics(requestCount metrics.Counter, requestLatency metrics.Histogram, countResult metrics.Histogram) *PrometheusMetrics {
	return &PrometheusMetrics{RequestCount: requestCount, RequestLatency: requestLatency, CountResult: countResult}
}

type InstrumentingMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	CountResult    metrics.Histogram
	Next           services.TasksService
}

func (mw *InstrumentingMiddleware) GetById(ctx context.Context, taskId string) (task *entities.Task, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetById", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Next.GetById(ctx, taskId)
}

func (mw *InstrumentingMiddleware) AddFile(ctx context.Context, taskId, fileId, fileName string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "AddFile", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Next.AddFile(ctx, taskId, fileId, fileName)
}

func (mw *InstrumentingMiddleware) Update(ctx context.Context, userId, taskId, name, description string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Update", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Next.Update(ctx, userId, taskId, name, description)
}

func (mw *InstrumentingMiddleware) Create(ctx context.Context, name, description, userId string) (output *entities.Task, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Create", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Next.Create(ctx, name, description, userId)
}

func (mw *InstrumentingMiddleware) GetByUserId(ctx context.Context, userId string, limit, offset int64) (output []entities.Task, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetByUserId", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Next.GetByUserId(ctx, userId, limit, offset)
}
