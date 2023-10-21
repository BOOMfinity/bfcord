package api

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
	"github.com/valyala/fasthttp"
)

var _ = (discord.EmojiQuery)(&EmojiQuery{})

type EmojiQuery struct {
	api *Client
	emptyOptions[discord.EmojiQuery]
	guild snowflake.ID
}

func (e EmojiQuery) Guild() snowflake.ID {
	return e.guild
}

func (e EmojiQuery) List() (emojis []discord.Emoji, err error) {
	req := e.api.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/emojis", FullApiUrl, e.guild))
	return emojis, e.api.DoResult(req, &emojis)
}

func (e EmojiQuery) Delete(id snowflake.ID) error {
	req := e.api.New(true)
	req.Header.SetMethod(fasthttp.MethodDelete)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/emojis/%v", FullApiUrl, e.guild, id))
	if e.reason != "" {
		req.Header.Set("reason", e.reason)
	}
	return e.api.DoNoResp(req)
}

func (e EmojiQuery) Get(id snowflake.ID) (emoji *discord.Emoji, err error) {
	req := e.api.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/emojis/%v", FullApiUrl, e.guild, id))
	return emoji, e.api.DoResult(req, &emoji)
}

func NewEmojiQuery(client *Client, id snowflake.ID) *EmojiQuery {
	d := &EmojiQuery{
		guild: id,
		api:   client,
	}
	d.emptyOptions = emptyOptions[discord.EmojiQuery]{data: d}
	return d
}
