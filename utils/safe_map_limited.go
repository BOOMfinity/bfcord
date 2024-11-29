package utils

import (
	"sort"
	"sync"
)

type limitedSafeMapImpl[K comparable, T any] struct {
	SafeMap[K, T]

	limit  int
	sorted []K
	sortBy LimitedSafeMapOrderBy[T]
	mut    sync.RWMutex
}

func (l *limitedSafeMapImpl[K, T]) sort() {
	sort.SliceStable(l.sorted, func(i, j int) bool {
		v1, _ := l.Get(l.sorted[i])
		v2, _ := l.Get(l.sorted[j])

		return l.sortBy(v1, v2)
	})
}

func (l *limitedSafeMapImpl[K, T]) Set(key K, value T) {
	l.mut.Lock()
	l.SafeMap.Set(key, value)
	if len(l.sorted) > l.limit {
		key, l.sorted = l.sorted[l.limit-1], l.sorted[:l.limit]
		l.SafeMap.Delete(key)
		l.sort()
	}
	l.mut.Unlock()
}

func (l *limitedSafeMapImpl[K, T]) Delete(key K) (ok bool) {
	l.mut.Lock()
	ok = l.SafeMap.Delete(key)
	if ok {
		index := -1
		for i, k := range l.sorted {
			if k == key {
				index = i
				break
			}
		}
		if index != -1 {
			l.sorted = append(l.sorted[:index], l.sorted[(index+1):]...)
			l.sort()
		}
	}
	l.mut.Unlock()
	return
}

func (l *limitedSafeMapImpl[K, T]) Sorted() (values []T) {
	l.mut.RLock()
	values = make([]T, 0, l.Size())
	for _, key := range l.sorted {
		value, ok := l.Get(key)
		if ok {
			values = append(values, value)
		}
	}
	l.mut.RUnlock()
	return
}

func NewLimitedSafeMap[K comparable, T any](limit int, sortBy LimitedSafeMapOrderBy[T]) LimitedSafeMap[K, T] {
	return &limitedSafeMapImpl[K, T]{
		SafeMap: NewSafeMap[K, T](),
		limit:   limit,
		sortBy:  sortBy,
	}
}
