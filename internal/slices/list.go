package slices

type List[V any] []V

func (x List[V]) Find(fn func(item V) bool) (data V, ok bool) {
	for i := range x {
		if fn(x[i]) {
			return x[i], true
		}
	}
	return
}
