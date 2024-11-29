package discord

import (
	"github.com/andersfylling/snowflake/v5"

	"github.com/BOOMfinity/bfcord/utils"
)

type Guild struct {
	ID                          snowflake.ID                       `json:"id,omitempty"`
	Name                        string                             `json:"name,omitempty"`
	Icon                        string                             `json:"icon,omitempty"`
	IconHash                    string                             `json:"icon_hash,omitempty"`
	Splash                      string                             `json:"splash,omitempty"`
	DiscoverySplash             string                             `json:"discovery_splash,omitempty"`
	Owner                       bool                               `json:"owner,omitempty"`
	OwnerID                     snowflake.ID                       `json:"owner_id,omitempty"`
	Permissions                 Permission                         `json:"permissions,omitempty"`
	Region                      string                             `json:"region,omitempty"`
	AFKChannelID                snowflake.ID                       `json:"afk_channel_id,omitempty"`
	AFKTimeout                  uint                               `json:"afk_timeout,omitempty"`
	WidgetEnabled               bool                               `json:"widget_enabled,omitempty"`
	WidgetChannelID             snowflake.ID                       `json:"widget_channel_id,omitempty"`
	VerificationLevel           GuildVerificationLevel             `json:"verification_level,omitempty"`
	DefaultMessageNotifications GuildMessageNotificationLevel      `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       GuildExplicitContentFilter         `json:"explicit_content_filter,omitempty"`
	Roles                       []Role                             `json:"roles,omitempty"`
	Emojis                      []Emoji                            `json:"emojis,omitempty"`
	Features                    []string                           `json:"features,omitempty"`
	MFALevel                    GuildMFALevel                      `json:"mfa_level,omitempty"`
	ApplicationID               snowflake.ID                       `json:"application_id,omitempty"`
	SystemChannelID             snowflake.ID                       `json:"system_channel_id,omitempty"`
	SystemChannelFlags          GuildSystemChannelFlag             `json:"system_channel_flags,omitempty"`
	RulesChannelID              snowflake.ID                       `json:"rules_channel_id,omitempty"`
	MaxPresences                uint                               `json:"max_presences,omitempty"`
	MaxMembers                  uint                               `json:"max_members,omitempty"`
	VanityURLCode               string                             `json:"vanity_url_code,omitempty"`
	Description                 string                             `json:"description,omitempty"`
	Banner                      string                             `json:"banner,omitempty"`
	PremiumTier                 GuildPremiumTier                   `json:"premium_tier,omitempty"`
	PremiumSubscriptionCount    uint                               `json:"premium_subscription_count,omitempty"`
	PreferredLocale             string                             `json:"preferred_locale,omitempty"`
	PublicUpdatesChannelID      snowflake.ID                       `json:"public_updates_channel_id,omitempty"`
	MaxVideoChannelUsers        uint                               `json:"max_video_channel_users,omitempty"`
	MaxStageVideoChannelUsers   uint                               `json:"max_stage_video_channel_users,omitempty"`
	ApproximateMemberCount      uint                               `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount    uint                               `json:"approximate_presence_count,omitempty"`
	WelcomeScreen               utils.Nullable[GuildWelcomeScreen] `json:"welcome_screen,omitempty"`
	NSFWLevel                   GuildNSFWLevel                     `json:"nsfw_level,omitempty"`
	PremiumProgressBarEnabled   bool                               `json:"premium_progress_bar_enabled,omitempty"`
	SafetyAlertsChannelID       snowflake.ID                       `json:"safety_alerts_channel_id,omitempty"`
	// TODO: implements stickers
}

func (g Guild) Role(id snowflake.ID) (r Role, _ bool) {
	for _, r = range g.Roles {
		if r.ID == id {
			return r, true
		}
	}
	return
}

type GuildNSFWLevel uint8

const (
	GuildNSFWLevelDefault GuildNSFWLevel = iota
	GuildNSFWLevelExplicit
	GuildNSFWLevelSafe
	GuildNSFWLevelAgeRestricted
)

type GuildWelcomeScreen struct {
	Description     string                 `json:"description,omitempty"`
	WelcomeChannels []WelcomeScreenChannel `json:"welcome_channels,omitempty"`
}

type WelcomeScreenChannel struct {
	ChannelID   snowflake.ID `json:"channel_id,omitempty"`
	EmojiID     snowflake.ID `json:"emoji_id,omitempty"`
	EmojiName   string       `json:"emoji_name,omitempty"`
	Description string       `json:"description,omitempty"`
}

type GuildPremiumTier uint8

const (
	ServerBoostTierNone GuildPremiumTier = iota
	ServerBoostTierTier1
	ServerBoostTierTier2
	ServerBoostTierTier3
)

type GuildSystemChannelFlag uint8
type GuildMFALevel uint8
type GuildExplicitContentFilter uint8
type GuildMessageNotificationLevel uint8
type GuildVerificationLevel uint8

const (
	GuildSystemChannelFlagSuppressJoinNotifications GuildSystemChannelFlag = 1 << iota
	GuildSystemChannelFlagSuppressPremiumSubscriptions
	GuildSystemChannelFlagSuppressGuildReminderNotifications
	GuildSystemChannelFlagSuppressJoinNotificationReplies
	GuildSystemChannelFlagSuppressRoleSubscriptionPurchaseNotifications
	GuildSystemChannelFlagSuppressRoleSubscriptionPurchaseNotificationReplies
)
const (
	GuildMFALevelNone GuildMFALevel = iota
	GuildMFALevelElevated
)
const (
	GuildExplicitContentFilterNone GuildExplicitContentFilter = iota
	GuildExplicitContentFilterMembersWithoutRoles
	GuildExplicitContentFilterAllMembers
)

const (
	MessageNotificationLevelAllMessages GuildMessageNotificationLevel = iota
	MessageNotificationLevelOnlyMentions
)

const (
	GuildVerificationLevelNone GuildVerificationLevel = iota
	GuildVerificationLevelLow
	GuildVerificationLevelMedium
	GuildVerificationLevelHigh
	GuildVerificationLevelVeryHigh
)
