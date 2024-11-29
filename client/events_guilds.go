package client

import (
	"fmt"
	"slices"

	"github.com/BOOMfinity/golog/v2"

	"github.com/BOOMfinity/bfcord/client/events"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/ws"
)

var guildCreateEventHandler = handle[ws.GuildCreateEvent](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, shard Shard, data *ws.GuildCreateEvent) {
	if sess.Cache() != nil {
		for _, obj := range data.Channels {
			if err := sess.Cache().Channels().Set(obj.ID, obj); err != nil {
				log.Error().Throw(fmt.Errorf("failed to save channel: %w", err))
			}
		}
		for _, obj := range data.Threads {
			if err := sess.Cache().Channels().Set(obj.ID, obj); err != nil {
				log.Error().Throw(fmt.Errorf("failed to save thread (channel): %w", err))
			}
		}
		for _, obj := range data.Members {
			if err := sess.Cache().Members().Get(data.ID).Set(obj.User.ID, obj.Member); err != nil {
				log.Error().Throw(fmt.Errorf("failed to save member: %w", err))
			}
			if err := sess.Cache().Users().Set(obj.User.ID, obj.User); err != nil {
				log.Error().Throw(fmt.Errorf("failed to save user: %w", err))
			}
		}
		for _, obj := range data.GuildScheduledEvents {
			if err := sess.Cache().ScheduledEvents().Get(data.ID).Set(obj.ID, obj); err != nil {
				log.Error().Throw(fmt.Errorf("failed to save scheduled event: %w", err))
			}
		}
		for _, obj := range data.Presences {
			if err := sess.Cache().Presences().Get(data.ID).Set(obj.User.ID, obj); err != nil {
				log.Error().Throw(fmt.Errorf("failed to save presence: %w", err))
			}
			if !obj.User.Partial() {
				if err := sess.Cache().Users().Set(obj.User.ID, obj.User); err != nil {
					log.Error().Throw(fmt.Errorf("failed to save user: %w", err))
				}
			}
		}
		for _, obj := range data.VoiceStates {
			obj.GuildID = data.ID
			if err := sess.Cache().VoiceStates().Get(data.ID).Set(obj.UserID, obj); err != nil {
				log.Error().Throw(fmt.Errorf("failed to save voice state: %w", err))
			}
		}
		if err := sess.Cache().Guilds().Set(data.ID, data.Guild); err != nil {
			log.Error().Throw(fmt.Errorf("failed to save guild: %w", err))
		}
	}

	if _, ok := shard.Unavailable().Get(data.ID); ok {
		shard.Unavailable().Delete(data.ID)
		if shard.Unavailable().Size() == 0 {
			log.Debug().Send("All unavailable guilds fetched")
		}
		return
	}

	sess.Events().GuildCreate().Sender(func(handler events.GuildCreateEvent) {
		handler(data)
	})
})

var guildUpdateEventHandler = handle[discord.Guild](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *discord.Guild) {
	var cached *discord.Guild

	if sess.Cache() != nil {
		if obj, err := sess.Cache().Guilds().Get(data.ID); err == nil {
			cached = &obj
		}
		if err := sess.Cache().Guilds().Set(data.ID, *data); err != nil {
			log.Error().Throw(fmt.Errorf("failed to save guild: %w", err))
		}
	}
	sess.Events().GuildUpdate().Sender(func(handler events.GuildUpdateEvent) {
		handler(data, cached)
	})
})

var guildDeleteEventHandler = handle[ws.UnavailableGuild](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *ws.UnavailableGuild) {
	if sess.Cache() != nil && !data.Unavailable {
		if err := sess.Cache().Guilds().Delete(data.ID); err != nil {
			log.Error().Throw(fmt.Errorf("failed to delete guild: %w", err))
		}
	}

	sess.Events().GuildDelete().Sender(func(handler events.GuildDeleteEvent) {
		handler(data.ID, data.Name)
	})
})

var guildBan = func(add bool) handleDispatchFn {
	return handle[ws.GuildBanEvent](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *ws.GuildBanEvent) {
		if add {
			sess.Events().GuildBanAdd().Sender(func(handler events.GuildBanAddEvent) {
				handler(data)
			})
		} else {
			sess.Events().GuildBanRemove().Sender(func(handler events.GuildBanRemoveEvent) {
				handler(data)
			})
		}
	})
}

var guildRoleEventHandler = func(update bool) handleDispatchFn {
	return handle[ws.GuildRoleEvent](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *ws.GuildRoleEvent) {
		var cached *discord.Role

		if sess.Cache() != nil {
			if guild, err := sess.Cache().Guilds().Get(data.GuildID); err == nil {
				index := -1
				if update {
					index = slices.IndexFunc(guild.Roles, func(role discord.Role) bool {
						return role.ID == data.Role.ID
					})
				}
				if index == -1 {
					guild.Roles = append(guild.Roles, data.Role)
				} else {
					guild.Roles[index] = data.Role
				}
				if err := sess.Cache().Guilds().Set(data.GuildID, guild); err != nil {
					log.Error().Throw(fmt.Errorf("failed to save guild: %w", err))
				}
			}
		}

		if update {
			sess.Events().GuildRoleUpdate().Sender(func(handler events.GuildRoleUpdateEvent) {
				handler(data, cached)
			})
		} else {
			sess.Events().GuildRoleAdd().Sender(func(handler events.GuildRoleAddEvent) {
				handler(data)
			})
		}
	})
}

var guildRoleDeleteEventHandler = handle[ws.GuildRoleDeleteEvent](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *ws.GuildRoleDeleteEvent) {
	var cached *discord.Role
	if sess.Cache() != nil {
		if guild, err := sess.Cache().Guilds().Get(data.GuildID); err == nil {
			index := slices.IndexFunc(guild.Roles, func(role discord.Role) bool {
				return role.ID == data.RoleID
			})
			if index != -1 {
				cached = &guild.Roles[index]
				guild.Roles = slices.Delete(guild.Roles, index, index+1)
				if err := sess.Cache().Guilds().Set(data.GuildID, guild); err != nil {
					log.Error().Throw(fmt.Errorf("failed to save guild: %w", err))
				}
			}
		}
	}

	sess.Events().GuildRoleDelete().Sender(func(handler events.GuildRoleDeleteEvent) {
		handler(data, cached)
	})
})

var guildScheduledEventHandler = func(t string) handleDispatchFn {
	return handle[discord.ScheduledEvent](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *discord.ScheduledEvent) {
		var cached *discord.ScheduledEvent
		if sess.Cache() != nil {
			switch t {
			case "create":
				if err := sess.Cache().ScheduledEvents().Get(data.GuildID).Set(data.ID, *data); err != nil {
					log.Error().Throw(fmt.Errorf("failed to save scheduled event: %w", err))
				}
			case "update":
				if scheduled, err := sess.Cache().ScheduledEvents().Get(data.GuildID).Get(data.ID); err == nil {
					cached = &scheduled
				}
				if err := sess.Cache().ScheduledEvents().Get(data.GuildID).Set(data.ID, *data); err != nil {
					log.Error().Throw(fmt.Errorf("failed to save scheduled event: %w", err))
				}
			case "delete":
				if scheduled, err := sess.Cache().ScheduledEvents().Get(data.GuildID).Get(data.ID); err == nil {
					cached = &scheduled
				}
				if err := sess.Cache().ScheduledEvents().Get(data.GuildID).Delete(data.ID); err != nil {
					log.Error().Throw(fmt.Errorf("failed to delete scheduled event: %w", err))
				}
			}
		}

		switch t {
		case "create":
			sess.Events().GuildScheduledCreate().Sender(func(handler events.GuildScheduledCreateEvent) {
				handler(data)
			})
		case "update":
			sess.Events().GuildScheduledUpdate().Sender(func(handler events.GuildScheduledUpdateEvent) {
				handler(data, cached)
			})
		case "delete":
			sess.Events().GuildScheduledDelete().Sender(func(handler events.GuildScheduledDeleteEvent) {
				handler(data)
			})
		}
	})
}

var guildScheduledUserEventHandler = func(removed bool) handleDispatchFn {
	return handle[ws.GuildScheduledUserEvent](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *ws.GuildScheduledUserEvent) {
		if removed {
			sess.Events().GuildScheduledUserRemove().Sender(func(handler events.GuildScheduledUserRemoveEvent) {
				handler(data)
			})
		} else {
			sess.Events().GuildScheduledUserAdd().Sender(func(handler events.GuildScheduledUserAddEvent) {
				handler(data)
			})
		}
	})
}

var guildMemberAddEventHandler = handle[ws.GuildMemberAddEvent](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *ws.GuildMemberAddEvent) {
	if sess.Cache() != nil {
		if err := sess.Cache().Members().Get(data.GuildID).Set(data.User.ID, data.Member); err != nil {
			log.Error().Throw(fmt.Errorf("failed to save member: %w", err))
		}
	}

	sess.Events().GuildMemberAdd().Sender(func(handler events.GuildMemberAddEvent) {
		handler(data)
	})
})

var guildMemberUpdateEventHandler = handle[ws.GuildMemberUpdateEvent](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *ws.GuildMemberUpdateEvent) {
	var cached *discord.Member
	if sess.Cache() != nil {
		if member, err := sess.Cache().Members().Get(data.GuildID).Get(data.User.ID); err == nil {
			cached = &member
		}

		if err := sess.Cache().Members().Get(data.GuildID).Set(data.User.ID, data.Member); err != nil {
			log.Error().Throw(fmt.Errorf("failed to save member: %w", err))
		}
		if err := sess.Cache().Users().Set(data.User.ID, data.User); err != nil {
			log.Error().Throw(fmt.Errorf("failed to save user: %w", err))
		}
	}

	sess.Events().GuildMemberUpdate().Sender(func(handler events.GuildMemberUpdateEvent) {
		handler(data, cached)
	})
})

var guildMemberRemoveEventHandler = handle[ws.GuildMemberRemoveEvent](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *ws.GuildMemberRemoveEvent) {
	var cached *discord.Member
	if sess.Cache() != nil {
		if member, err := sess.Cache().Members().Get(data.GuildID).Get(data.User.ID); err == nil {
			cached = &member
		}

		if err := sess.Cache().Members().Get(data.GuildID).Delete(data.User.ID); err != nil {
			log.Error().Throw(fmt.Errorf("failed to delete member: %w", err))
		}
	}

	sess.Events().GuildMemberRemove().Sender(func(handler events.GuildMemberRemoveEvent) {
		handler(data, cached)
	})
})
