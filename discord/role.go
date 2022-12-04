package discord

import (
	"github.com/BOOMfinity/bfcord/discord/permissions"
	"github.com/BOOMfinity/go-utils/inlineif"
	"github.com/andersfylling/snowflake/v5"
	"golang.org/x/exp/slices"
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

// ComparePosition compares this role's position to other one.
func (r Role) ComparePosition(other *Role) int {
	if r.Position == other.Position {
		return int(other.ID - r.ID)
	}

	return r.Position - other.Position
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

type RoleSlice []Role

func (rs RoleSlice) Highest() (highest *Role) {
	if len(rs) == 0 {
		return nil
	}

	for i := range rs {
		if highest == nil {
			highest = &rs[0]
			continue
		}

		if rs[i].ComparePosition(highest) > 0 {
			highest = &rs[i]
		}
	}
	return
}

func (rs RoleSlice) HighestWithin(member *Member) (highest *Role) {
	if len(rs) == 0 {
		return nil
	}

	for i := range rs {
		if highest == nil {
			highest = inlineif.IfElse(slices.Contains(member.Roles, rs[i].ID), &rs[i], nil)
			continue
		}

		if slices.Contains(member.Roles, rs[i].ID) && rs[i].ComparePosition(highest) > 0 {
			highest = &rs[i]
		}
	}
	return
}

func (rs RoleSlice) Find(id snowflake.ID) *Role {
	index := slices.IndexFunc(rs, func(r Role) bool {
		return r.ID == id
	})
	if index == -1 {
		return nil
	}
	return &rs[index]
}
