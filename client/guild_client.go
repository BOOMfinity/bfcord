package client

import (
	"github.com/BOOMfinity/bfcord/api"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
)

type guildClient struct {
	api.GuildClient
	id   snowflake.ID
	sess Session
}

func (c guildClient) Member(id snowflake.ID) api.MemberClient {
	return memberClient{
		MemberClient: c.GuildClient.Member(id),
		guild:        c.id,
		id:           id,
		sess:         c.sess,
	}
}

func (c guildClient) Get() (discord.Guild, error) {
	return getOrSet[discord.Guild](c.sess, func() (discord.Guild, error) {
		return c.sess.Cache().Guilds().Get(c.id)
	}, func() (discord.Guild, error) {
		return c.GuildClient.Get()
	}, func(data discord.Guild) error {
		return c.sess.Cache().Guilds().Set(c.id, data)
	})
}

func (c guildClient) Channels() ([]discord.Channel, error) {
	return getOrSet[[]discord.Channel](c.sess, func() ([]discord.Channel, error) {
		return c.sess.Cache().Channels().Search(func(channel discord.Channel) bool {
			return channel.GuildID == c.id
		})
	}, func() ([]discord.Channel, error) {
		return c.GuildClient.Channels()
	}, func(data []discord.Channel) error {
		for _, channel := range data {
			if err := c.sess.Cache().Channels().Set(c.id, channel); err != nil {
				return err
			}
		}
		return nil
	})
}
