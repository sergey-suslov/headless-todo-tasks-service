package endpoints

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"go.uber.org/dig"
	"headless-todo-tasks-service/internal/services"
	"log"
	"net/http"
)

type createTaskRequest struct {
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
	var request createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
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

	return httptransport.NewServer(
		makeCreateTaskEndpoint(service),
		decodeCreateRequest,
		encodeCreateRequest,
	)
}
