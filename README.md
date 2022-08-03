# go-tasks

A simple package for running tasks (functions) over channels.  Task handlers are simple functions.

## Basic Operation

Tasks are created by implementing a function and initialising a variable with either a `Buffered` or `Unbuffered` queue:

```golang
package foo

import "github.com/deltics/go-tasks"

var Task = task.Buffered(fooFunc)

func fooFunc(i int) error {
    // Do something magical with `i`
}
```
Tasks are invoked by enqueuing values to be processed by a task.  The result of enqueuing a request is a channel over which any error is received, or a result channel and an error channel:

```golang
    // When a Task yields a result and/or an error
    resultChan, errChan := foo.Task.Enqueue(42)

    // NOTE: select{} is potentially problematic here since it will yield
    //  ONLY the result OR the error, depending on which is received first.
    result := <-resultChan
    err := <-errChan
    if err != nil {
        log.Errorf("Unexpected error: %s", err)
    }

    // When a Task yields only an error
    errChan := foo.Task.Enqueue(42)
    err := <-errChan
    if err != nil {
        log.Errorf("Unexpected error: %s", err)
    }
```

**IMPORTANT:** `Enqueue()` is **non**-blocking.  If the request cannot be enqueued (e.g a buffered task queue is full) then `nil` channels are returned.

To queue a request and block until that request has been enqueued, use `MustEnqueue()`.

---

### Q: Why?
Why run a function via a channel rather than simply calling it directly?

1. Loose-coupling; e.g. between service ingress and functions that are performed in response to requests over that ingress

2. Improve testability; ingress behaviour can be tested using a mocked handler for any required tasks 

3. Scalability; optionally scaling-out functions by launching multiple listeners

### Required Go Version

This package uses generics and therefore requires GoLang 1.18 or later.


### Terminology


| Term | Meaning |
| ---- | ------- |
| Courier | A task performed by a function that returns only an error (or nil) |
| task | A task performed by a function that returns some value _and_ and error (or nil).<br><br>`Task` may also be used to refer to an `Courier` where the task has already been established to _be_ an Courier.|

In simple terms: an `Courier` is a `task` but not all `tasks` are `Couriers`. |

---
## Limitations

- Task functions must accept only a single argument and return an error (or nil) **or** some value _and_ an error (or nil).

- Separate packages are provided for Couriers and tasks
