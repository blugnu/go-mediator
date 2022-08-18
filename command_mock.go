package mediator

import "context"

type mockcommand[TRequest any] struct {
	requests []TRequest
	vfn      func(context.Context, TRequest) error
	fn       func(context.Context, TRequest) error
}

func MockCommand[TRequest any](cmd CommandFunc[TRequest]) (*mockcommand[TRequest], *reg) {
	h := &mockcommand[TRequest]{fn: cmd}
	r := RegisterCommandHandler[TRequest](h)
	return h, r
}

func MockCommandWithValidator[TRequest comparable](cmd CommandFunc[TRequest], validator ValidatorFunc[TRequest]) (*mockcommand[TRequest], *reg) {
	h := &mockcommand[TRequest]{
		fn:  cmd,
		vfn: validator,
	}
	r := RegisterCommandHandler[TRequest](h)
	return h, r
}

func MockSuccessfulCommand[TRequest comparable]() (*mockcommand[TRequest], *reg) {
	return MockCommand(func(context.Context, TRequest) error { return nil })
}

func MockCommandReturningError[TRequest comparable](err error) (*mockcommand[TRequest], *reg) {
	return MockCommand(func(context.Context, TRequest) error { return err })
}

func MockCommandWithValidatorError[TRequest comparable](err error) (*mockcommand[TRequest], *reg) {
	return MockCommandWithValidator(
		func(context.Context, TRequest) error { return nil },
		func(context.Context, TRequest) error { return err },
	)
}

func (mock *mockcommand[TRequest]) Execute(ctx context.Context, request TRequest) error {
	mock.requests = append(mock.requests, request)
	return mock.fn(ctx, request)
}

func (mock *mockcommand[TRequest]) Validate(ctx context.Context, request TRequest) error {
	if mock.vfn != nil {
		return mock.vfn(ctx, request)
	}
	return nil
}

func (mock *mockcommand[TRequest]) NumRequests() int {
	return len(mock.requests)
}

func (mock *mockcommand[TRequest]) Requests() []TRequest {
	return append([]TRequest{}, mock.requests...)
}

func (mock *mockcommand[TRequest]) WasCalled() bool {
	return len(mock.requests) > 0
}
