package timeconv

import (
	"strconv"
	"time"
)

// Seconds type is wrapper for time.Duration sent in seconds with fast (un)marshal functions
type Seconds time.Duration

func (s *Seconds) UnmarshalJSON(b []byte) error {
	var x int
	for _, c := range b {
		x = x*10 + int(c-'0')
	}
	*s = Seconds(time.Duration(x) * time.Second)
	return nil
}

func (s Seconds) MarshalJSON() ([]byte, error) {
	if s < 1 {
		return null, nil
	}
	return []byte(strconv.Itoa(int(s))), nil
}

func (s *Seconds) Duration() time.Duration {
	return time.Duration(*s)
}
