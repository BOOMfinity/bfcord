package discord

import (
	"github.com/BOOMfinity/bfcord/api/images"
	"github.com/andersfylling/snowflake/v5"
)

type Webhook struct {
	Avatar        string       `json:"avatar"`
	Token         string       `json:"token"`
	Name          string       `json:"name"`
	Url           string       `json:"url"`
	User          User         `json:"user"`
	SourceChannel Channel      `json:"source_channel"`
	SourceGuild   Guild        `json:"source_guild"`
	GuildID       snowflake.ID `json:"guild_id"`
	ID            snowflake.ID `json:"id"`
	Type          int          `json:"type"`
	ApplicationID snowflake.ID `json:"application_id"`
	ChannelID     snowflake.ID `json:"channel_id"`
}

type WebhookType uint8

const (
	WebhookIncoming WebhookType = iota + 1
	WebhookChannelFollower
	WebhookApplication
)

type WebhookExecute struct {
	MessageCreate
	Username  *string `json:"username,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type WebhookCreate struct {
	Name   string        `json:"name"`
	Avatar *images.Image `json:"avatar,omitempty"`
	Reason string        `json:"-"`
}
