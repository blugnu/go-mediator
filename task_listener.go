package tasks

type taskListener[T comparable, R any] struct {
	handler func(T) (R, error)
	task    *Task[T, R]
}

// execute is a function that performs a queued request, sending the
// result and any error to the result and error channel for that request.
func (l *taskListener[T, R]) execute(rq T) {
	// Extract the result and error channels for the request from the
	// task and remove them, since they will no longer be required once
	// the request is complete
	channels := l.task.result[rq]
	delete(l.task.result, rq)

	errChan := channels.err
	resultChan := channels.value

	// Call the handler
	result, err := l.handler(rq)

	// Send the result and any error to the appropriate channels
	resultChan <- result
	errChan <- err

}

// start initiates the listener processing queued requests
func (l *taskListener[T, R]) start() {
	for rq := range l.task.queue {
		go l.execute(rq)
	}
}
