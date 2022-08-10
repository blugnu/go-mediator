package mediator

import (
	"context"
	"testing"
)

type qryRequest struct{}
type qryHandler struct{}
type qryDuplicate struct{}
type qryInvalidResult struct{}

func (*qryDuplicate) Execute(context.Context, qryRequest) (string, error)   { return "", nil }
func (*qryHandler) Execute(context.Context, qryRequest) (string, error)     { return "", nil }
func (*qryInvalidResult) Execute(context.Context, qryRequest) (bool, error) { return false, nil }

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
	// ARRANGE (and ASSERT, since we're testing for a panic() :) )

	// Setup the panic test (deferred ASSERT)
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

	// ASSERT (deferred)
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
