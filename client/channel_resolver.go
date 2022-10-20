package client

import (
	"github.com/BOOMfinity/bfcord/api"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/errs"
)

type channelResolver struct {
	bot *client
	*api.ChannelQuery
	resolverOptions[discord.ChannelQuery]
}

func (mt channelResolver) Get() (ch discord.Channel, err error) {
	if mt.bot.Store() != nil && !mt.ignoreCache {
		guild := mt.bot.Store().ChannelGuild(mt.ID())
		if guild.Valid() {
			if store, ok := mt.bot.Store().Channels().Get(guild); ok {
				if _ch, ok := store.Get(mt.ID()); ok {
					return _ch, nil
				}
			}
		} else {
			_ch, ok := mt.bot.Store().Private().Get(mt.ID())
			if ok {
				return _ch, nil
			}
		}
	}
	if !mt.ignoreAPI {
		ch, err = mt.bot.API().Channel(mt.ID()).Get()
		if err != nil {
			return
		}
		if mt.bot.Store() != nil {
			if ch.GuildID.Valid() {
				mt.bot.Store().Channels().UnsafeGet(ch.GuildID).Set(ch.ID, ch)
				mt.bot.Store().SetChannelGuild(ch.ID, ch.GuildID)
			} else {
				mt.bot.Store().Private().Set(ch.ID, ch)
			}
		}
		return
	}
	return discord.Channel{}, errs.ItemNotFound
}
