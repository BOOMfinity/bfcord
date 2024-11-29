package discord

import "github.com/andersfylling/snowflake/v5"

type CommandChoice struct {
	Name              string            `json:"name,omitempty"`
	NameLocalizations map[string]string `json:"name_localizations,omitempty"`
	Value             any               `json:"value,omitempty"`
}

type Command struct {
	ID                       snowflake.ID                 `json:"id,omitempty"`
	Type                     CommandType                  `json:"type,omitempty"`
	ApplicationID            snowflake.ID                 `json:"application_id,omitempty"`
	GuildID                  snowflake.ID                 `json:"guild_id,omitempty"`
	Name                     string                       `json:"name,omitempty"`
	NameLocalized            string                       `json:"name_localized,omitempty"`
	NameLocalizations        map[string]string            `json:"name_localizations,omitempty"`
	Description              string                       `json:"description,omitempty"`
	DescriptionLocalized     string                       `json:"description_localized,omitempty"`
	DescriptionLocalizations map[string]string            `json:"description_localizations,omitempty"`
	Options                  []any                        `json:"options,omitempty"`
	DefaultMemberPermissions Permission                   `json:"default_member_permissions,omitempty"`
	NSFW                     bool                         `json:"nsfw,omitempty"`
	IntegrationTypes         []ApplicationIntegrationType `json:"integration_types,omitempty"`
	Version                  snowflake.ID                 `json:"version,omitempty"`
}

type CommandOption struct {
	Type                     CommandOptionType `json:"type,omitempty"`
	Name                     string            `json:"name,omitempty"`
	NameLocalized            string            `json:"name_localized,omitempty"`
	NameLocalizations        map[string]string `json:"name_localizations,omitempty"`
	Description              string            `json:"description,omitempty"`
	DescriptionLocalized     string            `json:"description_localized,omitempty"`
	DescriptionLocalizations map[string]string `json:"description_localizations,omitempty"`
	Required                 bool              `json:"required,omitempty"`
	Choices                  []CommandChoice   `json:"choices,omitempty"`
	Options                  []CommandOption   `json:"options,omitempty"`
	ChannelTypes             []ChannelType     `json:"channel_types,omitempty"`
	MinValue                 float64           `json:"min_value,omitempty"`
	MaxValue                 float64           `json:"max_value,omitempty"`
	MinLength                uint              `json:"min_length,omitempty"`
	MaxLength                uint              `json:"max_length,omitempty"`
	AutoComplete             bool              `json:"autocomplete,omitempty"`
}

type CreateCommand struct {
	Name                     string                       `json:"name,omitempty"`
	NameLocalizations        map[string]string            `json:"name_localizations,omitempty"`
	Description              string                       `json:"description,omitempty"`
	DescriptionLocalizations map[string]string            `json:"description_localizations,omitempty"`
	Options                  []CommandOption              `json:"options,omitempty"`
	DefaultMemberPermissions Permission                   `json:"default_member_permissions,omitempty"`
	IntegrationTypes         []ApplicationIntegrationType `json:"integration_types,omitempty"`
	Contexts                 []InteractionContextType     `json:"contexts,omitempty"`
	Type                     CommandType                  `json:"type,omitempty"`
	NSFW                     bool                         `json:"nsfw,omitempty"`
}

type CommandPermissions struct {
	ID            snowflake.ID        `json:"id,omitempty"`
	ApplicationID snowflake.ID        `json:"application_id,omitempty"`
	GuildID       snowflake.ID        `json:"guild_id,omitempty"`
	Permissions   []CommandPermission `json:"permissions,omitempty"`
}

type CommandPermission struct {
	ID         snowflake.ID          `json:"id,omitempty"`
	Type       CommandPermissionType `json:"type,omitempty"`
	Permission *bool                 `json:"permission,omitempty"`
}

type CommandPermissionType uint

const (
	CommandPermissionRole CommandPermissionType = iota + 1
	CommandPermissionUser
	CommandPermissionChannel
)
