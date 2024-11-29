package httpc

import (
	"errors"
	"fmt"
	"github.com/segmentio/encoding/json"
)

var (
	ErrMaxRetriesReached            = errors.New("reached the request retries limit")
	ErrFailedToParseRequestOptions  = errors.New("there was an error while executing request options")
	ErrExecutingRequest             = errors.New("something went wrong while executing this request")
	ErrFailedToParseResponseOptions = errors.New("there was an error while executing response options")
	ErrTooManyRequests              = errors.New("your have sent too many request to this bucket, wait and try again")
)

type DiscordError struct {
	Message string          `json:"message"`
	Errors  json.RawMessage `json:"errors"`
	Code    int             `json:"code"`
}

func (e DiscordError) Error() string {
	if len(e.Errors) > 0 {
		d, _ := json.MarshalIndent(e.Errors, "", "\t")
		return fmt.Sprintf("%s - %d\n\n%s\n", e.Message, e.Code, string(d))
	}
	if e.Code != 0 {
		return fmt.Sprintf("%s - %d", e.Message, e.Code)
	}
	return e.Message
}
