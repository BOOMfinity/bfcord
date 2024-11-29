package api

import (
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/internal/httpc"
	"github.com/andersfylling/snowflake/v5"
	"github.com/valyala/fasthttp"
)

type UserResolver struct {
	client *client
	ID     snowflake.ID
}

func (c UserResolver) Get() (u discord.User, _ error) {
	defer c.client.proxy.AddUser(u)
	return httpc.NewJSONRequest[discord.User](c.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("users", c.ID.String())
	})
}

func (c UserResolver) CreateDM(recipient snowflake.ID) (discord.Channel, error) {
	return httpc.NewJSONRequest[discord.Channel](c.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPost)
		b.Body(map[string]any{
			"recipient_id": recipient,
		})
		return b.Execute("users", c.ID.String(), "channels")
	})
}
