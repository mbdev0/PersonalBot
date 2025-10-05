package datastructures

import (
	"context"
	"sync"
	"time"
)

type Map[K comparable, V any] struct {
	sync.Mutex

	m    map[K]V
	subs map[K][]chan V
}

func NewMap[K comparable, V any]() Map[K, V] {
	return Map[K, V]{
		m:    map[K]V{},
		subs: map[K][]chan V{},
	}
}

func (m *Map[K, V]) Set(key K, value V) {
	m.Lock()
	defer m.Unlock()

	m.m[key] = value

	for _, sub := range m.subs[key] {
		sub <- value
	}

	delete(m.subs, key)
}

func (m *Map[K, V]) Get(key K) (value V, ok bool) {
	val, ok := m.m[key]
	if !ok {
		return val, false
	}

	return val, true
}

func (m *Map[K, V]) Delete(key K) {
	delete(m.m, key)
}

func (m *Map[K, V]) WaitForEntry(key K) (val V, ok bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	m.Lock()

	value, ok := m.m[key]
	if ok {
		m.Unlock()
		return value, true
	}

	ch := make(chan V, 1) // Buffered to prevent goroutine leak
	m.subs[key] = append(m.subs[key], ch)
	m.Unlock()

	select {
	case value := <-ch:
		return value, true
	case <-ctx.Done():
		// Clean up our subscription on timeout
		m.Lock()
		subs := m.subs[key]
		for i, sub := range subs {
			if sub == ch {
				m.subs[key] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
		if len(m.subs[key]) == 0 {
			delete(m.subs, key)
		}
		m.Unlock()

		var zero V
		return zero, false
	}
}
