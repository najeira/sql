package sql

import (
	"sync"

	metrics "github.com/rcrowley/go-metrics"
)

var (
	stringPool     sync.Pool
	intPool        sync.Pool
	floatPool      sync.Pool
	boolPool       sync.Pool
	valuesPoolPool sync.Pool

	poolCounter metrics.Counter
)

func init() {
	poolCounter = metrics.NewCounter()
}

type valuesPool struct {
	values []interface{}
}

func getValuesPool() *valuesPool {
	poolCounter.Inc(1)
	if v := valuesPoolPool.Get(); v != nil {
		return v.(*valuesPool)
	}
	return &valuesPool{
		values: make([]interface{}, 0, 1024),
	}
}

func (p *valuesPool) String() *NullString {
	v := getString()
	p.values = append(p.values, v)
	return v
}

func (p *valuesPool) Int64() *NullInt64 {
	v := getInt64()
	p.values = append(p.values, v)
	return v
}

func (p *valuesPool) Float64() *NullFloat64 {
	v := getFloat64()
	p.values = append(p.values, v)
	return v
}

func (p *valuesPool) Bool() *NullBool {
	v := getBool()
	p.values = append(p.values, v)
	return v
}

func (p *valuesPool) Close() error {
	if logv(logDebug) && len(p.values) > 0 {
		logf("sql: pool %d values", len(p.values))
	}
	for _, v := range p.values {
		poolValue(v)
	}
	p.values = p.values[:0]
	valuesPoolPool.Put(p)
	poolCounter.Dec(1)
	return nil
}

func getString() *NullString {
	poolCounter.Inc(1)
	if v := stringPool.Get(); v != nil {
		return v.(*NullString)
	}
	return &NullString{}
}

func getInt64() *NullInt64 {
	poolCounter.Inc(1)
	if v := intPool.Get(); v != nil {
		return v.(*NullInt64)
	}
	return &NullInt64{}
}

func getFloat64() *NullFloat64 {
	poolCounter.Inc(1)
	if v := floatPool.Get(); v != nil {
		return v.(*NullFloat64)
	}
	return &NullFloat64{}
}

func getBool() *NullBool {
	poolCounter.Inc(1)
	if v := boolPool.Get(); v != nil {
		return v.(*NullBool)
	}
	return &NullBool{}
}

func poolString(v *NullString) {
	stringPool.Put(v)
	poolCounter.Dec(1)
}

func poolInt64(v *NullInt64) {
	intPool.Put(v)
	poolCounter.Dec(1)
}

func poolFloat64(v *NullFloat64) {
	floatPool.Put(v)
	poolCounter.Dec(1)
}

func poolBool(v *NullBool) {
	boolPool.Put(v)
	poolCounter.Dec(1)
}

func poolValue(v interface{}) {
	if v == nil {
		return
	}
	switch x := v.(type) {
	case *NullString:
		poolString(x)
	case *NullInt64:
		poolInt64(x)
	case *NullFloat64:
		poolFloat64(x)
	case *NullBool:
		poolBool(x)
	}
}

func CountPool() int64 {
	n := poolCounter.Count()
	if n < 0 {
		panic("fatal")
	}
	return n
}
