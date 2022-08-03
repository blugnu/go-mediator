package tasks

import "context"

type Request[V any] struct {
	Ctx   context.Context
	Value V
}

type Courier[V comparable] func(context.Context, V) error
type Gofer[V comparable, R any] func(context.Context, V) (R, error)

type CourierProxy[V comparable] interface {
	CallWith(context.Context, V) error
}

type GoferProxy[V comparable, R any] interface {
	CallWith(context.Context, V) (R, error)
}

type CourierQueue[V comparable] interface {
	Enqueue(context.Context, V) chan error
	StartListener()
}

type GoferQueue[V comparable, R any] interface {
	Enqueue(context.Context, V) (chan R, chan error)
	StartListener()
}
