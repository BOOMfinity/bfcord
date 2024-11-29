package utils

import "sync"

type safeMapImpl[K comparable, T any] struct {
	data map[K]T
	mut  sync.RWMutex
}

func (m *safeMapImpl[K, T]) Get(key K) (T, bool) {
	m.mut.RLock()
	obj, ok := m.data[key]
	m.mut.RUnlock()
	return obj, ok
}

func (m *safeMapImpl[K, T]) Has(key K) (found bool) {
	m.mut.RLock()
	_, ok := m.data[key]
	m.mut.RUnlock()
	return ok
}

func (m *safeMapImpl[K, T]) Size() int {
	m.mut.RLock()
	size := len(m.data)
	m.mut.RUnlock()
	return size
}
func (m *safeMapImpl[K, T]) Set(key K, value T) {
	m.mut.Lock()
	m.data[key] = value
	m.mut.Unlock()
}

func (m *safeMapImpl[K, T]) Delete(key K) (ok bool) {
	m.mut.Lock()
	_, ok = m.data[key]
	delete(m.data, key)
	m.mut.Unlock()
	return
}

func (m *safeMapImpl[K, T]) Search(fn SafeMapLambda[T]) (values []T) {
	m.mut.RLock()
	values = make([]T, 0, len(m.data))
	for _, ptr := range m.data {
		if fn(ptr) {
			values = append(values, ptr)
		}
	}
	m.mut.RUnlock()
	return
}

func (m *safeMapImpl[K, T]) Each(fn SafeMapLambda[T]) {
	m.mut.RLock()
	for _, ptr := range m.data {
		if !fn(ptr) {
			break
		}
	}
	m.mut.RUnlock()
}

func (m *safeMapImpl[K, T]) Find(fn SafeMapLambda[T]) (obj T, ok bool) {
	m.mut.RLock()
	for _, v := range m.data {
		if fn(v) {
			return v, true
		}
	}
	m.mut.RUnlock()
	return
}

func NewSafeMap[K comparable, T any]() SafeMap[K, T] {
	return &safeMapImpl[K, T]{
		data: map[K]T{},
	}
}

type safeMapEmbedded[K comparable, T any] struct {
	data SafeMap[K, T]
	gen  SafeMapGenerator[T]
}

func (e *safeMapEmbedded[K, T]) Get(key K) T {
	value, ok := e.data.Get(key)
	if ok {
		return value
	}
	gen := e.gen()
	e.data.Set(key, gen)
	return gen
}

func (e *safeMapEmbedded[K, T]) Delete(key K) bool {
	return e.data.Delete(key)
}

func (e *safeMapEmbedded[K, T]) Size() int {
	return e.data.Size()
}

func (e *safeMapEmbedded[K, T]) Each(fn SafeMapLambda[T]) {
	e.data.Each(fn)
}

func NewSafeMapEmbedded[K comparable, T any](gen SafeMapGenerator[T]) SafeMapEmbedded[K, T] {
	return &safeMapEmbedded[K, T]{
		data: NewSafeMap[K, T](),
		gen:  gen,
	}
}
