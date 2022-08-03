package gofer

import (
	"context"
	"sync"
)

type listener[V comparable, R any] struct {
	task *queue[V, R]
	mx   sync.Mutex
}

// execute is a function that performs a queued request, sending the
// result and any error to the result and error channel for that request.
func (l *listener[V, R]) execute(ctx context.Context, rq V) {
	l.mx.Lock()
	defer l.mx.Unlock()

	// Extract the result and error channels for the request from the
	// task and remove them, since they will no longer be required once
	// the request is complete
	ch := l.task.result[rq]
	delete(l.task.result, rq)

	// Call the handler
	result, err := l.task.handler(ctx, rq)

	// If the err is non-nil send it back over the error channel, otherwise
	// send the result over the result channel.
	if err != nil {
		ch.err <- err
	} else {
		ch.val <- result
	}

}

// start initiates the listener processing queued requests
func (l *listener[V, R]) start() {
	for rq := range l.task.request {
		go l.execute(rq.Ctx, rq.Value)
	}
}
