package middlewares

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

func LoggerMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			_ = logger.Log("msg", "calling endpoint with request")
			defer func() {
				_ = logger.Log("msg", "called endpoint")
				err := recover()
				if err != nil {
					_ = logger.Log("unhandled error", err)
				}
			}()
			return next(ctx, request)
		}
	}
}
