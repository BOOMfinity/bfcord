package gateway

import (
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/internal/timeconv"
	"github.com/andersfylling/snowflake/v5"
)

type IntegrationDeleteEvent struct {
	ID            snowflake.ID `json:"id"`
	GuildID       snowflake.ID `json:"guild_id"`
	ApplicationID snowflake.ID `json:"application_id"`
}

type GuildScheduledUserEvent struct {
	EventID snowflake.ID `json:"guild_scheduled_event_id"`
	UserID  snowflake.ID `json:"user_id"`
	GuildID snowflake.ID `json:"guild_id"`
}

type TypingStartEvent struct {
	Timestamp timeconv.Timestamp `json:"timestamp"`
	Member    discord.Member     `json:"member"`
	ChannelID snowflake.ID       `json:"channel_id"`
	GuildID   snowflake.ID       `json:"guild_id"`
	UserID    snowflake.ID       `json:"user_id"`
}

type InviteCreateEvent struct {
	CreatedAt  timeconv.Timestamp `json:"created_at"`
	Code       string             `json:"code"`
	Inviter    discord.User       `json:"inviter"`
	TargetUser discord.User       `json:"target_user"`
	ChannelID  snowflake.ID       `json:"channel_id"`
	GuildID    snowflake.ID       `json:"guild_id"`
	MaxAge     int                `json:"max_age"`
	MaxUses    int                `json:"max_uses"`
	TargetType int                `json:"target_type"`
	Temporary  bool               `json:"temporary"`
	// TODO: TargetApplication discord.Application
}

type InviteDeleteEvent struct {
	Code      string       `json:"code"`
	ChannelID snowflake.ID `json:"channel_id"`
	GuildID   snowflake.ID `json:"guild_id"`
}

type GuildEmojisUpdateEvent struct {
	Emojis  []discord.Emoji `json:"emojis"`
	GuildID snowflake.ID    `json:"guild_id"`
}

type GuildBanEvent struct {
	User    discord.User `json:"user"`
	GuildID snowflake.ID `json:"guild_id"`
}

type GuildRoleDeleteEvent struct {
	GuildID snowflake.ID `json:"guild_id"`
	RoleID  snowflake.ID `json:"role_id"`
}

type GuildRoleEvent struct {
	Role    discord.Role `json:"role"`
	GuildID snowflake.ID `json:"guild_id"`
}

type ReadyEvent struct {
	SessionID        string `json:"session_id"`
	ResumeGatewayURL string `json:"resume_gateway_url"`
	Guilds           []struct {
		ID snowflake.ID `json:"id"`
	} `json:"guilds"`
	Shard []uint16     `json:"shard"`
	User  discord.User `json:"user"`
}

type MessageDeleteEvent struct {
	ChannelID snowflake.ID `json:"channel_id"`
	GuildID   snowflake.ID `json:"guild_id"`
	ID        snowflake.ID `json:"id"`
}

func (v MessageDeleteEvent) Message(api discord.ClientQuery) (discord.Message, error) {
	return api.Channel(v.ChannelID).Message(v.ID).Get()
}

type UnavailableGuild struct {
	ID          snowflake.ID `json:"id"`
	Unavailable bool         `json:"unavailable"`
}

type ThreadMembersUpdateEvent struct {
	AddedMembers     []discord.ThreadMember `json:"added_members"`
	RemovedMemberIDs []snowflake.ID         `json:"removed_member_ids"`
	ID               snowflake.ID           `json:"id"`
	GuildID          snowflake.ID           `json:"guild_id"`
	MemberCount      int                    `json:"member_count"`
}

type ThreadListSyncEvent struct {
	ChannelIDs []snowflake.ID         `json:"channel_ids"`
	Threads    []discord.Channel      `json:"threads"`
	Members    []discord.ThreadMember `json:"members"`
	GuildID    snowflake.ID           `json:"guild_id"`
}

type ChannelPinsUpdateEvent struct {
	LastPinTimestamp timeconv.Timestamp `json:"last_pin_timestamp"`
	GuildID          snowflake.ID       `json:"guild_id"`
	ChannelID        snowflake.ID       `json:"channel_id"`
}

type MessageReactionAddEvent struct {
	Member *discord.MemberWithUser `json:"member"`
	MessageReactionRemoveEvent
}

type MessageReactionRemoveEvent struct {
	Emoji     discord.Emoji `json:"emoji"`
	UserID    snowflake.ID  `json:"user_id"`
	ChannelID snowflake.ID  `json:"channel_id"`
	GuildID   snowflake.ID  `json:"guild_id"`
	MessageID snowflake.ID  `json:"message_id"`
}

func (v MessageReactionRemoveEvent) User(api discord.ClientQuery) (discord.User, error) {
	return api.User(v.UserID).Get()
}

func (v MessageReactionRemoveEvent) InGuild() bool {
	return v.GuildID.Valid()
}

type MessageReactionRemoveAllEvent struct {
	ChannelID snowflake.ID `json:"channel_id"`
	GuildID   snowflake.ID `json:"guild_id"`
	MessageID snowflake.ID `json:"message_id"`
}

type MessageReactionRemoveAllEmojiEvent struct {
	Emoji     discord.Emoji `json:"emoji"`
	ChannelID snowflake.ID  `json:"channel_id"`
	GuildID   snowflake.ID  `json:"guild_id"`
	MessageID snowflake.ID  `json:"message_id"`
}

type VoiceStateUpdateEvent = discord.VoiceState

type VoiceServerUpdateEvent struct {
	Token    string       `json:"token"`
	Endpoint string       `json:"endpoint"`
	GuildID  snowflake.ID `json:"guild_id"`
}

// TODO: Channel, Guild, Message
