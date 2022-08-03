package courier

import (
	"context"

	"github.com/deltics/go-tasks"
)

// Courier[V] represents a task satisfied by a handler function accepting a
// single argument of type `T` and returning `error`.
type queue[V comparable] struct {
	fn       tasks.Courier[V]
	handler  tasks.Courier[V]
	listener *listener[V]
	request  chan *tasks.Request[V]
	result   map[V]chan error
}

// Queue returns a Courier[V] with the specified handler and a buffered
// request channel of capacity `n`.
func Queue[V comparable](handler tasks.Courier[V], n int) tasks.CourierQueue[V] {
	return &queue[V]{
		fn:      handler,
		handler: handler,
		request: make(chan *tasks.Request[V], n),
		result:  make(map[V]chan error),
	}
}

func (task *queue[V]) Use(fn tasks.Courier[V]) {
	task.handler = fn
}

func (task *queue[V]) UseDefault() {
	task.handler = task.fn
}

func (task *queue[V]) Enqueue(ctx context.Context, v V) chan error {
	select {
	case task.request <- &tasks.Request[V]{Ctx: ctx, Value: v}:
		task.result[v] = make(chan error)
		return task.result[v]
	default:
		return nil
	}
}

// StartListener starts a listener using the task handler.
func (task *queue[V]) StartListener() {
	if task.listener != nil {
		panic("Listener already started")
	}
	task.listener = &listener[V]{
		task: task,
	}
	task.listener.start()
}
