package tasks

import (
	"testing"
)

func void(i int) error { return nil }

func TestThatABufferedTaskReturnsNilWhenBufferIsFull(t *testing.T) {
	// ARRANGE
	errand := BufferedErrand(void, 3)

	// ACT
	// With no listener started the Enqueue()ing will fill the buffer
	c1 := errand.Enqueue(1)
	c2 := errand.Enqueue(2)
	c3 := errand.Enqueue(3)
	c4 := errand.Enqueue(4)

	// ASSERT
	if c1 == nil || c2 == nil || c3 == nil {
		t.Errorf("Wanted 3 non-nil channels, got ch1: %x, ch1: %x, ch3: %x", c1, c2, c3)
	}

	if c4 != nil {
		t.Errorf("Wanted nil channel, got %x", c4)
	}
}
