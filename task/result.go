package task

type result[R any] struct {
	value chan R
	err   chan error
}
