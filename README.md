# go-mediator
A lightweight implementation of the Mediator Pattern for goLang, inspired by the MediatR framework for .net.

# Concepts
`go-mediator` treats `Commands` and `Queries` differently.  Whilst both accept a `context` and a `Requests` value, a `Command` returns only an `error` whilst a `Query` returns some result value in addition to, or instead of, an `error`.  If `Command` and `Query` were simple funcs their declarations would differ thus:

```go
func Command(context.Context, Request) error
func Query(context.Context, Request) (Result, error)
```

Treating them differently allows code that uses `go-mediator` to benefit from type inference to simply calls made to `Commands` in a way that isn't possible with a `Query`. 

# The Mediator Pattern
The Mediator Pattern is a simple pattern that uses a 3rd-party (the mediator) to facilitate communication between two other parties without those two other parties having knowledge of each other.

It is a powerful pattern for achieving loosely coupled code that might otherwise have tight-coupling in the form of direct function calls.

# How It Works
A Command or Query Handler is registered with the mediator.  Commands and Query handlers are generic interfaces, accepting a request `type` param that the handler is capable of responding to:

```golang
    type CommandHandler[TRequest any] interface {
        Execute(context.Context, TRequest) error
    }

    type QueryHandler[TRequest any, TResult any] interface {
        Execute(context.Context, TRequest) (TResult, error)
    }
```

Handlers are registered by passing an instantiated struct implementing the required interface to a registration function.  Registration functions are generic functions and the go type system therefore ensures that when registering handler for a specific type, the correct request type must be specified on the registration type parameter:

```golang
    RegisterCommandHandler[TRequest](handler CommandHandler[TRequest])

    RegisterQueryHandler[TRequest, TResult](handler QueryHandler[TRequest, TResult])
```

When registering a command handler, go is able to infer the `Request` type, so this need not be explicit:

```golang
    mediator.RegisterCommandHandler(&FooHandler{})
```

Unfortunately this is not (currently?) possible with Query handlers, so registration of those requires that the type parameters be specified on the registration function, e.g:

```golang
    mediator.RegisterQueryHandler[FooRequest, string](&FooHandler{})
```

Handlers are implemented as lightweight interfaces providing an `Execute()` method corresponding to the Command or Query signature accepting a `context` and that Request type.  The request and handler types for the above registration example might look similar to this:

```golang
type FooRequest struct {
    Foo string
}

type FooHandler struct {}

func (*FooHandler) Execute(ctx context.Context, req *FooRequest) error {
    err := FooTheRequest(ctx, req.Foo)
    return err
}
```

Code wishing to have a Command or Query performed, builds a `Request` of the appropriate type and sends that request to the mediator.  Continuing the example above, a requestor would do something similar to:

```golang
    err := mediator.Command(ctx, &FooRequest{ Foo: "something nice" })
```

The mediator identifies the appropriate handler for that request and passes the request on to that handler, returning the result back to the original requestor.

The above example illustrate that the go type system allows go to 

# What It Is NOT
- Mediator is not a message queue
- Mediator is not asynchronous
- Mediator is not complicated!