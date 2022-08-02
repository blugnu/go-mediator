package tasks

type result[R any] struct {
	value chan R
	err   chan error
}
