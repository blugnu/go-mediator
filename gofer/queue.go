package gofer

import (
	"context"

	"github.com/deltics/go-tasks"
)

type result[R any] struct {
	val chan R
	err chan error
}

type queue[V comparable, R any] struct {
	fn       tasks.Gofer[V, R]
	handler  tasks.Gofer[V, R]
	listener *listener[V, R]
	request  chan *tasks.Request[V]
	result   map[V]*result[R]
}

func Queue[V comparable, R any](handler tasks.Gofer[V, R], n int) *queue[V, R] {
	return &queue[V, R]{
		fn:      handler,
		handler: handler,
		request: make(chan *tasks.Request[V], n),
		result:  make(map[V]*result[R]),
	}
}

func (task *queue[V, R]) Enqueue(ctx context.Context, v V) (chan R, chan error) {
	select {
	case task.request <- &tasks.Request[V]{Ctx: ctx, Value: v}:
		ch := &result[R]{
			val: make(chan R),
			err: make(chan error),
		}
		task.result[v] = ch
		return ch.val, ch.err
	default:
		return nil, nil
	}
}

func (task *queue[V, R]) StartListener() {
	if task.listener != nil {
		panic("Listener already started")
	}
	task.listener = &listener[V, R]{
		task: task,
	}
	task.listener.start()
}
