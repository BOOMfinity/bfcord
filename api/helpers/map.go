package helpers

func Map[A, B any](arr []A, fn func(obj A) B) []B {
	result := make([]B, len(arr))
	for i, v := range arr {
		result[i] = fn(v)
	}
	return result
}
