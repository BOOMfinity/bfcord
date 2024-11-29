package discord

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/BOOMfinity/go-utils/ubytes"
)

const ISO8601Format = "2006-01-02T15:04:05Z07:00"

//easyjson:skip
type Timestamp struct {
	time.Time
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return ubytes.ToBytes(`"` + t.Format(ISO8601Format) + `"`), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		return nil
	}
	tm, err := time.Parse(ISO8601Format, ubytes.ToString(b[1:len(b)-1]))
	if err != nil {
		return fmt.Errorf("failed to parse iso8601: %w", err)
	}
	*t = Timestamp{tm}
	return nil
}

//easyjson:skip
type UnixTimestamp struct {
	time.Time
}

func (t UnixTimestamp) MarshalJSON() ([]byte, error) {
	return ubytes.ToBytes(strconv.FormatInt(t.UnixMilli(), 10)), nil
}

func (t *UnixTimestamp) UnmarshalJSON(bytes []byte) error {
	tm, err := strconv.ParseInt(ubytes.ToString(bytes), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse unix timestamp: %w", err)
	}
	t.Time = time.UnixMilli(tm)
	return nil
}

type BitField uint

func (b *BitField) Add(bits ...BitField) {
	for _, bit := range bits {
		*b |= BitField(1 << bit)
	}
}

func (b *BitField) Remove(bits ...BitField) {
	for _, bit := range bits {
		*b &= ^BitField(1 << bit)
	}
}

func (b *BitField) Toggle(bits ...BitField) {
	for _, bit := range bits {
		*b ^= BitField(1 << bit)
	}
}

func (b BitField) Has(bits ...BitField) bool {
	for _, bit := range bits {
		if (b & BitField(1<<bit)) == 0 {
			return false
		}
	}
	return true
}
