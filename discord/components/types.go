package components

import (
	"github.com/andersfylling/snowflake/v5"
)

type Component struct {
	Emoji       *PartialEmoji      `json:"emoji,omitempty"`
	Value       string             `json:"value,omitempty"`
	CustomID    string             `json:"custom_id,omitempty"`
	Placeholder string             `json:"placeholder,omitempty"`
	Label       string             `json:"label,omitempty"`
	Url         string             `json:"url,omitempty"`
	Options     []SelectMenuOption `json:"options,omitempty"`
	Components  []Component        `json:"components,omitempty"`
	Style       ButtonStyle        `json:"style,omitempty"`
	Disabled    bool               `json:"disabled,omitempty"`
	MinValues   uint8              `json:"min_values,omitempty"`
	MaxValues   uint8              `json:"max_values,omitempty"`
	Type        Type               `json:"type"`
}

type SelectMenuOptions = []SelectMenuOption

type SelectMenuOption struct {
	Emoji       *PartialEmoji `json:"emoji,omitempty"`
	Label       string        `json:"label"`
	Value       string        `json:"value"`
	Description string        `json:"description,omitempty"`
	Default     bool          `json:"default,omitempty"`
}

type PartialEmoji struct {
	Name     string       `json:"name"`
	ID       snowflake.ID `json:"id"`
	Animated bool         `json:"animated,omitempty"`
}
