package mediator

import "context"

type mockquery[TRequest any, TResult any] struct {
	requests []TRequest
	vfn      func(context.Context, TRequest) error
	fn       func(context.Context, TRequest) (TResult, error)
}

func MockQuery[TRequest any, TResult any](cmd QueryFunc[TRequest, TResult]) (*mockquery[TRequest, TResult], *reg) {
	h := &mockquery[TRequest, TResult]{fn: cmd}
	r := RegisterQueryHandler[TRequest, TResult](h)
	return h, r
}

func MockQueryWithValidator[TRequest any, TResult any](qry QueryFunc[TRequest, TResult], validator ValidatorFunc[TRequest]) (*mockquery[TRequest, TResult], *reg) {
	h := &mockquery[TRequest, TResult]{
		fn:  qry,
		vfn: validator,
	}
	r := RegisterQueryHandler[TRequest, TResult](h)
	return h, r
}

func MockQueryReturningValues[TRequest any, TResult any](result TResult, err error) (*mockquery[TRequest, TResult], *reg) {
	return MockQuery(func(context.Context, TRequest) (TResult, error) { return result, err })
}

func MockQueryWithValidatorError[TRequest any, TResult any](err error) (*mockquery[TRequest, TResult], *reg) {
	return MockQueryWithValidator(
		func(context.Context, TRequest) (TResult, error) { return *new(TResult), nil },
		func(context.Context, TRequest) error { return err },
	)
}

func (mock *mockquery[TRequest, TResult]) Execute(ctx context.Context, request TRequest) (TResult, error) {
	mock.requests = append(mock.requests, request)
	return mock.fn(ctx, request)
}

func (mock *mockquery[TRequest, TResult]) Validate(ctx context.Context, request TRequest) error {
	if mock.vfn != nil {
		return mock.vfn(ctx, request)
	}
	return nil
}

func (mock *mockquery[TRequest, TResult]) NumRequests() int {
	return len(mock.requests)
}

func (mock *mockquery[TRequest, TResult]) Requests() []TRequest {
	return append([]TRequest{}, mock.requests...)
}

func (mock *mockquery[TRequest, TResult]) WasCalled() bool {
	return len(mock.requests) > 0
}
