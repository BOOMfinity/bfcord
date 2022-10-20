package builders

import (
	"github.com/BOOMfinity/bfcord/api/images"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/discord/permissions"
	"github.com/andersfylling/snowflake/v5"
	"strings"
)

var _ = (discord.RoleBuilder)(&roleBuilder{})

type roleBuilder struct {
	guild snowflake.ID
	id    snowflake.ID
	data  discord.RoleCreate
}

func (r *roleBuilder) Name(str string) discord.RoleBuilder {
	r.data.Name = &str
	return r
}

func (r *roleBuilder) Permissions(perms permissions.Permission) discord.RoleBuilder {
	r.data.Permissions = &perms
	return r
}

func (r *roleBuilder) Color(c int64) discord.RoleBuilder {
	r.data.Color = &c
	return r
}

func (r *roleBuilder) ShowSeparately(b bool) discord.RoleBuilder {
	r.data.Hoist = &b
	return r
}

func (r *roleBuilder) Icon(i *images.MediaBuilder) discord.RoleBuilder {
	if i.Nil() {
		return r
	}
	data, err := i.ToBase64()
	if err != nil {
		panic(err)
	}
	r.data.Icon = &data
	return r
}

func (r *roleBuilder) UnicodeEmoji(str string) discord.RoleBuilder {
	r.data.UnicodeEmoji = &str
	return r
}

func (r *roleBuilder) Mentionable(b bool) discord.RoleBuilder {
	r.data.Mentionable = &b
	return r
}

func (r *roleBuilder) Execute(api discord.ClientQuery, reasons ...string) (role discord.Role, err error) {
	return api.LowLevel().Reason(strings.Join(reasons, " ")).CreateOrUpdate(r.guild, r.id, r.data)
}

func NewRoleBuilder(guild, id snowflake.ID, data discord.RoleCreate) *roleBuilder {
	return &roleBuilder{guild, id, data}
}
