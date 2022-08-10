package mediator

import (
	"context"
	"fmt"
	"reflect"
)

type QueryHandler[TRequest any, TResult any] interface {
	Execute(context.Context, TRequest) (TResult, error)
}

// RegisterQueryHandler registers the specified handler for a particular request type
// and returning a particular result type.
//
// If a handler is already registered for that request type the function will panic,
// otherwise the handler is registered.
func RegisterQueryHandler[TRequest any, TResult any](handler QueryHandler[TRequest, TResult]) *reg {
	var rq TRequest
	rqt := reflect.TypeOf(rq)

	_, exists := queryHandlers[rqt]
	if exists {
		panic(fmt.Sprintf("handler already registered for %T", rq))
	}

	queryHandlers[rqt] = handler

	return &reg{
		handlers: queryHandlers,
		rqt:      rqt,
	}
}

// Perform sends the specified request and context to the registered Query
// handler and returns the result and error from that handler.
func Query[TRequest any, TResult any](ctx context.Context, request TRequest) (TResult, error) {
	rqt := reflect.TypeOf(request)

	handlerReg, ok := queryHandlers[rqt]
	if !ok {
		return *new(TResult), &ErrNoHandler{request: request}
	}

	handler, ok := handlerReg.(QueryHandler[TRequest, TResult])
	if !ok {
		return *new(TResult), &ErrInvalidHandler{request: request, handler: handler}
	}

	response, err := handler.Execute(ctx, request)

	return response, err
}
