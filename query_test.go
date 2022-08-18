package mediator

import (
	"context"
	"errors"
	"testing"
)

func TestThatTheRegistrationInterfaceRemovesTheQueryHandler(t *testing.T) {

	if len(queryHandlers) > 0 {
		t.Fatal("invalid test: one or more query handlers are already registered")
	}

	// ARRANGE

	_, reg := MockQueryReturningValues[string]("", nil)

	// ACT

	wanted := 1
	got := len(queryHandlers)
	if wanted != got {
		t.Errorf("wanted %d handlers, got %d", wanted, got)
	}
	reg.Remove()

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
	_, reg := MockQueryReturningValues[string]("result", nil)
	defer reg.Remove()

	// ACT - attempt to register another handler for the same request type

	MockQueryReturningValues[string]("other", nil)

	// ASSERT (deferred, see above)
}

func TestThatQueryReturnsExpectedErrorWhenRequestHandlerIsNotRegistered(t *testing.T) {
	// ARRANGE
	// no-op

	// ACT

	_, err := Query[string, bool](context.Background(), "request")

	// ASSERT

	if _, ok := err.(*ErrNoHandler); !ok {
		t.Errorf("wanted *mediator.ErrNoHandler, got %T", err)
	}
}

func TestThatQueryReturnsExpectedErrorWhenRequestHandlerResultIsWrongType(t *testing.T) {
	// ARRANGE

	// Register a handler returning a string
	_, reg := MockQueryReturningValues[string]("string response", nil)
	defer reg.Remove()

	// ACT

	// Request a Query returning a bool
	_, err := Query[string, bool](context.Background(), "request")

	// ASSERT

	if _, ok := err.(*ErrInvalidHandler); !ok {
		t.Errorf("wanted *mediator.ErrInvalidHandler, got %T", err)
	}
}

func TestThatQueryValidatorErrorIsReturnedAsErrBadRequest(t *testing.T) {
	// ARRANGE

	_, reg := MockQueryWithValidatorError[string, string](errors.New("error"))
	defer reg.Remove()

	// ACT

	_, err := Query[string, string](context.Background(), "request")

	// ASSERT

	if _, ok := err.(*ErrBadRequest); !ok {
		t.Errorf("wanted %T, got %T (%[2]q)", new(ErrBadRequest), err)
	}
}

func TestThatQueryValidatorErrorsDoNotWrapErrBadRequestErrors(t *testing.T) {
	// ARRANGE

	badRequest := &ErrBadRequest{}
	_, reg := MockQueryWithValidatorError[string, string](badRequest)
	defer reg.Remove()

	// ACT

	_, err := Query[string, string](context.Background(), "request")

	// ASSERT

	bre, ok := err.(*ErrBadRequest)
	if !ok {
		wanted := badRequest
		got := err
		t.Errorf("wanted %T, got %T (%[2]q)", wanted, got)
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

	wanted := "result"
	_, reg := MockQueryReturningValues[string](wanted, nil)
	defer reg.Remove()

	// ACT

	result, err := Query[string, string](context.Background(), "request")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	// ASSERT

	got := result
	if wanted != got {
		t.Errorf("wanted %q, got %q", wanted, got)
	}
}
