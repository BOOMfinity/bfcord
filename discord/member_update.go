package discord

import (
	"github.com/BOOMfinity/bfcord/internal/timeconv"
	"github.com/BOOMfinity/go-utils/nullable"
	"github.com/andersfylling/snowflake/v5"
)

type MemberUpdate struct {
	Nick                       *string                          `json:"nick,omitempty"`
	Roles                      *[]snowflake.ID                  `json:"roles,omitempty"`
	Mute                       *bool                            `json:"mute,omitempty"`
	Deaf                       *bool                            `json:"deaf,omitempty"`
	ChannelID                  *nullable.Nullable[snowflake.ID] `json:"channel_id,omitempty"`
	CommunicationDisabledUntil *timeconv.Timestamp              `json:"communication_disabled_until,omitempty"`
}
