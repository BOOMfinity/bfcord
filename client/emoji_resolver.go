package client

import (
	"github.com/BOOMfinity/bfcord/api"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/errs"
	"github.com/BOOMfinity/bfcord/internal/slices"
	"github.com/andersfylling/snowflake/v5"
)

var _ = (discord.EmojiQuery)(&emojiResolver{})

type emojiResolver struct {
	bot *client
	*api.EmojiQuery
	resolverOptions[discord.EmojiQuery]
}

func (e emojiResolver) List() (emojis []discord.Emoji, err error) {
	if !e.ignoreCache && e.bot.Store() != nil {
		guild, _err := e.bot.Guild(e.Guild()).NoAPI().Get()
		if _err == nil {
			return guild.Emojis, nil
		}
		err = _err
	}
	if !e.ignoreAPI {
		guild, _err := e.bot.Guild(e.Guild()).NoCache().Get()
		if _err == nil {
			return guild.Emojis, nil
		}
		err = _err
	}
	if err == nil {
		err = errs.ItemNotFound
	}
	return
}

func (e emojiResolver) Get(id snowflake.ID) (emoji *discord.Emoji, err error) {
	if !e.ignoreCache && e.bot.Store() != nil {
		guild, _err := e.bot.Guild(e.Guild()).NoAPI().Get()
		if len(guild.Emojis) > 0 {
			emoji = slices.Find(guild.Emojis, func(item discord.Emoji) bool {
				return item.ID == id
			})
			if emoji != nil {
				return
			}
		}
		err = _err
	}
	if !e.ignoreAPI {
		guild, _err := e.bot.Guild(e.Guild()).NoCache().Get()
		if len(guild.Emojis) > 0 {
			emoji = slices.Find(guild.Emojis, func(item discord.Emoji) bool {
				return item.ID == id
			})
			if emoji != nil {
				return
			}
		}
		err = _err
	}
	if err == nil {
		err = errs.ItemNotFound
	}
	return
}
