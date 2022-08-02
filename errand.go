package tasks

// Errand[T] represents a task satisfied by a handler function accepting a
// single argument of type `T` and returning `error`.
type Errand[T comparable] struct {
	handler    func(T) error
	queue      chan T
	completion map[T]chan error
}

// Buffered returns a Task[T] with the specified handler and a buffered
// request channel of capacity `n`.
func BufferedErrand[T comparable](handler func(T) error, n int) Errand[T] {
	return Errand[T]{
		handler:    handler,
		queue:      make(chan T, n),
		completion: make(map[T]chan error),
	}
}

// Unbuffered returns a Task[T] with the specified handler and an
// unbuffered request channel.
func UnbufferedErrand[T comparable](handler func(T) error) Errand[T] {
	return Errand[T]{
		handler:    handler,
		queue:      make(chan T),
		completion: make(map[T]chan error),
	}
}

// Enqueue queues a request to process the value `rq`.  If successfully
// queued, a `chan error` is returned which will receive the `error`
// result from the task upon completion (including `nil` if the
// request completes without error).
//
// If the task cannot be queued a `nil` channel is returned.
func (task *Errand[T]) Enqueue(rq T) chan error {
	select {
	case task.queue <- rq:
		task.completion[rq] = make(chan error)
		return task.completion[rq]
	default:
		return nil
	}
}

// StartListener starts a listener using the task handler.
func (task *Errand[T]) StartListener() {
	task.StartListenerWith(task.handler)
}

// StartListenerWith starts a listener using the specified handler instead
//  of the task handler itself.
//
// This is typically used to supply mock handlers when testing code that
// makes requests of the task.
func (task *Errand[T]) StartListenerWith(handler func(T) error) {
	l := errandListener[T]{
		handler: handler,
		task:    task,
	}
	l.start()
}
