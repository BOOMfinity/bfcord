package discord

import (
	"github.com/BOOMfinity/bfcord/discord/permissions"
	"github.com/andersfylling/snowflake/v5"
)

// Role
//
// Reference: https://discord.com/developers/docs/topics/permissions#role-object
type Role struct {
	Name         string                 `json:"name"`
	Icon         string                 `json:"icon,omitempty"`
	UnicodeEmoji string                 `json:"unicode_emoji,omitempty"`
	Permissions  permissions.Permission `json:"permissions"`
	Tags         RoleTags               `json:"tags,omitempty"`
	ID           snowflake.ID           `json:"id"`
	Color        int                    `json:"color"`
	Position     int                    `json:"position"`
	GuildID      snowflake.ID           `json:"guild_id"`
	Hoist        bool                   `json:"hoist"`
	Managed      bool                   `json:"managed"`
	Mentionable  bool                   `json:"mentionable"`
}

func (r Role) Guild(api ClientQuery) GuildQuery {
	return api.Guild(r.GuildID)
}

// RoleTags
//
// Reference: https://discord.com/developers/docs/topics/permissions#role-object-role-tags-structure
type RoleTags struct {
	BotID             snowflake.ID `json:"bot_id,omitempty"`
	IntegrationID     snowflake.ID `json:"integration_id,omitempty"`
	PremiumSubscriber snowflake.ID `json:"premium_subscriber,omitempty"`
}

type RoleCreate struct {
	Name         *string                 `json:"name,omitempty"`
	Permissions  *permissions.Permission `json:"permissions,omitempty"`
	Color        *int64                  `json:"color,omitempty"`
	Hoist        *bool                   `json:"hoist,omitempty"`
	Icon         *string                 `json:"icon,omitempty"`
	UnicodeEmoji *string                 `json:"unicode_emoji,omitempty"`
	Mentionable  *bool                   `json:"mentionable,omitempty"`
}
