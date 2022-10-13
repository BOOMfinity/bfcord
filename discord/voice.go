package discord

import (
	"github.com/BOOMfinity/bfcord/internal/timeconv"
	"github.com/andersfylling/snowflake/v5"
)

// VoiceState
//
// Reference: https://discord.com/developers/docs/resources/voice#voice-state-object
type VoiceState struct {
	RequestToSpeakTimestamp timeconv.Timestamp `json:"request_to_speak_timestamp"`
	SessionID               string             `json:"session_id"`
	Member                  Member             `json:"member"`
	ChannelID               snowflake.ID       `json:"channel_id"`
	UserID                  snowflake.ID       `json:"user_id"`
	GuildID                 snowflake.ID       `json:"guild_id"`
	Deaf                    bool               `json:"deaf"`
	SelfDeaf                bool               `json:"self_deaf"`
	SelfMute                bool               `json:"self_mute"`
	SelfStream              bool               `json:"self_stream"`
	SelfVideo               bool               `json:"self_video"`
	Suppress                bool               `json:"suppress"`
	Mute                    bool               `json:"mute"`
}
