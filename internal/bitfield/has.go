package bitfield

import "golang.org/x/exp/constraints"

func Has[V constraints.Integer](base V, check V) bool {
	return base&check == check
}
