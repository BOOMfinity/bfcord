package slices

func FindIndex[V any](slice []V, fn func(item V) bool) int {
	for i := range slice {
		if fn(slice[i]) {
			return i
		}
	}
	return -1
}

func FindCopy[V any](slice []V, fn func(item V) bool) (item V, found bool) {
	for i := range slice {
		if fn(slice[i]) {
			return slice[i], true
		}
	}
	return
}

func Find[V any](slice []V, fn func(item V) bool) *V {
	for i := range slice {
		item := slice[i]
		if fn(item) {
			return &item
		}
	}
	return nil
}
