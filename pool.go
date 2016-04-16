package sql

import (
	"sync"
)

var (
	stringPool sync.Pool
	intPool    sync.Pool
	floatPool  sync.Pool
	boolPool   sync.Pool
	poolPool   sync.Pool
)

type Pool struct {
	inuse []interface{}
}

func NewPool() *Pool {
	if v := poolPool.Get(); v != nil {
		return v.(*Pool)
	}
	return &Pool{
		inuse: make([]interface{}, 0, 64),
	}
}

func (p *Pool) String() *NullString {
	v := getString()
	p.inuse = append(p.inuse, v)
	return v
}

func (p *Pool) Int64() *NullInt64 {
	v := getInt64()
	p.inuse = append(p.inuse, v)
	return v
}

func (p *Pool) Float64() *NullFloat64 {
	v := getFloat64()
	p.inuse = append(p.inuse, v)
	return v
}

func (p *Pool) Bool() *NullBool {
	v := getBool()
	p.inuse = append(p.inuse, v)
	return v
}

func (p *Pool) Close() {
	if len(p.inuse) > 0 {
		//debugf("sql: pool %d Pool", len(p.inuse))
		for _, v := range p.inuse {
			poolValue(v)
		}
		p.inuse = p.inuse[:0]
	}
	poolPool.Put(p)
}

func getString() *NullString {
	if v := stringPool.Get(); v != nil {
		return v.(*NullString)
	}
	return &NullString{}
}

func getInt64() *NullInt64 {
	if v := intPool.Get(); v != nil {
		return v.(*NullInt64)
	}
	return &NullInt64{}
}

func getFloat64() *NullFloat64 {
	if v := floatPool.Get(); v != nil {
		return v.(*NullFloat64)
	}
	return &NullFloat64{}
}

func getBool() *NullBool {
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
}

func poolInt64(v *NullInt64) {
	if v == nil {
		return
	}
	v.Valid = false
	v.Int64 = 0
	intPool.Put(v)
}

func poolFloat64(v *NullFloat64) {
	if v == nil {
		return
	}
	v.Valid = false
	v.Float64 = 0
	floatPool.Put(v)
}

func poolBool(v *NullBool) {
	if v == nil {
		return
	}
	v.Valid = false
	v.Bool = false
	boolPool.Put(v)
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
