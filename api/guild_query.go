package api

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/api/builders"
	"github.com/BOOMfinity/go-utils/inlineif"
	"net/url"

	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/internal/slices"
	"github.com/andersfylling/snowflake/v5"
	"github.com/segmentio/encoding/json"
	"github.com/valyala/fasthttp"
)

var _ = (discord.GuildQuery)(&GuildQuery{})

type GuildQuery struct {
	api *Client
	emptyOptions[discord.GuildQuery]
	guild snowflake.ID
}

func (v GuildQuery) UpdateRolePositions(roles discord.RolePositions) error {
	req := v.api.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/roles", FullApiUrl, v.guild))
	req.Header.SetMethod(fasthttp.MethodPatch)
	json, err := json.Marshal(roles.Map())
	if err != nil {
		return fmt.Errorf("failed to marshal json data: %w", err)
	}
	req.SetBody(json)
	return v.api.DoNoResp(req)
}

func (v GuildQuery) Role(id snowflake.ID) discord.RoleQuery {
	return NewRoleQuery(v.api, v.guild, id)
}

func (v GuildQuery) CreateRole() discord.RoleBuilder {
	return builders.NewRoleBuilder(v.guild, 0, discord.RoleCreate{})
}

func (v GuildQuery) VoiceStates() (states discord.Slice[discord.VoiceState], err error) {
	// TODO: Discord does not send voice states anywhere except guild_create event
	return
}

func (v GuildQuery) Edit() discord.GuildBuilder {
	return builders.NewGuildBuilder(v.guild)
}

func (v GuildQuery) Members(limit int, after snowflake.ID) (members []discord.MemberWithUser, err error) {
	var _limit uint16
	_members := make([]discord.MemberWithUser, 0, inlineif.IfElse(limit == -1, 1000, limit))
	for {
		if limit > 1000 {
			_limit = 1000
		}
		if limit == -1 {
			_limit = 1000
		} else if (len(members) + 1000) > limit {
			_limit = uint16(limit - len(members))
		}
		params := url.Values{}
		if limit > 0 {
			params.Set("limit", fmt.Sprint(_limit))
		}
		if after.Valid() {
			params.Set("after", after.String())
		}
		req := v.api.New(true)
		if len(params) != 0 {
			req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/members?%v", FullApiUrl, v.guild, params.Encode()))
		} else {
			req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/members", FullApiUrl, v.guild))
		}
		err = v.api.DoResult(req, &members)
		for i := range members {
			members[i].GuildID = v.ID()
			members[i].UserID = members[i].User.ID
		}
		if len(members) > 0 {
			after = members[len(members)-1].UserID
		}
		members = append(members, _members...)
		if len(members) < 1000 || len(members) == limit {
			break
		}
	}
	return
}

func (v GuildQuery) Bans(limit int, after snowflake.ID) (bans []discord.Ban, err error) {
	var _limit uint16
	var _bans []discord.Ban
	for {
		if limit > 1000 {
			_limit = 1000
		}
		if limit == -1 {
			_limit = 1000
		} else if (len(bans) + 1000) > limit {
			_limit = uint16(limit - len(bans))
		}
		params := url.Values{}
		if limit > 0 {
			params.Set("limit", fmt.Sprint(_limit))
		}
		if after.Valid() {
			params.Set("after", after.String())
		}
		req := v.api.New(true)
		if len(params) != 0 {
			req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/bans?%v", FullApiUrl, v.guild, params.Encode()))
		} else {
			req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/bans", FullApiUrl, v.guild))
		}
		err = v.api.DoResult(req, &bans)
		if len(bans) > 0 {
			after = bans[len(bans)-1].User.ID
		}
		bans = append(bans, _bans...)
		if len(bans) < 1000 || len(bans) == limit {
			break
		}
	}
	return
}

func (v GuildQuery) Get() (guild *discord.Guild, err error) {
	req := v.api.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v", FullApiUrl, v.guild))
	err = v.api.DoResult(req, &guild)
	if err != nil {
		return
	}
	guild.Patch()
	return
}

func (v GuildQuery) Delete() error {
	req := v.api.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v", FullApiUrl, v.guild))
	req.Header.SetMethod(fasthttp.MethodDelete)
	return v.api.DoNoResp(req)
}

func (v GuildQuery) Channels() (channels []discord.Channel, err error) {
	req := v.api.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/channels", FullApiUrl, v.guild))
	return channels, v.api.DoResult(req, &channels)
}

func (v GuildQuery) CreateChannel(name string) discord.GuildChannelBuilder {
	return builders.NewCreateChannelBuilder(v.guild, name)
}

func (v GuildQuery) UpdateChannelPositions(positions *discord.GuildChannelPositionsBuilder) error {
	req := v.api.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/channels", FullApiUrl, v.guild))
	req.Header.SetMethod(fasthttp.MethodPatch)
	data, err := json.Marshal(positions.Encode())
	if err != nil {
		return err
	}
	req.SetBody(data)
	return v.api.DoNoResp(req)
}

type activeThreadsData struct {
	Threads []discord.Channel      `json:"threads"`
	Members []discord.ThreadMember `json:"members"`
}

func (v GuildQuery) ActiveThreads() (threads []discord.Channel, err error) {
	req := v.api.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/threads/active", FullApiUrl, v.guild))
	var data activeThreadsData
	err = v.api.DoResult(req, &data)
	if err != nil {
		return nil, err
	}
	threads = data.Threads
	for i := range data.Members {
		data.Members[i].GuildID = v.guild
	}
	for i := range threads {
		index := slices.FindIndex(data.Members, func(item discord.ThreadMember) bool {
			return item.ID == threads[i].ID
		})
		if index == -1 {
			continue
		}
		threads[i].Member = &data.Members[index]
	}
	return
}

func (v GuildQuery) Member(id snowflake.ID) discord.GuildMemberQuery {
	return NewMemberQuery(v.api, v.guild, id)
}

func (v GuildQuery) Search(query string, limit uint16) (members []discord.MemberWithUser, err error) {
	req := v.api.New(true)
	params := url.Values{}
	params.Set("query", query)
	params.Set("limit", fmt.Sprint(limit))
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/members/search?%v", FullApiUrl, v.guild, params.Encode()))
	err = v.api.DoResult(req, &members)
	if err != nil {
		return nil, err
	}
	for i := range members {
		members[i].GuildID = v.guild
		members[i].UserID = members[i].User.ID
	}
	return
}

func (v GuildQuery) Roles() (roles discord.RoleSlice, err error) {
	req := v.api.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/roles", FullApiUrl, v.guild))
	return roles, v.api.DoResult(req, &roles)
}

func (v GuildQuery) SetCurrentNick(nick string) (err error) {
	req := v.api.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/members/@me", FullApiUrl, v.guild))
	req.Header.SetMethod(fasthttp.MethodPatch)
	data, err := json.Marshal(map[string]any{
		"nick": fastIf(nick != "", &nick, nil),
	})
	if err != nil {
		return err
	}
	req.SetBody(data)
	if v.reason != "" {
		req.Header.Set("X-Audit-Log-Reason", v.reason)
	}
	return v.api.DoNoResp(req)
}

func (v GuildQuery) Invites() (invites []discord.InviteWithMeta, err error) {
	req := v.api.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/guilds/%v/invites", FullApiUrl, v.guild))
	err = v.api.DoResult(req, &invites)
	return
}

func (v GuildQuery) ID() snowflake.ID {
	return v.guild
}

func NewGuildQuery(client *Client, id snowflake.ID) *GuildQuery {
	d := &GuildQuery{
		guild: id,
		api:   client,
	}
	d.emptyOptions = emptyOptions[discord.GuildQuery]{data: d}
	return d
}

func fastIf[T any](cond bool, ifTrue *T, ifFalse *T) *T {
	if cond {
		return ifTrue
	} else {
		return ifFalse
	}
}
