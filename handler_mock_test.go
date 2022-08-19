package mediator

import (
	"context"
	"reflect"
	"testing"
)

func TestMockHandler(t *testing.T) {

	if len(handlers) > 0 {
		t.Fatal("invalid test: one or more handlers are already registered")
	}

	mock, reg := MockHandler[string, string]()
	defer reg.Remove()

	t.Run("registers the handler", func(t *testing.T) {
		wanted := 1
		got := len(handlers)
		if wanted != got {
			t.Errorf("wanted %d, got %d", wanted, got)
		}
	})

	// ACT
	_, err := Perform[string, string](context.Background(), "test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	t.Run("captures number of requests handled", func(t *testing.T) {
		wanted := 1
		got := mock.NumRequests()
		if wanted != got {
			t.Errorf("wanted %d, got %d", wanted, got)
		}
	})

	t.Run("returns a copy of the handled requests", func(t *testing.T) {
		requests := mock.Requests()

		if reflect.ValueOf(requests).UnsafePointer() == reflect.ValueOf(mock.requests).UnsafePointer() {
			t.Error("got same slice")
		}

		if !reflect.DeepEqual(requests, mock.requests) {
			t.Errorf("wanted %v, got %v", mock.requests, requests)
		}
	})

	t.Run("captures that a handler was called", func(t *testing.T) {
		wantedWc := true
		wantedWnc := false
		gotWc := mock.WasCalled()
		gotWnc := mock.WasNotCalled()

		if wantedWc != gotWc || wantedWnc != gotWnc {
			t.Errorf("called / not called: wanted %v / %v, got %v / %v", wantedWc, wantedWnc, gotWc, gotWnc)
		}
	})

	t.Run("captures that a handler was not called", func(t *testing.T) {
		mock, reg := MockHandler[int, int]()
		defer reg.Remove()

		wantedWc := false
		wantedWnc := true
		gotWc := mock.WasCalled()
		gotWnc := mock.WasNotCalled()

		if wantedWc != gotWc || wantedWnc != gotWnc {
			t.Errorf("called / not called: wanted %v / %v, got %v / %v", wantedWc, wantedWnc, gotWc, gotWnc)
		}
	})

}
