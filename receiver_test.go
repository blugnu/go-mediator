package mediator

import (
	"context"
	"errors"
	"testing"
)

type cmdRequestWithResult struct {
	result string
}

func TestThatTheRegistrationInterfaceRemovesTheReceiver(t *testing.T) {
	// ARRANGE

	_, r := MockReceiver[string]()

	// ACT

	wanted := 1
	got := len(receivers)
	if wanted != got {
		t.Errorf("wanted %d handlers, got %d", wanted, got)
	}
	r.Remove()

	// ASSERT

	wanted = 0
	got = len(receivers)
	if wanted != got {
		t.Errorf("wanted %d handlers, got %d", wanted, got)
	}
}

func TestThatRegisterReceiverPanicsWhenDataTypeReceiverIsAlreadyRegistered(t *testing.T) {
	// ARRANGE

	// 'arrange' a deferred ASSERT since we're testing for a panic!
	defer func() {
		if r := recover(); r == nil {
			t.Error("did not panic")
		}
	}()

	// Register a handler and remove it when done
	_, r := MockReceiver[string]()
	defer r.Remove()

	// ACT

	// Attempt to register ANOTHER handler for the SAME request type
	MockReceiver[string]()

	// ASSERT (deferred, see above)
}

func TestSendErrorWhenNoReceiverIsRegistered(t *testing.T) {
	// ARRANGE
	// no-op

	// ACT

	err := Send(context.Background(), "test")

	// ASSERT

	wanted := NoReceiverError{}
	if !errors.As(err, &wanted) {
		t.Errorf("wanted %T, got %T", wanted, err)
	}
}

func TestThatResultsCanBeReturnedViaFieldsInAByRefRequestType(t *testing.T) {
	// ARRANGE

	original := "original"
	modified := "modified"

	_, reg := MockReceiverWithFunc(func(ctx context.Context, rq *cmdRequestWithResult) error {
		rq.result = modified
		return nil
	})
	defer reg.Remove()

	// ACT

	request := &cmdRequestWithResult{result: original}
	err := Send(context.Background(), request)

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

	_, reg := MockReceiverWithFunc(func(ctx context.Context, rq cmdRequestWithResult) error {
		rq.result = modified
		return nil
	})
	defer reg.Remove()

	// ACT

	request := cmdRequestWithResult{result: original}
	err := Send(context.Background(), request)

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

func TestThatReceiverValidatorErrorIsReturnedAsInputError(t *testing.T) {
	// ARRANGE

	_, reg := MockReceiverWithValidatorError[string](errors.New("error"))
	defer reg.Remove()

	// ACT

	err := Send(context.Background(), "test")

	// ASSERT

	wanted := ValidationError{}
	if !errors.As(err, &wanted) {
		t.Errorf("wanted %T, got %T (%[2]q)", wanted, err)
	}
}

func TestThatReceiverValidatorDoesNotWrapInputError(t *testing.T) {
	// ARRANGE

	_, reg := MockReceiverWithValidatorError[string](ValidationError{errors.New("inner error")})
	defer reg.Remove()

	// ACT

	err := Send(context.Background(), "")

	// ASSERT

	var wanted = ValidationError{}
	if !errors.As(err, &wanted) {
		t.Errorf("wanted %T, got %T (%[2]q)", wanted, err)
	}

	err = wanted.error
	if errors.As(err, &wanted) {
		t.Errorf("got %T wrapping %[1]T unnecessarily", err)
	}
}

func TestThatReceiverIsExecutedWhenRequestValidationIsSuccessful(t *testing.T) {
	// ARRANGE

	_, reg := MockReceiverWithValidatorError[string](nil)
	defer reg.Remove()

	// ACT

	err := Send(context.Background(), "request")

	// ASSERT

	wanted := error(nil)
	got := err
	if wanted != got {
		t.Errorf("wanted %q, got \"%v\"", wanted, got)
	}
}
