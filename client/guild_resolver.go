package client

import (
	"github.com/BOOMfinity/bfcord/api"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/errs"
	"github.com/andersfylling/snowflake/v5"
)

var _ = (discord.GuildQuery)(&guildResolver{})

type guildResolver struct {
	bot *client
	*api.GuildQuery
	resolverOptions[discord.GuildQuery]
}

func (gr guildResolver) VoiceStates() (states discord.Slice[discord.VoiceState], err error) {
	if gr.bot.Store() != nil && !gr.ignoreCache {
		guild, found := gr.bot.Store().VoiceStates().Get(gr.ID())
		if found {
			states = guild.ToSlice()
			return
		}
	}
	return nil, errs.ItemNotFound
}

func (gr guildResolver) Member(id snowflake.ID) discord.GuildMemberQuery {
	resolver := &memberResolver{MemberQuery: api.NewMemberQuery(gr.bot.API(), gr.ID(), id), bot: gr.bot}
	resolver.resolverOptions = resolverOptions[discord.GuildMemberQuery]{data: resolver}
	return resolver
}

func (gr guildResolver) Members(limit int, after snowflake.ID) (members []discord.MemberWithUser, err error) {
	members, err = gr.bot.Guild(gr.ID()).Members(limit, after)
	if err != nil {
		return nil, err
	}
	if gr.bot.Store() != nil {
		for i := range members {
			gr.bot.Store().Members().UnsafeGet(gr.ID()).Set(members[i].UserID, members[i].Member)
			gr.bot.Store().Users().Set(members[i].UserID, members[i].User)
		}
	}
	return
}

func (gr guildResolver) Get() (*discord.Guild, error) {
	if !gr.ignoreCache && gr.bot.Store() != nil {
		g, ok := gr.bot.Store().Guilds().Get(gr.ID())
		if ok {
			return &g, nil
		}
	}
	if !gr.ignoreAPI {
		guild, err := gr.GuildQuery.Get()
		if err != nil {
			return nil, err
		}
		if gr.bot.Store() != nil {
			gr.bot.Store().Guilds().Set(guild.ID, *guild)
		}
		return guild, nil
	}
	return nil, errs.ItemNotFound
}
