package discord

import (
	"github.com/BOOMfinity/bfcord/api/cdn"
	"github.com/andersfylling/snowflake/v5"
	"strconv"
	"strings"
)

// User
//
// Discord reference: https://discord.com/developers/docs/resources/user#user-object-user-structure
type User struct {
	// The user's banner hash, not url!
	Banner        string `json:"banner,omitempty"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	// Discord user's avatar HASH, NOT URL!
	//
	// To generate URL use User.AvatarURL method.
	Avatar      string          `json:"avatar,omitempty"`
	Locale      string          `json:"locale,omitempty"`
	AccentColor int64           `json:"accent_color,omitempty"`
	ID          snowflake.ID    `json:"id"`
	PublicFlags UserFlag        `json:"public_flags"`
	Flags       UserFlag        `json:"flags"`
	IsSystem    bool            `json:"system"`
	IsBot       bool            `json:"bot"`
	IsVerified  bool            `json:"verified"`
	PremiumType UserPremiumType `json:"premium_type"`
	HasMFA      bool            `json:"mfa_enabled,omitempty"`
}

func (v User) AvatarURL(size cdn.ImageSize, format cdn.ImageFormat, dynamic bool) string {
	if v.Avatar == "" {
		tag, _ := strconv.ParseInt(v.Discriminator, 10, 0)
		return cdn.Resolver.UserDefaultAvatar(uint32(tag))
	}
	if dynamic && strings.HasPrefix(v.Avatar, "a_") {
		format = cdn.ImageFormatGIF
	}
	return cdn.Resolver.UserAvatar(v.ID, v.Avatar, size, format)
}

func (v User) IsPartial() bool {
	if v.Username == "" {
		return true
	}
	if v.Discriminator == "" {
		return true
	}
	return false
}

func (v User) Tag() string {
	return v.Username + "#" + v.Discriminator
}

type UserFlag uint32

const (
	StaffUserFlag                     UserFlag = 1 << 0
	PartnerUserFlag                   UserFlag = 1 << 1
	HypeSquadUserFlag                 UserFlag = 1 << 2
	BugHunter1UserFlag                UserFlag = 1 << 3
	HypeSquadHouse1UserFlag           UserFlag = 1 << 6
	HypeSquadHouse2UserFlag           UserFlag = 1 << 7
	HypeSquadHouse3UserFlag           UserFlag = 1 << 8
	EarlySupporterUserFlag            UserFlag = 1 << 9
	PseudoTeamUserFlag                UserFlag = 1 << 10
	BugHunter2UserFlag                UserFlag = 1 << 14
	VerifiedBotUserFlag               UserFlag = 1 << 16
	VerifiedDeveloperUserFlag         UserFlag = 1 << 17
	CertifiedDiscordModeratorUserFlag UserFlag = 1 << 18
	HttpInteractionsOnlyUserFlag      UserFlag = 1 << 19
)

func (v UserFlag) Has(flag UserFlag) bool {
	return v&flag == flag
}

type UserPremiumType uint8

const (
	NonePremiumType UserPremiumType = iota
	NitroClassicPremiumType
	NitroPremiumType
)
