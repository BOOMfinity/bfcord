package api

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/api/builders"
	"github.com/BOOMfinity/bfcord/errs"

	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/discord/permissions"
	"github.com/andersfylling/snowflake/v5"
	"github.com/segmentio/encoding/json"
	"github.com/valyala/fasthttp"
)

var _ = (discord.GuildMemberQuery)(&MemberQuery{})

type MemberQuery struct {
	api *Client
	emptyOptions[discord.GuildMemberQuery]
	member snowflake.ID
	guild  snowflake.ID
}

func (v MemberQuery) VoiceState() (state *discord.VoiceState, err error) {
	// TODO: JEBANY DISCORD
	return nil, errs.HTTPNotFound
}

func (v MemberQuery) Get() (member *discord.MemberWithUser, err error) {
	req := v.api.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/members/%v", FullApiUrl, v.guild, v.member))
	err = v.api.DoResult(req, &member)
	if err != nil {
		return
	}
	member.UserID = v.member
	member.GuildID = v.guild
	return
}

func (v MemberQuery) Ban(days uint8) (err error) {
	req := v.api.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/bans/%v", FullApiUrl, v.guild, v.member))
	raw, err := json.Marshal(map[string]uint8{
		"delete_message_days": days,
	})
	if err != nil {
		return
	}
	req.SetBody(raw)
	req.Header.SetMethod(fasthttp.MethodPut)
	if v.reason != "" {
		req.Header.Set("X-Audit-Log-Reason", v.reason)
	}
	return v.api.DoNoResp(req)
}

func (v MemberQuery) Unban() (err error) {
	req := v.api.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/bans/%v", FullApiUrl, v.guild, v.member))
	req.Header.SetMethod(fasthttp.MethodDelete)
	if v.reason != "" {
		req.Header.Set("X-Audit-Log-Reason", v.reason)
	}
	return v.api.DoNoResp(req)
}

func (v MemberQuery) AddRole(role snowflake.ID) (err error) {
	req := v.api.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/members/%v/roles/%v", FullApiUrl, v.guild, v.member, role))
	req.Header.SetMethod(fasthttp.MethodPut)
	if v.reason != "" {
		req.Header.Set("X-Audit-Log-Reason", v.reason)
	}
	return v.api.DoNoResp(req)
}

func (v MemberQuery) RemoveRole(role snowflake.ID) (err error) {
	req := v.api.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/members/%v/roles/%v", FullApiUrl, v.guild, v.member, role))
	req.Header.SetMethod(fasthttp.MethodDelete)
	if v.reason != "" {
		req.Header.Set("X-Audit-Log-Reason", v.reason)
	}
	return v.api.DoNoResp(req)
}

func (v MemberQuery) Kick() (err error) {
	req := v.api.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/members/%v", FullApiUrl, v.guild, v.member))
	req.Header.SetMethod(fasthttp.MethodDelete)
	if v.reason != "" {
		req.Header.Set("X-Audit-Log-Reason", v.reason)
	}
	return v.api.DoNoResp(req)
}

func (v MemberQuery) Edit() discord.UpdateGuildMemberBuilder {
	return builders.NewUpdateGuildMemberBuilder(v.guild, v.member)
}

func (v MemberQuery) Permissions() (perm permissions.Permission, err error) {
	guild, err := v.api.Guild(v.guild).Get()
	if err != nil {
		return
	}
	member, err := v.api.Guild(v.guild).Member(v.member).Get()
	if err != nil {
		return
	}
	return discord.BasePermissions(member.Member, *guild), nil
}

func (v MemberQuery) PermissionsIn(channel snowflake.ID) (perm permissions.Permission, err error) {
	guild, err := v.api.Guild(v.guild).Get()
	if err != nil {
		return
	}
	member, err := v.api.Guild(v.guild).Member(v.member).Get()
	if err != nil {
		return
	}
	ch, err := v.api.Channel(channel).Get()
	if err != nil {
		return
	}
	return discord.ChannelPermissions(*guild, member.Member, ch.Overwrites), nil
}

func (v MemberQuery) ID() snowflake.ID {
	return v.member
}

func (v MemberQuery) GuildID() snowflake.ID {
	return v.guild
}

func NewMemberQuery(client *Client, guild snowflake.ID, id snowflake.ID) *MemberQuery {
	data := &MemberQuery{
		guild:  guild,
		member: id,
		api:    client,
	}
	data.emptyOptions = emptyOptions[discord.GuildMemberQuery]{data: data}
	return data
}
