package mediator

import (
	"context"
	"fmt"
	"reflect"
)

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

// Query sends the specified request and context to the registered Query
// handler and returns the result and error from that handler.  If the
// Query handler implements a RequestValidator, the Query is only executed
// if the request passes validation.
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

	// If the handler also provides a request validator, call that first
	// and return any error in an ErrBadRequest.
	if validator, ok := handlerReg.(RequestValidator[TRequest]); ok {
		err := validate(validator, ctx, request)
		if err != nil {
			return *new(TResult), err
		}
	}

	response, err := handler.Execute(ctx, request)

	return response, err
}
