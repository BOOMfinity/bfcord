package cache

import (
	"slices"
	"sync"
)

type Map[K comparable, V any] interface {
	Get(key K) (V, error)
	Set(key K, obj V) error
	Delete(key K) error
	// Has checks if specific key exists. Default implementation will return ErrNotFound error if key does not exist.
	Has(key K) error
	// Size
	//
	// It will not return an error in Default implementation.
	Size() (int, error)
	Each(fn MapLambda[V]) error
	Search(fn MapLambda[V]) ([]V, error)
	Clear() error
}

type SubMap[K comparable, V any] interface {
	Get(key K) V
	Delete(key K) error
	Size() (int, error)
	Each(fn MapLambda[V]) error
	Search(fn MapLambda[V]) ([]V, error)
	Clear() error
}

type MapLambda[T any] func(T) bool
type SubMapGen[T any] func() T

type subMapImpl[K comparable, V any] struct {
	Map[K, V]
	gen SubMapGen[V]
}

func (s *subMapImpl[K, V]) Get(key K) V {
	obj, err := s.Map.Get(key)
	if err == nil {
		return obj
	}
	obj = s.gen()
	if err = s.Set(key, obj); err != nil {
		panic("error when saving to the cache.Map")
	}
	return obj
}

type mapImpl[K comparable, V any] struct {
	data map[K]V
	sync.RWMutex
}

func (m *mapImpl[K, V]) Get(key K) (obj V, err error) {
	m.RLock()
	obj, ok := m.data[key]
	m.RUnlock()
	if !ok {
		err = ErrNotFound
	}
	return
}

func (m *mapImpl[K, V]) Set(key K, obj V) error {
	m.Lock()
	m.data[key] = obj
	m.Unlock()
	return nil
}

func (m *mapImpl[K, V]) Delete(key K) error {
	m.Lock()
	_, ok := m.data[key]
	delete(m.data, key)
	m.Unlock()
	if !ok {
		return ErrNotFound
	}
	return nil
}

func (m *mapImpl[K, V]) Has(key K) error {
	m.RLock()
	_, ok := m.data[key]
	m.RUnlock()
	if !ok {
		return ErrNotFound
	}
	return nil
}

func (m *mapImpl[K, V]) Size() (int, error) {
	m.RLock()
	size := len(m.data)
	m.RUnlock()
	return size, nil
}

func (m *mapImpl[K, V]) Each(fn MapLambda[V]) error {
	m.RLock()
	objects := m.data
	m.RUnlock()
	for _, obj := range objects {
		if !fn(obj) {
			break
		}
	}
	return nil
}

func (m *mapImpl[K, V]) Search(fn MapLambda[V]) ([]V, error) {
	size, _ := m.Size()
	data := make([]V, 0, size)
	return data, m.Each(func(obj V) bool {
		if fn(obj) {
			data = append(data, obj)
		}
		return true
	})
}

func (m *mapImpl[K, V]) Clear() error {
	m.Lock()
	clear(m.data)
	m.Unlock()
	return nil
}

type limitedMapImpl[K comparable, V any] struct {
	Map[K, V]
	sync.RWMutex
	keys  []K
	limit int
}

func (m *limitedMapImpl[K, V]) Set(key K, obj V) error {
	created := m.Has(key) == nil
	_ = m.Map.Set(key, obj)
	if created {
		m.Lock()
		m.keys = append(m.keys, key)
		if len(m.keys) >= m.limit {
			key, m.keys = m.keys[0], m.keys[1:]
			_ = m.Map.Delete(key)
		}
		m.Unlock()
	}
	return nil
}

func (m *limitedMapImpl[K, V]) Delete(key K) error {
	if err := m.Map.Delete(key); err == nil {
		m.Lock()
		m.keys = slices.DeleteFunc(m.keys, func(k K) bool {
			return k == key
		})
		m.Unlock()
	}
	return nil
}

func NewSubMap[K comparable, V any](fn SubMapGen[V]) SubMap[K, V] {
	return &subMapImpl[K, V]{
		Map: NewMap[K, V](0),
		gen: fn,
	}
}

func NewMap[K comparable, V any](preAllocation uint) Map[K, V] {
	return &mapImpl[K, V]{
		data: make(map[K]V, preAllocation),
	}
}

func NewLimitedMap[K comparable, V any](limit int, preAllocated uint) Map[K, V] {
	return &limitedMapImpl[K, V]{
		Map:   NewMap[K, V](preAllocated),
		limit: limit,
	}
}
