package client

import (
	"github.com/BOOMfinity/bfcord/client/events"
	"github.com/BOOMfinity/bfcord/ws"
	"github.com/BOOMfinity/golog/v2"
)

var inviteCreateEventHandler = handle[ws.InviteCreateEvent](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *ws.InviteCreateEvent) {
	sess.Events().InviteCreate().Sender(func(handler events.InviteCreateEvent) {
		handler(data)
	})
})
var inviteDeleteEventHandler = handle[ws.InviteDeleteEvent](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *ws.InviteDeleteEvent) {
	sess.Events().InviteDelete().Sender(func(handler events.InviteDeleteEvent) {
		handler(data)
	})
})
