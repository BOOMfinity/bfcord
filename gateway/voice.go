package gateway

import (
	"context"
	"github.com/andersfylling/snowflake/v5"
)

// ChangeVoiceState method will send OP4 Voice State Update to the gateway.
//
// You have to listen for VoiceStateUpdate and VoiceServerUpdate events by yourself or use a wrapper built-in to client package (client.ChangeVoiceState)
func (g *Gateway) ChangeVoiceState(opts ChangeVoiceStateOptions) error {
	return g.ws.WriteJSON(context.Background(), map[string]any{
		"op": VoiceStateUpdateOp,
		"d":  opts,
	})
}

type ChangeVoiceStateOptions struct {
	GuildID   snowflake.ID `json:"guild_id"`
	ChannelID snowflake.ID `json:"channel_id"`
	SelfMute  bool         `json:"self_mute"`
	SelfDeaf  bool         `json:"self_deaf"`
}
