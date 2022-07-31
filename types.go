package mediator

import "context"

type ReceiverFunc[TData any] func(context.Context, TData) error
type HandlerFunc[TRequest any, TResult any] func(context.Context, TRequest) (TResult, error)
type ValidatorFunc[TInput any] func(context.Context, TInput) error

// Receiver[TData] is the interface to be implemented by a receiver
type Receiver[TData any] interface {
	Execute(context.Context, TData) error
}

// Handler[TRequest, TResult] is the interface to be implemented by a handler
type Handler[TRequest any, TResult any] interface {
	Execute(context.Context, TRequest) (TResult, error)
}

// Validator[TInput] is an optional interface that may be implemented
// by both receivers and handlers, to separate the validation of the data
// or request (the input) from the execution the receiver or handler itself.
type Validator[TInput any] interface {
	Validate(context.Context, TInput) error
}
