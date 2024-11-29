package discord

import (
	"bytes"
	"github.com/andersfylling/snowflake/v5"
	"github.com/segmentio/encoding/json"
	"strconv"
)

type Interaction struct {
	ID             snowflake.ID           `json:"id,omitempty"`
	Data           InteractionData        `json:"data,omitempty"`
	ApplicationID  snowflake.ID           `json:"application_id,omitempty"`
	Type           InteractionType        `json:"type,omitempty"`
	Guild          Guild                  `json:"guild"`
	GuildID        snowflake.ID           `json:"guild_id,omitempty"`
	Channel        Channel                `json:"channel"`
	ChannelID      snowflake.ID           `json:"channel_id,omitempty"`
	Member         MemberWithUser         `json:"member"`
	User           User                   `json:"user"`
	Token          string                 `json:"token,omitempty"`
	Version        uint                   `json:"version,omitempty"`
	Message        Message                `json:"message"`
	AppPermissions Permission             `json:"app_permissions,omitempty"`
	Locale         string                 `json:"locale,omitempty"`
	GuildLocale    string                 `json:"guild_locale,omitempty"`
	Entitlements   []Entitlement          `json:"entitlements,omitempty"`
	Context        InteractionContextType `json:"context,omitempty"`
}

type InteractionContextType uint

const (
	InteractionContextTypeGuild InteractionContextType = iota
	InteractionContextTypeBotDM
	InteractionContextTypePrivateChannel
)

type Entitlement struct {
	ID            snowflake.ID    `json:"id,omitempty"`
	SkuID         snowflake.ID    `json:"sku_id,omitempty"`
	ApplicationID snowflake.ID    `json:"application_id,omitempty"`
	UserID        snowflake.ID    `json:"user_id,omitempty"`
	Type          EntitlementType `json:"type,omitempty"`
	Deleted       bool            `json:"deleted,omitempty"`
	StartsAt      Timestamp       `json:"starts_at"`
	EndsAt        Timestamp       `json:"ends_at"`
	GuildID       snowflake.ID    `json:"guild_id,omitempty"`
	Consumed      bool            `json:"consumed,omitempty"`
}

type EntitlementType uint

const (
	EntitlementTypePurchase EntitlementType = iota + 1
	EntitlementTypePremiumSubscription
	EntitlementTypeDeveloperGift
	EntitlementTypeTestModePurchase
	EntitlementTypeFreePurchase
	EntitlementTypeUserGift
	EntitlementTypePremiumPurchase
	EntitlementTypeApplicationSubscription
)

type InteractionOptions []InteractionCommandOption

func (opts InteractionOptions) GetString(name string) (_ string, _ bool) {
	return parseOption[string](opts, CommandOptionTypeString, name)
}

func (opts InteractionOptions) GetBoolean(name string) (_ bool, _ bool) {
	return parseOption[bool](opts, CommandOptionTypeBoolean, name)
}

func (opts InteractionOptions) GetInteger(name string) (val int, _ bool) {
	return parseOption[int](opts, CommandOptionTypeInteger, name)
}

func (opts InteractionOptions) GetFloat(name string) (val float64, _ bool) {
	return parseOption[float64](opts, CommandOptionTypeNumber, name)
}

func (opts InteractionOptions) GetUser(name string) (val snowflake.ID, _ bool) {
	return parseOption[snowflake.ID](opts, CommandOptionTypeUser, name)
}

func (opts InteractionOptions) GetChannel(name string) (val snowflake.ID, _ bool) {
	return parseOption[snowflake.ID](opts, CommandOptionTypeChannel, name)
}

func (opts InteractionOptions) GetMentionable(name string) (val snowflake.ID, _ bool) {
	return parseOption[snowflake.ID](opts, CommandOptionTypeMentionable, name)
}

func parseOption[T any](opts InteractionOptions, typ CommandOptionType, name string) (val T, ok bool) {
	for _, opt := range opts {
		if opt.Type == typ && opt.Name == name {
			if err := json.Unmarshal(opt.Value, &val); err != nil {
				return
			}
			return val, true
		}
		if val, ok = parseOption[T](opt.Options, typ, name); ok {
			return val, true
		}
	}
	return val, false
}

type InteractionData struct {
	ID            snowflake.ID       `json:"id,omitempty"`
	Name          string             `json:"name,omitempty"`
	Type          CommandType        `json:"type,omitempty"`
	Resolved      ResolvedData       `json:"resolved,omitempty"`
	Options       InteractionOptions `json:"options,omitempty"`
	GuildID       snowflake.ID       `json:"guild_id,omitempty"`
	TargetID      snowflake.ID       `json:"target_id,omitempty"`
	CustomID      string             `json:"custom_id,omitempty"`
	ComponentType ComponentType      `json:"component_type,omitempty"`
	Values        []SelectOption     `json:"values,omitempty"`
	Components    ActionRows         `json:"components,omitempty"`
}

type InteractionCommandOption struct {
	Name    string             `json:"name,omitempty"`
	Type    CommandOptionType  `json:"type,omitempty"`
	Value   json.RawMessage    `json:"value,omitempty"`
	Options InteractionOptions `json:"options,omitempty"`
	Focused bool               `json:"focused,omitempty"`
}

func (o InteractionCommandOption) Null() bool {
	return bytes.Equal(o.Value, []byte("null"))
}

func (o InteractionCommandOption) AsString() string {
	return string(json.Unescape(o.Value))
}

func (o InteractionCommandOption) AsBoolean() bool {
	return bytes.Equal(o.Value, []byte("true"))
}

func (o InteractionCommandOption) AsNumber() int {
	i, _ := strconv.Atoi(string(json.Unescape(o.Value)))
	return i
}

func (o InteractionCommandOption) AsFloat() float64 {
	f, _ := strconv.ParseFloat(string(json.Unescape(o.Value)), 64)
	return f
}

func (o InteractionCommandOption) AsSnowflake() snowflake.ID {
	s, _ := snowflake.ParseSnowflakeUint(o.AsString(), 64)
	return s
}

type CommandOptionType uint

const (
	CommandOptionTypeSubCommand CommandOptionType = iota + 1
	CommandOptionTypeSubCommandGroup
	CommandOptionTypeString
	CommandOptionTypeInteger
	CommandOptionTypeBoolean
	CommandOptionTypeUser
	CommandOptionTypeChannel
	CommandOptionTypeRole
	CommandOptionTypeMentionable
	CommandOptionTypeNumber
	CommandOptionTypeAttachment
)

type CommandType uint

const (
	CommandTypeChatInput CommandType = iota + 1
	CommandTypeUser
	CommandTypeMessage
)

type InteractionType BitField

const (
	InteractionTypePing InteractionType = iota + 1
	InteractionTypeApplicationCommand
	InteractionTypeMessageComponent
	InteractionTypeApplicationCommandAutoComplete
	InteractionTypeModalSubmit
)

type ApplicationIntegrationType uint8

const (
	ApplicationIntegrationTypeGuild ApplicationIntegrationType = iota
	ApplicationIntegrationTypeUser
)

type ResolvedData struct {
	Users       map[snowflake.ID]User       `json:"users,omitempty"`
	Members     map[snowflake.ID]Member     `json:"members,omitempty"`
	Roles       map[snowflake.ID]Role       `json:"roles,omitempty"`
	Channels    map[snowflake.ID]Channel    `json:"channels,omitempty"`
	Messages    map[snowflake.ID]Message    `json:"messages,omitempty"`
	Attachments map[snowflake.ID]Attachment `json:"attachments,omitempty"`
}

type InteractionCallback uint

const (
	InteractionCallbackPong           InteractionCallback = iota + 1
	InteractionCallbackChannelMessage                     = iota + 3
	InteractionCallbackDeferredChannelMessage
	InteractionCallbackDeferredUpdateMessage
	InteractionCallbackUpdateMessage
	InteractionCallbackAutoCompleteResult
	InteractionCallbackModal
)
