package ws

import (
	"github.com/andersfylling/snowflake/v5"

	"github.com/BOOMfinity/bfcord/discord"
)

type ReadyEvent struct {
	Version          uint                `json:"v"`
	User             discord.User        `json:"user"`
	Guilds           []UnavailableGuild  `json:"guilds"`
	SessionID        string              `json:"session_id"`
	Shard            [2]uint16           `json:"shard"`
	Application      discord.Application `json:"application"`
	ResumeGatewayURL string              `json:"resume_gateway_url"`
}

type GuildCreateEvent struct {
	discord.Guild
	JoinedAt             discord.Timestamp        `json:"joined_at,omitempty"`
	Large                bool                     `json:"large,omitempty"`
	Unavailable          bool                     `json:"unavailable,omitempty"`
	MemberCount          uint                     `json:"member_count,omitempty"`
	Members              []discord.MemberWithUser `json:"members,omitempty"`
	Channels             []discord.Channel        `json:"channels,omitempty"`
	Threads              []discord.Channel        `json:"threads,omitempty"`
	Presences            []discord.Presence       `json:"presences,omitempty"`
	GuildScheduledEvents []discord.ScheduledEvent `json:"guild_scheduled_events,omitempty"`
	StageInstances       []discord.StageInstance  `json:"stage_instances,omitempty"`
	VoiceStates          []discord.VoiceState     `json:"voice_states,omitempty"`
	// TODO: stage_instances field
}

type ThreadListSyncEvent struct {
	GuildID    snowflake.ID           `json:"guild_id,omitempty"`
	ChannelIDs []snowflake.ID         `json:"channel_i_ds,omitempty"`
	Threads    []discord.Channel      `json:"threads,omitempty"`
	Members    []discord.ThreadMember `json:"members,omitempty"`
}

type MessageCreateEvent struct {
	discord.Message
	GuildID  snowflake.ID   `json:"guild_id,omitempty"`
	Member   discord.Member `json:"member,omitempty"`
	Mentions []struct {
		discord.User
		Member discord.Member `json:"member,omitempty"`
	} `json:"mentions,omitempty"`
}

type MessageDeleteEvent struct {
	ID        snowflake.ID `json:"id,omitempty"`
	ChannelID snowflake.ID `json:"channel_id,omitempty"`
	GuildID   snowflake.ID `json:"guild_id,omitempty"`
}

type ChannelPinsUpdateEvent struct {
	GuildID          snowflake.ID      `json:"guild_id,omitempty"`
	ChannelID        snowflake.ID      `json:"channel_id,omitempty"`
	LastPinTimestamp discord.Timestamp `json:"last_pin_timestamp,omitempty"`
}

type GuildBanEvent struct {
	GuildID snowflake.ID `json:"guild_id,omitempty"`
	User    discord.User `json:"user,omitempty"`
}

type ThreadCreateEvent struct {
	discord.Channel
	NewlyCreated bool `json:"newly_created,omitempty"`
}

type ThreadDeleteEvent struct {
	ID       snowflake.ID        `json:"id,omitempty"`
	GuildID  snowflake.ID        `json:"guild_id,omitempty"`
	ParentID snowflake.ID        `json:"parent_id,omitempty"`
	Type     discord.ChannelType `json:"type,omitempty"`
}

type ThreadMembersUpdateEvent struct {
	ID               snowflake.ID           `json:"id,omitempty"`
	GuildID          snowflake.ID           `json:"guild_id,omitempty"`
	MemberCount      int                    `json:"member_count,omitempty"`
	AddedMembers     []discord.ThreadMember `json:"added_members,omitempty"`
	RemovedMemberIDs []snowflake.ID         `json:"removed_member_i_ds,omitempty"`
}

type GuildRoleEvent struct {
	GuildID snowflake.ID `json:"guild_id,omitempty"`
	Role    discord.Role `json:"role"`
}

type GuildRoleDeleteEvent struct {
	RoleID  snowflake.ID `json:"role_id,omitempty"`
	GuildID snowflake.ID `json:"guild_id,omitempty"`
}

type GuildScheduledUserEvent struct {
	GuildScheduledEventID snowflake.ID `json:"guild_scheduled_event_id,omitempty"`
	UserID                snowflake.ID `json:"user_id,omitempty"`
	GuildID               snowflake.ID `json:"guild_id,omitempty"`
}

type GuildMemberAddEvent struct {
	discord.MemberWithUser
	GuildID snowflake.ID `json:"guild_id,omitempty"`
}

type GuildMemberRemoveEvent struct {
	User    discord.User `json:"user,omitempty"`
	GuildID snowflake.ID `json:"guild_id,omitempty"`
}

type GuildMemberUpdateEvent struct {
	discord.MemberWithUser
	GuildID snowflake.ID `json:"guild_id,omitempty"`
}

type GuildMembersChunkEvent struct {
	GuildID    snowflake.ID             `json:"guild_id,omitempty"`
	Members    []discord.MemberWithUser `json:"members,omitempty"`
	ChunkIndex int                      `json:"chunk_index,omitempty"`
	ChunkCount int                      `json:"chunk_count,omitempty"`
	NotFound   []snowflake.ID           `json:"not_found,omitempty"`
	Presences  []discord.Presence       `json:"presences,omitempty"`
	Nonce      string                   `json:"nonce,omitempty"`
}

type InviteDeleteEvent struct {
	ChannelID snowflake.ID `json:"channel_id,omitempty"`
	GuildID   snowflake.ID `json:"guild_id,omitempty"`
	Code      string       `json:"code,omitempty"`
}

type InviteCreateEvent struct {
	ChannelID         snowflake.ID
	Code              string
	CreatedAt         discord.Timestamp
	GuildID           snowflake.ID
	Inviter           discord.User
	MaxAge            int
	MaxUses           int
	TargetType        discord.InviteTargetType
	TargetUser        discord.User
	TargetApplication discord.Application
	Temporary         bool
	Uses              int
}

type RequestGuildMembersParams struct {
	GuildID   snowflake.ID   `json:"guild_id,omitempty"`
	Query     string         `json:"query,omitempty"`
	Limit     uint           `json:"limit,omitempty"`
	Presences bool           `json:"presences,omitempty"`
	UserIDs   []snowflake.ID `json:"user_i_ds,omitempty"`
	Nonce     string         `json:"nonce,omitempty"`
}
