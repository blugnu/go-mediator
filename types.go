package mediator

import "context"

type handlerType int

const (
	command handlerType = iota
	query
)

func (t handlerType) Name() string {
	switch t {
	case command:
		return "command"
	case query:
		return "query"
	default:
		return "<undefined>"
	}
}

// CommandHandler[TRequest] is the interface to be implemented by Command handlers
type CommandHandler[TRequest any] interface {
	Execute(context.Context, TRequest) error
}

// QueryHandler[TRequest, TResult] is the interface to be implemented by Query handlers
type QueryHandler[TRequest any, TResult any] interface {
	Execute(context.Context, TRequest) (TResult, error)
}

// RequestValidator[TRequest] is an optional interface that may be implemented
// by Command and Query handlers, to separate the concerns of validating a
// request from those of executing or fulfilling the request.
type RequestValidator[TRequest any] interface {
	Validate(context.Context, TRequest) error
}
