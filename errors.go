package mediator

import (
	"fmt"
)

// NoReceiverError is returned by Perform if there is no handler
// registered for the request type being performed.
type NoHandlerError struct {
	request interface{}
}

func (e NoHandlerError) Error() string {
	return fmt.Sprintf("no handler for '%T'", e.request)
}

// NoReceiverError is returned by Send if there is no receiver registered
// for the data type being sent.
type NoReceiverError struct {
	data interface{}
}

func (e NoReceiverError) Error() string {
	return fmt.Sprintf("no receiver for '%T'", e.data)
}

// InvalidHandlerError is returned by Perform if the registered
// handler for the specified request type does not return then
// specified result type.
type InvalidHandlerError struct {
	handler interface{}
	request interface{}
	result  interface{}
}

func (e InvalidHandlerError) Error() string {
	return fmt.Sprintf("handler for %T (%T) does not return %T", e.request, e.handler, e.result)
}

// ValidationError is returned by Perform or Send if the handler or
// receiver implements a validator that has returned an error.
//
// ValidationError should also be used for "bad request" type errors
// by any receiver or handler that does not implement Validator but
// performs such validation in the receiver or handler Execute
// function itself.
type ValidationError struct {
	error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error: %v", e.error)
}

func (e ValidationError) Unwrap() error {
	return e.error
}
