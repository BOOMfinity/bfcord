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
	if !gr.ignoreAPI {
		states, err = gr.GuildQuery.VoiceStates()
		if gr.bot.Store() != nil {
			for i := range states {
				gr.bot.Store().VoiceStates().UnsafeGet(gr.ID()).Set(states[i].UserID, states[i])
			}
		}
		return
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
			gr.bot.Log().Debug().Any(gr.ID()).Send("Guild was returned from cache")
			return &discord.Guild{
				BaseGuild: g,
			}, nil
		}
	}
	if !gr.ignoreAPI {
		guild, err := gr.GuildQuery.Get()
		if err != nil {
			return nil, err
		}
		if gr.bot.Store() != nil {
			gr.bot.Store().Guilds().Set(guild.ID, guild.BaseGuild)
			for i := range guild.Members {
				gr.bot.Store().Members().UnsafeGet(guild.ID).Set(guild.Members[i].UserID, guild.Members[i].Member)
				gr.bot.Store().Users().Set(guild.Members[i].UserID, guild.Members[i].User)
			}
			for i := range guild.Channels {
				gr.bot.Store().Channels().UnsafeGet(guild.ID).Set(guild.Channels[i].ID, guild.Channels[i])
			}
			for i := range guild.Presences {
				gr.bot.Store().Presences().UnsafeGet(guild.ID).Set(guild.Presences[i].UserID, guild.Presences[i].BasePresence)
			}
			gr.bot.Store().Users().Set(guild.OwnerID, guild.Owner)
		}
		gr.bot.Log().Debug().Any(gr.ID()).Send("Guild was returned from API")
		return guild, nil
	}
	return nil, errs.ItemNotFound
}
