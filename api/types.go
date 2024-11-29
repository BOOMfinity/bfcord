package api

import (
	"github.com/BOOMfinity/bfcord/api/media"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
)

type BotGateway struct {
	URL    string            `json:"url"`
	Shards uint16            `json:"shards"`
	Limit  SessionStartLimit `json:"session_start_limit"`
}

type SessionStartLimit struct {
	Total          int `json:"total"`
	Remaining      int `json:"remaining"`
	ResetAfter     int `json:"reset_after"`
	MaxConcurrency int `json:"max_concurrency"`
}

type ModifyCurrentUserParams struct {
	Username string `json:"username,omitempty"`
	Avatar   []byte `json:"avatar,omitempty"`
	Banner   []byte `json:"banner,omitempty"`
}

type AllowedMentions struct {
	Parse       []string       `json:"parse,omitempty"`
	Roles       []snowflake.ID `json:"roles,omitempty"`
	Users       []snowflake.ID `json:"users,omitempty"`
	RepliedUser bool           `json:"replied_user,omitempty"`
}

type MessageFile struct {
	ID          snowflake.ID               `json:"id,omitempty"`
	Filename    string                     `json:"filename,omitempty"`
	Title       string                     `json:"title,omitempty"`
	Description string                     `json:"description,omitempty"`
	Ephemeral   bool                       `json:"ephemeral,omitempty"`
	Resolver    media.AttachmentResolverFn `json:"-"`
}

type EditMessageParams struct {
	Content         *string                `json:"content,omitempty"`
	Embeds          []discord.MessageEmbed `json:"embeds,omitempty"`
	Flags           *discord.MessageFlag   `json:"flags,omitempty"`
	AllowedMentions *AllowedMentions       `json:"allowed_mentions,omitempty"`
	Components      discord.ActionRows     `json:"components,omitempty"`
	Attachments     []MessageFile          `json:"attachments,omitempty"`
}

type CreateMessageParams struct {
	Content          string                    `json:"content,omitempty"`
	Nonce            string                    `json:"nonce,omitempty"`
	TTS              bool                      `json:"tts,omitempty"`
	Embeds           []discord.MessageEmbed    `json:"embeds,omitempty"`
	AllowedMentions  *AllowedMentions          `json:"allowed_mentions,omitempty"`
	MessageReference *discord.MessageReference `json:"message_reference,omitempty"`
	Components       discord.ActionRows        `json:"components,omitempty"`
	Attachments      []MessageFile             `json:"attachments,omitempty"`
	Flags            *discord.MessageFlag      `json:"flags,omitempty"`
	EnforceNonce     bool                      `json:"enforce_nonce,omitempty"`
	Poll             *PollCreateParams         `json:"poll,omitempty"`
}

type PollCreateParams struct {
	Question         discord.PollMedia      `json:"question,omitempty"`
	Answers          []discord.PollAnswer   `json:"answers,omitempty"`
	Duration         uint                   `json:"duration,omitempty"`
	AllowMultiselect bool                   `json:"allow_multiselect,omitempty"`
	LayoutType       discord.PollLayoutType `json:"layout_type,omitempty"`
}

type StartThreadParams struct {
	Name                string `json:"name,omitempty"`
	AutoArchiveDuration uint   `json:"auto_archive_duration,omitempty"`
	RateLimitPerUser    uint   `json:"rate_limit_per_user,omitempty"`
}

type StartThreadWithoutMessageParams struct {
	StartThreadParams
	Type      discord.ChannelType `json:"type,omitempty"`
	Invitable *bool               `json:"invitable,omitempty"`
}

type StartForumOrMediaThreadParams struct {
	StartThreadParams
	Message     EditMessageParams `json:"message,omitempty"`
	AppliedTags []snowflake.ID    `json:"applied_tags,omitempty"`
}

type ModifyGroupChannelParams struct {
	Name string `json:"name,omitempty"`
	Icon []byte `json:"icon,omitempty"`
}

type GuildChannelParams struct {
	Name                          *string                         `json:"name,omitempty"`
	Type                          *discord.ChannelType            `json:"type,omitempty"`
	Position                      *uint                           `json:"position,omitempty"`
	Topic                         *string                         `json:"topic,omitempty"`
	NSFW                          *bool                           `json:"nsfw,omitempty"`
	RateLimitPerUser              *uint                           `json:"rate_limit_per_user,omitempty"`
	Bitrate                       *uint                           `json:"bitrate,omitempty"`
	UserLimit                     *uint                           `json:"user_limit,omitempty"`
	PermissionOverwrite           []discord.PermissionOverwrite   `json:"permission_overwrite,omitempty"`
	ParentID                      *snowflake.ID                   `json:"parent_id,omitempty"`
	RTCRegion                     *string                         `json:"rtc_region,omitempty"`
	DefaultAutoArchiveDuration    *uint                           `json:"default_auto_archive_duration,omitempty"`
	Flags                         *discord.ChannelFlag            `json:"flags,omitempty"`
	AvailableTags                 []discord.ChannelTag            `json:"available_tags,omitempty"`
	DefaultReactionEmoji          *discord.ChannelDefaultReaction `json:"default_reaction_emoji,omitempty"`
	DefaultThreadRateLimitPerUser *uint                           `json:"default_thread_rate_limit_per_user,omitempty"`
	DefaultSortOrder              *discord.ChannelSortOrder       `json:"default_sort_order,omitempty"`
	DefaultForumLayout            *discord.ChannelForumLayout     `json:"default_forum_layout,omitempty"`
}

type ModifyThreadChannelParams struct {
	Name                *string              `json:"name,omitempty"`
	Archived            *bool                `json:"archived,omitempty"`
	AutoArchiveDuration *uint                `json:"auto_archive_duration,omitempty"`
	Locked              *bool                `json:"locked,omitempty"`
	Invitable           *bool                `json:"invitable,omitempty"`
	RateLimitPerUser    *uint                `json:"rate_limit_per_user,omitempty"`
	Flags               *discord.ChannelFlag `json:"flags,omitempty"`
	AppliedTags         []snowflake.ID       `json:"applied_tags,omitempty"`
}

type ModifyChannelParams struct {
	GroupDM *ModifyGroupChannelParams
	Guild   *GuildChannelParams
	Thread  *ModifyThreadChannelParams
}

type UpdateChannelPermissionsParams struct {
	Allow discord.Permission              `json:"allow,omitempty"`
	Deny  discord.Permission              `json:"deny,omitempty"`
	Type  discord.PermissionOverwriteType `json:"type,omitempty"`
}

type CreateChannelInviteParams struct {
	MaxAge              *uint                     `json:"max_age,omitempty"`
	MaxUses             *uint                     `json:"max_uses,omitempty"`
	Temporary           bool                      `json:"temporary,omitempty"`
	Unique              bool                      `json:"unique,omitempty"`
	TargetType          *discord.InviteTargetType `json:"target_type,omitempty"`
	TargetUserID        snowflake.ID              `json:"target_user_id,omitempty"`
	TargetApplicationID snowflake.ID              `json:"target_application_id,omitempty"`
}

type ArchivedThreadListParams struct {
	Before discord.Timestamp
	Limit  uint
}

type JoinedArchivedThreadListParams struct {
	Before snowflake.ID
	Limit  uint
}

type ArchivedThreadList struct {
	Threads []discord.Channel      `json:"threads,omitempty"`
	Members []discord.ThreadMember `json:"members,omitempty"`
	HasMore bool                   `json:"has_more,omitempty"`
}

type MessageReactionsParams struct {
	Type  discord.ReactionType
	After snowflake.ID
	Limit uint
}

type InteractionMessageParams struct {
	TTS             bool                   `json:"tts,omitempty"`
	Content         string                 `json:"content,omitempty"`
	Embeds          []discord.MessageEmbed `json:"embeds,omitempty"`
	AllowedMentions *AllowedMentions       `json:"allowed_mentions,omitempty"`
	Flags           *discord.MessageFlag   `json:"flags,omitempty"`
	Components      discord.ActionRows     `json:"components,omitempty"`
	Attachments     []MessageFile          `json:"attachments,omitempty"`
	Poll            *PollCreateParams      `json:"poll,omitempty"`
}

type FollowUpParams struct {
	CreateMessageParams
	AppliedTags []snowflake.ID `json:"applied_tags,omitempty"`
	ThreadName  string         `json:"thread_name,omitempty"`
}

type WebhookExecuteParams struct {
	CreateMessageParams
	AppliedTags []snowflake.ID `json:"applied_tags,omitempty"`
	ThreadName  string         `json:"thread_name,omitempty"`
	Username    string         `json:"username,omitempty"`
	AvatarURL   string         `json:"avatar_url,omitempty"`
	ThreadID    snowflake.ID   `json:"thread_id,omitempty"`
	Wait        bool           `json:"wait,omitempty"`
}

type CreateWebhookParams struct {
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

type GuildChannelPosition struct {
	ID              snowflake.ID  `json:"id,omitempty"`
	Position        uint          `json:"position,omitempty"`
	LockPermissions bool          `json:"lock_permissions,omitempty"`
	ParentID        *snowflake.ID `json:"parent_id,omitempty"`
}

type GuildRolePosition struct {
	ID       snowflake.ID `json:"id,omitempty"`
	Position uint         `json:"position,omitempty"`
}

type ModifyGuildParams struct {
	Name                        string                                 `json:"name,omitempty"`
	VerificationLevel           *discord.GuildVerificationLevel        `json:"verification_level,omitempty"`
	DefaultMessageNotifications *discord.GuildMessageNotificationLevel `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *discord.GuildExplicitContentFilter    `json:"explicit_content_filter,omitempty"`
	AFKChannelID                *snowflake.ID                          `json:"afk_channel_id,omitempty"`
	AFKTimeout                  *AFKChannelTimeout                     `json:"afk_timeout,omitempty"`
	Icon                        []byte                                 `json:"icon,omitempty"`
	OwnerID                     *snowflake.ID                          `json:"owner_id,omitempty"`
	Splash                      []byte                                 `json:"splash,omitempty"`
	DiscoverySplash             []byte                                 `json:"discovery_splash,omitempty"`
	Banner                      []byte                                 `json:"banner,omitempty"`
	SystemChannelID             *snowflake.ID                          `json:"system_channel_id,omitempty"`
	SystemChannelFlags          *discord.GuildSystemChannelFlag        `json:"system_channel_flags,omitempty"`
	RulesChannelID              *snowflake.ID                          `json:"rules_channel_id,omitempty"`
	PublicUpdatesChannelID      *snowflake.ID                          `json:"public_updates_channel_id,omitempty"`
	PreferredLocale             *string                                `json:"preferred_locale,omitempty"`
	Features                    []string                               `json:"features,omitempty"`
	Description                 *string                                `json:"description,omitempty"`
	PremiumProgressBarEnabled   *bool                                  `json:"premium_progress_bar_enabled,omitempty"`
	SafetyAlertsChannelID       *snowflake.ID                          `json:"safety_alerts_channel_id,omitempty"`
}

type AFKChannelTimeout uint

var (
	AFKChannelTimeout60   = AFKChannelTimeout(60)
	AFKChannelTimeout300  = AFKChannelTimeout(300)
	AFKChannelTimeout900  = AFKChannelTimeout(900)
	AFKChannelTimeout1800 = AFKChannelTimeout(1800)
	AFKChannelTimeout3600 = AFKChannelTimeout(3600)
)

type GuildActiveThreads struct {
	Threads []discord.Channel
	Members []discord.ThreadMember
}

type GuildMembersParams struct {
	Limit uint
	After snowflake.ID
}

type ModifyGuildMemberParams struct {
	Nick                       *string             `json:"nick,omitempty"`
	Roles                      []snowflake.ID      `json:"roles,omitempty"`
	Mute                       *bool               `json:"mute,omitempty"`
	Deaf                       *bool               `json:"deaf,omitempty"`
	ChannelID                  *snowflake.ID       `json:"channel_id,omitempty"`
	CommunicationDisabledUntil *discord.Timestamp  `json:"communication_disabled_until,omitempty"`
	Flags                      *discord.MemberFlag `json:"flags,omitempty"`
}

type GuildBanAddResponse struct {
	BannedUsers []snowflake.ID `json:"banned_users,omitempty"`
	FailedUsers []snowflake.ID `json:"failed_users,omitempty"`
}

type CreateRoleParams struct {
	Name         string              `json:"name,omitempty"`
	Permissions  *discord.Permission `json:"permissions,omitempty"`
	Color        *uint               `json:"color,omitempty"`
	Hoist        *bool               `json:"hoist,omitempty"`
	Icon         []byte              `json:"icon,omitempty"`
	UnicodeEmoji *string             `json:"unicode_emoji,omitempty"`
	Mentionable  *bool               `json:"mentionable,omitempty"`
}

type MessagePollVoters struct {
	Users []discord.User `json:"users,omitempty"`
}

type MessagePollVotersParams struct {
	After snowflake.ID
	Limit uint
}

type CreateStageInstanceParams struct {
	ChannelID             snowflake.ID              `json:"channel_id,omitempty"`
	Topic                 string                    `json:"topic,omitempty"`
	PrivacyLevel          discord.StagePrivacyLevel `json:"privacy_level,omitempty"`
	SendStartNotification bool                      `json:"send_start_notification,omitempty"`
	GuildScheduledEventID snowflake.ID              `json:"guild_scheduled_event_id,omitempty"`
}

type ModifyStageInstanceParams struct {
	Topic        *string                    `json:"topic,omitempty"`
	PrivacyLevel *discord.StagePrivacyLevel `json:"privacy_level,omitempty"`
}

type ModifyUserVoiceStateParams struct {
	ChannelID *snowflake.ID `json:"channel_id,omitempty"`
	Suppress  *bool         `json:"suppress,omitempty"`
}

type ModifyCurrentUserVoiceStateParams struct {
	ModifyUserVoiceStateParams
	RequestToSpeakTimestamp discord.Timestamp `json:"request_to_speak_timestamp,omitempty"`
}

type ModifyEmojiParams struct {
	Name  string         `json:"name,omitempty"`
	Roles []snowflake.ID `json:"roles,omitempty"`
}

type CreateEmojiParams struct {
	Name  string         `json:"name,omitempty"`
	Roles []snowflake.ID `json:"roles,omitempty"`
	Image []byte         `json:"image,omitempty"`
}

type CreateScheduledEventParams struct {
	ChannelID          snowflake.ID                       `json:"channel_id,omitempty"`
	EntityMetadata     *discord.ScheduledEventEntity      `json:"entity_metadata,omitempty"`
	Name               string                             `json:"name,omitempty"`
	PrivacyLevel       discord.ScheduledEventPrivacyLevel `json:"privacy_level,omitempty"`
	ScheduledStartTime *discord.Timestamp                 `json:"scheduled_start_time,omitempty"`
	ScheduledEndTime   *discord.Timestamp                 `json:"scheduled_end_time,omitempty"`
	Description        *string                            `json:"description,omitempty"`
	EntityType         discord.ScheduledEventEntityType   `json:"entity_type,omitempty"`
	Image              []byte                             `json:"image,omitempty"`
	// TODO
	RecurrenceRule any `json:"recurrence_rule,omitempty"`
}

type ModifyWebhookParams struct {
	Name      string       `json:"name,omitempty"`
	Avatar    []byte       `json:"avatar,omitempty"`
	ChannelID snowflake.ID `json:"channel_id,omitempty"`
}

type ModifyScheduledEventParams struct {
	CreateScheduledEventParams
	Status discord.ScheduledEventStatus `json:"status,omitempty"`
}

type ScheduledEventUser struct {
	GuildScheduledEventID snowflake.ID   `json:"guild_scheduled_event_id,omitempty"`
	User                  discord.User   `json:"user,omitempty"`
	Member                discord.Member `json:"member,omitempty"`
}

type ScheduledEventUserParams struct {
	Limit      uint
	WithMember bool
	Before     snowflake.ID
	After      snowflake.ID
}

type TextInputParams struct {
	CustomID   string              `json:"custom_id,omitempty"`
	Title      string              `json:"title,omitempty"`
	Components []discord.Component `json:"components,omitempty"`
}
