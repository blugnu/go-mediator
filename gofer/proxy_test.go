package gofer

import (
	"context"
	"testing"
)

type nothing struct{}

func void(ctx context.Context, i int) (*nothing, error) { return nil, nil }

func TestThatABufferedTaskReturnsNilWhenBufferIsFull(t *testing.T) {
	// ARRANGE
	Courier := Queue(void, 3)
	ctx := context.Background()

	// ACT
	// With no listener started the Enqueue()ing will fill the buffer
	c1, _ := Courier.Enqueue(ctx, 1)
	c2, _ := Courier.Enqueue(ctx, 2)
	c3, _ := Courier.Enqueue(ctx, 3)
	c4, _ := Courier.Enqueue(ctx, 4)

	// ASSERT
	if c1 == nil || c2 == nil || c3 == nil {
		t.Errorf("Wanted 3 non-nil channels, got ch1: %x, ch1: %x, ch3: %x", c1, c2, c3)
	}

	if c4 != nil {
		t.Errorf("Wanted nil channel, got %x", c4)
	}
}
