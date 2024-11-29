package api

import (
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/internal/httpc"
	"github.com/andersfylling/snowflake/v5"
	"github.com/valyala/fasthttp"
)

type SlashResolver struct {
	client        *client
	CurrentUserID snowflake.ID
}

func (s SlashResolver) Command(id snowflake.ID) SlashCommandClient {
	return SlashCommandResolver{
		client:        s.client,
		ID:            id,
		CurrentUserID: s.CurrentUserID,
	}
}

func (s SlashResolver) Guild(id snowflake.ID) GuildSlashCommandClient {
	return GuildSlashCommandResolver{
		client:        s.client,
		ID:            id,
		CurrentUserID: s.CurrentUserID,
	}
}

func (s SlashResolver) List() ([]discord.Command, error) {
	return httpc.NewJSONRequest[[]discord.Command](s.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("applications", s.CurrentUserID.String(), "commands")
	})
}

func (s SlashResolver) Create(params discord.CreateCommand) (discord.Command, error) {
	return httpc.NewJSONRequest[discord.Command](s.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPost)
		b.Body(params)
		return b.Execute("applications", s.CurrentUserID.String(), "commands")
	})
}

func (s SlashResolver) Bulk(cmds []discord.CreateCommand) ([]discord.Command, error) {
	return httpc.NewJSONRequest[[]discord.Command](s.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPut)
		b.Body(cmds)
		return b.Execute("applications", s.CurrentUserID.String(), "commands")
	})
}

type GuildSlashCommandResolver struct {
	client        *client
	ID            snowflake.ID
	CurrentUserID snowflake.ID
}

func (g GuildSlashCommandResolver) PermissionList() ([]discord.CommandPermissions, error) {
	return httpc.NewJSONRequest[[]discord.CommandPermissions](g.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("applications", g.CurrentUserID.String(), "guilds", g.ID.String(), "permissions")
	})
}

func (g GuildSlashCommandResolver) List() ([]discord.Command, error) {
	return httpc.NewJSONRequest[[]discord.Command](g.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("applications", g.CurrentUserID.String(), "guilds", g.ID.String(), "commands")
	})
}

func (g GuildSlashCommandResolver) Create(params discord.CreateCommand) (discord.Command, error) {
	return httpc.NewJSONRequest[discord.Command](g.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPost)
		b.Body(params)
		return b.Execute("applications", g.CurrentUserID.String(), "guilds", g.ID.String(), "commands")
	})
}

func (g GuildSlashCommandResolver) Bulk(cmds []discord.CreateCommand) ([]discord.Command, error) {
	return httpc.NewJSONRequest[[]discord.Command](g.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPut)
		b.Body(cmds)
		return b.Execute("applications", g.CurrentUserID.String(), "guilds", g.ID.String(), "commands")
	})
}

func (g GuildSlashCommandResolver) Command(id snowflake.ID) SlashCommandClient {
	return SlashCommandResolver{
		client:        g.client,
		ID:            id,
		Guild:         g.ID,
		CurrentUserID: g.CurrentUserID,
	}
}

type SlashCommandResolver struct {
	client        *client
	ID            snowflake.ID
	Guild         snowflake.ID
	CurrentUserID snowflake.ID
}

func (c SlashCommandResolver) Get() (discord.Command, error) {
	return httpc.NewJSONRequest[discord.Command](c.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("applications", c.CurrentUserID.String(), "commands", c.ID.String())
	})
}

func (c SlashCommandResolver) Delete() error {
	return httpc.NewRequest(c.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		return b.Execute("applications", c.CurrentUserID.String(), "commands", c.ID.String())
	})
}

func (c SlashCommandResolver) Modify(params discord.CreateCommand) (discord.Command, error) {
	return httpc.NewJSONRequest[discord.Command](c.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPatch)
		b.Body(params)
		return b.Execute("applications", c.CurrentUserID.String(), "commands", c.ID.String())
	})
}

func (c SlashCommandResolver) ModifyPermissions(params discord.CommandPermissions) (discord.CommandPermissions, error) {
	return httpc.NewJSONRequest[discord.CommandPermissions](c.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPatch)
		b.Body(params)
		return b.Execute("applications", c.CurrentUserID.String(), "guilds", c.Guild.String(), "commands", c.ID.String(), "permissions")
	})
}

func (c SlashCommandResolver) Permissions() (discord.CommandPermissions, error) {
	return httpc.NewJSONRequest[discord.CommandPermissions](c.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("applications", c.CurrentUserID.String(), "guilds", c.Guild.String(), "commands", c.ID.String(), "permissions")
	})
}
