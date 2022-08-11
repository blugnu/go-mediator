# go-mediator
A lightweight implementation of the Mediator Pattern for goLang, inspired by the MediatR framework for .net.

## Mediator Pattern
[The Mediator Pattern](https://en.wikipedia.org/wiki/Mediator_pattern) is a simple pattern that uses a 3rd-party (the mediator) to facilitate communication between two other parties without them requiring knowledge of each other.

It is a powerful pattern for achieving loosely-coupled code.

There are many ways to implement a mediator, from simple `func` pointers to sophisticated and complex messaging systems.

`go-mediator` sits firmly at the *simple* end of that spectrum.

## What go-mediator Is NOT
- `go-mediator` is not a message queue
- `go-mediator` is not asynchronous
- `go-mediator` is not complicated!

# Concepts
`go-mediator` maintains a registry of handlers for specific request types.

Requests are submitted to the `mediator` which consults a registry to locate the handler for the type of request involved which is called with the request and the results passed back to the original caller.

`go-mediator` maintains separate handler registries for `Command` handlers and `Query` handlers.

A `Command` returns only an error, while a `Query` returns a value *and* an error:

```go
func Command(context.Context, Request) error
func Query(context.Context, Request) (Result, error)
```

Treating them differently allows code that uses `go-mediator` to benefit from type inference to simplify calls made to `Commands` in a way that isn't (currently?) possible with a `Query`.

Handlers may (optionally) choose to implement validation of requests as a separate concern from the primary handler execution.

# Getting Started

## 1. Define a Request Type
Only one handler can be registered for any request type, so even if you have multiple requests that accept the same values you will need a distinct request type for each one.

```go
    type FooRequest struct {
        Foo string
    }
```

An exception to the "*one handler per request type*" rule is that you can have separate `Command` and `Query` handlers sharing the same `request` type.  This is possible since there are separate registries for each, which cannot be confused.

**NOTE: A future update is being considered to also allow Request types to be shared by Query Handlers having different result types.** 

## 2. Implement a Handler Interface
Handlers are generic interfaces with a single `Execute` method accepting a `context` and the `request` value.

Both `Command` and `Query` handlers accept a `TRequest` type parameter identifying the *request* type.  A `Query` handler interface additionally requires a *`TResult`* type parameter:

```golang
    // The CommandHandler interface...
    type CommandHandler[TRequest any] interface {
        Execute(context.Context, TRequest) error
    }

    // and the QueryHandler interface...
    type QueryHandler[TRequest any, TResult any] interface {
        Execute(context.Context, TRequest) (TResult, error)
    }
```

An example command handler implementation might look similar to:

```golang
    // FooHandler is a CommandHandler (returns `error`)
    type FooHandler struct {}

    func (*FooHandler) Execute(ctx context.Context, req FooRequest) error {
        // Do some foo'ing
        return nil
    }
```

**NOTE: Since handlers are interface implementations with an underlying struct, this may be used to hold services used by the handler. This cab be useful for substituting fake or mock services when registering handlers in test code.  See more on testing further below...**

## 3. Register the Handler
Handlers are registered by passing an implementation of the handler interface to the appropriate registration function.

The registration functions are generic functions with the same `TRequest` and `TResult` type parameters (for `Query` handlers) as the corresponding handlers.  The go type system is then able to ensure that the request types correspond:

```golang
    RegisterCommandHandler[FooRequest](&FooHandler{})

    RegisterQueryHandler[FooRequest, string](&FooQueryHandler{})
```

## 4. (Optional) Implement a RequestValidator
In addition to the `Execute()` method of the `Command` or `Query` handler interface, handlers may also choose to implement the `RequestValidator` interface:

```golang
    type RequestValidator[TRequest any] interface {
	    Validate(context.Context, TRequest) error
    }
```

**If** implemented by a handler, this will be called by `mediator` *before* the `Execute()` method.  If the `Validate()` method returns an error then the `Execute` method will not be called; and the validation error is returned to the caller as n `ErrBadRequest`.

## 5. Send Requests to Handlers Via Mediator
To call a `Command` or `Query`, simply construct a request of the type required and pass it to mediator using either the `Command()` or `Query()` function, according to the nature of the handler you expect to respond to the request:

```golang
    err := mediator.Command(ctx, &FooRequest{ Foo: "something nice" })

    err := mediator.Query[FooRequest, string](ctx, &FooRequest{ Foo: "something nice" })
```

Notice that for `Command` requests the request type does not need to be specified - go is able to infer the type from the request parameter itself.

This is not possible for `Query` requests since the result type is not represented in the call.  Hence for `Query` requests bot the request and result type parameters must be identified.

# Alternative Result Handling
Results other than errors may be returned by a `Command` handler by using a pointed to a struct for the request type.

The handler may then manipulate the members of the struct, including ones provided explicitly for the purpose of "returning" a value.

This avoids the need to specify request and result types when calling a Query via the mediator but at the expense of losing the usual go result, error pattern when calling functions:

```golang
    result, err := mediator.Query[FooRequest, string](FooRequest{})
    if err != nil {
        log.Error(err)
        return err
    }
```

Might become something similar to:

```golang
    request := &FooRequest{}
    err := mediator.Query[*FooRequest, string](request)
    if err != nil {
        log.Error(err)
        return err
    }
    result := request.Result
```

This is not to either recommend or condemn such an approach, merely to highlight that it is possible.

Note however that careful attention must be paid to such handlers receiving their requests *by reference*.  If passed *by value*, any updates to the request members will be applied to a *copy* of the request, **not** the one held by the original caller. 

# Testing With Mediator

## Handler Dependency Injection
As mentioned above, since a handler is implemented as a struct, this may be used to hold references to services required by the handler implementation.  For example a repository interface:

```golang
    type FooHandler struct {
        Repository   FooRepository
    }
```

Since handlers are usually registered as part of the initialisation of your application or service (e.g. in or by `main.go`), they are not *usually* registered when test code runs and so must be explicitly registered in tests.

This provides an opportunity to '*inject*' alternative services for a handler to use when running under test conditions.

For example subsituting an in-memory repository to remove dependency on a physical database when running tests.

## De-and Re-Registering Handlers
It is good practice for tests to be self-contained and independent.

This means that when code under test makes mediator requests, an appropriate handler must be registered by that test and *removed* when done (so that other tests can register their own handlers for the same request type if necessary).

To achieve this, the handler registration functions return a registration reference which provides a single method: `Remove()`.

```golang
func TestSomethingThatMakesMediatorRequests(t *testing.T) {
    // ARRANGE
    reg := RegisterCommandHandler[*FooRequest](&FakeFooHandler{})
    defer reg.Remove()

    // ACT
    ...

    // ASSERT
    ...
}
```

## Fake Handlers
Ss well as providing for the injection of test services into a production handler, the decoupling of handlers and requestors also makes it very easy to use fakes, spys and other mock handlers in place of 'production' handlers for the purposes of testing.
