package endpoints

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kitnats "github.com/go-kit/kit/transport/nats"
	"github.com/nats-io/nats.go"
	"go.uber.org/dig"
	"headless-todo-tasks-service/internal/adapters/middlewares"
	"headless-todo-tasks-service/internal/services"
	"log"
)

const GetTaskById = "tasks.getById"

type getTaskByIdRequest struct {
	TaskId string `json:"taskId"`
}

func makeGetTaskByIdEndpoint(service services.TasksService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getTaskByIdRequest)
		task, err := service.GetById(ctx, req.TaskId)
		if err != nil {
			return nil, err
		}
		return task, nil
	}
}

func CreateGetTaskByIdHandler(c *dig.Container) *nats.Subscription {
	var nc *nats.Conn
	err := c.Invoke(func(n *nats.Conn) {
		nc = n
	})
	if err != nil {
		log.Fatal(err)
	}

	var service services.TasksService
	err = c.Invoke(func(s services.TasksService) {
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
	ep := makeGetTaskByIdEndpoint(service)

	subscriber := kitnats.NewSubscriber(ep, func(ctx context.Context, msg *nats.Msg) (request interface{}, err error) {
		var getTaskByIdRequest getTaskByIdRequest
		marshErr := json.Unmarshal(msg.Data, &getTaskByIdRequest)
		if marshErr != nil {
			return nil, err
		}
		return getTaskByIdRequest, nil
	}, kitnats.EncodeJSONResponse)

	subscription, err := nc.QueueSubscribe(GetTaskById, CommonNatsGroup, subscriber.ServeMsg(nc))
	if err != nil {
		log.Fatal(err)
	}
	return subscription
}
