package utils

import "sync"

type SafeMapLambda[T any] func(obj T) bool
type SafeMapGenerator[T any] func() T

type SafeMap[K comparable, T any] interface {
	Get(key K) (T, bool)
	Delete(key K) bool
	Has(key K) bool
	Set(key K, value T)
	Size() int

	Find(fn SafeMapLambda[T]) (T, bool)
	Search(fn SafeMapLambda[T]) (values []T)
	Each(fn SafeMapLambda[T])
}

type SafeMapEmbedded[K comparable, T any] interface {
	Get(key K) T
	Delete(key K) (ok bool)
	Each(fn SafeMapLambda[T])
	Size() int
}

type LimitedSafeMapOrderBy[T any] func(a, b T) bool

type LimitedSafeMap[K comparable, T any] interface {
	SafeMap[K, T]

	Sorted() []T
}

type RWLocker interface {
	sync.Locker
	RLock()
	RUnlock()
}
