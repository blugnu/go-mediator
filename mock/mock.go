package mock

import "github.com/deltics/go-tasks"

type Courier[V comparable] interface {
	Use(tasks.Courier[V])
	UseDefault()
}
