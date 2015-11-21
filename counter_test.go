package sql

import (
	"testing"

	mt "github.com/rcrowley/go-metrics"
)

func TestCounter(t *testing.T) {
	c := mt.NewCounter()
	if c.Count() != 0 {
		t.Errorf("counter.count: expected %d, got %d", 0, c.Count())
	}

	c.Inc(1)
	if c.Count() != 1 {
		t.Errorf("counter.Inc: expected %d, got %d", 1, c.Count())
	}

	c.Inc(-2)
	if c.Count() != -1 {
		t.Errorf("counter.Inc: expected %d, got %d", -1, c.Count())
	}

	c.Dec(1)
	if c.Count() != -2 {
		t.Errorf("counter.Dec: expected %d, got %d", -2, c.Count())
	}

	c.Dec(-3)
	if c.Count() != 1 {
		t.Errorf("counter.Dec: expected %d, got %d", 1, c.Count())
	}

	c.Clear()
	if c.Count() != 0 {
		t.Errorf("counter.Clear: expected %d, got %d", 0, c.Count())
	}
}
