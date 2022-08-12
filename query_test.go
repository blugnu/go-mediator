package mediator

import (
	"context"
	"errors"
	"testing"
)

const qryHandlerResultValue = "result!"

type qryRequest struct{}
type qryHandler struct{}
type qryDuplicate struct{}

func (*qryDuplicate) Execute(context.Context, qryRequest) (string, error) { return "", nil }
func (*qryHandler) Execute(context.Context, qryRequest) (string, error) {
	return qryHandlerResultValue, nil
}

type qryRequestHandlerWithValidatorReturningError struct{}

func (*qryRequestHandlerWithValidatorReturningError) Execute(context.Context, qryRequest) (string, error) {
	return "ok", nil
}
func (*qryRequestHandlerWithValidatorReturningError) Validate(context.Context, qryRequest) error {
	return errors.New("validation failed")
}

type qryRequestHandlerWithValidatorReturningErrBadRequest struct{}

func (*qryRequestHandlerWithValidatorReturningErrBadRequest) Execute(context.Context, qryRequest) (string, error) {
	return "ok", nil
}
func (*qryRequestHandlerWithValidatorReturningErrBadRequest) Validate(context.Context, qryRequest) error {
	return &ErrBadRequest{err: errors.New("already a bad request")}
}

func TestThatTheRegistrationInterfaceRemovesTheQueryHandler(t *testing.T) {
	// ARRANGE

	r := RegisterQueryHandler[qryRequest, string](&qryHandler{})

	// ACT

	wanted := 1
	got := len(queryHandlers)
	if wanted != got {
		t.Errorf("wanted %d handlers, got %d", wanted, got)
	}
	r.Remove()

	// ASSERT

	wanted = 0
	got = len(queryHandlers)
	if wanted != got {
		t.Errorf("wanted %d handlers, got %d", wanted, got)
	}
}

func TestThatRegisterQueryHandlerPanicsWhenHandlerIsAlreadyRegisteredForAType(t *testing.T) {
	// ARRANGE

	// 'arrange' the deferred ALERT since we're testing for a panic!
	defer func() {
		if r := recover(); r == nil {
			t.Error("did not panic")
		}
	}()

	// Register a handler and remove it when done
	r := RegisterQueryHandler[qryRequest, string](&qryHandler{})
	defer r.Remove()

	// ACT - attempt to register another handler for the same request type

	RegisterQueryHandler[qryRequest, string](&qryDuplicate{})

	// ASSERT (deferred, see above)
}

func TestThatQueryReturnsExpectedErrorWhenRequestHandlerIsNotRegistered(t *testing.T) {
	// ARRANGE
	// no-op

	// ACT

	_, err := Query[qryRequest, bool](context.Background(), qryRequest{})

	// ASSERT

	if _, ok := err.(*ErrNoHandler); !ok {
		t.Errorf("wanted *mediator.ErrNoHandler, got %T", err)
	}
}

func TestThatQueryReturnsExpectedErrorWhenRequestHandlerResultIsWrongType(t *testing.T) {
	// ARRANGE

	// Register a handler returning a string
	r := RegisterQueryHandler[qryRequest, string](&qryHandler{})
	defer r.Remove()

	// ACT

	// Request a Query returning a bool
	_, err := Query[qryRequest, bool](context.Background(), qryRequest{})

	// ASSERT

	if _, ok := err.(*ErrInvalidHandler); !ok {
		t.Errorf("wanted *mediator.ErrInvalidHandler, got %T", err)
	}
}

func TestThatQueryValidatorErrorIsReturnedAsErrBadRequest(t *testing.T) {
	// ARRANGE

	reg := RegisterQueryHandler[qryRequest, string](&qryRequestHandlerWithValidatorReturningError{})
	defer reg.Remove()

	// ACT

	_, err := Query[qryRequest, string](context.Background(), qryRequest{})

	// ASSERT

	if _, ok := err.(*ErrBadRequest); !ok {
		t.Errorf("wanted %T, got %T (%[1]q)", new(ErrBadRequest), err)
	}
}

func TestThatQueryValidatorErrorsDoNotWrapErrBadRequestErrors(t *testing.T) {
	// ARRANGE

	badRequest := &ErrBadRequest{}
	reg := RegisterQueryHandler[qryRequest, string](&qryRequestHandlerWithValidatorReturningErrBadRequest{})
	defer reg.Remove()

	// ACT

	_, err := Query[qryRequest, string](context.Background(), qryRequest{})

	// ASSERT

	bre, ok := err.(*ErrBadRequest)
	if !ok {
		wanted := badRequest
		got := bre
		t.Errorf("wanted %T, got %T (%[1]q)", wanted, got)
	}

	if bre.InnerError() != nil {
		got := bre.InnerError()
		if _, ok := got.(*ErrBadRequest); ok {
			t.Errorf("got %T wrapping %[1]T unnecessarily", badRequest)
		}
	}
}

func TestThatQueryResultIsReturnedToCaller(t *testing.T) {
	// ARRANGE

	reg := RegisterQueryHandler[qryRequest, string](&qryHandler{})
	defer reg.Remove()

	// ACT

	result, err := Query[qryRequest, string](context.Background(), qryRequest{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	// ASSERT

	wanted := qryHandlerResultValue
	got := result
	if wanted != got {
		t.Errorf("wanted %q, got %q", wanted, got)
	}
}
