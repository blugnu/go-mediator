package mediator

import (
	"context"
	"reflect"
)

type mockreceiver[TData any] struct {
	received []TData
	validate func(context.Context, TData) error
	execute  func(context.Context, TData) error
}

func MockReceiver[TData any]() (*mockreceiver[TData], *reg) {
	return MockReceiverReturningError[TData](nil)
}

func MockReceiverWithFunc[TData any](executor ReceiverFunc[TData]) (*mockreceiver[TData], *reg) {
	h := &mockreceiver[TData]{execute: executor}
	r := RegisterReceiver[TData](h)
	return h, r
}

func MockReceiverWithValidator[TData any](executor ReceiverFunc[TData], validator ValidatorFunc[TData]) (*mockreceiver[TData], *reg) {
	h := &mockreceiver[TData]{
		execute:  executor,
		validate: validator,
	}
	r := RegisterReceiver[TData](h)
	return h, r
}

func MockReceiverReturningError[TData any](err error) (*mockreceiver[TData], *reg) {
	return MockReceiverWithFunc(func(context.Context, TData) error { return err })
}

func MockReceiverWithValidatorError[TData any](err error) (*mockreceiver[TData], *reg) {
	return MockReceiverWithValidator(
		func(context.Context, TData) error { return nil },
		func(context.Context, TData) error { return err },
	)
}

func (mock *mockreceiver[TData]) Execute(ctx context.Context, request TData) error {
	mock.received = append(mock.received, request)
	return mock.execute(ctx, request)
}

func (mock *mockreceiver[TData]) Validate(ctx context.Context, request TData) error {
	if mock.validate != nil {
		return mock.validate(ctx, request)
	}
	return nil
}

func (mock *mockreceiver[TData]) Received(data TData) bool {
	for _, received := range mock.received {
		if reflect.DeepEqual(received, data) {
			return true
		}
	}
	return false
}

func (mock *mockreceiver[TData]) DataReceived() []TData {
	return append([]TData{}, mock.received...)
}

func (mock *mockreceiver[TData]) WasCalled() bool {
	return len(mock.received) > 0
}

func (mock *mockreceiver[TData]) WasNotCalled() bool {
	return len(mock.received) == 0
}
