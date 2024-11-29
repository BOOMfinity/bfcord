package client

import (
	"github.com/BOOMfinity/bfcord/api"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
)

type channelClient struct {
	api.ChannelClient
	id   snowflake.ID
	sess Session
}

func (c channelClient) Get() (discord.Channel, error) {
	return getOrSet[discord.Channel](c.sess, func() (discord.Channel, error) {
		return c.sess.Cache().Channels().Get(c.id)
	}, func() (discord.Channel, error) {
		return c.ChannelClient.Get()
	}, func(data discord.Channel) error {
		return c.sess.Cache().Channels().Set(c.id, data)
	})
}
