package mediator

import "fmt"

// ErrNoHandler is returned by Perform or Query if there is no handler registered
// for the request being made.
type ErrNoHandler struct {
	handler handlerType
	request interface{}
}

func (e *ErrNoHandler) Error() string {
	return fmt.Sprintf("no %s handler for '%T'",
		e.handler.Name(),
		e.request,
	)
}

// ErrInvalidHandler is returned by Perform or Query if the registered handler for
// a request is of the wrong type (which should be impossible!).
type ErrInvalidHandler struct {
	handlerType handlerType
	handler     interface{}
	request     interface{}
}

func (e *ErrInvalidHandler) Error() string {
	return fmt.Sprintf("%T is not a valid %s handler for '%T' requests",
		e.handler,
		e.handlerType.Name(),
		e.request,
	)
}

// ErrBadRequest is returned by Perform or Query if the handler implements a
// validator that has returned an error.
type ErrBadRequest struct {
	err error
}

func (e *ErrBadRequest) Error() string {
	return fmt.Sprintf("bad request: %v", e.err)
}

func (e *ErrBadRequest) InnerError() error {
	return e.err
}
