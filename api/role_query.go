package api

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/api/builders"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
	"github.com/valyala/fasthttp"
)

var _ = (discord.RoleQuery)(&RoleQuery{})

type RoleQuery struct {
	api *Client
	emptyOptions[discord.RoleQuery]
	guild snowflake.ID
	role  snowflake.ID
}

func (r RoleQuery) Edit() discord.RoleBuilder {
	return builders.NewRoleBuilder(r.guild, r.role, discord.RoleCreate{})
}

func (r RoleQuery) Delete() error {
	req := r.api.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/roles/%v", FullApiUrl, r.guild, r.role))
	req.Header.SetMethod(fasthttp.MethodDelete)
	return r.api.DoNoResp(req)
}

func NewRoleQuery(client *Client, guild snowflake.ID, id snowflake.ID) *RoleQuery {
	data := &RoleQuery{
		guild: guild,
		role:  id,
		api:   client,
	}
	data.emptyOptions = emptyOptions[discord.RoleQuery]{data: data}
	return data
}
