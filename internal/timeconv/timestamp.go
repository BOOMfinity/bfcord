package timeconv

import (
	"github.com/segmentio/encoding/json"
	"github.com/unxcepted/iso8601/v2"
	"time"
)

var null = []byte(`null`)

// Timestamp type is time.Time wrapper with null support and optimized (un)marshal functions
type Timestamp struct {
	time.Time
}

var _ json.Marshaler = (*Timestamp)(nil)
var _ json.Unmarshaler = (*Timestamp)(nil)

func (t Timestamp) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return null, nil
	}

	return []byte(`"` + t.Format(time.RFC3339) + `"`), nil
}

func (t *Timestamp) UnmarshalJSON(date []byte) (err error) {
	if date[0] == 'n' && date[3] == 'l' { // null
		return nil
	}
	date = date[1 : len(date)-1] //trim quotes

	t.Time, err = iso8601.Parse(date)
	return err
}
