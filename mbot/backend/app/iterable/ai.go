package iterable

import (
	"sync/atomic"
)

// var Itr Iterable

type Iterable struct {
	id atomic.Int64
}

func NewIterable() *Iterable {
	return &Iterable{}
}

func (i *Iterable) ID() int64 {
	return i.id.Add(1)
}

func (i *Iterable) SetIterable(start int64) {
	i.id.Store(start)
}
