package voice

import (
	"github.com/andersfylling/snowflake/v5"

	"github.com/BOOMfinity/bfcord/discord"
)

type StateUpdateEvent = discord.VoiceState

type ServerUpdateEvent struct {
	Token    string       `json:"token,omitempty"`
	GuildID  snowflake.ID `json:"guild_id,omitempty"`
	Endpoint string       `json:"endpoint,omitempty"`
}
