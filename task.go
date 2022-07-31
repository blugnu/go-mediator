package tasks

type Task[T comparable] struct {
	handler    func(T) error
	queue      chan T
	completion map[T]chan error
}

func Buffered[T comparable](handler func(T) error, n int) Task[T] {
	return Task[T]{
		handler:    handler,
		queue:      make(chan T, n),
		completion: make(map[T]chan error),
	}
}

func Unbuffered[T comparable](handler func(T) error) Task[T] {
	return Task[T]{
		handler:    handler,
		queue:      make(chan T),
		completion: make(map[T]chan error),
	}
}

func (w *Task[T]) Enqueue(rq T) chan error {
	w.completion[rq] = make(chan error)
	w.queue <- rq
	return w.completion[rq]
}

func (w *Task[T]) execute(rq T) {
	completion := w.completion[rq]
	delete(w.completion, rq)
	completion <- w.handler(rq)
}

func (w *Task[T]) StartListening() {
	for rq := range w.queue {
		go w.execute(rq)
	}
}

func (w *Task[T]) StartListeningWith(fn func(T) error) {
	w.handler = fn
	go w.StartListening()
}
