package discord

import (
	"bytes"

	"github.com/andersfylling/snowflake/v5"
)

type Role struct {
	ID           snowflake.ID `json:"id,omitempty"`
	Name         string       `json:"name,omitempty"`
	Permissions  Permission   `json:"permissions,omitempty"`
	Color        int          `json:"color,omitempty"`
	Hoist        bool         `json:"hoist,omitempty"`
	Icon         string       `json:"icon,omitempty"`
	UnicodeEmoji string       `json:"unicode_emoji,omitempty"`
	Position     int          `json:"position,omitempty"`
	Mentionable  bool         `json:"mentionable,omitempty"`
	Managed      bool         `json:"managed,omitempty"`
	Tags         RoleTags     `json:"tags,omitempty"`
	Flags        RoleFlag     `json:"flags,omitempty"`
}

type RoleTags struct {
	BotID                 snowflake.ID `json:"bot_id,omitempty"`
	IntegrationID         snowflake.ID `json:"integration_id,omitempty"`
	PremiumSubscriber     TagNullType  `json:"premium_subscriber,omitempty"`
	SubscriptionListingID snowflake.ID `json:"subscription_listing_id,omitempty"`
	AvailableForPurchase  TagNullType  `json:"available_for_purchase,omitempty"`
	GuildConnections      TagNullType  `json:"guild_connections,omitempty"`
}

type RoleFlag uint8

const (
	RoleFlagInPrompt RoleFlag = 1 << 0
)

type TagNullType bool

func (d TagNullType) MarshalJSON() ([]byte, error) {
	if d {
		return []byte("null"), nil
	} else {
		return nil, nil
	}
}

func (d *TagNullType) UnmarshalJSON(data []byte) error {
	*d = TagNullType(bytes.Equal(data, []byte("null")))
	return nil
}
