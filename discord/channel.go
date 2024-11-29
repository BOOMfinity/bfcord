package discord

import (
	"github.com/andersfylling/snowflake/v5"

	"github.com/BOOMfinity/bfcord/utils"
)

type Channel struct {
	ID                            snowflake.ID                   `json:"id,omitempty"`
	Type                          ChannelType                    `json:"type,omitempty"`
	GuildID                       snowflake.ID                   `json:"guild_id,omitempty"`
	Position                      uint                           `json:"position,omitempty"`
	PermissionOverwrites          PermissionOverwrites           `json:"permission_overwrites,omitempty"`
	Name                          string                         `json:"name,omitempty"`
	Topic                         string                         `json:"topic,omitempty"`
	NSFW                          bool                           `json:"nsfw,omitempty"`
	LastMessageID                 snowflake.ID                   `json:"last_message_id,omitempty"`
	Bitrate                       uint                           `json:"bitrate,omitempty"`
	UserLimit                     uint                           `json:"user_limit,omitempty"`
	RateLimitPerUser              uint                           `json:"rate_limit_per_user,omitempty"`
	Recipients                    []User                         `json:"recipients,omitempty"`
	Icon                          string                         `json:"icon,omitempty"`
	OwnerID                       snowflake.ID                   `json:"owner_id,omitempty"`
	ApplicationID                 snowflake.ID                   `json:"application_id,omitempty"`
	Managed                       bool                           `json:"managed,omitempty"`
	ParentID                      snowflake.ID                   `json:"parent_id,omitempty"`
	LastPinTimestamp              Timestamp                      `json:"last_pin_timestamp"`
	RTCRegion                     string                         `json:"rtc_region,omitempty"`
	VideoQualityMode              uint                           `json:"video_quality_mode,omitempty"`
	MessageCount                  uint                           `json:"message_count,omitempty"`
	MemberCount                   uint                           `json:"member_count,omitempty"`
	ThreadMetadata                utils.Nullable[ThreadMetadata] `json:"thread_metadata,omitempty"`
	ThreadMember                  utils.Nullable[ThreadMember]   `json:"member,omitempty"`
	DefaultAutoArchiveDuration    uint                           `json:"default_auto_archive_duration,omitempty"`
	Permissions                   Permission                     `json:"permissions,omitempty"`
	Flags                         ChannelFlag                    `json:"flags,omitempty"`
	TotalMessageSent              uint                           `json:"total_message_sent,omitempty"`
	AvailableTags                 []ChannelTag                   `json:"available_tags,omitempty"`
	AppliedTags                   []string                       `json:"applied_tags,omitempty"`
	DefaultReactionEmoji          utils.Nullable[Emoji]          `json:"default_reaction_emoji,omitempty"`
	DefaultThreadRateLimitPerUser uint                           `json:"default_thread_rate_limit_per_user,omitempty"`
	DefaultSortOrder              ChannelSortOrder               `json:"default_sort_order,omitempty"`
	DefaultForumLayout            ChannelForumLayout             `json:"default_forum_layout,omitempty"`
}

func (ch Channel) Thread() bool {
	return ch.Type == ChannelTypePublicThread ||
		ch.Type == ChannelTypePrivateThread ||
		ch.Type == ChannelTypeAnnouncementThread
}

type ChannelTag struct {
	ID        snowflake.ID `json:"id,omitempty"`
	Name      string       `json:"name,omitempty"`
	Moderated bool         `json:"moderated,omitempty"`
	EmojiID   snowflake.ID `json:"emoji_id,omitempty"`
	EmojiName string       `json:"emoji_name,omitempty"`
}

type ChannelDefaultReaction struct {
	EmojiID   snowflake.ID `json:"emoji_id,omitempty"`
	EmojiName string       `json:"emoji_name,omitempty"`
}

type ThreadMember struct {
	ThreadID      snowflake.ID           `json:"id,omitempty"`
	UserID        snowflake.ID           `json:"user_id,omitempty"`
	JoinTimestamp Timestamp              `json:"join_timestamp"`
	Flags         uint                   `json:"flags,omitempty"`
	Member        utils.Nullable[Member] `json:"member,omitempty"`
}

type ThreadMetadata struct {
	Archived bool `json:"archived,omitempty"`
	// AutoArchiveDuration can be set to: 60, 1440, 4320, 10080
	AutoArchiveDuration uint      `json:"auto_archive_duration,omitempty"`
	ArchiveTimestamp    Timestamp `json:"archive_timestamp"`
	Locked              bool      `json:"locked,omitempty"`
	Invitable           bool      `json:"invitable,omitempty"`
	CreateTimestamp     Timestamp `json:"create_timestamp"`
}

type ChannelForumLayout uint

const (
	ChannelForumLayoutNotSet = iota
	ChannelForumLayoutListView
	ChannelForumLayoutGalleryView
)

type ChannelSortOrder uint

const (
	ChannelSortOrderLatestActivity = iota
	ChannelSortOrderCreationDate
)

type ChannelType uint

const (
	ChannelTypeText ChannelType = iota
	ChannelTypeDM
	ChannelTypeVoice
	ChannelTypeGroup
	ChannelTypeCategory
	ChannelTypeAnnouncement
	ChannelTypeAnnouncementThread ChannelType = iota + 4
	ChannelTypePublicThread
	ChannelTypePrivateThread
	ChannelTypeStageVoice
	ChannelTypeDirectory
	ChannelTypeForum
	ChannelTypeMedia
)

type ChannelFlag uint

const (
	ChannelFlagPinned               ChannelFlag = 1 << 1
	ChannelFlagRequireTag           ChannelFlag = 1 << 4
	ChannelHideMediaDownloadOptions ChannelFlag = 1 << 15
)

type FollowedChannel struct {
	ChannelID snowflake.ID `json:"channel_id,omitempty"`
	WebhookID snowflake.ID `json:"webhook_id,omitempty"`
}
