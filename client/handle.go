package client

import (
	"fmt"
	"github.com/BOOMfinity/go-utils/sets"

	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/discord/interactions"
	"github.com/BOOMfinity/bfcord/gateway"
	"github.com/BOOMfinity/bfcord/gateway/events"
	"github.com/andersfylling/snowflake/v5"
)

var unavailableGuilds = sets.Safe[snowflake.ID, snowflake.ID]{}

func (v *client) handle(data *gateway.Payload) {
	shard := v.Get(data.Shard)
	switch data.Event {
	case events.ShardReady:
		ev, err := gateway.PayloadTo[gateway.ReadyEvent](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		for i := range ev.Guilds {
			unavailableGuilds.PushEnd(ev.Guilds[i].ID)
		}
		v.current = ev.User.ID
		if v.Store() != nil {
			v.Store().Users().Set(ev.User.ID, ev.User)
		}
		shard.SetStatus(gateway.ShardStatusConnected)
		Execute(v.manager, func(_h ShardReadyEvent) {
			_h(v, shard, ev)
		})
	case events.MessageCreate:
		ev, err := gateway.PayloadTo[discord.Message](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		ev.Patch()
		if v.Store() != nil {
			if !ev.Author.IsPartial() {
				v.Store().Users().Set(ev.AuthorID, ev.Author)
			}
			if ev.Type == discord.MessageTypeDefault {
				v.Store().Members().UnsafeGet(ev.GuildID).Set(ev.AuthorID, ev.Member)
				v.Store().Messages().UnsafeGet(ev.ChannelID).PushStart(ev.BaseMessage)
			}
		}
		Execute(v.manager, func(_h MessageCreateEvent) {
			_h(v, shard, ev)
		})
	case events.MessageDelete:
		ev, err := gateway.PayloadTo[gateway.MessageDeleteEvent](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		var old *discord.BaseMessage
		if v.Store() != nil {
			if store, ok := v.Store().Messages().Get(ev.ChannelID); ok {
				old = store.Get(ev.ID)
			}
		}
		Execute(v.manager, func(_h MessageDeleteEvent) {
			_h(v, shard, ev, old)
		})
	case events.MessageDeleteBulk:
		// TODO: Cache
		ev, err := gateway.PayloadTo[struct {
			IDs       []snowflake.ID `json:"ids"`
			GuildID   snowflake.ID   `json:"guild_id"`
			ChannelID snowflake.ID   `json:"channel_id"`
		}](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		Execute(v.manager, func(_h MessageBulkDeleteEvent) {
			_h(v, shard, ev.IDs, ev.GuildID, ev.ChannelID)
		})
	case events.MessageUpdate:
		ev, err := gateway.PayloadTo[discord.Message](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		ev.Patch()
		var old *discord.BaseMessage
		if v.Store() != nil {
			if !ev.Author.IsPartial() {
				v.Store().Users().Set(ev.AuthorID, ev.Author)
			}
			v.Store().Members().UnsafeGet(ev.GuildID).Set(ev.AuthorID, ev.Member)
			if store, ok := v.Store().Messages().Get(ev.ChannelID); ok {
				old = store.Get(ev.ID)
			}
		}
		Execute(v.manager, func(_h MessageUpdateEvent) {
			_h(v, shard, ev, old)
		})
	case events.GuildCreate:
		ev, err := gateway.PayloadTo[discord.Guild](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		ev.Patch()
		if v.Store() != nil {
			v.Store().Guilds().Set(ev.ID, ev.BaseGuild)
			for i := range ev.Presences {
				user := ev.Presences[i].User
				if !user.IsPartial() {
					v.Store().Users().Set(user.ID, user)
				}
				v.Store().Presences().UnsafeGet(ev.ID).Set(user.ID, ev.Presences[i].BasePresence)
			}
			for i := range ev.Members {
				member := ev.Members[i]
				v.Store().Members().UnsafeGet(ev.ID).Set(member.UserID, member.Member)
				if !member.User.IsPartial() {
					v.Store().Users().Set(member.User.ID, member.User)
				}
			}
			for i := range ev.Channels {
				channel := ev.Channels[i]
				v.Store().Channels().UnsafeGet(ev.ID).Set(channel.ID, channel)
				v.Store().SetChannelGuild(channel.ID, channel.GuildID)
			}
			for i := range ev.Threads {
				thread := ev.Threads[i]
				v.Store().Channels().UnsafeGet(ev.ID).Set(thread.ID, thread)
			}
			for i := range ev.VoiceStates {
				v.Store().VoiceStates().UnsafeGet(ev.ID).Set(ev.VoiceStates[i].UserID, ev.VoiceStates[i])
			}
		}
		if unavailableGuilds.Exists(ev.ID) {
			return
		}
		unavailableGuilds.PushEnd(ev.ID)
		Execute(v.manager, func(_h GuildCreateEvent) {
			_h(v, shard, ev)
		})
	case events.MessageReactionAdd:
		ev, err := gateway.PayloadTo[gateway.MessageReactionAddEvent](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		if v.Store() != nil {
			if reactions, found := v.Store().Reactions().Get(ev.MessageID); found {
				if reaction, found := reactions.Get(ev.Emoji.ToString()); found {
					reaction.Count++
					reactions.Set(ev.Emoji.ToString(), reaction)
				}
			}
		}
		Execute(v.manager, func(_h MessageReactionAddEvent) {
			_h(v, shard, ev)
		})
	case events.MessageReactionRemove:
		ev, err := gateway.PayloadTo[gateway.MessageReactionRemoveEvent](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		if v.Store() != nil {
			if reactions, found := v.Store().Reactions().Get(ev.MessageID); found {
				if reaction, found := reactions.Get(ev.Emoji.ToString()); found {
					reaction.Count--
					reactions.Set(ev.Emoji.ToString(), reaction)
				}
			}
		}
		Execute(v.manager, func(_h MessageReactionRemoveEvent) {
			_h(v, shard, ev)
		})
	case events.MessageReactionRemoveEmoji:
		ev, err := gateway.PayloadTo[gateway.MessageReactionRemoveAllEmojiEvent](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		if v.Store() != nil {
			if reactions, found := v.Store().Reactions().Get(ev.MessageID); found {
				reactions.Delete(ev.Emoji.ToString())
			}
		}
		Execute(v.manager, func(_h MessageReactionRemoveEmojiEvent) {
			_h(v, shard, ev)
		})
	case events.MessageReactionRemoveAll:
		ev, err := gateway.PayloadTo[gateway.MessageReactionRemoveAllEvent](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		if v.Store() != nil {
			v.Store().Reactions().Delete(ev.MessageID)
		}
		Execute(v.manager, func(_h MessageReactionRemoveAllEvent) {
			_h(v, shard, ev)
		})
	case events.ChannelCreate:
		ev, err := gateway.PayloadTo[discord.Channel](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		if v.Store() != nil {
			if ev.GuildID.Valid() {
				v.Store().Channels().UnsafeGet(ev.ID).Set(ev.ID, ev)
				v.Store().SetChannelGuild(ev.ID, ev.GuildID)
			} else {
				v.Store().Private().Set(ev.ID, ev)
			}
		}
		Execute(v.manager, func(_h ChannelCreateEvent) {
			_h(v, shard, ev)
		})
	case events.ChannelUpdate:
		ev, err := gateway.PayloadTo[discord.Channel](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		patch := func(src discord.Channel, dst discord.Channel) discord.Channel {
			dst.Name = src.Name
			dst.NSFW = src.NSFW
			dst.Type = src.Type
			dst.Bitrate = src.Bitrate
			dst.Topic = src.Topic
			dst.DefaultAutoArchiveDuration = src.DefaultAutoArchiveDuration
			dst.Icon = src.Icon
			dst.LastPinTimestamp = src.LastPinTimestamp
			dst.ParentID = src.ParentID
			if src.ThreadMetadata != nil {
				dst.ThreadMetadata = src.ThreadMetadata
			}
			if src.Member != nil {
				dst.Member = src.Member
			}
			dst.RTCRegion = src.RTCRegion
			dst.UserLimit = src.UserLimit
			dst.Position = src.Position
			return dst
		}
		var old *discord.Channel
		if v.Store() != nil {
			if ev.GuildID.Valid() {
				if channel, found := v.Store().Channels().UnsafeGet(ev.GuildID).Get(ev.ID); found {
					old = &channel
					v.Store().Channels().UnsafeGet(ev.GuildID).Set(channel.ID, patch(ev, channel))
				} else {
					v.Store().Channels().UnsafeGet(ev.GuildID).Set(ev.ID, ev)
				}
			} else {
				if channel, found := v.Store().Private().Get(ev.ID); found {
					old = &channel
					v.Store().Private().Set(channel.ID, patch(ev, channel))
				} else {
					v.Store().Private().Set(ev.ID, ev)
				}
			}
		}
		Execute(v.manager, func(_h ChannelUpdateEvent) {
			_h(v, shard, ev, old)
		})
	case events.InteractionCreate:
		ev, err := gateway.PayloadTo[*interactions.Interaction](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		ev.Message.Patch()
		Execute(v.manager, func(_h InteractionEvent) {
			_h(v, shard, ev)
		})
	case events.ChannelDelete:
		ev, err := gateway.PayloadTo[discord.Channel](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		if v.Store() != nil {
			if ev.GuildID.Valid() {
				if guild, ok := v.Store().Channels().Get(ev.GuildID); ok {
					guild.Delete(ev.ID)
				}
			} else {
				v.Store().Private().Delete(ev.ID)
			}
		}
		Execute(v.manager, func(_h ChannelDeleteEvent) {
			_h(v, shard, ev)
		})
	case events.ChannelPinsUpdate:
		ev, err := gateway.PayloadTo[gateway.ChannelPinsUpdateEvent](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		Execute(v.manager, func(_h ChannelPinsUpdateEvent) {
			_h(v, shard, ev)
		})
	case events.GuildUpdate:
		ev, err := gateway.PayloadTo[discord.BaseGuild](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		ev.Patch()
		var old *discord.BaseGuild
		if v.Store() != nil {
			if _g, ok := v.Store().Guilds().Get(ev.ID); ok {
				old = &_g
			}
			v.Store().Guilds().Set(ev.ID, ev)
		}
		Execute(v.manager, func(_h GuildUpdateEvent) {
			_h(v, shard, ev, old)
		})
	case events.GuildDelete:
		ev, err := gateway.PayloadTo[gateway.UnavailableGuild](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		var old *discord.BaseGuild
		if v.Store() != nil {
			if _g, ok := v.Store().Guilds().Get(ev.ID); ok {
				old = &_g
			}
			v.Store().Guilds().Delete(ev.ID)
			v.Store().Channels().Delete(ev.ID)
			v.Store().Members().Delete(ev.ID)
			v.Store().Presences().Delete(ev.ID)
		}
		Execute(v.manager, func(_h GuildDeleteEvent) {
			_h(v, shard, ev, old)
		})
	case events.ThreadCreate:
		ev, err := gateway.PayloadTo[discord.Channel](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		if v.Store() != nil {
			if g, ok := v.Store().Channels().Get(ev.GuildID); ok {
				g.Set(ev.ID, ev)
			}
		}
		Execute(v.manager, func(_h ThreadCreateEvent) {
			_h(v, shard, ev)
		})
	case events.ThreadDelete:
		ev, err := gateway.PayloadTo[discord.Channel](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		if v.Store() != nil {
			if g, ok := v.Store().Channels().Get(ev.GuildID); ok {
				g.Delete(ev.ID)
			}
		}
		Execute(v.manager, func(_h ThreadDeleteEvent) {
			_h(v, shard, ev)
		})
	case events.ThreadUpdate:
		ev, err := gateway.PayloadTo[discord.Channel](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		var old *discord.Channel
		if v.Store() != nil {
			if g, ok := v.Store().Channels().Get(ev.GuildID); ok {
				if ch, ok := g.Get(ev.ID); ok {
					old = &ch
				}
				g.Set(ev.ID, ev)
			}
		}
		Execute(v.manager, func(_h ThreadUpdateEvent) {
			_h(v, shard, ev, old)
		})
	case events.GuildRoleCreate:
		ev, err := gateway.PayloadTo[gateway.GuildRoleEvent](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		ev.Role.GuildID = ev.GuildID
		if v.Store() != nil {
			if guild, ok := v.Store().Guilds().Get(ev.GuildID); ok {
				guild.Roles = append(guild.Roles, ev.Role)
				v.Store().Guilds().Set(ev.GuildID, guild)
			}
		}
		Execute(v.manager, func(_h GuildRoleCreateEvent) {
			_h(v, shard, ev.Role)
		})
	case events.GuildRoleDelete:
		ev, err := gateway.PayloadTo[gateway.GuildRoleDeleteEvent](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		var cache *discord.Role
		if v.Store() != nil {
			if guild, ok := v.Store().Guilds().Get(ev.GuildID); ok {
				index := -1
				for i := range guild.Roles {
					if guild.Roles[i].ID == ev.RoleID {
						index = i
						break
					}
				}
				if index != -1 {
					r := guild.Roles[index]
					cache = &r
					guild.Roles = append(guild.Roles[:index], guild.Roles[index+1:]...)
					v.Store().Guilds().Set(ev.GuildID, guild)
				}
			}
		}
		Execute(v.manager, func(_h GuildRoleDeleteEvent) {
			_h(v, shard, ev, cache)
		})
	case events.GuildRoleUpdate:
		ev, err := gateway.PayloadTo[gateway.GuildRoleEvent](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		ev.Role.GuildID = ev.GuildID
		var old *discord.Role
		if v.Store() != nil {
			if guild, ok := v.Store().Guilds().Get(ev.GuildID); ok {
				index := -1
				for i := range guild.Roles {
					if guild.Roles[i].ID == ev.Role.ID {
						index = i
						break
					}
				}
				if index != -1 {
					r := guild.Roles[index]
					old = &r
					guild.Roles[index] = ev.Role
					v.Store().Guilds().Set(ev.GuildID, guild)
				}
			}
		}
		Execute(v.manager, func(_h GuildRoleUpdateEvent) {
			_h(v, shard, ev.Role, old)
		})
	case events.GuildBanAdd, events.GuildBanRemove:
		ev, err := gateway.PayloadTo[gateway.GuildBanEvent](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		if data.Event == events.GuildBanAdd {
			Execute(v.manager, func(_h GuildBanAddEvent) {
				_h(v, shard, ev)
			})
		} else {
			Execute(v.manager, func(_h GuildBanRemoveEvent) {
				_h(v, shard, ev)
			})
		}
	case events.GuildEmojisUpdate:
		ev, err := gateway.PayloadTo[gateway.GuildEmojisUpdateEvent](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		for i := range ev.Emojis {
			ev.Emojis[i].GuildID = ev.GuildID
		}
		if v.Store() != nil {
			if guild, ok := v.Store().Guilds().Get(ev.GuildID); ok {
				guild.Emojis = ev.Emojis
				v.Store().Guilds().Set(ev.GuildID, guild)
			}
		}
		Execute(v.manager, func(_h GuildEmojisUpdateEvent) {
			_h(v, shard, ev)
		})
	case events.GuildIntegrationsUpdate:
		ev, err := gateway.PayloadTo[struct {
			GuildID snowflake.ID `json:"guild_id"`
		}](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		Execute(v.manager, func(_h GuildIntegrationsUpdateEvent) {
			_h(v, shard, ev.GuildID)
		})
	case events.GuildMemberAdd, events.GuildMemberUpdate:
		ev, err := gateway.PayloadTo[discord.MemberWithUser](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		ev.UserID = ev.User.ID
		if v.Store() != nil {
			v.Store().Users().Set(ev.User.ID, ev.User)
		}
		switch data.Event {
		case events.GuildMemberUpdate:
			var old *discord.Member
			if v.Store() != nil {
				if store, ok := v.Store().Members().Get(ev.GuildID); ok {
					if member, ok := store.Get(ev.UserID); ok {
						old = &member
					}
				}
				v.Store().Members().UnsafeGet(ev.GuildID).Set(ev.UserID, ev.Member)
			}
			Execute(v.manager, func(_h GuildMemberUpdateEvent) {
				_h(v, shard, ev, old)
			})
		case events.GuildMemberAdd:
			if v.Store() != nil {
				v.Store().Members().UnsafeGet(ev.GuildID).Set(ev.UserID, ev.Member)
			}
			Execute(v.manager, func(_h GuildMemberAddEvent) {
				_h(v, shard, ev)
			})
		}
	case events.GuildMemberRemove:
		ev, err := gateway.PayloadTo[discord.MemberWithUser](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		var old *discord.Member
		if v.Store() != nil {
			if store, ok := v.Store().Members().Get(ev.GuildID); ok {
				if member, ok := store.Get(ev.User.ID); ok {
					old = &member
				}
			}
		}
		Execute(v.manager, func(_h GuildMemberRemoveEvent) {
			_h(v, shard, ev.GuildID, ev.User, old)
		})
	case events.InviteCreate:
		ev, err := gateway.PayloadTo[gateway.InviteCreateEvent](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		Execute(v.manager, func(_h InviteCreateEvent) {
			_h(v, shard, ev)
		})
	case events.InviteDelete:
		ev, err := gateway.PayloadTo[gateway.InviteDeleteEvent](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		Execute(v.manager, func(_h InviteDeleteEvent) {
			_h(v, shard, ev)
		})
	case events.PresenceUpdate:
		ev, err := gateway.PayloadTo[discord.Presence](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		ev.UserID = ev.User.ID
		if v.Store() != nil {
			if !ev.User.IsPartial() {
				v.Store().Users().Set(ev.UserID, ev.User)
			} else if user, ok := v.Store().Users().Get(ev.UserID); ok {
				ev.User = user
			}
			v.Store().Presences().UnsafeGet(ev.GuildID).Set(ev.UserID, ev.BasePresence)
		}
		Execute(v.manager, func(_h PresenceUpdate) {
			_h(v, shard, ev)
		})
	case events.TypingStart:
		fmt.Println(string(data.Data))
		ev, err := gateway.PayloadTo[gateway.TypingStartEvent](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		Execute(v.manager, func(_h TypingStartEvent) {
			_h(v, shard, ev)
		})
	case events.UserUpdate:
		ev, err := gateway.PayloadTo[discord.User](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		var old *discord.User
		if v.Store() != nil {
			if user, ok := v.Store().Users().Get(ev.ID); ok {
				old = &user
			}
			v.Store().Users().Set(ev.ID, ev)
		}
		Execute(v.manager, func(_h UserUpdateEvent) {
			_h(v, shard, ev, old)
		})
	case events.GuildStickersUpdate:
		ev, err := gateway.PayloadTo[struct {
			Stickers []discord.GuildSticker `json:"stickers"`
			GuildID  snowflake.ID           `json:"guild_id"`
		}](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		if v.Store() != nil {
			if guild, ok := v.Store().Guilds().Get(ev.GuildID); ok {
				guild.Stickers = ev.Stickers
				v.Store().Guilds().Set(ev.GuildID, guild)
			}
		}
		Execute(v.manager, func(_h GuildStickersUpdateEvent) {
			_h(v, shard, ev.GuildID, ev.Stickers)
		})
	case events.GuildScheduledEventCreate, events.GuildScheduledEventDelete, events.GuildScheduledEventUpdate:
		ev, err := gateway.PayloadTo[discord.GuildScheduledEvent](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		switch data.Event {
		case events.GuildScheduledEventUpdate:
			Execute(v.manager, func(_h GuildScheduledUpdateEvent) {
				_h(v, shard, ev)
			})
		case events.GuildScheduledEventDelete:
			Execute(v.manager, func(_h GuildScheduledDeleteEvent) {
				_h(v, shard, ev)
			})
		case events.GuildScheduledEventCreate:
			Execute(v.manager, func(_h GuildScheduledCreateEvent) {
				_h(v, shard, ev)
			})
		}
	case events.GuildScheduledEventUserAdd, events.GuildScheduledEventUserRemove:
		ev, err := gateway.PayloadTo[gateway.GuildScheduledUserEvent](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		switch data.Event {
		case events.GuildScheduledEventUserAdd:
			Execute(v.manager, func(_h GuildScheduledUserAddEvent) {
				_h(v, shard, ev.EventID, ev.UserID, ev.GuildID)
			})
		case events.GuildScheduledEventUserRemove:
			Execute(v.manager, func(_h GuildScheduledUserRemoveEvent) {
				_h(v, shard, ev.EventID, ev.UserID, ev.GuildID)
			})
		}
	case events.IntegrationCreate, events.IntegrationUpdate:
		ev, err := gateway.PayloadTo[discord.Integration](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		switch data.Event {
		case events.IntegrationCreate:
			Execute(v.manager, func(_h IntegrationCreateEvent) {
				_h(v, shard, ev)
			})
		case events.IntegrationUpdate:
			Execute(v.manager, func(_h IntegrationUpdateEvent) {
				_h(v, shard, ev)
			})
		}
	case events.IntegrationDelete:
		ev, err := gateway.PayloadTo[gateway.IntegrationDeleteEvent](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		Execute(v.manager, func(_h IntegrationDeleteEvent) {
			_h(v, shard, ev)
		})
	case events.VoiceStateUpdate:
		ev, err := gateway.PayloadTo[gateway.VoiceStateUpdateEvent](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		if v.Store() != nil {
			if !ev.ChannelID.Valid() {
				if store, found := v.Store().VoiceStates().Get(ev.GuildID); found {
					store.Delete(ev.UserID)
				}
			} else {
				v.Store().VoiceStates().UnsafeGet(ev.GuildID).Set(ev.UserID, ev)
			}
		}
		Execute(v.manager, func(_h VoiceStateUpdateEvent) {
			_h(v, shard, ev)
		})
	case events.VoiceServerUpdate:
		ev, err := gateway.PayloadTo[gateway.VoiceServerUpdateEvent](data)
		if err != nil {
			v.Log().Error().Send("failed unmarshalling %v event: %v", data.Event, err.Error())
			return
		}
		Execute(v.manager, func(_h VoiceServerUpdateEvent) {
			_h(v, shard, ev)
		})
	}
}
