package sql

import (
	"testing"
)

func TestCounter(t *testing.T) {
	c := &counter{}
	if c.count != 0 {
		t.Errorf("counter.count: expected %d, got %d", 0, c.count)
	}

	if c.count != c.Count() {
		t.Errorf("counter.Count: expected %d, got %d", c.count, c.Count())
	}

	c.Inc(1)
	if c.count != 1 {
		t.Errorf("counter.Inc: expected %d, got %d", 1, c.count)
	}

	c.Inc(-2)
	if c.count != -1 {
		t.Errorf("counter.Inc: expected %d, got %d", -1, c.count)
	}

	c.Dec(1)
	if c.count != -2 {
		t.Errorf("counter.Dec: expected %d, got %d", -2, c.count)
	}

	c.Dec(-3)
	if c.count != 1 {
		t.Errorf("counter.Dec: expected %d, got %d", 1, c.count)
	}

	c.Clear()
	if c.count != 0 {
		t.Errorf("counter.Clear: expected %d, got %d", 0, c.count)
	}
}
