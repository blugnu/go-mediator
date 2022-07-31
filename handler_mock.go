package mediator

import "context"

type mockhandler[TRequest any, TResult any] struct {
	requests []TRequest
	validate func(context.Context, TRequest) error
	execute  func(context.Context, TRequest) (TResult, error)
}

func MockHandler[TRequest any, TResult any]() (*mockhandler[TRequest, TResult], *reg) {
	return MockHandlerReturningError[TRequest, TResult](nil)
}

func MockHandlerWithFunc[TRequest any, TResult any](cmd HandlerFunc[TRequest, TResult]) (*mockhandler[TRequest, TResult], *reg) {
	h := &mockhandler[TRequest, TResult]{execute: cmd}
	r := RegisterHandler[TRequest, TResult](h)
	return h, r
}

func MockHandlerWithValidator[TRequest any, TResult any](qry HandlerFunc[TRequest, TResult], validator ValidatorFunc[TRequest]) (*mockhandler[TRequest, TResult], *reg) {
	h := &mockhandler[TRequest, TResult]{
		execute:  qry,
		validate: validator,
	}
	r := RegisterHandler[TRequest, TResult](h)
	return h, r
}

func MockHandlerReturningError[TRequest any, TResult any](err error) (*mockhandler[TRequest, TResult], *reg) {
	return MockHandlerWithFunc(func(context.Context, TRequest) (TResult, error) { return *new(TResult), err })
}

func MockHandlerReturningValues[TRequest any, TResult any](result TResult, err error) (*mockhandler[TRequest, TResult], *reg) {
	return MockHandlerWithFunc(func(context.Context, TRequest) (TResult, error) { return result, err })
}

func MockHandlerWithValidatorError[TRequest any, TResult any](err error) (*mockhandler[TRequest, TResult], *reg) {
	return MockHandlerWithValidator(
		func(context.Context, TRequest) (TResult, error) { return *new(TResult), nil },
		func(context.Context, TRequest) error { return err },
	)
}

func (mock *mockhandler[TRequest, TResult]) Execute(ctx context.Context, request TRequest) (TResult, error) {
	return mock.execute(ctx, request)
}

func (mock *mockhandler[TRequest, TResult]) Validate(ctx context.Context, request TRequest) error {
	mock.requests = append(mock.requests, request)
	if mock.validate != nil {
		return mock.validate(ctx, request)
	}
	return nil
}

func (mock *mockhandler[TRequest, TResult]) NumRequests() int {
	return len(mock.requests)
}

func (mock *mockhandler[TRequest, TResult]) Requests() []TRequest {
	return append([]TRequest{}, mock.requests...)
}

func (mock *mockhandler[TRequest, TResult]) WasCalled() bool {
	return len(mock.requests) > 0
}

func (mock *mockhandler[TRequest, TResult]) WasNotCalled() bool {
	return len(mock.requests) == 0
}
