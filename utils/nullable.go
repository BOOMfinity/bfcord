package utils

import (
	"bytes"

	"github.com/segmentio/encoding/json"
)

type Nullable[T any] struct {
	data *T
}

func (n *Nullable[T]) Set(v T) {
	n.data = &v
}

func (n *Nullable[T]) Clear() {
	n.data = nil
}

func (n Nullable[T]) Get() (v T) {
	if n.data != nil {
		return *n.data
	}
	return
}

func (n Nullable[T]) Valid() bool {
	return n.data != nil
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.data)
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		n.data = nil
		return nil
	}
	if bytes.Equal(data, []byte("null")) {
		n.data = nil
		return nil
	}
	n.data = new(T)
	return json.Unmarshal(data, n.data)
}
