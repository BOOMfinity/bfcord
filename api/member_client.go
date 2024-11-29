package api

import (
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/internal/httpc"
	"github.com/andersfylling/snowflake/v5"
	"github.com/valyala/fasthttp"
)

type MemberResolver struct {
	client *client
	Guild  snowflake.ID
	Member snowflake.ID
}

func (m MemberResolver) Get() (discord.MemberWithUser, error) {
	return httpc.NewJSONRequest[discord.MemberWithUser](m.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("guilds", m.Guild.String(), "members", m.Member.String())
	})
}

func (m MemberResolver) Modify(params ModifyGuildMemberParams, reason ...string) (discord.Member, error) {
	return httpc.NewJSONRequest[discord.Member](m.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPatch)
		b.Reason(reason...)
		b.Body(params)
		return b.Execute("guilds", m.Guild.String(), "members", m.Member.String())
	})
}

func (m MemberResolver) AddRole(id snowflake.ID, reason ...string) error {
	return httpc.NewRequest(m.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPut)
		b.Reason(reason...)
		return b.Execute("guilds", m.Guild.String(), "members", m.Member.String(), "roles", id.String())
	})
}

func (m MemberResolver) RemoveRole(id snowflake.ID, reason ...string) error {
	return httpc.NewRequest(m.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		b.Reason(reason...)
		return b.Execute("guilds", m.Guild.String(), "members", m.Member.String(), "roles", id.String())
	})
}

func (m MemberResolver) Kick(reason ...string) error {
	return httpc.NewRequest(m.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		b.Reason(reason...)
		return b.Execute("guilds", m.Guild.String(), "members", m.Member.String())
	})
}

func (m MemberResolver) CreateBan(seconds uint, reason ...string) error {
	return httpc.NewRequest(m.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPut)
		b.Reason(reason...)
		if seconds > 0 {
			b.Body(map[string]any{
				"delete_message_seconds": seconds,
			})
		}
		return b.Execute("guilds", m.Guild.String(), "bans", m.Member.String())
	})
}

func (m MemberResolver) RemoveBan(reason ...string) error {
	return httpc.NewRequest(m.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		b.Reason(reason...)
		return b.Execute("guilds", m.Guild.String(), "bans", m.Member.String())
	})
}

func (m MemberResolver) VoiceState() (discord.VoiceState, error) {
	return httpc.NewJSONRequest[discord.VoiceState](m.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("guilds", m.Guild.String(), "voice-states", m.Member.String())
	})
}

func (m MemberResolver) ModifyVoiceState(params ModifyUserVoiceStateParams) (discord.VoiceState, error) {
	return httpc.NewJSONRequest[discord.VoiceState](m.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPatch)
		b.Body(params)
		return b.Execute("guilds", m.Guild.String(), "voice-states", m.Member.String())
	})
}
