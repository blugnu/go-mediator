package courier

import (
	"context"
	"testing"

	"github.com/deltics/go-tasks/mock"
)

func TestThatEnqueuingACourierReturnsNilWhenTheQueueIsFull(t *testing.T) {
	// ARRANGE
	c := Queue(void, 3)
	ctx := context.Background()

	// ACT
	// With no listener started the Enqueue()ing will fill the buffer
	c1 := c.Enqueue(ctx, 1)
	c2 := c.Enqueue(ctx, 2)
	c3 := c.Enqueue(ctx, 3)
	c4 := c.Enqueue(ctx, 4)

	// ASSERT
	if c1 == nil || c2 == nil || c3 == nil {
		t.Errorf("Wanted 3 non-nil channels, got ch1: %x, ch1: %x, ch3: %x", c1, c2, c3)
	}

	if c4 != nil {
		t.Errorf("Wanted nil channel, got %x", c4)
	}
}

func TestThatAQueueingCourierCanBeMockedMultipleTimes(t *testing.T) {
	// ARRANGE
	x = 99
	a = 0
	b = 0
	c := Queue(incX, 3)
	ct := c.(mock.Courier[int])
	ctx := context.Background()

	go c.StartListener()

	// ACT
	ec := c.Enqueue(ctx, 1)
	<-ec
	if x != 100 {
		t.Errorf("Default handler not invoked on first call")
	}
	ct.Use(incA)
	ec = c.Enqueue(ctx, 1)
	<-ec
	if a != 1 {
		t.Errorf("StartListenerWith(incA) did not install incA handler")
	}

	ct.Use(incB)
	ec = c.Enqueue(ctx, 1)
	<-ec
	if b != 1 {
		t.Errorf("StartListenerWith(incB) did not install incB handler")
	}

	ct.UseDefault()
	ec = c.Enqueue(ctx, -1)
	<-ec
	if x != 99 {
		t.Errorf("StartListener() did not install default handler")
	}
}
