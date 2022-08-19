package mediator

import (
	"context"
	"fmt"
	"reflect"
)

// RegisterHandler registers a handler for the specified request type
// returning the specified result type.
//
// If a handler is already registered for the request type, the
// function will panic, otherwise the handler is registered.
func RegisterHandler[TRequest any, TResult any](handler Handler[TRequest, TResult]) *reg {
	dummyrequest := *new(TRequest)
	requesttype := reflect.TypeOf(dummyrequest)

	_, exists := handlers[requesttype]
	if exists {
		panic(fmt.Sprintf("handler already registered for %T", dummyrequest))
	}

	handlers[requesttype] = handler

	return &reg{
		registry:       handlers,
		registeredtype: requesttype,
	}
}

// Perform sends the specified request and context to the registered handler
// for the request type and returns the result and error from that handler.
//
// If the handler implements Validator and the validator returns an error,
// then handler is not called and the error returned by Perform will be a
// ValidationError, wrapping the error returned by the validator.
func Perform[TRequest any, TResult any](ctx context.Context, request TRequest) (TResult, error) {
	requesttype := reflect.TypeOf(request)
	zeroresult := *new(TResult)

	reg, ok := handlers[requesttype]
	if !ok {
		return zeroresult, &NoReceiverError{data: request}
	}

	handler, ok := reg.(Handler[TRequest, TResult])
	if !ok {
		return zeroresult, &InvalidHandlerError{handler: handler, request: request, result: zeroresult}
	}

	// If the handler implements validator, call that first
	// and return any error
	if validator, ok := reg.(Validator[TRequest]); ok {
		err := validate(validator, ctx, request)
		if err != nil {
			return zeroresult, err
		}
	}

	response, err := handler.Execute(ctx, request)

	return response, err
}
