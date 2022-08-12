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

	handler, _ := handlerReg.(CommandHandler[TRequest])

	// You may be thinking that we should test that the handler we found is
	// of the correct type, but the magic of generics and the strict type
	// system takes care of that for us, so there's no need.  \o/

	// If the handler also provides a request validator call that first
	if validator, ok := handlerReg.(RequestValidator[TRequest]); ok {
		err := validate(validator, ctx, request)
		if err != nil {
			return err
		}
	}

	return handler.Execute(ctx, request)
}
