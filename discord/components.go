package discord

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/utils"
	"github.com/andersfylling/snowflake/v5"
	"github.com/segmentio/encoding/json"
)

type ComponentType uint

const (
	ComponentTypeActionRow ComponentType = iota + 1
	ComponentTypeButton
	ComponentTypeStringSelect
	ComponentTypeTextInput
	ComponentTypeUserSelect
	ComponentTypeRoleSelect
	ComponentTypeMentionableSelect
	ComponentTypeChannelSelect
)

type ButtonStyle uint

const (
	ButtonStylePrimary ButtonStyle = iota + 1
	ButtonStyleSecondary
	ButtonStyleSuccess
	ButtonStyleDanger
	ButtonStyleLink
	ButtonStylePremium
)

const (
	DefaultValueUser    SelectDefaultValueType = "user"
	DefaultValueRole    SelectDefaultValueType = "role"
	DefaultValueChannel SelectDefaultValueType = "channel"
)

type TextInputStyle uint

const (
	TextInputShort TextInputStyle = iota + 1
	TextInputLong
)

type SelectDefaultValueType string

type ActionRow []Component

type ActionRows []ActionRow

func (r *ActionRows) Get(id string) (_ Component, _ bool) {
	for _, row := range *r {
		if comp, ok := row.Get(id); ok {
			return comp, true
		}
	}
	return
}

func (r *ActionRow) Get(id string) (_ Component, _ bool) {
	for _, component := range *r {
		if component.ComponentID() == id {
			return component, true
		}
	}
	return
}

type actionRow struct {
	Type       ComponentType     `json:"type,omitempty"`
	Components []json.RawMessage `json:"components,omitempty"`
}

func (r *ActionRow) MarshalJSON() ([]byte, error) {
	var row actionRow
	row.Type = ComponentTypeActionRow
	row.Components = make([]json.RawMessage, 0, len(*r))
	for _, comp := range *r {
		data, err := json.Marshal(comp)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal component: %w", err)
		}
		row.Components = append(row.Components, data)
	}
	return json.Marshal(row)
}

func (r *ActionRow) UnmarshalJSON(data []byte) error {
	row := new(actionRow)
	if err := json.Unmarshal(data, row); err != nil {
		return fmt.Errorf("failed to unmarshal components: %w", err)
	}
	comp := new(BaseComponent)
	for _, obj := range row.Components {
		if err := json.Unmarshal(obj, comp); err != nil {
			return fmt.Errorf("failed to unmarshal component custom id: %w", err)
		}
		switch comp.Type {
		case ComponentTypeButton:
			var btn Button
			if err := json.Unmarshal(obj, &btn); err != nil {
				return fmt.Errorf("failed to unmarshal button: %w", err)
			}
			*r = append(*r, btn)
		case ComponentTypeRoleSelect, ComponentTypeMentionableSelect, ComponentTypeChannelSelect, ComponentTypeUserSelect, ComponentTypeStringSelect:
			var sel SelectMenu
			if err := json.Unmarshal(obj, &sel); err != nil {
				return fmt.Errorf("failed to unmarshal select menu: %w", err)
			}
			*r = append(*r, sel)
		case ComponentTypeTextInput:
			var input TextInput
			if err := json.Unmarshal(obj, &input); err != nil {
				return fmt.Errorf("failed to unmarshal text input: %w", err)
			}
			*r = append(*r, input)
		default:
			*r = append(*r, comp)
		}
	}
	return nil
}

type Component interface {
	ComponentID() string
	ComponentType() ComponentType
	SelectMenu() bool
}

type BaseComponent struct {
	CustomID string        `json:"custom_id,omitempty"`
	Type     ComponentType `json:"type,omitempty"`
}

func (b BaseComponent) ComponentID() string {
	return b.CustomID
}

func (b BaseComponent) ComponentType() ComponentType {
	return b.Type
}

func (b BaseComponent) SelectMenu() bool {
	return b.Type == ComponentTypeStringSelect ||
		b.Type == ComponentTypeChannelSelect ||
		b.Type == ComponentTypeRoleSelect ||
		b.Type == ComponentTypeMentionableSelect ||
		b.Type == ComponentTypeUserSelect
}

type Button struct {
	BaseComponent
	Style    ButtonStyle           `json:"style,omitempty"`
	Label    string                `json:"label,omitempty"`
	Emoji    utils.Nullable[Emoji] `json:"emoji,omitempty"`
	SkuID    snowflake.ID          `json:"sku_id,omitempty"`
	URL      string                `json:"url,omitempty"`
	Disabled bool                  `json:"disabled,omitempty"`
}

func (b Button) MarshalJSON() ([]byte, error) {
	type cpy Button
	var v cpy
	b.Type = ComponentTypeButton
	v = cpy(b)
	return json.Marshal(v)
}

type SelectMenu struct {
	BaseComponent
	Options       []SelectOption       `json:"options,omitempty"`
	ChannelTypes  []ChannelType        `json:"channel_types,omitempty"`
	Placeholder   string               `json:"placeholder,omitempty"`
	DefaultValues []SelectDefaultValue `json:"default_values,omitempty"`
	MinValues     utils.Nullable[uint] `json:"min_values,omitempty"`
	MaxValues     uint                 `json:"max_values,omitempty"`
	Disabled      bool                 `json:"disabled,omitempty"`
}

type TextInput struct {
	BaseComponent
	Style       TextInputStyle       `json:"style,omitempty"`
	Label       string               `json:"label,omitempty"`
	MinLength   utils.Nullable[uint] `json:"min_length,omitempty"`
	MaxLength   uint                 `json:"max_length,omitempty"`
	Required    utils.Nullable[uint] `json:"required,omitempty"`
	Value       string               `json:"value,omitempty"`
	Placeholder string               `json:"placeholder,omitempty"`
}

func (b TextInput) MarshalJSON() ([]byte, error) {
	type cpy TextInput
	var v cpy
	b.Type = ComponentTypeTextInput
	v = cpy(b)
	return json.Marshal(v)
}

type SelectDefaultValue struct {
	ID   snowflake.ID           `json:"id,omitempty"`
	Type SelectDefaultValueType `json:"type,omitempty"`
}

type SelectOption struct {
	Label       string                `json:"label,omitempty"`
	Value       string                `json:"value,omitempty"`
	Description string                `json:"description,omitempty"`
	Emoji       utils.Nullable[Emoji] `json:"emoji,omitempty"`
	Default     bool                  `json:"default,omitempty"`
}
