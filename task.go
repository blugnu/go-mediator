package tasks

type Task[T comparable, R any] struct {
	handler func(T) (R, error)
	queue   chan T
	result  map[T]*result[R]
}

func BufferedTask[T comparable, R any](handler func(T) (R, error), n int) Task[T, R] {
	return Task[T, R]{
		handler: handler,
		queue:   make(chan T, n),
		result:  make(map[T]*result[R]),
	}
}

func UnbufferedTask[T comparable, R any](handler func(T) (R, error)) Task[T, R] {
	return Task[T, R]{
		handler: handler,
		queue:   make(chan T),
	}
}

func (task *Task[T, R]) Enqueue(rq T) (chan R, chan error) {
	result := &result[R]{
		value: make(chan R),
		err:   make(chan error),
	}
	task.result[rq] = result
	task.queue <- rq

	return result.value, result.err
}

func (task *Task[T, R]) StartListener() {
	task.StartListenerWith(task.handler)
}

func (task *Task[T, R]) StartListenerWith(handler func(T) (R, error)) {
	l := taskListener[T, R]{
		handler: handler,
		task:    task,
	}
	l.start()
}
