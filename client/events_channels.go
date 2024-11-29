package client

import (
	"fmt"
	"slices"

	"github.com/BOOMfinity/golog/v2"

	"github.com/BOOMfinity/bfcord/client/events"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/ws"
)

var channelCreateEventHandler = handle[discord.Channel](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *discord.Channel) {
	if sess.Cache() != nil {
		if err := sess.Cache().Channels().Set(data.ID, *data); err != nil {
			log.Error().Throw(fmt.Errorf("failed to save channel: %w", err))
		}
	}

	sess.Events().ChannelCreate().Sender(func(handler events.ChannelCreateEvent) {
		handler(data)
	})
})

var channelUpdateEventHandler = handle[discord.Channel](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *discord.Channel) {
	var cached *discord.Channel
	if sess.Cache() != nil {
		cpy := *data
		obj, err := sess.Cache().Channels().Get(data.ID)
		if err == nil {
			cached = &obj
			cpy.LastMessageID = obj.LastMessageID
		}
		if err = sess.Cache().Channels().Set(cpy.ID, cpy); err != nil {
			log.Error().Throw(fmt.Errorf("failed to update channel: %w", err))
		}
	}

	sess.Events().ChannelUpdate().Sender(func(handler events.ChannelUpdateEvent) {
		handler(data, cached)
	})
})

var channelDeleteEventHandler = handle[discord.Channel](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *discord.Channel) {
	if sess.Cache() != nil {
		if err := sess.Cache().Channels().Delete(data.ID); err != nil {
			log.Error().Throw(fmt.Errorf("failed to delete channel: %w", err))
		}
	}

	sess.Events().ChannelDelete().Sender(func(handler events.ChannelDeleteEvent) {
		handler(data)
	})
})

var channelPinsUpdateEventHandler = handle[ws.ChannelPinsUpdateEvent](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *ws.ChannelPinsUpdateEvent) {
	if sess.Cache() != nil {
		if obj, err := sess.Cache().Channels().Get(data.ChannelID); err == nil {
			obj.LastPinTimestamp = data.LastPinTimestamp
			if err = sess.Cache().Channels().Set(obj.ID, obj); err != nil {
				log.Error().Throw(fmt.Errorf("failed to save channel: %w", err))
			}
		}
	}

	sess.Events().ChannelPinsUpdate().Sender(func(handler events.ChannelPinsUpdateEvent) {
		handler(data)
	})
})

var threadCreateEventHandler = handle[ws.ThreadCreateEvent](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *ws.ThreadCreateEvent) {
	if sess.Cache() != nil {
		if err := sess.Cache().Channels().Set(data.ID, data.Channel); err != nil {
			log.Error().Throw(fmt.Errorf("failed to save thread (channel): %w", err))
		}
	}

	sess.Events().ThreadCreate().Sender(func(handler events.ThreadCreateEvent) {
		handler(data)
	})
})

var threadUpdateEventHandler = handle[discord.Channel](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *discord.Channel) {
	var cached *discord.Channel
	if sess.Cache() != nil {
		if thread, err := sess.Cache().Channels().Get(data.ID); err == nil {
			cached = &thread
			data.LastMessageID = thread.LastMessageID
		}
		if err := sess.Cache().Channels().Set(data.ID, *data); err != nil {
			log.Error().Throw(fmt.Errorf("failed to save thread (channel): %w", err))
		}
	}

	sess.Events().ThreadUpdate().Sender(func(handler events.ThreadUpdateEvent) {
		handler(data, cached)
	})
})

var threadDeleteEventHandler = handle[ws.ThreadDeleteEvent](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *ws.ThreadDeleteEvent) {
	var cached *discord.Channel
	if sess.Cache() != nil {
		if thread, err := sess.Cache().Channels().Get(data.ID); err == nil {
			cached = &thread
		}
		if err := sess.Cache().Channels().Delete(data.ID); err != nil {
			log.Error().Throw(fmt.Errorf("failed to delete thread (channel): %w", err))
		}
	}

	sess.Events().ThreadDelete().Sender(func(handler events.ThreadDeleteEvent) {
		handler(data, cached)
	})
})

var threadListSyncEventHandler = handle[ws.ThreadListSyncEvent](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *ws.ThreadListSyncEvent) {
	if sess.Cache() != nil {
		if len(data.ChannelIDs) > 0 {
			for _, thread := range data.Threads {
				if !slices.Contains(data.ChannelIDs, thread.ParentID) {
					if err := sess.Cache().Channels().Delete(thread.ID); err != nil {
						log.Error().Throw(fmt.Errorf("failed to delete thread (channel): %w", err))
					}
				}
			}
		} else {
			threads, err := sess.Cache().Channels().Search(func(obj discord.Channel) bool {
				return obj.Thread()
			})
			if err != nil {
				log.Error().Throw(fmt.Errorf("failed to search over threads (channel): %w", err))
			}
			for _, thread := range threads {
				if err := sess.Cache().Channels().Delete(thread.ID); err != nil {
					log.Error().Throw(fmt.Errorf("failed to delete thread (channel): %w", err))
				}
			}
		}

		for _, thread := range data.Threads {
			i := slices.IndexFunc(data.Members, func(member discord.ThreadMember) bool {
				return member.ThreadID == thread.ID
			})
			if i != -1 {
				thread.ThreadMember.Set(data.Members[i])
			}
			if err := sess.Cache().Channels().Set(thread.ID, thread); err != nil {
				log.Error().Throw(fmt.Errorf("failed to save thread (channel): %w", err))
			}
		}
	}

	sess.Events().ThreadListSync().Sender(func(handler events.ThreadListSyncEvent) {
		handler(data)
	})
})

var threadMembersUpdateEventHandler = handle[ws.ThreadMembersUpdateEvent](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *ws.ThreadMembersUpdateEvent) {
	sess.Events().ThreadMembersUpdate().Sender(func(handler events.ThreadMembersUpdateEvent) {
		handler(data)
	})
})
