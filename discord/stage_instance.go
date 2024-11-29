package discord

import "github.com/andersfylling/snowflake/v5"

type StageInstance struct {
	ID                    snowflake.ID      `json:"id,omitempty"`
	GuildID               snowflake.ID      `json:"guild_id,omitempty"`
	ChannelID             snowflake.ID      `json:"channel_id,omitempty"`
	Topic                 string            `json:"topic,omitempty"`
	PrivacyLevel          StagePrivacyLevel `json:"privacy_level,omitempty"`
	DiscoverableDisabled  bool              `json:"discoverable_disabled,omitempty"`
	GuildScheduledEventID snowflake.ID      `json:"guild_scheduled_event_id,omitempty"`
}

type StagePrivacyLevel uint

const (
	StagePrivacyPublic StagePrivacyLevel = iota + 1
	StagePrivacyGuild
)
