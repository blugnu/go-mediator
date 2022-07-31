package mediator

import (
	"context"
	"errors"
)

// validate calls the supplied Validator for the context and input specified,
// wrapping an any returned error that is not a ValidationError.
func validate[TInput any](v Validator[TInput], ctx context.Context, input TInput) error {
	if err := v.Validate(ctx, input); err != nil {
		if errors.As(err, &ValidationError{}) {
			return err
		}
		return ValidationError{err}
	}
	return nil
}
