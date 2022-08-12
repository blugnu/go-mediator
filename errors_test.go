package mediator

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestThatErrNoHandlerYieldsCorrectString(t *testing.T) {
	request := "a string value"
	tests := map[string]struct {
		handler handlerType
		wanted  string
	}{
		"command error": {handler: command, wanted: fmt.Sprintf("no command handler for '%T'", request)},
		"query error":   {handler: query, wanted: fmt.Sprintf("no query handler for '%T'", request)},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := &ErrNoHandler{handler: test.handler, request: request}
			got := err.Error()
			if got != test.wanted {
				t.Errorf("wanted %q, got %q", test.wanted, got)
			}
		})
	}
}

type boolCmd struct{}

func (*boolCmd) Execute(context.Context, bool) error { return nil }

type boolQry struct{}

func (*boolQry) Execute(context.Context, bool) (bool, error) { return true, nil }

func TestThatErrInvalidHandlerYieldsCorrectString(t *testing.T) {
	request := "a string value"
	tests := map[string]struct {
		handlerType handlerType
		handler     interface{}
		wanted      string
	}{
		"command error": {handlerType: command, handler: &boolCmd{}, wanted: fmt.Sprintf("%T is not a valid command handler for '%T' requests", &boolCmd{}, request)},
		"query error":   {handlerType: query, handler: &boolQry{}, wanted: fmt.Sprintf("%T is not a valid query handler for '%T' requests", &boolQry{}, request)},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := &ErrInvalidHandler{handlerType: test.handlerType, handler: test.handler, request: request}
			got := err.Error()
			if got != test.wanted {
				t.Errorf("wanted %q, got %q", test.wanted, got)
			}
		})
	}
}

func TestErrBadRequest(t *testing.T) {
	// ARRANGE

	inner := errors.New("inner error")
	err := &ErrBadRequest{err: inner}

	// ACT

	result := err.Error()

	// ASSERT

	t.Run("yields correct String()", func(t *testing.T) {
		wanted := fmt.Sprintf("bad request: %s", inner)
		got := result
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})

	t.Run("returns InnerError() correctly", func(t *testing.T) {
		wanted := inner
		got := err.InnerError()
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)
		}
	})
}
