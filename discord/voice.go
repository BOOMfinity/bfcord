package discord

import "github.com/andersfylling/snowflake/v5"

type VoiceState struct {
	GuildID                 snowflake.ID `json:"guild_id,omitempty"`
	ChannelID               snowflake.ID `json:"channel_id,omitempty"`
	UserID                  snowflake.ID `json:"user_id,omitempty"`
	Member                  Member       `json:"member,omitempty"`
	SessionID               string       `json:"session_id,omitempty"`
	Deaf                    bool         `json:"deaf,omitempty"`
	Mute                    bool         `json:"mute,omitempty"`
	SelfDeaf                bool         `json:"self_deaf,omitempty"`
	SelfMute                bool         `json:"self_mute,omitempty"`
	SelfStream              bool         `json:"self_stream,omitempty"`
	SelfVideo               bool         `json:"self_video,omitempty"`
	Suppress                bool         `json:"suppress,omitempty"`
	RequestToSpeakTimestamp Timestamp    `json:"request_to_speak_timestamp,omitempty"`
}
