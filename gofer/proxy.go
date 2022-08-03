package gofer

import (
	"context"

	"github.com/deltics/go-tasks"
)

type proxy[V comparable, R any] struct {
	fn      tasks.Gofer[V, R]
	handler tasks.Gofer[V, R]
}

func Proxy[V comparable, R any](handler tasks.Gofer[V, R]) *proxy[V, R] {
	return &proxy[V, R]{
		fn:      handler,
		handler: handler,
	}
}

func (task *proxy[V, R]) CallWith(ctx context.Context, v V) (R, error) {
	return task.handler(ctx, v)
}
