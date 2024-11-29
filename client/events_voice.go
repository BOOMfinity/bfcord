package client

import (
	"github.com/BOOMfinity/golog/v2"

	"github.com/BOOMfinity/bfcord/client/events"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/voice"
	"github.com/BOOMfinity/bfcord/ws"
)

var handleVoiceStateUpdate = handle[voice.StateUpdateEvent](func(log golog.Logger, sess Session, _ *ws.Event, _ Shard, data *discord.VoiceState) {
	if sess.Cache() != nil {
		if data.ChannelID == 0 {
			sess.Cache().VoiceStates().Get(data.GuildID).Delete(data.UserID)
		} else {
			sess.Cache().VoiceStates().Get(data.GuildID).Set(data.UserID, *data)
		}
	}
	sess.Events().VoiceStateUpdate().Sender(func(handler events.VoiceStateUpdateEvent) {
		handler(data)
	})
})

var handleVoiceServerUpdate = handle[voice.ServerUpdateEvent](func(log golog.Logger, sess Session, _ *ws.Event, _ Shard, data *voice.ServerUpdateEvent) {
	sess.Events().VoiceServerUpdate().Sender(func(handler events.VoiceServerUpdateEvent) {
		handler(data)
	})
})
