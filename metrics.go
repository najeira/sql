package sql

import (
	"sync"
	"time"

	mt "github.com/rcrowley/go-metrics"
)

var (
	metrics *metricsDB
	timers  *timersDB
)

func init() {
	metrics = &metricsDB{
		connections: mt.NewHistogram(mt.NewExpDecaySample(1028, 0.015)),
		queries:     mt.NewMeter(),
		executes:    mt.NewMeter(),
		rows:        mt.NewMeter(),
		affects:     mt.NewMeter(),
	}
	timers = &timersDB{timers: make(map[string]mt.Timer)}
}

func GetMetrics() map[string]float64 {
	return metrics.Get()
}

func GetTimers() map[string]map[string]float64 {
	return timers.Get()
}

type metricsDB struct {
	connections mt.Histogram
	queries     mt.Meter
	executes    mt.Meter
	rows        mt.Meter
	affects     mt.Meter
}

func (m *metricsDB) MarkQueries(v int) {
	if v != 0 {
		m.queries.Mark(int64(v))
	}
}

func (m *metricsDB) MarkExecutes(v int) {
	if v != 0 {
		m.executes.Mark(int64(v))
	}
}

func (m *metricsDB) MarkRows(v int) {
	if v != 0 {
		m.rows.Mark(int64(v))
	}
}

func (m *metricsDB) MarkAffects(v int) {
	if v != 0 {
		m.affects.Mark(int64(v))
	}
}

func (m *metricsDB) MarkConnections(v int) {
	m.connections.Update(int64(v))
}

func (m *metricsDB) Get() map[string]float64 {
	return map[string]float64{
		"connections_min": float64(m.connections.Min()),
		"connections_max": float64(m.connections.Max()),
		"connections_avg": m.connections.Mean(),
		"queries_count":   float64(m.queries.Count()),
		"queries_rate":    m.queries.Rate1(),
		"executes_count":  float64(m.executes.Count()),
		"executes_rate":   m.executes.Rate1(),
		"rows_count":      float64(m.rows.Count()),
		"rows_rate":       m.rows.Rate1(),
		"affects_count":   float64(m.affects.Count()),
		"affects_rate":    m.affects.Rate1(),
	}
}

type timersDB struct {
	timers map[string]mt.Timer
	mu     sync.RWMutex
}

func (m *timersDB) Measure(key string, start time.Time) {
	elapsed := time.Now().Sub(start)

	m.mu.RLock()
	t, ok := m.timers[key]
	m.mu.RUnlock()

	if !ok {
		m.mu.Lock()
		t, ok = m.timers[key]
		if !ok {
			t = mt.NewTimer()
			m.timers[key] = t
		}
		m.mu.Unlock()
	}

	t.Update(elapsed)
}

func (m *timersDB) Get() map[string]map[string]float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string]map[string]float64)
	for query, timer := range m.timers {
		result[query] = map[string]float64{
			"count": float64(timer.Count()),
			"min":   float64(timer.Min()) / float64(time.Millisecond),
			"max":   float64(timer.Max()) / float64(time.Millisecond),
			"avg":   timer.Mean() / float64(time.Millisecond),
			"rate":  timer.Rate1(),
			"p50":   timer.Percentile(0.5) / float64(time.Millisecond),
			"p75":   timer.Percentile(0.75) / float64(time.Millisecond),
			"p95":   timer.Percentile(0.95) / float64(time.Millisecond),
			"p99":   timer.Percentile(0.99) / float64(time.Millisecond),
		}
	}
	return result
}
