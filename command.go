package mediator

import (
	"context"
	"fmt"
	"reflect"
)

// RegisterCommandHandler registers the specified handler for a particular request type.
//
// If a handler is already registered for that type the function will panic, otherwise
// the handler is registered.
func RegisterCommandHandler[TRequest any](handler CommandHandler[TRequest]) *reg {
	var rq TRequest
	rqt := reflect.TypeOf(rq)

	_, exists := commandHandlers[rqt]
	if exists {
		panic(fmt.Sprintf("handler already registered for %T", rq))
	}

	commandHandlers[rqt] = handler

	return &reg{
		handlers: commandHandlers,
		rqt:      rqt,
	}
}

// Perform sends the specified request and context to the registered Perform
// handler for the request type.   If the Command handler implements a
// RequestValidator, the Command is only executed if the request passes
// validation.
func Perform[TRequest any](ctx context.Context, request TRequest) error {
	rqt := reflect.TypeOf(request)

	handlerReg, ok := commandHandlers[rqt]
	if !ok {
		return &ErrNoHandler{request: request}
	}

	handler, ok := handlerReg.(CommandHandler[TRequest])
	if !ok {
		// This shouldn't actually be possible as the type system should
		// prevent the registration of a handler of the wrong type.
		// But just in case...
		return &ErrInvalidHandler{handler: handler, request: request}
	}

	// If the handler also provides a request validator, call that first
	// and return any error in an ErrBadRequest.
	if validator, ok := handlerReg.(RequestValidator[TRequest]); ok {
		return validate(validator, ctx, request)
	}

	return handler.Execute(ctx, request)
}
