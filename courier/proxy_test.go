package courier

import (
	"context"
	"testing"

	"github.com/deltics/go-tasks/mock"
)

func void(ctx context.Context, i int) error { return nil }

var (
	a int = 0
	b int = 0
	x int = 99
)

func incA(ctx context.Context, n int) error {
	a += n
	return nil
}

func incB(ctx context.Context, n int) error {
	b += n
	return nil
}

func incX(ctx context.Context, n int) error {
	x += n
	return nil
}

func TestThatHandlerFuncsCanBeReplacedDuringATestRun(t *testing.T) {
	// ARRANGE
	x = 99
	a = 0
	b = 0

	c := Proxy(incX)
	ct := c.(mock.Courier[int])
	ctx := context.Background()

	// ACT
	c.CallWith(ctx, 1)
	if x != 100 {
		t.Errorf("Default handler not invoked on first call")
	}

	ct.Use(incA)
	c.CallWith(ctx, 1)
	if a != 1 {
		t.Errorf("Use(incA) did not install incA handler")
	}

	ct.Use(incB)
	c.CallWith(ctx, 1)
	if b != 1 {
		t.Errorf("Use(incB) did not install incB handler")
	}

	ct.UseDefault()
	c.CallWith(ctx, -1)
	if x != 99 {
		t.Errorf("UseDefalut() not restore default handler")
	}
}
