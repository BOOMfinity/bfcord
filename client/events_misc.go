package client

import (
	"fmt"

	"github.com/BOOMfinity/bfcord/client/events"
	"github.com/BOOMfinity/bfcord/ws"
	"github.com/BOOMfinity/golog/v2"
)

var readyEventHandler = handle[ws.ReadyEvent](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, shard Shard, data *ws.ReadyEvent) {
	if sess.Cache() != nil {
		if err := sess.Cache().Users().Set(data.User.ID, data.User); err != nil {
			log.Error().Throw(fmt.Errorf("failed to save user: %w", err))
		}
	}

	for _, v := range data.Guilds {
		shard.Unavailable().Set(v.ID, v)
	}

	log.Debug().Send("%d unavailable guilds added waiting to be loaded from GUILD_CREATE event", shard.Unavailable().Size())

	sess.Events().Ready().Sender(func(handler events.ReadyEvent) {
		handler(sess.Shards(), sess.ShardCount(), data)
	})
})
