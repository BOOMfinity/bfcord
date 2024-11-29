package utils

import "sync"

type simpleMap[K comparable, T any] struct {
	data map[K]T
	sync.RWMutex
}

func (self *simpleMap[K, T]) Set(k K, v T) {
	self.Lock()
	self.data[k] = v
	self.Unlock()
}

func (self *simpleMap[K, T]) Delete(k K) bool {
	self.Lock()
	_, ok := self.data[k]
	delete(self.data, k)
	self.Unlock()
	return ok
}

func (self *simpleMap[K, T]) Get(k K) (T, bool) {
	self.RLock()
	data, ok := self.data[k]
	self.RUnlock()
	return data, ok
}

func (self *simpleMap[K, T]) Size() int {
	self.RLock()
	size := len(self.data)
	self.RUnlock()
	return size
}

func (self *simpleMap[K, T]) Each(fn func(key K, value T)) {
	self.RLock()
	for key, value := range self.data {
		fn(key, value)
	}
	self.RUnlock()
}

type SimpleMap[K comparable, T any] interface {
	Set(key K, value T)
	Delete(key K) (ok bool)
	Get(key K) (value T, ok bool)
	Size() int
	Each(fn func(key K, value T))
}

func NewSimpleMap[K comparable, T any]() SimpleMap[K, T] {
	return &simpleMap[K, T]{
		data: make(map[K]T),
	}
}
