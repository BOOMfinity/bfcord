package api

import (
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/internal/httpc"
	"github.com/andersfylling/snowflake/v5"
	"github.com/valyala/fasthttp"
)

type GuildEventResolver struct {
	client *client
	Guild  snowflake.ID
	Event  snowflake.ID
}

func (g GuildEventResolver) Get(withUserCount bool) (discord.ScheduledEvent, error) {
	return httpc.NewJSONRequest[discord.ScheduledEvent](g.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("guilds", g.Guild.String(), "scheduled-events", g.Event.String())
	})
}

func (g GuildEventResolver) Modify(params ModifyScheduledEventParams, reason ...string) (discord.ScheduledEvent, error) {
	return httpc.NewJSONRequest[discord.ScheduledEvent](g.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPatch)
		b.Reason(reason...)
		b.Body(params)
		return b.Execute("guilds", g.Guild.String(), "scheduled-events", g.Event.String())
	})
}

func (g GuildEventResolver) Delete(reason ...string) error {
	return httpc.NewRequest(g.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		b.Reason(reason...)
		return b.Execute("guilds", g.Guild.String(), "scheduled-events", g.Event.String())
	})
}

func (g GuildEventResolver) Users() GuildEventQuery {
	return GuildEventQueryResolver{
		client: g.client,
		Guild:  g.Guild,
		Event:  g.Event,
	}
}
