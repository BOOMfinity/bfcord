package client

import (
	"fmt"

	"github.com/BOOMfinity/golog/v2"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/segmentio/encoding/json"

	"github.com/BOOMfinity/bfcord/client/events"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/ws"
)

var messageCreateEventHandler = handle[ws.MessageCreateEvent](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *ws.MessageCreateEvent) {
	if sess.Cache() != nil {
		if err := sess.Cache().Messages().Get(data.ChannelID).Set(data.ID, data.Message); err != nil {
			log.Error().Throw(fmt.Errorf("failed to save message: %w", err))
		}
		if !data.Member.Partial() {
			if err := sess.Cache().Members().Get(data.GuildID).Set(data.Author.ID, data.Member); err != nil {
				log.Error().Throw(fmt.Errorf("failed to save member: %w", err))
			}
		}
		if !data.Author.Partial() {
			if err := sess.Cache().Users().Set(data.Author.ID, data.Author); err != nil {
				log.Error().Throw(fmt.Errorf("failed to save message author (user): %w", err))
			}
		}
		for _, mention := range data.Mentions {
			if !mention.User.Partial() {
				if err := sess.Cache().Users().Set(mention.User.ID, mention.User); err != nil {
					log.Error().Throw(fmt.Errorf("failed to save message mention (user): %w", err))
				}
			}
			if !mention.Member.Partial() {
				if err := sess.Cache().Members().Get(data.GuildID).Set(mention.User.ID, mention.Member); err != nil {
					log.Error().Throw(fmt.Errorf("failed to save message mention (member): %w", err))
				}
			}
		}

		if obj, err := sess.Cache().Channels().Get(data.ID); err == nil {
			obj.LastMessageID = data.ID
			if err = sess.Cache().Channels().Set(data.ID, obj); err != nil {
				log.Error().Throw(fmt.Errorf("failed to update channel: %w", err))
			}
		}
	}

	sess.Events().MessageCreate().Sender(func(handler events.MessageCreateEvent) {
		handler(data)
	})
})

var messageUpdateEventHandler = handle[ws.MessageCreateEvent](func(log golog.Logger, sess Session, raw ws.InternalDispatchEvent, _ Shard, data *ws.MessageCreateEvent) {
	var cached *discord.Message
	if sess.Cache() != nil {
		if obj, err := sess.Cache().Messages().Get(data.ChannelID).Get(data.ID); err == nil {
			cached = &obj

			b, err := json.Marshal(cached)
			if err != nil {
				panic(fmt.Errorf("failed to marshal cached message: %w", err))
			}

			modified, err := jsonpatch.MergePatch(b, raw.Data)
			if err != nil {
				panic(fmt.Errorf("failed to merge messages: %w", err))
			}

			var msg discord.Message
			if err = json.Unmarshal(modified, &msg); err != nil {
				panic(fmt.Errorf("failed to unmarshal modified message: %w", err))
			}
			if err := sess.Cache().Messages().Get(data.ChannelID).Set(data.ID, msg); err != nil {
				log.Error().Throw(fmt.Errorf("failed to update message: %w", err))
			}
		} else {
			if err := sess.Cache().Messages().Get(data.ChannelID).Set(data.ID, data.Message); err != nil {
				log.Error().Throw(fmt.Errorf("failed to save message: %w", err))
			}
		}
	}
	sess.Events().MessageUpdate().Sender(func(handler events.MessageUpdateEvent) {
		handler(data, cached)
	})
})

var messageDeleteEventHandler = handle[ws.MessageDeleteEvent](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *ws.MessageDeleteEvent) {
	var cached *discord.Message
	if sess.Cache() != nil {
		if obj, err := sess.Cache().Messages().Get(data.ChannelID).Get(data.ID); err == nil {
			cached = &obj
		}
		if err := sess.Cache().Messages().Get(data.ChannelID).Delete(data.ID); err != nil {
			log.Error().Throw(fmt.Errorf("failed to delete member: %w", err))
		}
	}

	sess.Events().MessageDelete().Sender(func(handler events.MessageDeleteEvent) {
		handler(data, cached)
	})
})
