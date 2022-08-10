package mediator

import "fmt"

// ErrNoHandler is returned by Perform or Query if there is no handler registered
// for the request being made.
type ErrNoHandler struct {
	request interface{}
}

func (e *ErrNoHandler) Error() string {
	return fmt.Sprintf("no handler for %T", e.request)
}

// ErrInvalidHandler is returned by Perform or Query if the registered handeler for
// a request is of the wrong type (which should be impossible!).
type ErrInvalidHandler struct {
	handler interface{}
	request interface{}
}

func (e *ErrInvalidHandler) Error() string {
	return fmt.Sprintf("%T is not a valid handler of %T", e.handler, e.request)
}
