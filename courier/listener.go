package courier

import (
	"context"
	"sync"
)

type listener[V comparable] struct {
	task *queue[V]
	mx   sync.Mutex
}

func (l *listener[V]) execute(ctx context.Context, v V) {
	l.mx.Lock()
	defer l.mx.Unlock()

	completion := l.task.result[v]
	delete(l.task.result, v)
	err := l.task.handler(ctx, v)
	completion <- err
}

func (l *listener[V]) start() {
	for rq := range l.task.request {
		go l.execute(rq.Ctx, rq.Value)
	}
}
