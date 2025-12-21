package iterable

import (
	"sync/atomic"
)

var Itr Iterable

type Iterable struct {
	id atomic.Int64
}

// this should start the id at the max of what's in the tasks table
func (i *Iterable) SetCurrentId() {
	//get max id from tasks table
	// set a.id to that
}

func (i *Iterable) ID() int64 {
	return i.id.Add(1)
}
