package endpoints

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	"github.com/nats-io/stan.go"
	"go.uber.org/dig"
	"headless-todo-tasks-service/internal/adapters/middlewares"
	"headless-todo-tasks-service/internal/services"
	"log"
	"time"
)

const FileAddedSubjectName = "tasks.files.added"
const FileAddedQueueGroup = "tasks.files.added.group"

type addFileRequest struct {
	TaskId   string `json:"taskId"`
	FileId   string `json:"fileId"`
	FileName string `json:"fileName"`
}

func makeAddFileEndpoint(service services.TasksService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addFileRequest)
		err := service.AddFile(ctx, req.TaskId, req.FileId, req.FileName)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func CreateAddFileHandler(c *dig.Container) {
	var sc stan.Conn
	err := c.Invoke(func(s stan.Conn) {
		sc = s
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
	ep := makeAddFileEndpoint(service)

	_, err = sc.QueueSubscribe(FileAddedSubjectName, FileAddedQueueGroup, func(msg *stan.Msg) {
		var addFileRequest addFileRequest
		err := json.Unmarshal(msg.Data, &addFileRequest)
		if err != nil {
			return
		}

		ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFunc()
		_, err = ep(ctx, addFileRequest)
		if err != nil {
			return
		}

		_ = msg.Ack()
	}, stan.AckWait(8*time.Second))
	if err != nil {
		log.Fatal(err)
	}

}
