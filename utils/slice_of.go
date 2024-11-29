package utils

func SliceOf[A any, B any](orig []A) (dst []B) {
	dst = make([]B, len(orig))
	for i, obj := range orig {
		if v, ok := any(obj).(B); ok {
			dst[i] = v
		}
	}
	return
}
