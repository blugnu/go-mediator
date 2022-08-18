package mediator

import (
	"context"
	"errors"
	"testing"
)

type cmdRequestWithResult struct {
	result string
}

func TestThatTheRegistrationInterfaceRemovesTheHandler(t *testing.T) {
	// ARRANGE

	_, r := MockSuccessfulCommand[string]()

	// ACT

	wanted := 1
	got := len(commandHandlers)
	if wanted != got {
		t.Errorf("wanted %d handlers, got %d", wanted, got)
	}
	r.Remove()

	// ASSERT

	wanted = 0
	got = len(commandHandlers)
	if wanted != got {
		t.Errorf("wanted %d handlers, got %d", wanted, got)
	}
}

func TestThatRegisterCommandHandlerPanicsWhenHandlerIsAlreadyRegisteredForAType(t *testing.T) {
	// ARRANGE

	// 'arrange' a deferred ASSERT since we're testing for a panic!
	defer func() {
		if r := recover(); r == nil {
			t.Error("did not panic")
		}
	}()

	// Register a handler and remove it when done
	_, r := MockSuccessfulCommand[string]()
	defer r.Remove()

	// ACT

	// Attempt to register ANOTHER handler for the SAME request type
	MockSuccessfulCommand[string]()

	// ASSERT (deferred, see above)
}

func TestThatPerformReturnsExpectedErrorWhenRequestHandlerIsNotRegistered(t *testing.T) {
	// ARRANGE
	// no-op

	// ACT

	err := Perform(context.Background(), "test")

	// ASSERT

	if _, ok := err.(*ErrNoHandler); !ok {
		t.Errorf("wanted *mediator.ErrNoHandler, got %T", err)
	}
}

func TestThatResultsCanBeReturnedViaFieldsInAByRefRequestType(t *testing.T) {
	// ARRANGE

	original := "original"
	modified := "modified"

	_, reg := MockCommand(func(ctx context.Context, rq *cmdRequestWithResult) error { rq.result = modified; return nil })
	defer reg.Remove()

	// ACT

	request := &cmdRequestWithResult{result: original}
	err := Perform(context.Background(), request)

	// ASSERT

	if err != nil {
		t.Errorf("unexpected error Perform()ing request: %v", err)
	}

	wanted := modified
	got := request.result
	if wanted != got {
		t.Errorf("wanted %q in request.Result, got %q", wanted, got)
	}
}

func TestThatResultsCannotBeReturnedViaFieldsInAByValueRequestType(t *testing.T) {
	// ARRANGE
	original := "original"
	modified := "modified"

	_, reg := MockCommand(func(ctx context.Context, rq cmdRequestWithResult) error { rq.result = modified; return nil })
	defer reg.Remove()

	// ACT

	request := cmdRequestWithResult{result: original}
	err := Perform(context.Background(), request)

	// ASSERT

	if err != nil {
		t.Errorf("unexpected error Perform()ing request: %v", err)
	}

	wanted := original
	got := request.result
	if wanted != got {
		t.Errorf("wanted %q in request.Result, got %q", wanted, got)
	}
}

func TestThatCommandValidatorErrorIsReturnedAsErrBadRequest(t *testing.T) {
	// ARRANGE

	_, reg := MockCommandWithValidatorError[string](errors.New("error"))
	defer reg.Remove()

	// ACT

	err := Perform(context.Background(), "test")

	// ASSERT

	if _, ok := err.(*ErrBadRequest); !ok {
		t.Errorf("wanted %T, got %T (%q)", new(ErrBadRequest), err, err)
	}
}

func TestThatCommandValidatorErrorsDoNotWrapErrBadRequestErrors(t *testing.T) {
	// ARRANGE

	badRequest := &ErrBadRequest{err: errors.New("inner error")}
	_, reg := MockCommandWithValidatorError[string](badRequest)
	defer reg.Remove()

	// ACT

	err := Perform(context.Background(), "")

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

func TestThatCommandHandlerIsExecutedWhenRequestValidationIsSuccessful(t *testing.T) {
	// ARRANGE

	_, reg := MockCommandWithValidatorError[string](nil)
	defer reg.Remove()

	// ACT

	err := Perform(context.Background(), "request")

	// ASSERT

	wanted := error(nil)
	got := err
	if wanted != got {
		t.Errorf("wanted %q, got \"%v\"", wanted, got)
	}
}
