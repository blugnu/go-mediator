package tasks

type errandListener[T comparable] struct {
	handler func(T) error
	task    *Errand[T]
}

// execute is a function that performs a queued request and sends the
// result to the completion channel allocated for that request.
//
// It is performed as a goroutine by a Listening goroutine.
func (l *errandListener[T]) execute(rq T) {
	completion := l.task.completion[rq]
	delete(l.task.completion, rq)
	completion <- l.handler(rq)
}

func (l *errandListener[T]) start() {
	for rq := range l.task.queue {
		go l.execute(rq)
	}
}
