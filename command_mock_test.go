package mediator

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestMockCommand(t *testing.T) {

	if len(commandHandlers) > 0 {
		t.Fatal("invalid test: one or more command handlers are already registered")
	}

	mock, reg := MockCommand(func(ctx context.Context, s string) error { return nil })
	defer reg.Remove()

	t.Run("registers the handler", func(t *testing.T) {
		wanted := 1
		got := len(commandHandlers)
		if wanted != got {
			t.Errorf("wanted %d, got %d", wanted, got)
		}
	})

	// ACT
	err := Perform(context.Background(), "test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	t.Run("NumRequests() captures number of requests processed", func(t *testing.T) {
		wanted := 1
		got := mock.NumRequests()
		if wanted != got {
			t.Errorf("wanted %d, got %d", wanted, got)
		}
	})

	t.Run("Requests() returns a copy of the processed requests", func(t *testing.T) {
		requests := mock.Requests()

		if reflect.ValueOf(requests).UnsafePointer() == reflect.ValueOf(mock.requests).UnsafePointer() {
			t.Error("got same slice")
		}

		if !reflect.DeepEqual(requests, mock.requests) {
			t.Errorf("wanted %v, got %v", mock.requests, requests)
		}
	})

	t.Run("WasCalled() captures whether the mock was called", func(t *testing.T) {
		wanted := true
		got := mock.WasCalled()

		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})
}

func TestMockCommandReturningError(t *testing.T) {
	// ARRANGE
	wanted := errors.New("error")
	_, reg := MockCommandReturningError[string](wanted)
	defer reg.Remove()

	// ACT
	got := Perform(context.Background(), "test")

	// ASSERT
	if wanted != got {
		t.Errorf("wanted %v, got %v", wanted, got)
	}
}
