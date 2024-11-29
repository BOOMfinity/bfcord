package bfcord

func PointerOf[T any](v T) *T {
	return &v
}

func ReadPointer[T any](v *T, returnWhenNil T) T {
	if v == nil {
		return returnWhenNil
	}
	return *v
}
