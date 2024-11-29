package client

import (
	"github.com/BOOMfinity/bfcord/api"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
)

type memberClient struct {
	api.MemberClient
	guild snowflake.ID
	id    snowflake.ID
	sess  Session
}

func (c memberClient) Get() (discord.MemberWithUser, error) {
	return getOrSet[discord.MemberWithUser](c.sess, func() (data discord.MemberWithUser, err error) {
		data.Member, err = c.sess.Cache().Members().Get(c.guild).Get(c.id)
		if err != nil {
			return
		}
		data.User, err = c.sess.Cache().Users().Get(c.id)
		return
	}, func() (discord.MemberWithUser, error) {
		return c.MemberClient.Get()
	}, func(data discord.MemberWithUser) error {
		err := c.sess.Cache().Members().Get(c.guild).Set(c.id, data.Member)
		if err != nil {
			return err
		}
		return c.sess.Cache().Users().Set(c.id, data.User)
	})
}
