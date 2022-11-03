package client

import (
	"github.com/BOOMfinity/bfcord/api"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/errs"
)

type resolverOptions[V any] struct {
	data        V
	reason      string
	ignoreCache bool
	ignoreAPI   bool
}

func (r *resolverOptions[V]) NoCache() V {
	r.ignoreCache = true
	return r.data
}

func (r *resolverOptions[V]) NoAPI() V {
	r.ignoreAPI = true
	return r.data
}

func (r *resolverOptions[V]) Reason(str string) V {
	r.reason = str
	return r.data
}

type userResolver struct {
	*api.UserQuery
	bot *client
	resolverOptions[discord.UserQuery]
}

func (u userResolver) Get() (user *discord.User, err error) {
	if !u.ignoreCache && u.bot.Store() != nil {
		cache, ok := u.bot.Store().Users().Get(u.ID())
		if ok {
			return &cache, nil
		}
	}
	if !u.ignoreAPI {
		user, err = u.bot.API().User(u.ID()).Get()
		if err != nil {
			return
		}
		if u.bot.Store() != nil {
			u.bot.Store().Users().Set(user.ID, *user)
		}
		return
	}
	return nil, errs.ItemNotFound
}

func (u userResolver) Send() (msg discord.CreateMessageBuilder, err error) {
	dm, err := u.CreateDM()
	if err != nil {
		return nil, err
	}
	return u.bot.API().LowLevel().SendDM(dm.ID), nil
}

func (u userResolver) CreateDM() (ch *discord.Channel, err error) {
	if !u.ignoreCache && u.bot.Store() != nil {
		ch, ok := u.bot.Store().Private().Get(u.ID())
		if ok {
			return &ch, nil
		}
	}
	ch, err = u.UserQuery.CreateDM()
	if err == nil && !u.ignoreCache && u.bot.Store() != nil {
		u.bot.Store().Private().Set(u.ID(), *ch)
	}
	return
}
