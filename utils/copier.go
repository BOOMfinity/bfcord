package utils

type DeepCopier[T any] interface {
	Copy() *T
}

func CopyPointer[T any](orig *T, usePointerCopier ...bool) *T {
	if orig == nil {
		return nil
	}

	useBuiltInCopier := false

	if len(usePointerCopier) > 0 {
		useBuiltInCopier = usePointerCopier[0]
	}

	var copy T

	if copier, ok := any(orig).(DeepCopier[T]); ok && useBuiltInCopier {
		copy = *copier.Copy()
	} else {
		copy = *orig
	}

	return &copy
}

func CopySlice[T any](orig []T, usePointerCopier ...bool) []T {
	if orig == nil {
		return nil
	}

	useBuiltInCopier := false

	if len(usePointerCopier) > 0 {
		useBuiltInCopier = usePointerCopier[0]
	}

	copy := make([]T, len(orig))

	for i := range orig {
		if copier, ok := any(orig[i]).(DeepCopier[T]); ok && useBuiltInCopier {
			copy[i] = *copier.Copy()
		} else {
			copy[i] = orig[i]
		}
	}

	return copy
}

func CopyMap[K comparable, T any](orig map[K]T, usePointerCopier ...bool) map[K]T {
	if orig == nil {
		return nil
	}

	useBuiltInCopier := false

	if len(usePointerCopier) > 0 {
		useBuiltInCopier = usePointerCopier[0]
	}

	copy := make(map[K]T, len(orig))

	for key, value := range orig {
		if copier, ok := any(value).(DeepCopier[T]); ok && useBuiltInCopier {
			copy[key] = *copier.Copy()
		} else {
			copy[key] = value
		}
	}

	return copy
}
