<div align="center" style="margin-bottom:20px">
  <img src=".assets/banner.png" alt="go-mediator" />
  <div align="center">
    <a href="https://github.com/deltics/go-mediator/actions/workflows/qa.yml"><img alt="build-status" src="https://github.com/deltics/go-mediator/actions/workflows/qa.yml/badge.svg?branch=master&style=flat-square"/></a>
    <a href="https://goreportcard.com/report/github.com/deltics/go-mediator" ><img alt="go report" src="https://goreportcard.com/badge/github.com/deltics/go-mediator"/></a>
    <a><img alt="go version >= 1.18" src="https://img.shields.io/badge/go%20version-%3E=1.18-61CFDD.svg?style=flat-square"/></a>
    <a href="https://github.com/deltics/go-mediator/blob/master/LICENCE"><img alt="MIT License" src="https://img.shields.io/github/license/deltics/go-mediator?color=%234275f5&style=flat-square"/></a>
    <a href="https://coveralls.io/github/deltics/go-mediator?branch=master"><img alt="coverage" src="https://img.shields.io/coveralls/github/deltics/go-mediator?style=flat-square"/></a>
    <a href="https://pkg.go.dev/github.com/deltics/go-mediator"><img alt="docs" src="https://pkg.go.dev/badge/github.com/deltics/go-mediator"/></a>
  </div>
</div>

<br/>

# go-mediator

A light-weight implementation of the [Mediator Pattern](https://en.wikipedia.org/wiki/Mediator_pattern) for `goLang`, inspired by [jbogard's MediatR framework for .net](https://github.com/jbogard/MediatR) but with far more limited ambition (_for now at least_).

<br/>

## Mediator Pattern
[The Mediator](https://en.wikipedia.org/wiki/Mediator_pattern) is a simple [pattern](https://en.wikipedia.org/wiki/Software_design_pattern) that uses a 3rd-party (the mediator) to facilitate communication between two other parties without them requiring knowledge of each other.

It is a powerful pattern for achieving loosely-coupled code.

There are many ways to implement a mediator, from simple `func` pointers to sophisticated and complex messaging systems.

`go-mediator` sits firmly at the *simple* end of that spectrum and intends staying there!

<br/>

## What go-mediator Is NOT
- `go-mediator` is not a message queue
- `go-mediator` is not asynchronous
- `go-mediator` is not complicated!

<br/>

# Terminology and Concepts

`go-mediator` maintains a registry of receivers and handlers that receive and respond to values of a specific type.

Values are submitted to the `mediator` which locates the appropriate `Receiver` or `Handler` for the type of Value submitted.  The submitted value is passed to the `Receiver` or `Handler` and any result (or `error`) returned is then passed back to the original caller.

Validation of submitted values may be performed 'in-line' with the execution of the `Receiver` or `Handler` or via the implementation of a `Validator` interface, for more complex validation needs.

`mediator` calls are synchronous.

<br/>

## Receivers vs Handlers

`go-mediator` makes a formal distinction between a `Receiver` and `Handler` as follows:

- A `Receiver` accepts `TData` and returns _only_ an `error` (or nil)
- A `Handler` accepts `TRequest` and returns *a `TResult`* **and** an `error` (or nil)

```go
// Receiver[TData] is the interface implemented by a receiver
type Receiver[TData any] interface {
    Execute(context.Context, TData) error
}

// Handler[TRequest, TResult] is the interface implemented by a handler
type Handler[TRequest any, TResult any] interface {
    Execute(context.Context, TRequest) (TResult, error)
}
```

Since receivers and handlers have different type parameters, separate functions are provided for registering implementations:

```go
func RegisterReceiver[TData any](Receiver[TData]) *reg
func RegisterHandler[TRequest any, TResult any](Handler[TRequest, TResult]) *reg
```

This distinction simplifies code that uses `mediator` to send data to a `Receiver`, thanks to `golang` type inference:

```go
    err := mediator.Send(ctx, &SomeData{ Id: id })

    // vs

    result, err := mediator.Perform[GetProductRequest, GetProductResult](ctx, &GetProductRequest{ Id: productId})
```

It also means that `Receiver` implementations are not required to return a placeholder result value and callers of a `Receiver` are not required to ignore that result.  i.e. it more clearly expresses the contract between the caller and receiver, mediated by `mediator`.

It also makes it apparent that when calling a `Handler` with a `Request`, there is a result, in addition to any error, which should _not_ be ignored.

<br/>

## Validators
Both a `Receiver` or a `Handler` may (*optionally*) implement the `Validator` interface, to separate validation from execution of a known valid input:

```go
// Validator[TInput] is an optional interface that may be implemented
// by both receivers and handlers.
type Validator[TInput any] interface {
    Validate(context.Context, TInput) error
}
```

For a `Receiver` implementating `Validator`, `TInput` is the _same type_ as `TData`.  For a `Handler`, the _same type_ as `TRequest`.

Any `error` returned from a `Validator` will be wrapped and returned by `mediator` to the caller in a `ValidationError`.  If `Validator` returns a `ValidationError`, this will **not** be wrapped.

<br>

><br>_Since it is impossible for `mediator` to differentiate between an error returned from `Execute()` which relates to validation rather than execution, any validation errors returned from `Execute()` should explicitly be of type `ValidationError`._<br><br>


<br/>

# Getting Started

For the purposes of this section, only a `Receiver` will be considered.  The steps are essentially the same for a `Handler`, with the addition of a `TResult` type, but where there are significant differences these will be mentioned.

## 1. Declare a Data Type

The data type for values sent via `mediator` is effectively the unique "address" of a handler, as there can be only one `Receiver` for a given data type.  Therefore you should declare a specific type for each `Receiver` to receive and should **not* use built-in types.

It is common to use a `struct` for `Receiver` data types, even where only a single value is passed:

```go
    type FooData struct {
        Foo string
    }
```

><br>_Because data types have a 1:1 relationship with a specific `Receiver`, if you have different receivers that accept effectively the same types of values (e.g. a single string value), you need separate and distinct request types for each one, as illustrated below._<br><br>

```go
    // Two separate Data types for different receivers,
    // each receiving a single string:

    // Data sent to a FooReceiver
    type FooData struct {
        Foo string
    }

    // Data sent to a BarReceiver
    type BarData struct {
        Bar string
    }
```

<br/>

## 2. Implement the Receiver Interface
A `Receiver` is an implementation of a generic interface with a single `Execute` method accepting a `context` and a `TData` value.

```golang
    type Receiver[TData any] interface {
        Execute(context.Context, TData) error
    }
```

An implementation for receiving `FooData` might be similar to this example:

```golang
    // FooReceiver
    type FooReceiver struct {}

    // FooReceiver implements the Receiver interface
    func (*FooReceiver) Execute(ctx context.Context, data FooData) error {
        // Check that data holds a valid Foo
        if len(data.Foo) == 0 {
            return mediator.ValidationError{errors.New("missing Foo in data")}
        }

        // Do some foo'ing with data.Foo
        return nil
    }
```

><br>_An underlying `struct` of a `Receiver` implementation may be used to hold services needed by the `Receiver`. This may be useful in test code for injecting fake or mock services, for example.  See more on testing further below..._<br><br>

<br/>

## 3. Register the Receiver
A `Receiver` implementation must be registered with `mediator` for the `TData` type that it receives.

```go
    RegisterReceiver[FooData](&FooReceiver{})
```

<br/>

## 4. Implement RequestValidator (*optional*)
If a `Receiver` has more complex validation needs, or if you simply wish to strictly separate concerns, implement the `Validator` interface on your `Receiver`:

```go
    // FooReceiver
    type FooReceiver struct {}

    // FooReceiver implements the Receiver interface
    func (*FooReceiver) Execute(ctx context.Context, data FooData) error {
        // Do some foo'ing with data.Foo - we can be sure it is valid
        return nil
    }

    // FooReceiver also implements the Validator interface
    func (*FooReceiver) Validate(ctx context.Context, data FooData) error {
        // Check that data holds a valid Foo
        if len(data.Foo) == 0 {
            return errors.New("missing Foo in data")
        }

        // The data is valid
        return nil
    }
```

If `Validate()` returns an error:

 - the error is automatically wrapped in a `ValidationError` if necessary, before being returned to the caller
 - if the error is already a `ValidationError` it will not be wrapped
 - the `Execute()` function _**will not be called**_ for that data

<br/>

## 5. Send Data to Receiver via Mediator
To send data to your `Receiver`, initialise a value of the required data type and `Send` it:

```go
    err := mediator.Send(ctx, &FooData{ Foo: "do something for me" })
```

For a `Handler`, instead of `Send()`ing data, you instead ask `mediator` to `Perform()` a request:

```go
    err := mediator.Perform[FooRequest, string](ctx, &FooRequest{ Foo: "get me something nice" })
```

When `Send()`ing to a `Receiver` the type parameter on the `Send()` generic function is not required as it is inferred from the data value parameter.

This is not possible when `Perform()`ing requests since the required `TResult` type is not expressed in the args, so there is nothing from which it can be inferred.  As a result, for `Perform()` call, both `TRequest` and `TResult` must be specified.

Sorry.

<br/>

# Alternative Result Handling

Normally a `Receiver` can return only an `error` (or nil).  It may be tempting to return values other than an `error` using a _by reference_ type for the data (e.g. pointer to struct).

The receiver may then manipulate the members of the received data, perhaps even ones provided explicitly for the purpose of "returning" a value (or values).

This avoids the need to specify request and result types when `Perform()`ing, but introduces its own problems, not to mention violating the idiomatic `result, error := ` pattern.

Consider:

```go
    // Using a `Handler` with explicit `TResult` (string):

    result, err := mediator.Perform[FooRequest, string](FooRequest{})
    if err != nil {
        log.Error(err)
        return err
    }

    // Process result...
```

Versus:

```go
    // Using a `Receiver` with by-reference data and side effect(s):

    request := &FooRequest{}
    err := mediator.Send(request)
    if err != nil {
        log.Error(err)
        return err
    }
    result := request.Result    // a 'string' member of the request

    // Process result...
```

><br> In case it was not already clear... **DO NOT DO THIS**<br><br>**The fact that it is even possible is only mentioned in order to highlight the reasons for _not_ doing it, should you be tempted.  :)**<br><br>

This does not mean you should not use by-reference types at all (e.g. for efficiency by avoiding the copying of a by-value struct type), only that you should not exploit the ability to mutate values in the request.

<br/>

# Testing With Mediator

The loose-coupling that can be achieved with a mediator has obvious utility when it comes to testing code.

Most obviously it enables code under test to make requests that are picked up by *test* handlers, rather than *production* handlers, without the code under test having to do anything to achieve this or even being aware that it is happening!

There are some other benefits of the particular implementation of `go-mediator`, as well as a couple of gotchas to watch out for...

<br/>

## Receiver/Handler Dependency Injection
A receiver or handler implemented as a `struct` may hold references to services *required by* the implementation, injected at the time that the implementation is initialised (usually at the time of registration).

For example a receiver might employ a repository interface:

```golang
    type FooReceiver struct {
        Repository   FooRepository
    }
```

Implementation registered by your application or service (e.g. in `main.go`) will be registered using the 'production' services required (or other concrete services appropraite to the runtime environment, if you use different environments for integration or acceptance tests for example).

Unless you are using package initialisation, these registrations will not usually be present when your unit tests run.  Your unit tests must register an implementation for any `mediator` calls that are made when exercise code under test.

This might mean registering the usual *`production`* handler, but injecting mock services or other test dependencies for the implementation to use.

For example, injecting an in-memory repository to remove dependency on a physical database when running unit tests.

```go
func TestSomeHigherLevelFunctionUsingFoo(t *testing T) {
    // ARRANGE
    fooReceiver := &FooReceiver{ 
        Repository: InMemoryRepository,
    }
    reg := mediator.RegisterReceiver[FooData](fooReceiver)
    defer reg.Remove()

    // ACT
    err := SomeHigherLevelFunction()

    // ASSERT
    ..etc
}
```

Or it might mean injecting an entirely different implementation of a receiver or handler, such as a fake, stub, spy or other mock.

Which brings us to...

<br/>

## Fake/Mock Handlers
You can of course implement test `Receiver` and `Handler` implementations as needed.  But out-of-the-box, `go-mediator` provides mock implementations that can be used in most - if not all - common use cases. 

Factory methods are provided to create mocks that can simulate specific return values and/or to determine how many times the handlers are called and with what data or requests, over the course of execution of code under test.

For example, to fake a `FooReceiver` that simply indicates successful completion for _any and all_ data it receives:

```go
    mock, reg := mediator.MockReceiver[FooData]()
```

Other factory functions provide for different use cases:

```go
    // A receiver that returns a given error, for any data it receives
    MockReceiverReturningError[TData](error) (mock, *reg)

    // A receiver that implements Validtor and returns a given error
    // from that validator for any data it receives
    MockReceiverWithValidatorError[TData](error) (mock, *reg)
    
    // A receiver that runs a provided func when Executing() data
    MockReceiverWithFunc[TData](ReceiverFunc[TData]) (mock, *reg)

    // A receiver that implements Validator, using a provided func,
    // and a func that runs when Executing() data
    MockReceiverWithValidator[TData](ReceiverFunc[TData], ValidatorFunc[TData]) (mock, *reg)
```

> <br>
> Similar mocks and factories are provided for `Handler` implementations. 
> <br>
> <br>
<br>

In all cases, the `Mock...()` functions create and register the mock, returning a reference to the mock and the registration reference.  The mock reference can be used to access spy functions in subsequents tests, otherwise it can be ignored.

The registration reference is provided so that the mock can be de-registered:

```go
    // ARRANGE
    mock, reg := mediator.MockReceiver[FooData]()
    defer reg.Remove()

    // ACT
    .. exercise code under test ..

    // ASSERT
    if mock.WasNotCalled() {
        t.Error("'FooData' receiver was not called")
    }
```

<br/>

## De-and Re-Registering Handlers
It is good practice for tests to be self-contained and independent.

This means that when code under test makes `mediator` requests, an appropriate handler *must* be registered by that test *and **removed*** when done, so that other tests can register their own receiver or handler for the same type.

The registration reference returned by registration functions and mock factories is for just this purpose, providing a single method: `Remove()` which removes the registration it references.

So a typical test would start something like this:

```golang
func TestSomethingThatMakesMediatorRequests(t *testing.T) {
    // ARRANGE
    reg := RegisterReceiver[*FooData](&testFoo{})
    defer reg.Remove()

    // ACT

      /* test some code that Perform()s a `FooRequest` that will
         be handled for this test by the testFoo implementation
         (assuming a MockReceiver() wasn't enough in this case) */

    // ASSERT
    
      /* etc... */
}
```

Or, if using a mock receiver:

```golang
func TestSomethingThatMakesMediatorRequests(t *testing.T) {
    // ARRANGE
    _, reg := MockReceiver[FooData]()
    defer reg.Remove()

    // etc...
}
```

<br/>

# Structuring Handler and Receiver Code
><br>_This section is not intended to be prescriptive, only illustrative.  Different use cases might call for different approaches.
<br><br>In particular, more complex registration, for example conditionally registering different handlers based on runtime conditions, would not fit very comfortably within the pattern described here._<br><br>

You may find it useful to organise your code into separate packages for each handler or receiver.  You can then leverage the package scoping syntax of `golang` to provide consistent and concise naming for `TData`, `TRequest`, `TResult`, `Receiver` and `Handler` types throughout your code.

So you might have a "services" package in which each service (receiver or handler) is implemented in its own package:

```
   <myproject root>
   > services
       > createProduct
            createProduct.go
       > getProduct
            getProduct.go
       > placeOrder
            placeOrder.go
       configure.go
```

One of those services might then look something like this:

```golang
package getProduct

import (
    "context"
    "database/sql"

    "github.com/deltics/go-mediator"

    model "myproject/database/models"
)

type Request *struct {
    ProductId string
}

type Result *model.Product

type Handler {
    DB: *sql.DB
}

func (*Handler) Execute(ctx context.Context, Request) (Result, error) {
    // fetch the requested product from the DB
    return product, nil
}

func Register(handler Handler) {
    _ := mediator.RegisterHandler[Request, Result](handler)
}
```

The reference returned by the `mediator` registration function can usually be ignored in production registrations, which do not typically need to be removed or replaced.

><br>_If the handler does not require any services to be injected, the `Register()` function might even initialise the handler itself and not require it to be supplied._<br><br>

The handler could then be registered in a `Configure()` function (in the `configure.go` file of the service package) by calling the `Register()` func exported by each service package:

```golang
    package services

    import (
        "database/sql"

        "myproject/service/createProduct"
        "myproject/service/getProduct"
        "myproject/service/placeOrder"
    )

    func Configure(db *sql.DB) {
        createProduct.Register(&createProduct.Handler{DB: db})
        getProduct.Register(&getProduct.Handler{DB: db})
        placeOrder.Register(&placeOrder.Handler{DB: db})
    }
```

And in main:

```go
    services.Configure()
```

Finally, the calls made through `mediator` would use code similar to:

```golang
    rq := &getProduct.Request{ 
        ProductId: id,
    }
    product, err := mediator.Perform[*getProduct.Request, *getProduct.Result](ctx, rq)
    if err != nil {
        ..etc..
    }
```

Completing the picture, a test involving a scenario where a product did not exist could use one of the mock utilities to create and register a mock handler to return a "not found" error for any product:

```golang
    _, reg := mediator.MockHandlerReturningError[*getProduct.Request, *getProduct.Result](errors.New("not found"))
    defer reg.Remove()
```
