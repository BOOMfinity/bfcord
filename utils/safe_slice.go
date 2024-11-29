package utils

import (
	"sync"

	"github.com/segmentio/encoding/json"
)

type SafeSliceUnlock func()
type SafeSliceIndex[T any] func(v T) bool

type SafeSlice[T any] interface {
	json.Marshaler
	json.Unmarshaler

	Read() []T
	Remove(i int)
	Size() int
	Get(i int) T
	Set(i int, v T)
	Add(T)
	Index(fn SafeSliceIndex[T]) int
	Lock() (unlock SafeSliceUnlock)
	Replace([]T)
}

type safeSliceImpl[T any] struct {
	data []T
	mut  sync.RWMutex
}

func (s *safeSliceImpl[T]) UnmarshalJSON(data []byte) error {
	s.mut.Lock()
	defer s.mut.Unlock()
	return json.Unmarshal(data, &s.data)
}

func (s *safeSliceImpl[T]) MarshalJSON() ([]byte, error) {
	s.mut.Lock()
	defer s.mut.Unlock()
	return json.Marshal(s.data)
}

func (s *safeSliceImpl[T]) Read() []T {
	s.mut.RLock()
	data := s.data[:]
	s.mut.RUnlock()
	return data
}

func (s *safeSliceImpl[T]) Remove(i int) {
	s.mut.Lock()
	s.data = append(s.data[:i], s.data[i+1:]...)
	s.mut.Unlock()
}

func (s *safeSliceImpl[T]) Index(fn SafeSliceIndex[T]) (index int) {
	index = -1
	s.mut.RLock()
	for i, v := range s.data {
		if fn(v) {
			index = i
			break
		}
	}
	s.mut.RUnlock()
	return
}

func (s *safeSliceImpl[T]) Size() int {
	s.mut.RLock()
	l := len(s.data)
	s.mut.RUnlock()
	return l
}

func (s *safeSliceImpl[T]) Get(i int) T {
	s.mut.RLock()
	v := s.data[i]
	s.mut.RUnlock()
	return v
}

func (s *safeSliceImpl[T]) Set(i int, v T) {
	s.mut.Lock()
	s.data[i] = v
	s.mut.Unlock()
}

func (s *safeSliceImpl[T]) Lock() SafeSliceUnlock {
	s.mut.Lock()
	return func() {
		s.mut.Unlock()
	}
}

func (s *safeSliceImpl[T]) Add(v T) {
	s.mut.Lock()
	s.data = append(s.data, v)
	s.mut.Unlock()
}

func (s *safeSliceImpl[T]) Replace(v []T) {
	s.mut.Lock()
	s.data = v
	s.mut.Unlock()
}

func NewSafeSlice[T any]() SafeSlice[T] {
	return &safeSliceImpl[T]{}
}
