package discord

import (
	"github.com/andersfylling/snowflake/v5"
)

type Webhook struct {
	ID            snowflake.ID `json:"id,omitempty"`
	Token         string       `json:"token,omitempty"`
	Type          WebhookType  `json:"type,omitempty"`
	GuildID       snowflake.ID `json:"guild_id,omitempty"`
	ChannelID     snowflake.ID `json:"channel_id,omitempty"`
	User          User         `json:"user,omitempty"`
	Name          string       `json:"name,omitempty"`
	Avatar        string       `json:"avatar,omitempty"`
	ApplicationID snowflake.ID `json:"application_id,omitempty"`
	SourceGuild   Guild        `json:"source_guild,omitempty"`
	SourceChannel Channel      `json:"source_channel,omitempty"`
	Url           string       `json:"url,omitempty"`
}

type WebhookType uint

const (
	WebhookTypeIncoming WebhookType = iota + 1
	WebhookTypeChannelFollower
	WebhookTypeApplication
)
