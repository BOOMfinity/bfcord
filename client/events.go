package client

import (
	"sync"

	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/discord/interactions"
	"github.com/BOOMfinity/bfcord/gateway"
	"github.com/andersfylling/snowflake/v5"
	"go.uber.org/atomic"
)

type EventSubscriber interface {
	GuildCreate(event GuildCreateEvent) *EventHandler
	ShardReady(event ShardReadyEvent) *EventHandler
	TypingStart(event TypingStartEvent) *EventHandler

	MessageCreate(event MessageCreateEvent) *EventHandler
	MessageUpdate(event MessageUpdateEvent) *EventHandler
	MessageDelete(event MessageDeleteEvent) *EventHandler
	MessageBulkDelete(event MessageBulkDeleteEvent) *EventHandler
	ReactionAdd(event MessageReactionAddEvent) *EventHandler
	ReactionRemove(event MessageReactionRemoveEvent) *EventHandler
	ReactionRemoveAll(event MessageReactionRemoveAllEvent) *EventHandler
	ReactionRemoveEmoji(event MessageReactionRemoveEmojiEvent) *EventHandler

	ShardResumed(event ShardResumedEvent) *EventHandler
	ChannelCreate(event ChannelCreateEvent) *EventHandler
	ChannelUpdate(event ChannelUpdateEvent) *EventHandler
	ChannelPinsUpdate(event ChannelPinsUpdateEvent) *EventHandler

	ThreadCreate(event ThreadCreateEvent) *EventHandler
	ThreadUpdate(event ThreadUpdateEvent) *EventHandler
	ThreadDelete(event ThreadDeleteEvent) *EventHandler
	ThreadMemberUpdate(event ThreadMemberUpdateEvent) *EventHandler
	ThreadMembersUpdate(event ThreadMembersUpdateEvent) *EventHandler

	GuildUpdate(event GuildUpdateEvent) *EventHandler
	GuildDelete(event GuildDeleteEvent) *EventHandler
	GuildRoleAdd(event GuildRoleCreateEvent) *EventHandler
	GuildRoleUpdate(event GuildRoleUpdateEvent) *EventHandler
	GuildRoleDelete(event GuildRoleDeleteEvent) *EventHandler
	GuildBanAdd(event GuildBanAddEvent) *EventHandler
	GuildBanRemove(event GuildBanRemoveEvent) *EventHandler
	GuildEmojisUpdate(event GuildEmojisUpdateEvent) *EventHandler
	GuildStickersUpdate(event GuildStickersUpdateEvent) *EventHandler

	GuildMemberAdd(event GuildMemberAddEvent) *EventHandler
	GuildMemberRemove(event GuildMemberRemoveEvent) *EventHandler
	GuildMemberUpdate(event GuildMemberUpdateEvent) *EventHandler

	InviteCreate(event InviteCreateEvent) *EventHandler
	InviteDelete(event InviteDeleteEvent) *EventHandler

	PresenceUpdate(event PresenceUpdate) *EventHandler

	StageCreate(event StageEvent) *EventHandler
	StageDelete(event StageEvent) *EventHandler
	StageUpdate(event StageEvent) *EventHandler

	UserUpdate(event UserUpdateEvent) *EventHandler

	IntegrationCreate(event IntegrationCreateEvent) *EventHandler
	IntegrationUpdate(event IntegrationUpdateEvent) *EventHandler
	IntegrationDelete(event IntegrationDeleteEvent) *EventHandler
}

type StageEvent func(bot Client, shard *gateway.Shard, ev discord.StageInstance)

type ReadyEvent func(bot Client)
type InteractionEvent func(bot Client, shard *gateway.Shard, ev *interactions.Interaction)
type GuildCreateEvent func(bot Client, shard *gateway.Shard, ev discord.GuildWithData)

type MessageCreateEvent func(bot Client, shard *gateway.Shard, ev discord.Message)
type MessageUpdateEvent func(bot Client, shard *gateway.Shard, ev discord.Message, old *discord.BaseMessage)
type MessageDeleteEvent func(bot Client, shard *gateway.Shard, ev gateway.MessageDeleteEvent, cache *discord.BaseMessage)

// MessageBulkDeleteEvent
//
// TODO: Partial message type instead of ids
type MessageBulkDeleteEvent func(bot Client, shard *gateway.Shard, ids []snowflake.ID, guild snowflake.ID, channel snowflake.ID)
type MessageReactionAddEvent func(bot Client, shard *gateway.Shard, ev gateway.MessageReactionAddEvent)
type MessageReactionRemoveEvent func(bot Client, shard *gateway.Shard, ev gateway.MessageReactionRemoveEvent)
type MessageReactionRemoveAllEvent func(bot Client, shard *gateway.Shard, ev gateway.MessageReactionRemoveAllEvent)
type MessageReactionRemoveEmojiEvent func(bot Client, shard *gateway.Shard, data gateway.MessageReactionRemoveAllEmojiEvent)

type ShardReadyEvent func(bot Client, shard *gateway.Shard, ev gateway.ReadyEvent)
type ShardResumedEvent func(bot Client, shard *gateway.Shard)
type ChannelCreateEvent func(bot Client, shard *gateway.Shard, ev discord.Channel)
type ChannelUpdateEvent func(bot Client, shard *gateway.Shard, new discord.Channel, old *discord.Channel)

type GuildRoleCreateEvent func(bot Client, shard *gateway.Shard, role discord.Role)
type GuildRoleUpdateEvent func(bot Client, shard *gateway.Shard, new discord.Role, old *discord.Role)
type GuildRoleDeleteEvent func(bot Client, shard *gateway.Shard, data gateway.GuildRoleDeleteEvent, cache *discord.Role)
type GuildBanAddEvent func(bot Client, shard *gateway.Shard, data gateway.GuildBanEvent)
type GuildBanRemoveEvent func(bot Client, shard *gateway.Shard, data gateway.GuildBanEvent)
type GuildEmojisUpdateEvent func(bot Client, shard *gateway.Shard, data gateway.GuildEmojisUpdateEvent)
type GuildIntegrationsUpdateEvent func(bot Client, shard *gateway.Shard, guild snowflake.ID)
type GuildStickersUpdateEvent func(bot Client, shard *gateway.Shard, guild snowflake.ID, stickers []discord.GuildSticker)

type GuildMemberAddEvent func(bot Client, shard *gateway.Shard, member discord.MemberWithUser)
type GuildMemberUpdateEvent func(bot Client, shard *gateway.Shard, new discord.MemberWithUser, old *discord.Member)
type GuildMemberRemoveEvent func(bot Client, shard *gateway.Shard, guild snowflake.ID, user discord.User, cache *discord.Member)

type IntegrationCreateEvent func(bot Client, shard *gateway.Shard, ev discord.Integration)
type IntegrationUpdateEvent func(bot Client, shard *gateway.Shard, ev discord.Integration)
type IntegrationDeleteEvent func(bot Client, shard *gateway.Shard, ev gateway.IntegrationDeleteEvent)

type GuildScheduledCreateEvent func(bot Client, shard *gateway.Shard, data discord.GuildScheduledEvent)
type GuildScheduledDeleteEvent func(bot Client, shard *gateway.Shard, data discord.GuildScheduledEvent)
type GuildScheduledUpdateEvent func(bot Client, shard *gateway.Shard, data discord.GuildScheduledEvent)
type GuildScheduledUserAddEvent func(bot Client, shard *gateway.Shard, scheduled snowflake.ID, user snowflake.ID, guild snowflake.ID)
type GuildScheduledUserRemoveEvent func(bot Client, shard *gateway.Shard, scheduled snowflake.ID, user snowflake.ID, guild snowflake.ID)

type InviteCreateEvent func(bot Client, shard *gateway.Shard, data gateway.InviteCreateEvent)
type InviteDeleteEvent func(bot Client, shard *gateway.Shard, data gateway.InviteDeleteEvent)

type ChannelDeleteEvent func(bot Client, shard *gateway.Shard, ev discord.Channel)
type ChannelPinsUpdateEvent func(bot Client, shard *gateway.Shard, ev gateway.ChannelPinsUpdateEvent)
type ThreadCreateEvent func(bot Client, shard *gateway.Shard, ev discord.Channel)
type ThreadUpdateEvent func(bot Client, shard *gateway.Shard, ev discord.Channel, old *discord.Channel)

// TODO: Threads

// ThreadDeleteEvent
//
// Reference: https://discord.com/developers/docs/topics/gateway#thread-delete
type ThreadDeleteEvent func(bot Client, shard *gateway.Shard, ev discord.Channel)
type ThreadMemberUpdateEvent func(bot Client, shard *gateway.Shard, ev discord.ThreadMember)
type ThreadMembersUpdateEvent func(bot Client, shard *gateway.Shard, ev gateway.ThreadMembersUpdateEvent)
type GuildUpdateEvent func(bot Client, shard *gateway.Shard, ev discord.Guild, old *discord.Guild)
type GuildDeleteEvent func(bot Client, shard *gateway.Shard, ev gateway.UnavailableGuild, cache *discord.Guild)

type PresenceUpdate func(bot Client, shard *gateway.Shard, ev discord.Presence)

type TypingStartEvent func(bot Client, shard *gateway.Shard, ev gateway.TypingStartEvent)
type UserUpdateEvent func(bot Client, shard *gateway.Shard, new discord.User, old *discord.User)

type VoiceStateUpdateEvent func(bot Client, shard *gateway.Shard, ev gateway.VoiceStateUpdateEvent)
type VoiceServerUpdateEvent func(bot Client, shard *gateway.Shard, ev gateway.VoiceServerUpdateEvent)

type EventHandler struct {
	handler any
	manager *EventManager
	id      uint64
}

func (v EventHandler) Close() {
	v.manager.Remove(v.id)
}

type EventManager struct {
	id       *atomic.Uint64
	handlers []*EventHandler
	mut      sync.Mutex
}

func (v *EventManager) IntegrationCreate(event IntegrationCreateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) IntegrationUpdate(event IntegrationUpdateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) IntegrationDelete(event IntegrationDeleteEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) StickersUpdate(event GuildStickersUpdateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) UserUpdate(event UserUpdateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) TypingStart(event TypingStartEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) StageCreate(event StageEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) StageUpdate(event StageEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) StageDelete(event StageEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) InviteCreate(event InviteCreateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) PresenceUpdate(event PresenceUpdate) *EventHandler {
	return On(v, event)
}

func (v *EventManager) InviteDelete(event InviteDeleteEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) GuildMemberAdd(event GuildMemberAddEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) GuildMemberRemove(event GuildMemberRemoveEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) GuildMemberUpdate(event GuildMemberUpdateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) GuildEmojisUpdate(event GuildEmojisUpdateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) GuildBanAdd(event GuildBanAddEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) GuildBanRemove(event GuildBanRemoveEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) GuildRoleCreate(event GuildRoleCreateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) GuildRoleUpdate(event GuildRoleUpdateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) GuildRoleDelete(event GuildRoleDeleteEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) ChannelPinsUpdate(event ChannelPinsUpdateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) ThreadCreate(event ThreadCreateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) ThreadUpdate(event ThreadUpdateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) ThreadDelete(event ThreadDeleteEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) ThreadMemberUpdate(event ThreadMemberUpdateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) ThreadMembersUpdate(event ThreadMembersUpdateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) GuildUpdate(event GuildUpdateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) GuildDelete(event GuildDeleteEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) Interaction(event InteractionEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) ReactionAdd(event MessageReactionAddEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) ReactionRemove(event MessageReactionRemoveEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) ReactionRemoveAll(event MessageReactionRemoveAllEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) ReactionRemoveEmoji(event MessageReactionRemoveEmojiEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) Ready(event ReadyEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) GuildCreate(event GuildCreateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) ShardReady(event ShardReadyEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) MessageCreate(event MessageCreateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) MessageDelete(event MessageDeleteEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) MessageBulkDelete(event MessageBulkDeleteEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) MessageUpdate(event MessageUpdateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) ShardResumed(event ShardResumedEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) ChannelCreate(event ChannelCreateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) ChannelUpdate(event ChannelUpdateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) VoiceStateUpdate(event VoiceStateUpdateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) VoiceServerUpdate(event VoiceServerUpdateEvent) *EventHandler {
	return On(v, event)
}

func (v *EventManager) Remove(id uint64) {
	v.mut.Lock()
	index := -1
	for i, h := range v.handlers {
		if h.id == id {
			index = i
			break
		}
	}
	if index == -1 {
		v.mut.Unlock()
		return
	}
	copy(v.handlers[index:], v.handlers[index+1:])
	v.handlers[len(v.handlers)-1] = nil
	v.handlers = v.handlers[:len(v.handlers)-1]
	v.mut.Unlock()
}

func On[V any](manager *EventManager, fn V) *EventHandler {
	manager.mut.Lock()
	h := &EventHandler{id: manager.id.Inc(), manager: manager, handler: fn}
	manager.handlers = append(manager.handlers, h)
	manager.mut.Unlock()
	return h
}

func Execute[V any](manager *EventManager, exec func(handler V)) {
	manager.mut.Lock()
	cpy := manager.handlers[:]
	manager.mut.Unlock()
	for i := range cpy {
		h, ok := cpy[i].handler.(V)
		if ok {
			go exec(h)
		}
	}
}
