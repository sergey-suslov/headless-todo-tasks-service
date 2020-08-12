package endpoints

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/sony/gobreaker"
	"go.uber.org/dig"
	"headless-todo-tasks-service/internal/adapters/middlewares"
	"headless-todo-tasks-service/internal/services"
	"log"
	"net/http"
)

type createTaskRequest struct {
	UserClaim
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *createTaskRequest) SetUserClaim(claim UserClaim) {
	c.UserClaim = claim
}

func makeCreateTaskEndpoint(service services.TasksService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*createTaskRequest)
		task, err := service.Create(ctx, req.Name, req.Description, req.UserClaim.ID)
		if err != nil {
			return nil, err
		}
		return task, nil
	}
}

func CreateTaskHandler(c *dig.Container) http.Handler {
	var service services.TasksService
	err := c.Invoke(func(s services.TasksService) {
		service = s
	})
	if err != nil {
		log.Fatal(err)
	}

	var logger kitlog.Logger
	err = c.Invoke(func(log kitlog.Logger) {
		logger = log
	})
	if err != nil {
		log.Fatal(err)
	}

	service = &middlewares.LoggerMiddleware{Logger: kitlog.With(logger), Next: service}
	taskEndpoint := makeCreateTaskEndpoint(service)
	breaker := circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "Create-Task",
		MaxRequests: 100,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	}))
	taskEndpoint = breaker(taskEndpoint)

	return httptransport.NewServer(
		taskEndpoint,
		DefaultRequestDecoder(func(r *http.Request) (UserClaimable, error) {
			var request createTaskRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				return nil, err
			}
			return &request, nil
		}),
		DefaultRequestEncoder,
	)
}
