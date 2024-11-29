package ws

import (
	"errors"
	"fmt"

	"github.com/andersfylling/snowflake/v5"
)

var (
	ErrGatewayNotConnected     = errors.New("gateway is disconnected from Discord")
	ErrFetchingMembersTimedOut = errors.New("could not wait longer for guild members chunk")
)

type ErrNotFound []snowflake.ID

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%d users were not found", len(e))
}
