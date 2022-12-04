package discord

import (
	"github.com/BOOMfinity/bfcord/discord/permissions"
	"github.com/BOOMfinity/bfcord/internal/timeconv"
	"github.com/BOOMfinity/go-utils/nullable"
	"github.com/andersfylling/snowflake/v5"
)

type GuildWithData struct {
	Members         []MemberWithUser      `json:"members"`
	Channels        []Channel             `json:"channels"`
	Threads         []Channel             `json:"threads"`
	Presences       []Presence            `json:"presences"`
	VoiceStates     []VoiceState          `json:"voice_states"`
	StageInstances  []StageInstance       `json:"stage_instances"`
	ScheduledEvents []GuildScheduledEvent `json:"guild_scheduled_events"`
	Owner           User                  `json:"owner"`
	Guild
}

func (g *GuildWithData) Patch() {
	g.Guild.Patch()
	for i := range g.Members {
		g.Members[i].UserID = g.Members[i].User.ID
		g.Members[i].GuildID = g.ID
	}
	for i := range g.Channels {
		g.Channels[i].GuildID = g.ID
	}
	for i := range g.Threads {
		g.Threads[i].GuildID = g.ID
	}
	for i := range g.Presences {
		g.Presences[i].UserID = g.Presences[i].User.ID
		g.Presences[i].GuildID = g.ID
	}
	for i := range g.StageInstances {
		g.StageInstances[i].GuildID = g.ID
	}
	for i := range g.ScheduledEvents {
		g.ScheduledEvents[i].GuildID = g.ID
	}
}

type Guild struct {
	JoinedAt                    timeconv.Timestamp         `json:"joined_at"`
	Banner                      string                     `json:"banner"`
	Name                        string                     `json:"name"`
	IconHash                    string                     `json:"icon_hash"`
	Splash                      string                     `json:"splash"`
	VanityUrlCode               string                     `json:"vanity_url_code"`
	PreferredLocale             string                     `json:"preferred_locale"`
	DiscoverySplash             string                     `json:"discovery_splash"`
	Icon                        string                     `json:"icon"`
	Permissions                 string                     `json:"permissions"`
	Description                 string                     `json:"description"`
	Features                    []string                   `json:"features"`
	Stickers                    []GuildSticker             `json:"stickers"`
	Roles                       RoleSlice                  `json:"roles"`
	Emojis                      []Emoji                    `json:"emojis"`
	AFKTimeout                  timeconv.Seconds           `json:"afk_timeout"`
	MFALevel                    int                        `json:"mfa_level"`
	ApplicationID               snowflake.ID               `json:"application_id"`
	SystemChannelID             snowflake.ID               `json:"system_channel_id"`
	SystemChannelFlags          int                        `json:"system_channel_flags"`
	RulesChannelID              snowflake.ID               `json:"rules_channel_id"`
	OwnerID                     snowflake.ID               `json:"owner_id"`
	PublicUpdatesChannelID      snowflake.ID               `json:"public_updates_channel_id"`
	AFKChannelID                snowflake.ID               `json:"afk_channel_id"`
	PremiumSubscriptionCount    int                        `json:"premium_subscription_count"`
	MaxPresences                int                        `json:"max_presences"`
	MaxMembers                  int                        `json:"max_members"`
	WidgetChannelID             snowflake.ID               `json:"widget_channel_id"`
	MaxVideoChannelUsers        int                        `json:"max_video_channel_users"`
	MemberCount                 int                        `json:"member_count"`
	ID                          snowflake.ID               `json:"id"`
	PremiumTier                 GuildPremiumTier           `json:"premium_tier"`
	ExplicitContentFilter       GuildExplicitContentFilter `json:"explicit_content_filter"`
	Unavailable                 bool                       `json:"unavailable"`
	DefaultMessageNotifications GuildDefaultNotifications  `json:"default_message_notifications"`
	NSFWLevel                   GuildNSFWLevel             `json:"nsfw_level"`
	VerificationLevel           GuildVerificationLevel     `json:"verification_level"`
	WidgetEnabled               bool                       `json:"widget_enabled"`
	PremiumProgressBarEnabled   bool                       `json:"premium_progress_bar_enabled"`
	Large                       bool                       `json:"large"`
}

func (v Guild) Patch() {
	for i := range v.Emojis {
		v.Emojis[i].GuildID = v.ID
	}
	for i := range v.Roles {
		v.Roles[i].GuildID = v.ID
	}
	for i := range v.Stickers {
		v.Stickers[i].GuildID = v.ID
	}
}

func (v Guild) MemberPermissions(api ClientQuery, member snowflake.ID) (perm permissions.Permission, err error) {
	if v.OwnerID == member {
		return permissions.All, nil
	}
	return api.Guild(v.ID).Member(member).Permissions()
}

type GuildUpdate struct {
	Name                        *string                          `json:"name,omitempty"`
	VerificationLevel           *GuildVerificationLevel          `json:"verification_level,omitempty"`
	DefaultMessageNotifications *GuildDefaultNotifications       `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *GuildExplicitContentFilter      `json:"explicit_content_filter,omitempty"`
	AFKChannelID                *nullable.Nullable[snowflake.ID] `json:"afk_channel_id,omitempty"`
	AFKTimeout                  *uint32                          `json:"afk_timeout,omitempty"`
	Icon                        *nullable.Nullable[string]       `json:"icon,omitempty"`
	OwnerID                     *snowflake.ID                    `json:"owner_id,omitempty"`
	Splash                      *nullable.Nullable[string]       `json:"splash,omitempty"`
	DiscoverySplash             *nullable.Nullable[string]       `json:"discovery_splash,omitempty"`
	Banner                      *nullable.Nullable[string]       `json:"banner,omitempty"`
	SystemChannelID             *nullable.Nullable[snowflake.ID] `json:"system_channel_id,omitempty"`
	SystemChannelFlags          *SystemChannelFlag               `json:"system_channel_flags,omitempty"`
	RulesChannelID              *snowflake.ID                    `json:"rules_channel_id,omitempty"`
	PublicUpdatesChannelID      *snowflake.ID                    `json:"public_updates_channel_id,omitempty"`
	PreferredLocale             *string                          `json:"preferred_locale,omitempty"`
	Description                 *string                          `json:"description,omitempty"`
	PremiumProgressBarEnabled   *bool                            `json:"premium_progress_bar_enabled,omitempty"`
}

type SystemChannelFlag uint8

const (
	ChannelFlagDisableMemberJoin SystemChannelFlag = 1 << iota
	ChannelFlagDisableBoost
	ChannelFlagDisableSetupTips
	ChannelFlagDisableReplyButton
)

type GuildSticker struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Tags        string       `json:"tags"`
	User        User         `json:"user"`
	ID          snowflake.ID `json:"id"`
	GuildID     snowflake.ID `json:"guild_id"`
	PackID      snowflake.ID `json:"pack_id"`
	SortValue   int          `json:"sort_value"`
	Available   bool         `json:"available"`
}

// StageInstance
//
// Reference: https://discord.com/developers/docs/resources/stage-instance#stage-instance-object-stage-instance-structure
type StageInstance struct {
	Topic                string       `json:"topic"`
	ID                   snowflake.ID `json:"id"`
	GuildID              snowflake.ID `json:"guild_id"`
	ChannelID            snowflake.ID `json:"channel_id"`
	PrivacyLevel         int          `json:"privacy_level"`
	DiscoverableDisabled bool         `json:"discoverable_disabled"`
}

type GuildNSFWLevel uint8

const (
	GuildNSFWDefault GuildNSFWLevel = iota
	GuildNSFWExplicit
	GuildNSFWSafe
	GuildNSFWAgeRestricted
)

type GuildPremiumTier uint8

const (
	GuildPremiumNone GuildPremiumTier = iota
	GuildPremiumTier1
	GuildPremiumTier2
	GuildPremiumTier3
)

type GuildExplicitContentFilter uint8

const (
	GuildExplicitContentFilterDisabled GuildExplicitContentFilter = iota
	GuildExplicitContentFilterWithoutRoles
	GuildExplicitContentFilterAll
)

type GuildDefaultNotifications uint8

const (
	GuildDefaultNotificationsAll GuildDefaultNotifications = iota
	GuildDefaultNotificationsMentions
)

type GuildVerificationLevel uint8

const (
	GuildVerificationNone GuildVerificationLevel = iota
	GuildVerificationLow
	GuildVerificationMedium
	GuildVerificationHigh
	GuildVerificationVeryHigh
)

// GuildScheduledEvent
//
// Reference: https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-structure
type GuildScheduledEvent struct {
	ScheduledStartTime timeconv.Timestamp   `json:"scheduled_start_time"`
	ScheduledEndTime   timeconv.Timestamp   `json:"scheduled_end_time"`
	EntityMetadata     ScheduledEventMeta   `json:"entity_metadata"`
	Name               string               `json:"name"`
	Description        string               `json:"description"`
	Creator            User                 `json:"creator"`
	ID                 snowflake.ID         `json:"id"`
	CreatorID          snowflake.ID         `json:"creator_id"`
	PrivacyLevel       int                  `json:"privacy_level"`
	ChannelID          snowflake.ID         `json:"channel_id"`
	GuildID            snowflake.ID         `json:"guild_id"`
	EntityID           snowflake.ID         `json:"entity_id"`
	UserCount          int                  `json:"user_count"`
	EntityType         ScheduledEventType   `json:"entity_type"`
	Status             ScheduledEventStatus `json:"status"`
}

type ScheduledEventMeta struct {
	Location string `json:"location"`
}

type ScheduledEventType uint8

const (
	ScheduledEventStage ScheduledEventType = iota + 1
	ScheduledEventVoice
	ScheduledEventExternal
)

type ScheduledEventStatus uint8

const (
	ScheduledEventStatusScheduled ScheduledEventStatus = iota + 1
	ScheduledEventStatusActive
	ScheduledEventStatusCompleted
	ScheduledEventStatusCanceled
)
