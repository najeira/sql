package sql

import (
	"sync/atomic"
)

type counter struct {
	count int64
}

func (c *counter) Clear() {
	atomic.StoreInt64(&c.count, 0)
}

func (c *counter) Count() int64 {
	return atomic.LoadInt64(&c.count)
}

func (c *counter) Dec(i int64) {
	atomic.AddInt64(&c.count, -i)
}

func (c *counter) Inc(i int64) {
	atomic.AddInt64(&c.count, i)
}
