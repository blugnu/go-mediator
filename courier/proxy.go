package courier

import (
	"context"

	"github.com/deltics/go-tasks"
)

// Proxy[V] represents a task satisfied by a handler function accepting a
// single argument of type `T` and returning `error`.
type proxy[V comparable] struct {
	fn      tasks.Courier[V]
	handler tasks.Courier[V]
}

// Unbuffered returns a Courier[V] with the specified handler and an
// unbuffered request channel.
func Proxy[V comparable](h tasks.Courier[V]) tasks.CourierProxy[V] {
	return &proxy[V]{
		fn:      h,
		handler: h,
	}
}

func (task *proxy[V]) CallWith(ctx context.Context, v V) error {
	return task.handler(ctx, v)
}

func (task *proxy[V]) Use(fn tasks.Courier[V]) {
	task.handler = fn
}

func (task *proxy[V]) UseDefault() {
	if task.fn == nil {
		panic("Proxy has no default function")
	}
	task.handler = task.fn
}
