package slash

import (
	"github.com/BOOMfinity/bfcord/discord"
)

type Option struct {
	MaxValue                 any                   `json:"max_value,omitempty"`
	MinValue                 any                   `json:"min_value,omitempty"`
	DescriptionLocalizations map[string]string     `json:"description_localizations,omitempty"`
	NameLocalizations        map[string]string     `json:"name_localizations,omitempty"`
	Description              string                `json:"description"`
	Name                     string                `json:"name"`
	Choices                  []Choice              `json:"choices,omitempty"`
	Options                  []Option              `json:"options,omitempty"`
	ChannelTypes             []discord.ChannelType `json:"channel_types,omitempty"`
	MinLength                uint16                `json:"min_length,omitempty"`
	MaxLength                uint16                `json:"max_length,omitempty"`
	Required                 bool                  `json:"required"`
	Type                     OptionType            `json:"type"`
	Autocomplete             bool                  `json:"autocomplete,omitempty"`
}

type OptionType uint8

const (
	OptionTypeSubCommand OptionType = iota + 1
	OptionTypeSubCommandGroup
	OptionTypeString
	OptionTypeInteger
	OptionTypeBoolean
	OptionTypeUser
	OptionTypeChannel
	OptionTypeRole
	OptionTypeMentionable
	OptionTypeDouble
	OptionTypeAttachment
)
