package endpoints

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"go.uber.org/dig"
	"headless-todo-tasks-service/internal/adapters/middlewares"
	"headless-todo-tasks-service/internal/services"
	"log"
	"net/http"
)

type getTasksRequest struct {
	UserClaim
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

func (g *getTasksRequest) SetUserClaim(claim UserClaim) {
	g.UserClaim = claim
}

func makeGetTasksEndpoint(service services.TasksService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*getTasksRequest)
		task, err := service.GetByUserId(ctx, req.ID, req.Limit, req.Offset)
		if err != nil {
			return nil, err
		}
		return task, nil
	}
}

func GetTasksHandler(c *dig.Container) http.Handler {
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
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here
	service = &middlewares.InstrumentingMiddleware{RequestCount: requestCount, RequestLatency: requestLatency, CountResult: countResult, Next: service}
	taskEndpoint := makeGetTasksEndpoint(service)

	return httptransport.NewServer(
		taskEndpoint,
		DefaultRequestDecoder(func(r *http.Request) (
			UserClaimable, error) {
			var request getTasksRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				return nil, err
			}
			return &request, nil
		}),
		DefaultRequestEncoder,
	)
}
