package discord

import "github.com/andersfylling/snowflake/v5"

// Emoji
//
// Reference: https://discord.com/developers/docs/resources/emoji#emoji-object
type Emoji struct {
	Name          string         `json:"name"`
	Roles         []snowflake.ID `json:"roles,omitempty"`
	User          User           `json:"user,omitempty"`
	ID            snowflake.ID   `json:"id"`
	GuildID       snowflake.ID   `json:"-"`
	RequireColons bool           `json:"require_colons,omitempty"`
	Managed       bool           `json:"managed,omitempty"`
	Animated      bool           `json:"animated,omitempty"`
	Available     bool           `json:"available,omitempty"`
}

func (v Emoji) Guild(client ClientQuery) GuildQuery {
	return client.Guild(v.GuildID)
}

func (v Emoji) ToString() string {
	if v.ID.Valid() {
		return v.Name + ":" + v.ID.String()
	}
	return v.Name
}
