package api

import "github.com/BOOMfinity/bfcord/discord"

type CacheProxy interface {
	AddUser(users ...discord.User)
	AddChannel(channels ...discord.Channel)
}

type noopProxy struct{}

func (n noopProxy) AddUser(users ...discord.User) {
	return
}

func (n noopProxy) AddChannel(channels ...discord.Channel) {
	return
}
