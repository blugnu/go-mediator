package mediator

import (
	"context"
	"reflect"
	"testing"
)

func TestMockQuery(t *testing.T) {

	if len(queryHandlers) > 0 {
		t.Fatal("invalid test: one or more query handlers are already registered")
	}

	mock, reg := MockQuery(func(ctx context.Context, s string) (string, error) { return "", nil })
	defer reg.Remove()

	t.Run("registers the handler", func(t *testing.T) {
		wanted := 1
		got := len(queryHandlers)
		if wanted != got {
			t.Errorf("wanted %d, got %d", wanted, got)
		}
	})

	// ACT
	_, err := Query[string, string](context.Background(), "test")
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
