package api

import (
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/internal/httpc"
	"github.com/andersfylling/snowflake/v5"
	"github.com/valyala/fasthttp"
)

type EmojiResolver struct {
	client *client
	Guild  snowflake.ID
	Emoji  snowflake.ID
}

func (e EmojiResolver) Get() (discord.Emoji, error) {
	return httpc.NewJSONRequest[discord.Emoji](e.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("guilds", e.Guild.String(), "emojis", e.Emoji.String())
	})
}

func (e EmojiResolver) Modify(params ModifyEmojiParams, reason ...string) (discord.Emoji, error) {
	return httpc.NewJSONRequest[discord.Emoji](e.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPatch)
		b.Reason(reason...)
		b.Body(params)
		return b.Execute("guilds", e.Guild.String(), "emojis", e.Emoji.String())
	})
}

func (e EmojiResolver) Delete(reason ...string) error {
	return httpc.NewRequest(e.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		b.Reason(reason...)
		return b.Execute("guilds", e.Guild.String(), "emojis", e.Emoji.String())
	})
}
