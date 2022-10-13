package slices

import "golang.org/x/exp/constraints"

func Count[V any](slice []V, filter func(item V) bool) (x int) {
	for i := range slice {
		if filter(slice[i]) {
			x++
		}
	}
	return
}

func Sum[V constraints.Integer](slice []V) (x int) {
	for i := range slice {
		x += int(slice[i])
	}
	return
}

func SumCustom[V any](slice []V, fn func(item V) int) (x int) {
	for i := range slice {
		x += fn(slice[i])
	}
	return
}
