package client

import (
	"github.com/BOOMfinity/bfcord/client/cache"
	"github.com/BOOMfinity/bfcord/discord"
)

type proxyImpl struct {
	store cache.Store
}

func (p proxyImpl) AddUser(users ...discord.User) {
	if p.store == nil {
		return
	}
	for _, user := range users {
		if user.ID.Valid() {
			_ = p.store.Users().Set(user.ID, user)
		}
	}
}

func (p proxyImpl) AddChannel(channels ...discord.Channel) {
	if p.store == nil {
		return
	}
	for _, ch := range channels {
		if ch.ID.Valid() {
			_ = p.store.Channels().Set(ch.ID, ch)
		}
	}
}
