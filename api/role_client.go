package api

import (
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/internal/httpc"
	"github.com/andersfylling/snowflake/v5"
	"github.com/valyala/fasthttp"
)

type RoleResolver struct {
	Guild  snowflake.ID
	ID     snowflake.ID
	client *client
}

func (r RoleResolver) Get() (discord.Role, error) {
	return httpc.NewJSONRequest[discord.Role](r.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("guilds", r.Guild.String(), "roles", r.ID.String())
	})
}

func (r RoleResolver) Modify(params CreateRoleParams, reason ...string) (discord.Role, error) {
	return httpc.NewJSONRequest[discord.Role](r.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPatch)
		b.Body(params)
		b.Reason(reason...)
		return b.Execute("guilds", r.Guild.String(), "roles", r.ID.String())
	})
}

func (r RoleResolver) Delete(reason ...string) error {
	return httpc.NewRequest(r.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		b.Reason(reason...)
		return b.Execute("guilds", r.Guild.String(), "roles", r.ID.String())
	})
}
