package mediator

import (
	"context"
	"testing"
)

type cmdRequest struct{}
type cmdHandler struct{}
type cmdDuplicate struct{}

func (*cmdDuplicate) Execute(context.Context, cmdRequest) error { return nil }
func (*cmdHandler) Execute(context.Context, cmdRequest) error   { return nil }

func TestThatTheRegistrationInterfaceRemovesTheHandler(t *testing.T) {
	// ARRANGE
	r := RegisterCommandHandler[cmdRequest](&cmdHandler{})

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
	// ARRANGE (and ASSERT, since we're testing for a panic() :) )

	// Setup the panic test (deferred ASSERT)
	defer func() {
		if r := recover(); r == nil {
			t.Error("did not panic")
		}
	}()

	// Register a handler and remove it when done
	r := RegisterCommandHandler[cmdRequest](&cmdHandler{})
	defer r.Remove()

	// ACT - attempt to register another handler for the same request type
	RegisterCommandHandler[cmdRequest](&cmdDuplicate{})

	// ASSERT (deferred)
}

func TestThatPerformReturnsExpectedErrorWhenRequestHandlerIsNotRegistered(t *testing.T) {
	// ARRANGE
	// no-op

	// ACT
	err := Perform(context.Background(), cmdRequest{})

	// ASSERT
	if _, ok := err.(*ErrNoHandler); !ok {
		t.Errorf("wanted *mediator.ErrNoHandler, got %T", err)
	}
}
