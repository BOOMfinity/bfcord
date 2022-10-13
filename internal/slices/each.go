package slices

func Each[V any](slice []V, fn func(item V) bool) {
	for i := range slice {
		if !fn(slice[i]) {
			break
		}
	}
}
