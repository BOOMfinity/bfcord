package slash

import (
	"github.com/BOOMfinity/bfcord/discord/permissions"
	"github.com/andersfylling/snowflake/v5"
)

// Ref: https://discord.com/developers/docs/interactions/application-commands#application-command-object
type Command struct {
	NameLocalizations        map[string]string      `json:"name_localizations,omitempty"`
	DescriptionLocalizations map[string]string      `json:"description_localizations,omitempty"`
	Description              string                 `json:"description,omitempty"`
	Name                     string                 `json:"name,omitempty"`
	Options                  []Option               `json:"options,omitempty"`
	GuildID                  snowflake.ID           `json:"guild_id,omitempty"`
	ID                       snowflake.ID           `json:"id,omitempty"`
	ApplicationID            snowflake.ID           `json:"application_id,omitempty"`
	DefaultMemberPermissions permissions.Permission `json:"default_member_permissions,omitempty"`
	Version                  snowflake.ID           `json:"version,omitempty"`
	Type                     CommandType            `json:"type,omitempty"`
	DM                       bool                   `json:"dm_permission,omitempty"`
	DefaultPermission        bool                   `json:"default_permission,omitempty"`
}

type CommandType uint8

const (
	CommandTypeChatInput CommandType = iota + 1
	CommandTypeUser
	CommandTypeMessage
)
