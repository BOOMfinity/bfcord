package errs

import (
	"fmt"
	"github.com/segmentio/encoding/json"
)

type DiscordError struct {
	Message string          `json:"message"`
	Details json.RawMessage `json:"errors"`
	Code    int             `json:"code"`
}

func (v DiscordError) Error() string {
	if len(v.Details) > 0 {
		return fmt.Sprintf("discord error(%v): %v\n%v", v.Code, v.Message, string(v.Details))
	}
	return fmt.Sprintf("discord error(%v): %v", v.Code, v.Message)
}

/*type DiscordError struct {
	Message string          `json:"message"`
	Code    int             `json:"code"`
	Errors  json.RawMessage `json:"errors"`
}

func (v DiscordError) Error() string {
	data := fmt.Sprintf("%v (%v)", v.Message, v.Code)
	if len(v.Errors) > 0 {
		data += fmt.Sprintf("\n%v", string(v.Errors))
	}
	return data
}

var (
	DiscordUnauthorized = &DiscordError{Message: "Unauthorized", Code: 401}
	DiscordNotFound     = &DiscordError{Message: "NotFound", Code: 404}
)

func Is(err error, code int) bool {
	if dcErr, ok := err.(*DiscordError); ok {
		if dcErr.Code == code {
			return true
		}
	}
	return false
}
*/

func Is(err error, code int) bool {
	if dcErr, ok := err.(*DiscordError); ok {
		if dcErr.Code == code {
			return true
		}
	}
	return false
}
