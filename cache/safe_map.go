package cache

import (
	"sort"
	"sync"
)

type SafeMap[K comparable, V any] struct {
	indexes    map[K]int
	defaultNew func() V
	data       []V
	nils       []int
	m          sync.RWMutex
}

func (v *SafeMap[K, V]) UnsafeGet(key K) V {
	return v.GetOrSet(key, v.defaultNew)
}

func (v *SafeMap[K, V]) Get(key K) (value V, ok bool) {
	v.m.RLock()
	defer v.m.RUnlock()
	if !v.Has(key) {
		return
	}
	index := v.indexes[key]
	value = v.data[index]
	ok = true
	return
}

type sortable[V any] struct {
	less func(a, b V) bool
	data []V
}

func (a sortable[V]) Len() int {
	return len(a.data)
}

func (a sortable[V]) Less(i, j int) bool {
	return a.less(a.data[i], a.data[j])
}

func (a sortable[V]) Swap(i, j int) {
	a.data[i], a.data[j] = a.data[j], a.data[i]
}

func (v *SafeMap[K, V]) Sort(fn func(a, b V) bool) (sorted []V) {
	sorted = make([]V, 0, v.Size())
	v.m.RLock()
	for i := range v.data {
		sorted = append(sorted, v.data[i])
	}
	v.m.RUnlock()
	sort.Sort(sortable[V]{data: sorted, less: fn})
	return
}

func (v *SafeMap[K, V]) Each(fn func(item V)) {
	v.m.RLock()
	defer v.m.RUnlock()
	for key := range v.indexes {
		index := v.indexes[key]
		val := v.data[index]
		fn(val)
	}
	return
}

func (v *SafeMap[K, V]) Find(fn func(item V) bool) (_ V, found bool) {
	v.m.RLock()
	defer v.m.RUnlock()
	for i := range v.data {
		val := v.data[i]
		if fn(val) {
			return val, true
		}
	}
	return
}

func (v *SafeMap[K, V]) ToSlice() (data []V) {
	data = make([]V, 0, v.Size())
	v.m.RLock()
	defer v.m.RUnlock()
	for i := range v.data {
		data = append(data, v.data[i])
	}
	return
}

func (v *SafeMap[K, V]) Filter(fn func(item V) bool) (data []V) {
	v.m.RLock()
	defer v.m.RUnlock()
	for i := range v.data {
		val := v.data[i]
		if fn(val) {
			data = append(data, val)
		}
	}
	return
}

func (v *SafeMap[K, V]) GetOrSet(key K, set func() V) V {
	val, ok := v.Get(key)
	if ok {
		return val
	}
	x := set()
	v.Set(key, x)
	return x
}

func (v *SafeMap[K, V]) Set(key K, value V) {
	v.m.Lock()
	defer v.m.Unlock()
	if index, ok := v.indexes[key]; ok {
		v.data[index] = value
		return
	}
	var index int
	if len(v.nils) > 0 {
		index, v.nils = v.nils[0], v.nils[1:]
	} else {
		index = len(v.data)
	}
	v.data = append(v.data, value)
	v.indexes[key] = index
}

func (v *SafeMap[K, V]) Has(key K) bool {
	v.m.RLock()
	defer v.m.RUnlock()
	index, ok := v.indexes[key]
	return ok && index != -1
}

func (v *SafeMap[K, V]) Delete(key K) bool {
	if !v.Has(key) {
		return false
	}
	v.m.Lock()
	defer v.m.Unlock()
	if ckey, ok := v.indexes[key]; ok && ckey != -1 {
		var x V
		delete(v.indexes, key)
		v.nils = append(v.nils, ckey)
		v.data[ckey] = x
		return true
	}
	return false
}

func (v *SafeMap[K, V]) Size() int {
	v.m.RLock()
	defer v.m.RUnlock()
	return len(v.data)
}

func (v *SafeMap[K, V]) Update(key K, fn func(value V) V) (ok bool) {
	v.m.Lock()
	defer v.m.Unlock()
	index, ok := v.indexes[key]
	if !ok {
		return false
	}
	v.data[index] = fn(v.data[index])
	return true
}

func NewSafeMap[K comparable, V any](prealloc int) *SafeMap[K, V] {
	return &SafeMap[K, V]{
		data:    make([]V, 0, prealloc),
		indexes: make(map[K]int, prealloc),
	}
}

func NewSafeMapWithInitializer[K comparable, V any](prealloc int, init func() V) *SafeMap[K, V] {
	safe := NewSafeMap[K, V](prealloc)
	safe.defaultNew = init
	return safe
}
