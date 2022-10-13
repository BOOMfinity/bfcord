package discord

import (
	"errors"
	"sort"

	"github.com/BOOMfinity/bfcord/discord/permissions"
	"github.com/BOOMfinity/bfcord/internal/slices"
	"github.com/BOOMfinity/bfcord/internal/timeconv"
	"github.com/andersfylling/snowflake/v5"
)

type ChannelWithMessage struct {
	Channel
	Message Message `json:"message"`
}

// Channel
//
// Reference: https://discord.com/developers/docs/resources/channel#channel-object-channel-structure
type Channel struct {
	LastPinTimestamp           timeconv.Timestamp      `json:"last_pin_timestamp"`
	ThreadMetadata             *ThreadMetadata         `json:"thread_metadata"`
	Member                     *ThreadMember           `json:"member"`
	Topic                      string                  `json:"topic"`
	RTCRegion                  string                  `json:"rtc_region"`
	Name                       string                  `json:"name"`
	Icon                       string                  `json:"icon"`
	Recipients                 []User                  `json:"recipients"`
	Overwrites                 []permissions.Overwrite `json:"permission_overwrites"`
	ID                         snowflake.ID            `json:"id"`
	UserLimit                  int                     `json:"user_limit"`
	RateLimitPerUser           int                     `json:"rate_limit_per_user"`
	Bitrate                    int                     `json:"bitrate"`
	LastMessageID              snowflake.ID            `json:"last_message_id"`
	OwnerID                    snowflake.ID            `json:"owner_id"`
	ApplicationID              snowflake.ID            `json:"application_id"`
	ParentID                   snowflake.ID            `json:"parent_id"`
	Position                   int                     `json:"position"`
	GuildID                    snowflake.ID            `json:"guild_id"`
	DefaultAutoArchiveDuration timeconv.Seconds        `json:"default_auto_archive_duration"`
	NSFW                       bool                    `json:"nsfw"`
	Flags                      ChannelFlag             `json:"flags"`
	TotalMessageSent           uint32                  `json:"total_message_sent"`
	Tags                       []snowflake.ID          `json:"applied_tags"`
	AvailableTags              []ForumTag              `json:"available_tags"`
	DefaultSortOrder           ForumSortOrder          `json:"default_sort_order"`
	Permissions                permissions.Permission  `json:"permissions"`
	MessageCount               uint32                  `json:"message_count"`
	Type                       ChannelType             `json:"type"`
}

type ForumSortOrder uint8

const (
	ForumSortOrderLatestActivity ForumSortOrder = iota
	ForumSortOrderCreationDate
)

func (v Channel) IsForum() bool {
	return v.Type == ChannelTypeGuildForum
}

func (v Channel) IsThread() bool {
	return v.Type == ChannelTypeNewsThread || v.Type == ChannelTypePublicThread || v.Type == ChannelTypePrivateThread
}

func (v Channel) IsStage() bool {
	return v.Type == ChannelTypeStage
}

func (v Channel) Guild(api ClientQuery) GuildQuery {
	if v.GuildID.IsZero() {
		return nil
	}
	return api.Guild(v.GuildID)
}

func (v Channel) Parent(api ClientQuery) ChannelQuery {
	if v.ParentID.IsZero() {
		return nil
	}
	return api.Channel(v.ParentID)
}

func (v Channel) MemberPermissions(api ClientQuery, member snowflake.ID) (perm permissions.Permission, err error) {
	if !v.GuildID.Valid() || !member.Valid() {
		return 0, errors.New("invalid id")
	}
	return api.Guild(v.GuildID).Member(member).PermissionsIn(v.ID)
}

func (v Channel) Edit(api ClientQuery) UpdateChannelTypeSelector {
	return api.Channel(v.ID).Edit()
}

func (v Channel) Thread() bool {
	return v.Type == ChannelTypePublicThread || v.Type == ChannelTypeNewsThread || v.Type == ChannelTypePrivateThread
}

type ChannelFlag uint8

const (
	ChannelFlagPinned     ChannelFlag = 1 << 1
	ChannelFlagRequireTag ChannelFlag = 1 << 4
)

type ForumTag struct {
	ID        snowflake.ID `json:"id"`
	Name      string       `json:"name"`
	Moderated bool         `json:"moderated"`
	EmojiID   snowflake.ID `json:"emoji_id"`
	EmojiName string       `json:"emoji_name"`
}

type ThreadMember struct {
	JoinTimestamp timeconv.Timestamp `json:"join_timestamp"`
	ID            snowflake.ID       `json:"id"`
	UserID        snowflake.ID       `json:"user_id"`
	GuildID       snowflake.ID       `json:"guild_id"`
}

type ThreadMetadata struct {
	ArchiveTimestamp    timeconv.Timestamp `json:"archive_timestamp"`
	AutoArchiveDuration timeconv.Seconds   `json:"auto_archive_duration"`
	Archived            bool               `json:"archived"`
	Locked              bool               `json:"locked"`
	Invitable           bool               `json:"invitable"`
}

type ThreadCreate struct {
	Name                *string                `json:"name,omitempty"`
	AutoArchiveDuration *ThreadArchiveDuration `json:"auto_archive_duration,omitempty"`
	Type                *ChannelType           `json:"type,omitempty"`
	Invitable           *bool                  `json:"invitable,omitempty"`
	RateLimitPerUser    uint32                 `json:"rate_limit_per_user,omitempty"`
}

type ChannelUpdate struct {
	Name                   *string                  `json:"name,omitempty"`
	Type                   *ChannelType             `json:"type,omitempty"`
	Position               *int                     `json:"position,omitempty"`
	Topic                  *string                  `json:"topic,omitempty"`
	Nsfw                   *bool                    `json:"nsfw,omitempty"`
	RateLimitPerUser       *uint32                  `json:"rate_limit_per_user,omitempty"`
	Bitrate                *uint64                  `json:"bitrate,omitempty"`
	UserLimit              *uint16                  `json:"user_limit,omitempty"`
	ParentID               *snowflake.ID            `json:"parent_id,omitempty"`
	DefaultArchiveDuration *ThreadArchiveDuration   `json:"default_archive_duration,omitempty"`
	Overwrites             *[]permissions.Overwrite `json:"overwrites,omitempty"`

	Archived            *bool                  `json:"archived,omitempty"`
	AutoArchiveDuration *ThreadArchiveDuration `json:"auto_archive_duration,omitempty"`
	Locked              *bool                  `json:"locked,omitempty"`
	Invitable           *bool                  `json:"invitable,omitempty"`
	//Icon any    `json:"icon"`
}

type ThreadArchiveDuration uint16

const (
	ThreadArchiveHour  ThreadArchiveDuration = 60
	ThreadArchiveDay   ThreadArchiveDuration = 1440
	ThreadArchive3Days ThreadArchiveDuration = 4320
	ThreadArchiveWeek  ThreadArchiveDuration = 10080
)

type ChannelType uint8

const (
	ChannelTypeText ChannelType = iota
	ChannelTypeDM
	ChannelTypeVoice
	ChannelTypeGroup
	ChannelTypeCategory
	ChannelTypeNews
	ChannelTypeStore
	ChannelTypeNewsThread ChannelType = iota + 3
	ChannelTypePublicThread
	ChannelTypePrivateThread
	ChannelTypeStage
	ChannelTypeGuildDirectory
	ChannelTypeGuildForum
)

type GuildChannelPositions struct {
	Position        *int          `json:"-"`
	LockPermissions *bool         `json:"-"`
	ParentID        *snowflake.ID `json:"-"`
}

type GuildChannelPositionsBuilder struct {
	data     map[snowflake.ID]GuildChannelPositions
	channels []Channel
}

func (x GuildChannelPositionsBuilder) Encode() (res []map[string]any) {
	for i := range x.data {
		channel := x.data[i]
		data := map[string]any{"id": i}
		if channel.Position != nil {
			data["position"] = *channel.Position
		}
		if channel.LockPermissions != nil {
			data["lock_permissions"] = *channel.LockPermissions
		}
		if channel.ParentID != nil {
			if *channel.ParentID != 69 {
				data["parent_id"] = *channel.ParentID
			} else {
				data["parent_id"] = nil
			}
		}
		res = append(res, data)
	}
	return
}

func (x GuildChannelPositionsBuilder) filter(fn func(item Channel) bool) (data []Channel) {
	for i := range x.channels {
		if fn(x.channels[i]) {
			data = append(data, x.channels[i])
		}
	}
	return
}

func (x GuildChannelPositionsBuilder) find(id snowflake.ID) (index int) {
	index = -1
	for i := range x.channels {
		if x.channels[i].ID == id {
			return i
		}
	}
	return
}

func (x GuildChannelPositionsBuilder) Pos(id snowflake.ID, pos int) {
	x.Set(id, 0, false, pos, false)
}

func (x GuildChannelPositionsBuilder) Set(id snowflake.ID, parent snowflake.ID, lock bool, pos int, removeParent bool) {
	index := x.find(id)
	if index == -1 {
		return
	}
	ch := &x.channels[index]
	if !parent.IsZero() {
		ch.ParentID = parent
	}
	if pos != -1 {
		ch.Position = pos
	}
	if removeParent {
		ch.ParentID = 0
	}
	channels := x.filter(func(item Channel) bool {
		return item.ParentID == ch.ParentID &&
			(item.Type == ch.Type ||
				(ch.Type == ChannelTypeText &&
					(item.Type == ChannelTypeNews || item.Type == ChannelTypeStore)) ||
				(ch.Type == ChannelTypeStage &&
					(item.Type == ChannelTypeStage || item.Type == ChannelTypeVoice)))
	})
	sort.SliceStable(channels, func(i, j int) bool {
		return channels[i].Position < channels[j].Position
	})
	if parent.Valid() && pos == -1 {
		pos = -2
		ch.Position = len(channels)
	}
	if ch.Position > len(channels) {
		ch.Position = len(channels)
	}
	if ch.Position < 1 {
		ch.Position = 1
	}
	_index := slices.FindIndex(channels, func(item Channel) bool {
		return item.ID == id
	})
	if _index != -1 {
		channels = append(channels[:_index], channels[_index+1:]...)
	}
	channels = append(channels[:ch.Position-1], append([]Channel{*ch}, channels[ch.Position-1:]...)...)
	for i := range channels {
		channels[i].Position = i + 1
		_ch := x.find(channels[i].ID)
		x.channels[_ch].Position = channels[i].Position
		if channels[i].ID == id {
			_data := GuildChannelPositions{
				Position:        &channels[i].Position,
				LockPermissions: &lock,
			}
			if removeParent {
				_data.ParentID = snowPtr(69)
				x.channels[_ch].ParentID = 0
			}
			if channels[i].ParentID.Valid() && !removeParent {
				_data.ParentID = &channels[i].ParentID
				x.channels[_ch].ParentID = channels[i].ParentID
			}
			x.data[channels[i].ID] = _data
		} else {
			x.data[channels[i].ID] = GuildChannelPositions{
				Position: &channels[i].Position,
			}
		}
	}
}

func snowPtr(x snowflake.ID) *snowflake.ID {
	return &x
}

func NewGuildChannelPositionsBuilder(channels []Channel) (x *GuildChannelPositionsBuilder) {
	data := new(GuildChannelPositionsBuilder)
	data.data = map[snowflake.ID]GuildChannelPositions{}
	data.channels = channels
	return data
}
