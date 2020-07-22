package endpoints

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
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

func makeCreateTaskEndpoint(service *services.TasksService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createTaskRequest)
		task, err := service.Create(ctx, req.Name, req.Description)
		if err != nil {
			return nil, err
		}
		return task, nil
	}
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userClaim, err := GetUserClaimFromRequest(r)
	if err != nil {
		return nil, err
	}

	var request createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	request.UserClaim = *userClaim
	return request, nil
}

func encodeCreateRequest(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func CreateTaskHandler(c *dig.Container) http.Handler {
	var service *services.TasksService
	err := c.Invoke(func(s *services.TasksService) {
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

	taskEndpoint := makeCreateTaskEndpoint(service)
	loggerMiddleware := middlewares.LoggerMiddleware(kitlog.With(logger, "method", "create-task"))
	taskEndpoint = loggerMiddleware(taskEndpoint)

	return httptransport.NewServer(
		taskEndpoint,
		decodeCreateRequest,
		encodeCreateRequest,
	)
}
