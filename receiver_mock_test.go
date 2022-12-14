package mediator

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestMockReceiver(t *testing.T) {

	if len(receivers) > 0 {
		t.Fatal("invalid test: one or more receivers are already registered")
	}

	dataSent := "sent"
	dataNotSent := "not sent"

	mock, reg := MockReceiver[string]()
	defer reg.Remove()

	t.Run("registers the receiver", func(t *testing.T) {
		wanted := 1
		got := len(receivers)
		if wanted != got {
			t.Errorf("wanted %d, got %d", wanted, got)
		}
	})

	// ACT
	err := Send(context.Background(), dataSent)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	t.Run("captures data received", func(t *testing.T) {
		wanted := 1
		got := len(mock.DataReceived())
		if wanted != got {
			t.Errorf("wanted %d, got %d", wanted, got)
		}
	})

	t.Run("captures a copy of received data", func(t *testing.T) {
		received := mock.DataReceived()

		if reflect.ValueOf(received).UnsafePointer() == reflect.ValueOf(mock.received).UnsafePointer() {
			t.Error("got same slice")
		}

		if !reflect.DeepEqual(received, mock.received) {
			t.Errorf("wanted %v, got %v", mock.received, received)
		}
	})

	t.Run("captures that data was received", func(t *testing.T) {
		wanted := true
		got := mock.Received(dataSent)

		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})

	t.Run("captures that data was not received", func(t *testing.T) {
		wanted := false
		got := mock.Received(dataNotSent)

		if wanted != got {
			t.Errorf("wanted %v, got %v", wanted, got)
		}
	})

	t.Run("captures that a receiver was called", func(t *testing.T) {
		wantedWc := true
		wantedWnc := false
		gotWc := mock.WasCalled()
		gotWnc := mock.WasNotCalled()
		if wantedWc != gotWc || wantedWnc != gotWnc {
			t.Errorf("called / not called: wanted %v / %v, got %v / %v", wantedWc, wantedWnc, gotWc, gotWnc)
		}
	})

	t.Run("captures that a receiver was not called", func(t *testing.T) {
		mock, reg := MockReceiver[int]()
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

func TestMockReceiverReturningError(t *testing.T) {
	// ARRANGE
	wanted := errors.New("error")
	_, reg := MockReceiverReturningError[string](wanted)
	defer reg.Remove()

	// ACT
	got := Send(context.Background(), "test")

	// ASSERT
	if wanted != got {
		t.Errorf("wanted %v, got %v", wanted, got)
	}
}
