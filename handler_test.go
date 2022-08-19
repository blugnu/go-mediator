package mediator

import (
	"context"
	"errors"
	"testing"
)

func TestThatTheRegistrationInterfaceRemovesTheHandlerHandler(t *testing.T) {

	if len(handlers) > 0 {
		t.Fatal("invalid test: one or more handlers are already registered")
	}

	// ARRANGE

	_, reg := MockHandlerReturningValues[string]("", nil)

	// ACT

	wanted := 1
	got := len(handlers)
	if wanted != got {
		t.Errorf("wanted %d handlers, got %d", wanted, got)
	}
	reg.Remove()

	// ASSERT

	wanted = 0
	got = len(handlers)
	if wanted != got {
		t.Errorf("wanted %d handlers, got %d", wanted, got)
	}
}

func TestThatRegisterHandlerHandlerPanicsWhenHandlerIsAlreadyRegisteredForAType(t *testing.T) {
	// ARRANGE

	// 'arrange' the deferred ALERT since we're testing for a panic!
	defer func() {
		if r := recover(); r == nil {
			t.Error("did not panic")
		}
	}()

	// Register a handler and remove it when done
	_, reg := MockHandlerReturningValues[string]("result", nil)
	defer reg.Remove()

	// ACT - attempt to register another handler for the same request type

	MockHandlerReturningValues[string]("other", nil)

	// ASSERT (deferred, see above)
}

func TestThatHandlerReturnsExpectedErrorWhenHandlerIsNotRegistered(t *testing.T) {
	// ARRANGE
	// no-op

	// ACT

	_, err := Perform[string, bool](context.Background(), "request")

	// ASSERT

	if _, ok := err.(*NoReceiverError); !ok {
		t.Errorf("wanted *mediator.ErrNoHandler, got %T", err)
	}
}

func TestThatHandlerReturnsExpectedErrorWhenHandlerResultIsWrongType(t *testing.T) {
	// ARRANGE

	// Register a handler returning a string
	_, reg := MockHandlerReturningValues[string]("string response", nil)
	defer reg.Remove()

	// ACT

	// Request a Handler returning a bool
	_, err := Perform[string, bool](context.Background(), "request")

	// ASSERT

	if _, ok := err.(*InvalidHandlerError); !ok {
		t.Errorf("wanted *mediator.ErrInvalidHandler, got %T", err)
	}
}

func TestThatHandlerValidatorErrorIsReturnedAsInputError(t *testing.T) {
	// ARRANGE

	_, reg := MockHandlerWithValidatorError[string, string](errors.New("error"))
	defer reg.Remove()

	// ACT

	_, err := Perform[string, string](context.Background(), "request")

	// ASSERT

	wanted := ValidationError{}
	if !errors.As(err, &wanted) {
		t.Errorf("wanted %T, got %T (%[2]q)", wanted, err)
	}
}

func TestThatHandlerValidatorDoesNotWrapInputError(t *testing.T) {
	// ARRANGE

	_, reg := MockHandlerWithValidatorError[string, string](ValidationError{})
	defer reg.Remove()

	// ACT

	_, err := Perform[string, string](context.Background(), "request")

	// ASSERT

	var wanted = ValidationError{}
	if !errors.As(err, &wanted) {
		t.Errorf("wanted %T, got %T (%[2]q)", wanted, err)
	}

	got := wanted.error
	if errors.As(got, &wanted) {
		t.Errorf("got %T wrapping %[1]T unnecessarily (%[1]v)", err)
	}
}

func TestThatHandlerResultIsReturnedToCaller(t *testing.T) {
	// ARRANGE

	wanted := "result"
	_, reg := MockHandlerReturningValues[string](wanted, nil)
	defer reg.Remove()

	// ACT

	result, err := Perform[string, string](context.Background(), "request")
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
