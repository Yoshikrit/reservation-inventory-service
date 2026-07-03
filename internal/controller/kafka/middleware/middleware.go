package middleware

import (
	"context"

	kafka "github.com/segmentio/kafka-go"
)

type Handler func(ctx context.Context, msg kafka.Message) error

type Middleware func(Handler) Handler

func Chain(h Handler, middlewares ...Middleware) Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
