package mediator

import (
	"context"
	"fmt"
	"reflect"
)

// RegisterReceiver registers the specified handler for a particular request type.
//
// If a handler is already registered for that type the function will panic, otherwise
// the handler is registered.
func RegisterReceiver[TData any](handler Receiver[TData]) *reg {
	var data TData
	datatype := reflect.TypeOf(data)

	_, exists := receivers[datatype]
	if exists {
		panic(fmt.Sprintf("receiver already registered for %T", data))
	}

	receivers[datatype] = handler

	return &reg{
		registry:       receivers,
		registeredtype: datatype,
	}
}

// Send sends the specified data and context to the registered receiver
// for the data type and returns any error returned by the recevier.
//
// If the receiver implements Validator and the validator returns an error,
// then receiver is not called and the error returned by Send will be a
// ValidationError, wrapping the error returned by the validator.
func Send[TData any](ctx context.Context, data TData) error {
	datatype := reflect.TypeOf(data)

	receiver, ok := receivers[datatype].(Receiver[TData])
	if !ok {
		return NoReceiverError{data: data}
	}

	// You may be thinking that we should test that the handler we found is
	// of the correct type, but the magic of generics and the strict type
	// system takes care of that for us, so there's no need.  \o/

	// If the handler also provides a request validator call that first
	if validator, ok := receiver.(Validator[TData]); ok {
		err := validate(validator, ctx, data)
		if err != nil {
			return err
		}
	}

	return receiver.Execute(ctx, data)
}
