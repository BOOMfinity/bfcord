package events

import (
	"github.com/BOOMfinity/bfcord/api"
	"github.com/andersfylling/snowflake/v5"

	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/voice"
	"github.com/BOOMfinity/bfcord/ws"
)

// Custom events

type ReadyEvent func(shards []uint16, shardCount uint16, ready *ws.ReadyEvent)

// Invite events

type InviteCreateEvent func(event *ws.InviteCreateEvent)
type InviteDeleteEvent func(event *ws.InviteDeleteEvent)

// Guild events

type GuildCreateEvent func(guild *ws.GuildCreateEvent)
type GuildDeleteEvent func(id snowflake.ID, name string)
type GuildUpdateEvent func(new, old *discord.Guild)

type GuildBanAddEvent func(event *ws.GuildBanEvent)
type GuildBanRemoveEvent func(event *ws.GuildBanEvent)

type GuildRoleAddEvent func(event *ws.GuildRoleEvent)
type GuildRoleUpdateEvent func(event *ws.GuildRoleEvent, cached *discord.Role)
type GuildRoleDeleteEvent func(event *ws.GuildRoleDeleteEvent, cached *discord.Role)

type GuildScheduledCreateEvent func(event *discord.ScheduledEvent)
type GuildScheduledUpdateEvent func(event, cached *discord.ScheduledEvent)
type GuildScheduledDeleteEvent func(event *discord.ScheduledEvent)
type GuildScheduledUserAddEvent func(event *ws.GuildScheduledUserEvent)
type GuildScheduledUserRemoveEvent func(event *ws.GuildScheduledUserEvent)

type GuildMemberAddEvent func(event *ws.GuildMemberAddEvent)
type GuildMemberUpdateEvent func(event *ws.GuildMemberUpdateEvent, cached *discord.Member)
type GuildMemberRemoveEvent func(event *ws.GuildMemberRemoveEvent, cached *discord.Member)

// Channel events

type ChannelCreateEvent func(channel *discord.Channel)
type ChannelUpdateEvent func(old, new *discord.Channel)
type ChannelDeleteEvent func(deleted *discord.Channel)
type ChannelPinsUpdateEvent func(event *ws.ChannelPinsUpdateEvent)

// Thread events

type ThreadCreateEvent func(event *ws.ThreadCreateEvent)
type ThreadUpdateEvent func(new, old *discord.Channel)
type ThreadDeleteEvent func(event *ws.ThreadDeleteEvent, cached *discord.Channel)
type ThreadListSyncEvent func(event *ws.ThreadListSyncEvent)
type ThreadMembersUpdateEvent func(event *ws.ThreadMembersUpdateEvent)

// Message events

type MessageCreateEvent func(event *ws.MessageCreateEvent)
type MessageUpdateEvent func(new *ws.MessageCreateEvent, old *discord.Message)
type MessageDeleteEvent func(event *ws.MessageDeleteEvent, cached *discord.Message)

type InteractionCreateEvent func(i *discord.Interaction, respond api.InteractionClient)

// Voice events
type VoiceStateUpdateEvent func(event *voice.StateUpdateEvent)
type VoiceServerUpdateEvent func(event *voice.ServerUpdateEvent)
