package sql

import (
	"sync"
)

var (
	stringPool sync.Pool
	intPool    sync.Pool
	floatPool  sync.Pool
	boolPool   sync.Pool
	valuesPool sync.Pool

	poolCounter *counter
)

func init() {
	poolCounter = &counter{}
}

type values struct {
	inuse []interface{}
}

func getValues() *values {
	poolCounter.Inc(1)
	if v := valuesPool.Get(); v != nil {
		return v.(*values)
	}
	return &values{
		inuse: make([]interface{}, 0, 64),
	}
}

func (p *values) String() *NullString {
	v := getString()
	p.inuse = append(p.inuse, v)
	return v
}

func (p *values) Int64() *NullInt64 {
	v := getInt64()
	p.inuse = append(p.inuse, v)
	return v
}

func (p *values) Float64() *NullFloat64 {
	v := getFloat64()
	p.inuse = append(p.inuse, v)
	return v
}

func (p *values) Bool() *NullBool {
	v := getBool()
	p.inuse = append(p.inuse, v)
	return v
}

func (p *values) Close() error {
	if len(p.inuse) > 0 {
		if logv(logDebug) {
			logf("sql: pool %d values", len(p.inuse))
		}
		for _, v := range p.inuse {
			poolValue(v)
		}
		p.inuse = p.inuse[:0]
	}
	valuesPool.Put(p)
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
	if v == nil {
		return
	}
	v.Valid = false
	v.String = ""
	stringPool.Put(v)
	poolCounter.Dec(1)
}

func poolInt64(v *NullInt64) {
	if v == nil {
		return
	}
	v.Valid = false
	v.Int64 = 0
	intPool.Put(v)
	poolCounter.Dec(1)
}

func poolFloat64(v *NullFloat64) {
	if v == nil {
		return
	}
	v.Valid = false
	v.Float64 = 0
	floatPool.Put(v)
	poolCounter.Dec(1)
}

func poolBool(v *NullBool) {
	if v == nil {
		return
	}
	v.Valid = false
	v.Bool = false
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
