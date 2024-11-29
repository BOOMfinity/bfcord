package discord

import (
	"github.com/BOOMfinity/bfcord/utils"
	"github.com/andersfylling/snowflake/v5"
)

type Member struct {
	Nick                       string                           `json:"nick,omitempty"`
	Avatar                     string                           `json:"avatar,omitempty"`
	Roles                      []snowflake.ID                   `json:"roles,omitempty"`
	JoinedAt                   Timestamp                        `json:"joined_at"`
	PremiumSince               Timestamp                        `json:"premium_since"`
	Deaf                       bool                             `json:"deaf,omitempty"`
	Mute                       bool                             `json:"mute,omitempty"`
	Flags                      MemberFlag                       `json:"flags,omitempty"`
	Pending                    bool                             `json:"pending,omitempty"`
	Permissions                Permission                       `json:"permissions,omitempty"`
	CommunicationDisabledUntil Timestamp                        `json:"communication_disabled_until"`
	AvatarDecorationData       utils.Nullable[AvatarDecoration] `json:"avatar_decoration_data"`
}

func (m Member) Partial() bool {
	return m.JoinedAt.IsZero()
}

type AvatarDecoration struct {
	Asset string       `json:"asset,omitempty"`
	SkuID snowflake.ID `json:"sku_id,omitempty"`
}

type MemberFlag BitField

const (
	MemberFlagDidRejoin MemberFlag = 1 << iota
	MemberFlagCompletedOnBoarding
	MemberFlagBypassedVerification
	MemberFlagStartedOnBoarding
)
