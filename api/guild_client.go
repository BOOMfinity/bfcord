package api

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/internal/httpc"
	"github.com/andersfylling/snowflake/v5"
	"github.com/valyala/fasthttp"
	"net/url"
)

type GuildResolver struct {
	client *client
	ID     snowflake.ID
}

func (u GuildResolver) UpdateChannelPositions(positions []GuildChannelPosition, reason ...string) error {
	return httpc.NewRequest(u.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPatch)
		b.Reason(reason...)
		b.Body(positions)
		return b.Execute("guilds", u.ID.String(), "channels")
	})
}

func (u GuildResolver) UpdateRolePositions(positions []GuildRolePosition, reason ...string) error {
	return httpc.NewRequest(u.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPatch)
		b.Reason(reason...)
		b.Body(positions)
		return b.Execute("guilds", u.ID.String(), "roles")
	})
}

func (u GuildResolver) CreateChannel(params GuildChannelParams, reason ...string) (discord.Channel, error) {
	return httpc.NewJSONRequest[discord.Channel](u.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPost)
		b.Reason(reason...)
		b.Body(params)
		return b.Execute("guilds", u.ID.String(), "channels")
	})
}

func (u GuildResolver) Get() (discord.Guild, error) {
	return httpc.NewJSONRequest[discord.Guild](u.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("guilds", u.ID.String())
	})
}

func (g GuildResolver) Modify(params ModifyGuildParams, reason ...string) (discord.Guild, error) {
	return httpc.NewJSONRequest[discord.Guild](g.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPatch)
		b.Reason(reason...)
		b.Body(params)
		return b.Execute("guilds", g.ID.String())
	})
}

func (g GuildResolver) Delete(reason ...string) error {
	return httpc.NewRequest(g.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		b.Reason(reason...)
		return b.Execute("guilds", g.ID.String())
	})
}

func (g GuildResolver) Channels() ([]discord.Channel, error) {
	return httpc.NewJSONRequest[[]discord.Channel](g.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("guilds", g.ID.String(), "channels")
	})
}

func (g GuildResolver) ActiveThreads() ([]discord.Channel, error) {
	return httpc.NewJSONRequest[[]discord.Channel](g.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("guilds", g.ID.String(), "threads", "active")
	})
}

func (g GuildResolver) Member(id snowflake.ID) MemberClient {
	return MemberResolver{
		client: g.client,
		Guild:  g.ID,
		Member: id,
	}
}

func (g GuildResolver) Members(params GuildMembersParams) ([]discord.Member, error) {
	return httpc.NewJSONRequest[[]discord.Member](g.client.http, func(b httpc.RequestBuilder) error {
		values := url.Values{}
		if params.Limit > 0 {
			values.Set("limit", fmt.Sprint(params.Limit))
		}
		if params.After.Valid() {
			values.Set("after", params.After.String())
		}
		return b.Execute("guilds", g.ID.String(), "members?"+values.Encode())
	})
}

func (g GuildResolver) Search(query string, limit ...uint) ([]discord.Member, error) {
	return httpc.NewJSONRequest[[]discord.Member](g.client.http, func(b httpc.RequestBuilder) error {
		values := url.Values{}
		if len(limit) > 0 {
			values.Set("limit", fmt.Sprint(limit[0]))
		}
		values.Set("query", query)
		return b.Execute("guilds", g.ID.String(), "members", "search?"+values.Encode())
	})
}

func (g GuildResolver) ModifyCurrentMember(nick string, reason ...string) (discord.Member, error) {
	return httpc.NewJSONRequest[discord.Member](g.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPatch)
		b.Body(map[string]any{
			"nick": nick,
		})
		b.Reason(reason...)
		return b.Execute("guilds", g.ID.String(), "members", "@me")
	})
}

func (g GuildResolver) BulkBan(ids []snowflake.ID, seconds uint, reason ...string) (GuildBanAddResponse, error) {
	return httpc.NewJSONRequest[GuildBanAddResponse](g.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPost)
		b.Body(map[string]any{
			"user_ids":               ids,
			"delete_message_seconds": seconds,
		})
		b.Reason(reason...)
		return b.Execute("guilds", g.ID.String(), "bulk-ban")
	})
}

func (g GuildResolver) Roles() ([]discord.Role, error) {
	return httpc.NewJSONRequest[[]discord.Role](g.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("guilds", g.ID.String(), "roles")
	})
}

func (g GuildResolver) Role(id snowflake.ID) RoleClient {
	return RoleResolver{
		client: g.client,
		Guild:  g.ID,
		ID:     id,
	}
}

func (g GuildResolver) CreateRole(params CreateRoleParams, reason ...string) (discord.Role, error) {
	return httpc.NewJSONRequest[discord.Role](g.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPost)
		b.Body(params)
		b.Reason(reason...)
		return b.Execute("guilds", g.ID.String(), "roles")
	})
}

func (g GuildResolver) CurrentUserVoiceState() (discord.VoiceState, error) {
	return httpc.NewJSONRequest[discord.VoiceState](g.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("guilds", g.ID.String(), "voice-states", "@me")
	})
}

func (g GuildResolver) ModifyCurrentUserVoiceState(params ModifyCurrentUserVoiceStateParams) (discord.VoiceState, error) {
	return httpc.NewJSONRequest[discord.VoiceState](g.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPatch)
		b.Body(params)
		return b.Execute("guilds", g.ID.String(), "voice-states", "@me")
	})
}

func (g GuildResolver) Emojis() ([]discord.Emoji, error) {
	return httpc.NewJSONRequest[[]discord.Emoji](g.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("guilds", g.ID.String(), "emojis")
	})
}

func (g GuildResolver) Emoji(id snowflake.ID) EmojiClient {
	return EmojiResolver{
		client: g.client,
		Guild:  g.ID,
		Emoji:  id,
	}
}

func (g GuildResolver) CreateEmoji(params CreateEmojiParams, reason ...string) (discord.Emoji, error) {
	return httpc.NewJSONRequest[discord.Emoji](g.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPost)
		b.Reason(reason...)
		b.Body(params)
		return b.Execute("guilds", g.ID.String(), "emojis")
	})
}

func (g GuildResolver) Events() ([]discord.ScheduledEvent, error) {
	return httpc.NewJSONRequest[[]discord.ScheduledEvent](g.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("guilds", g.ID.String(), "scheduled-events")
	})
}

func (g GuildResolver) CreateEvent(params CreateScheduledEventParams, reason ...string) (discord.ScheduledEvent, error) {
	return httpc.NewJSONRequest[discord.ScheduledEvent](g.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPost)
		b.Reason(reason...)
		b.Body(params)
		return b.Execute("guilds", g.ID.String(), "scheduled-events")
	})
}

func (g GuildResolver) Event(id snowflake.ID) GuildEventClient {
	return GuildEventResolver{
		client: g.client,
		Guild:  g.ID,
		Event:  id,
	}
}
