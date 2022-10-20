package client

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/api"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/errs"
)

var _ = (discord.GuildMemberQuery)(&memberResolver{})

type memberResolver struct {
	bot *client
	*api.MemberQuery
	resolverOptions[discord.GuildMemberQuery]
}

func (mr memberResolver) VoiceState() (state discord.VoiceState, err error) {
	states, _err := mr.bot.Guild(mr.GuildID()).NoAPI().VoiceStates()
	if _err != nil && !errs.IsNotFound(_err) {
		err = fmt.Errorf("could not fetch voice states from cache: %w", _err)
		return
	}
	_state, found := states.Find(func(item discord.VoiceState) bool {
		return item.UserID == mr.ID()
	})
	if found {
		return _state, nil
	}
	return discord.VoiceState{}, errs.ItemNotFound
}

func (mr memberResolver) Get() (member discord.MemberWithUser, err error) {
	if !mr.ignoreCache && mr.bot.Store() != nil {
		m, ok := mr.bot.Store().Members().UnsafeGet(mr.GuildID()).Get(mr.ID())
		if ok {
			member.Member = m
			user, _err := mr.bot.User(mr.ID()).Get()
			if _err != nil {
				err = fmt.Errorf("could not fetch user: %w", err)
				return
			}
			member.User = user
			return
		}
	}
	if !mr.ignoreAPI {
		member, err = mr.MemberQuery.Get()
		if err != nil {
			return
		}
		if mr.bot.Store() != nil {
			mr.bot.Store().Members().UnsafeGet(mr.GuildID()).Set(member.UserID, member.Member)
		}
		return
	}
	return discord.MemberWithUser{}, errs.ItemNotFound
}
