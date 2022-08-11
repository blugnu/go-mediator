package mediator

import "context"

// validate calls the supplied RequestValidator for the context and request specified,
// wrapping an ErrBadRequest around any returned error that is not already an ErrBadRequest.
func validate[TRequest any](v RequestValidator[TRequest], ctx context.Context, rq TRequest) error {
	if err := v.Validate(ctx, rq); err != nil {
		// If the error is not ErrBadRequest, wrap it
		if _, ok := err.(*ErrBadRequest); !ok {
			err = &ErrBadRequest{err: err}
		}
		return err
	}
	return nil
}
